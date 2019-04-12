package ctx

import (
	"context"
	"time"

	"github.com/arschles/crathens/pkg/resp"
)

// Coordinator is a context plus utilities to coordinate a goroutine running
// on a ticker and sending resp.ModuleAndVersions somewhere
type Coordinator interface {
	context.Context
	Ticker() *time.Ticker
	Ch() chan resp.ModuleAndVersion
	StopCtx()
	DoneCh() <-chan struct{}
}

type coordinator struct {
	context.Context
	ticker  *time.Ticker
	ch      chan resp.ModuleAndVersion
	ctxDone func()
}

func (c *coordinator) Ticker() *time.Ticker {
	return c.ticker
}

func (c *coordinator) Ch() chan resp.ModuleAndVersion {
	return c.ch
}

func (c *coordinator) StopCtx() {
	c.ctxDone()
}

func (c *coordinator) DoneCh() <-chan struct{} {
	return c.Done()
}

// CoordinatorFromCtx creates a new coordinator from a given context
func CoordinatorFromCtx(ctx context.Context, dur time.Duration) *coordinator {
	ctx, done := context.WithCancel(ctx)
	return &coordinator{
		ticker:  time.NewTicker(dur),
		ch:      make(chan resp.ModuleAndVersion),
		Context: ctx,
		ctxDone: done,
	}
}
