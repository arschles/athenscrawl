package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/arschles/crathens/pkg/resp"
	"github.com/google/go-github/github"
)

type inMemoryCrawler struct {
	ghFetchTicker    *time.Ticker
	ghFetchCh        <-chan string
	ghFetch          *ctxCancelPair
	athensWarmTicker *time.Ticker
	athensWarmCh     <-chan string
	athensWarm       *ctxCancelPair
}

// InMemory creates a new crawler implementation that works only in memory
func InMemory(
	ctx context.Context,
	ghCl *github.Client,
	ghTickDur time.Duration,
	athensTickDur time.Duration,
) Crawler {
	ghFetchTicker := time.NewTicker(ghTickDur)
	ghFetchCh := make(chan string)
	ghFetchCtx, ghFetchCtxDone := context.WithCancel(ctx)
	go ghFetcher(gitFetchCtx, ghCl, modCh, athensWarmCh, ghFetchTicker)

	athensWarmTicker := time.NewTicker(athensTickDr)
	athensWarmCh := make(chan resp.ModuleAndVersion)
	athensWarmCtx, athensWarmCtxdone := context.WithCancel(ctx)
	go athensWarmer(athensWarmCtx, athensWarmCh, athensWarmTicker)
	return &inMemoryCrawler{
		ghTicker:  ghTicker,
		ghFetchCh: ghFetchCh,
		ghFetch: &ctxCancelPair{
			ctx:  ghFetchCtx,
			done: ghFetchCtxDone,
		},
		athensWarmTicker: athensWarmTicker,
		athensWarmCh:     athensWarmCh,
		athensWarm: &ctxCancelPair{
			ctx:  athensFetchCtx,
			done: athensFechCtxDone,
		},
	}
}

func (i *inMemoryCrawler) Crawl(
	ctx context.Context,
	mav resp.ModuleAndVersion,
) error {
	select {
	case i.githubFetchCh <- mav.Module:
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
	case <-i.athensWarm.ctx.Done():
		i.stopTickers()
		i.stopContexts()
		return fmt.Errorf("Athens fetcher stopped")
	case <-i.githubFetch.ctx.Done():
		i.stopTickers()
		i.stopContexts()
		return fmt.Errorf("Github fetcher stopped")
	}
	return nil
}

func (i *inMemoryCrawler) stopTickers() {
	i.ghFetchTicker.Stop()
	i.athensWarmTicker.Stop()
}

func (i *inMemoryCrawler) stopContextx() {
	i.ghFetch.done()
	i.athensWarm.done()
}
