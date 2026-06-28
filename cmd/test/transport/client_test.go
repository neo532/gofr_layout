package transport

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/smallnest/rpcx/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
)

const (
	retryTimes       = 2          // max additional retries
	retryDuration    = time.Microsecond  // initial backoff
	retryMaxDuration = 20 * time.Microsecond  // backoff ceiling
)

// ---------------------------------------------------------------------------
// HTTP middleware: logging + retry via http.RoundTripper chain
// ---------------------------------------------------------------------------

type loggingTripper struct {
	t    *testing.T
	next http.RoundTripper
}

func (l *loggingTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	l.t.Logf("[HTTP] %s %s", req.Method, req.URL)
	resp, err := l.next.RoundTrip(req)
	if err != nil {
		l.t.Logf("[HTTP] error: %v", err)
	}
	return resp, err
}

// retryTripper matches gokit/transport/http retry rules:
//   - 4xx (400-407): not retryable (cancelRetry)
//   - 5xx+ / network error: retryable with exponential backoff
type retryTripper struct {
	t        *testing.T
	next     http.RoundTripper
	maxRetry int
}

func (r *retryTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var bodyBuf []byte
	if req.Body != nil {
		var err error
		bodyBuf, err = io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, err
		}
	}
	backoff := retryDuration
	var lastErr error
	for i := 0; i <= r.maxRetry; i++ {
		if i > 0 {
			r.t.Logf("[HTTP retry] attempt %d/%d", i, r.maxRetry)
			req.Body = io.NopCloser(bytes.NewReader(bodyBuf))
			time.Sleep(backoff)
			if backoff < retryMaxDuration {
				backoff *= 2
			}
		}
		resp, err := r.next.RoundTrip(req)
		if err != nil {
			lastErr = err
			continue
		}
		// 4xx client errors: not retryable (same as gokit DefaultErrorDecoder)
		if resp.StatusCode >= 400 && resp.StatusCode <= 407 {
			resp.Body.Close()
			return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		if resp.StatusCode >= http.StatusInternalServerError {
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
			resp.Body.Close()
			continue
		}
		return resp, nil
	}
	return nil, lastErr
}

// ---------------------------------------------------------------------------
// gRPC middleware: logging + retry via unary interceptor chain
// ---------------------------------------------------------------------------

// grpcRetryInterceptor retries on retryable gRPC status codes:
// Unavailable, DeadlineExceeded, ResourceExhausted, etc.
// Non-retryable codes (InvalidArgument, NotFound, PermissionDenied, etc.) short-circuit.
func grpcRetryInterceptor(t *testing.T, maxRetry int) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any,
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		backoff := retryDuration
		var lastErr error
		for i := 0; i <= maxRetry; i++ {
			if i > 0 {
				t.Logf("[gRPC retry] attempt %d/%d", i, maxRetry)
				time.Sleep(backoff)
				if backoff < retryMaxDuration {
					backoff *= 2
				}
			}
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}
			if st, ok := status.FromError(err); ok {
				switch st.Code() {
				case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted,
					codes.Aborted, codes.Internal, codes.DataLoss:
					lastErr = err
				default:
					return err // non-retryable
				}
			} else {
				lastErr = err // network error, retryable
			}
		}
		return lastErr
	}
}

// ---------------------------------------------------------------------------
// rpcx middleware: logging via plugin system (Failtry mode = built-in retry)
// ---------------------------------------------------------------------------

type loggingPlugin struct{ t *testing.T }

func (p *loggingPlugin) PreCall(_ context.Context, svcPath, svcMethod string, args any) error {
	p.t.Logf("[rpcx] %s.%s args=%+v", svcPath, svcMethod, args)
	return nil
}
func (p *loggingPlugin) PostCall(_ context.Context, svcPath, svcMethod string, args, reply any) error {
	p.t.Logf("[rpcx] %s.%s done", svcPath, svcMethod)
	return nil
}

// ---------------------------------------------------------------------------
// ClientPool — shared connection pool across protocols
// ---------------------------------------------------------------------------

// ClientPool holds one generated client per protocol, each backed by a
// properly pooled connection and middleware chain.
type ClientPool struct {
	HTTP *pb.UserApiClient
	GRPC *pb.UserApiClient
	RPCX *pb.UserApiClient
	WS   *pb.UserApiClient

	// kept alive for cleanup
	clientGRPC *grpc.ClientConn
	clientHTTP *http.Client
	clientRPCX client.XClient
	clientWS   *websocket.Dialer
}

// NewClientPool creates a ClientPool with connection pooling and middleware.
func NewClientPool(t *testing.T) (clt *ClientPool) {

	clt = &ClientPool{}
	var err error

	// ── HTTP: transport-level connection pool + retry + logging ──
	clt.clientHTTP = &http.Client{
		Transport: &retryTripper{
			t: t, maxRetry: retryTimes,
			next: &loggingTripper{t: t, next: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
			}},
		},
		Timeout: 10 * time.Second,
	}
	clt.HTTP = pb.NewUserApiHTTPClient("http://127.0.0.1:8500", clt.clientHTTP)

	// ── gRPC: built-in connection pool + retry + logging ──
	if clt.clientGRPC, err = grpc.NewClient("127.0.0.1:8501",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			func(ctx context.Context, method string, req, reply any,
				cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
				opts ...grpc.CallOption) error {
				t.Logf("[gRPC] %s req=%+v", method, req)
				return invoker(ctx, method, req, reply, cc, opts...)
			},
			grpcRetryInterceptor(t, retryTimes),
		),
	); err != nil {
		t.Fatal(err)
	}
	clt.GRPC = pb.NewUserApiGRPCClient(clt.clientGRPC)

	// ── rpcx: Failtry = built-in retry + plugin middleware ──
	d, err := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8502", "")
	if err != nil {
		t.Fatal(err)
	}
	clt.clientRPCX = client.NewXClient("user.v1.UserApi", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	pc := client.NewPluginContainer()
	pc.Add(&loggingPlugin{t: t})
	clt.clientRPCX.SetPlugins(pc)
	clt.RPCX = pb.NewUserApiRPCXClient(clt.clientRPCX)

	// ── WebSocket: shared dialer (retry at application level, see TestClientWS) ──
	clt.clientWS = &websocket.Dialer{HandshakeTimeout: 5 * time.Second}
	clt.WS = pb.NewUserApiWSClient("ws://127.0.0.1:8503", clt.clientWS)

	return
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestClientHTTP(t *testing.T) {
	pool := NewClientPool(t)
	t.Run("GetById", func(t *testing.T) {
		user, err := pool.HTTP.GetById(t.Context(), &pb.GetByIdRequest{Id: 1})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("GetById: id=%d name=%s", user.Id, user.Name)
	})
	t.Run("Post", func(t *testing.T) {
		_, err := pool.HTTP.Post(t.Context(), &pb.User{Id: 2, Name: "alice"})
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Post ok")
	})
}

func TestClientGRPC(t *testing.T) {
	pool := NewClientPool(t)
	t.Run("GetById", func(t *testing.T) {
		user, err := pool.GRPC.GetById(t.Context(), &pb.GetByIdRequest{Id: 1})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("GetById: id=%d name=%s", user.Id, user.Name)
	})
	t.Run("Post", func(t *testing.T) {
		_, err := pool.GRPC.Post(t.Context(), &pb.User{Id: 2, Name: "alice"})
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Post ok")
	})
}

func TestClientRPCX(t *testing.T) {
	pool := NewClientPool(t)
	t.Run("GetById", func(t *testing.T) {
		user, err := pool.RPCX.GetById(t.Context(), &pb.GetByIdRequest{Id: 1})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("GetById: id=%d name=%s", user.Id, user.Name)
	})
	t.Run("Post", func(t *testing.T) {
		_, err := pool.RPCX.Post(t.Context(), &pb.User{Id: 2, Name: "alice"})
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Post ok")
	})
}

// TestClientWS uses WebSocket. Each call opens a new connection, so retry
// is handled at the application level rather than in the transport layer.
func TestClientWS(t *testing.T) {
	pool := NewClientPool(t)
	retry := retryTimes

	t.Run("GetById", func(t *testing.T) {
		backoff := retryDuration
		var lastErr error
		for i := 0; i <= retry; i++ {
			if i > 0 {
				t.Logf("[WS retry] attempt %d/%d", i, retry)
				time.Sleep(backoff)
				if backoff < retryMaxDuration {
					backoff *= 2
				}
			}
			user, err := pool.WS.GetById(t.Context(), &pb.GetByIdRequest{Id: 1})
			if err == nil {
				t.Logf("GetById: id=%d name=%s", user.Id, user.Name)
				return
			}
			lastErr = err
		}
		t.Fatal(lastErr)
	})
	t.Run("Post", func(t *testing.T) {
		backoff := retryDuration
		var lastErr error
		for i := 0; i <= retry; i++ {
			if i > 0 {
				t.Logf("[WS retry] attempt %d/%d", i, retry)
				time.Sleep(backoff)
				if backoff < retryMaxDuration {
					backoff *= 2
				}
			}
			_, err := pool.WS.Post(t.Context(), &pb.User{Id: 2, Name: "alice"})
			if err == nil {
				t.Log("Post ok")
				return
			}
			lastErr = err
		}
		t.Fatal(lastErr)
	})
}
