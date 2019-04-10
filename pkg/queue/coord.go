package queue

import (
	"context"
	"time"

	"github.com/arschles/crathens/pkg/resp"
)

type coordinator struct {
	ticker  *time.Ticker
	ch      chan resp.ModuleAndVersion
	ctx     context.Context
	ctxDone func()
}

func newCoordinator(ctx context.Context, dur time.Duration) *coordinator {
	ctx, done := context.WithCancel(ctx)
	return &coordinator{
		ticker:  time.NewTicker(dur),
		ch:      make(chan resp.ModuleAndVersion),
		ctx:     ctx,
		ctxDone: done,
	}
}
