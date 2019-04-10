package queue

import (
	"context"
	"fmt"
)

type ctxCancelPair struct {
	ctx  context.context
	done func()
}

func send(ctx context.Context, ch chan<- string, val string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed sending value %s due to context cancel", val)
	case ch <- val:
		return nil
	}
}
