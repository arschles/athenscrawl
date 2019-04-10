package queue

type ctxCancelPair struct {
	ctx  context.context
	done func()
}
