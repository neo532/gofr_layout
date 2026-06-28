package log

import (
	"context"

	"github.com/neo532/gofr/transport"
	"github.com/neo532/gokit/queue"
)

// ProtocolFromContext extracts the transport protocol kind from context for logging.
// Falls back to the "protocol" queue header for consumer (Kafka) paths.
func ProtocolFromContext(c context.Context) (string, any) {
	if tr, ok := transport.FromServerContext(c); ok {
		return KeyProtocol, string(tr.Kind())
	}
	if h, ok := queue.GetHeaderFromContext(c); ok {
		if v := h.Value(KeyProtocol); v != "" {
			return KeyProtocol, v
		}
	}
	return "", ""
}
