package dispatcher

import (
	"context"
	"time"

	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
)

type retryHandler struct {
	inner    pubsub.IPubSubMessageHandler
	attempts int
	delay    time.Duration
}

func NewRetryHandler(inner pubsub.IPubSubMessageHandler, attempts int, delay time.Duration) pubsub.IPubSubMessageHandler {
	return &retryHandler{inner: inner, attempts: attempts, delay: delay}
}

func (r *retryHandler) Process(ctx context.Context, msg *pubsub.Message) error {
	var err error
	for i := 1; i <= r.attempts; i++ {
		err = r.inner.Process(ctx, msg)
		if err == nil {
			return nil
		}
		log.DefaultLogger().Warnf("retry %d/%d for webhook dispatch: %v", i, r.attempts, err)
		time.Sleep(r.delay)
	}
	return err
}
