package log

import (
	"context"
	"strconv"

	gofrTrace "github.com/neo532/gofr/middleware/trace"
	"github.com/neo532/gofr/transport"
	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/queue"
)

// KindFromContext extracts the transport kind from context for logging.
// Falls back to the "kind" queue header for consumer (Kafka) paths.
func KindFromContext(c context.Context) (string, any) {
	if tr, ok := transport.FromServerContext(c); ok {
		return KeyKind, string(tr.Kind())
	}
	if h, ok := queue.GetHeaderFromContext(c); ok {
		if v := h.Value(KeyKind); v != "" {
			return KeyKind, v
		}
	}
	return "", ""
}

// TraceIDFromContext extracts the active trace ID for logging: the OTel request
// trace ID when a span is active, falling back to the process boot trace ID so
// general logs emitted outside any request (startup, goroutines, DB pool) still
// carry an identifier.
func TraceIDFromContext(c context.Context) (string, any) {
	if id := gofrTrace.TraceID(c); id != "" {
		return KeyTraceID, id
	}
	return KeyTraceID, ""
}

// TraceparentFromContext extracts the raw inbound traceparent header of the
// current request span, so a log entry carries the caller's trace context
// (trace id + parent span id + flags). Chained across services: service A logs
// spanId=S, service B logs traceparent=...-S-...
func TraceparentFromContext(c context.Context) (string, any) {
	if tp := gofrTrace.Traceparent(c); tp != "" {
		return KeyTraceparent, tp
	}
	return KeyTraceparent, ""
}

func FileFromContext(c context.Context) (string, any) {
	file, line := logger.GetSourceByDepth(4)
	return "file", file + ":" + strconv.Itoa(line)
}
