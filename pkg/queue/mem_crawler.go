package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/arschles/crathens/pkg/resp"
	"github.com/google/go-github/github"
)

type inMemoryCrawler struct {
	ghFetchCoord    *coordinator
	athensWarmCoord *coordinator
}

// InMemory creates a new crawler implementation that works only in memory
func InMemory(
	ctx context.Context,
	ghCl *github.Client,
	ghTickDur time.Duration,
	athensTickDur time.Duration,
) Crawler {
	ghFetchCoordinator := newCoordinator(ctx, ghTickDur)
	athensWarmCoordinator := newCoordinator(ctx, athensTickDur)

	go ghFetcher(ghFetchCoordinator, ghCl, athensWarmCoordinator.ch)

	go athensWarmer(athensWarmCoordinator)
	return &inMemoryCrawler{
		ghFetchCoord:    ghFetchCoordinator,
		athensWarmCoord: athensWarmCoordinator,
	}
}

func (i *inMemoryCrawler) Enqueue(
	ctx context.Context,
	mav resp.ModuleAndVersion,
) error {
	select {
	case i.ghFetchCoord.ch <- mav:
		return nil
	case <-ctx.Done():
		return fmt.Errorf(
			"Failed to start crawling module %s due to context timeout",
			mav.Module,
		)
	}
}

func (i *inMemoryCrawler) Wait(context.Context) error {
	select {
	case <-i.athensWarmCoord.ctx.Done():
		i.stopTickers()
		i.stopContexts()
		return fmt.Errorf("Athens fetcher stopped")
	case <-i.ghFetchCoord.ctx.Done():
		i.stopTickers()
		i.stopContexts()
		return fmt.Errorf("Github fetcher stopped")
	}
	return nil
}

func (i *inMemoryCrawler) stopTickers() {
	i.ghFetchCoord.ticker.Stop()
	i.athensWarmCoord.ticker.Stop()
}

func (i *inMemoryCrawler) stopContexts() {
	i.ghFetchCoord.ctxDone()
	i.athensWarmCoord.ctxDone()
}
