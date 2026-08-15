package log

import (
	"context"

	"github.com/neo532/gofr/transport"
	"github.com/neo532/gokit/queue"
)

func ProducerHeaderSet(ph queue.ProducerHandler) queue.ProducerHandler {
	return func(c context.Context, message any) (err error) {
		if tr, ok := transport.FromServerContext(c); ok {
			c = queue.AppendHeaderToContext(c, KeyKind, string(tr.Kind()))
		}
		return ph(c, message)
	}
}
