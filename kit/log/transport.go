package log

import (
	"context"

	"github.com/neo532/gofr/transport"
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
