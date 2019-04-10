package queue

import "context"

type ctxCancelPair struct {
	ctx  context.Context
	done func()
}
