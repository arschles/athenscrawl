package queue

import (
	"context"
	"fmt"
	"time"

	pkgctx "github.com/arschles/crathens/pkg/ctx"
	"github.com/arschles/crathens/pkg/log"
	"github.com/arschles/crathens/pkg/resp"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type inMemoryCrawler struct {
	ghFetchCoord    pkgctx.Coordinator
	athensWarmCoord pkgctx.Coordinator
}

// InMemory creates a new crawler implementation that works only in memory
func InMemory(
	ctx context.Context,
	endpoint string,
	ghCl *github.Client,
	ghTickDur time.Duration,
	athensTickDur time.Duration,
) Crawler {
	ghFetchCoordinator := pkgctx.CoordinatorFromCtx(ctx, ghTickDur)
	athensWarmCoordinator := pkgctx.CoordinatorFromCtx(ctx, athensTickDur)

	go ghFetcher(ghFetchCoordinator, ghCl, athensWarmCoordinator.Ch())

	go athensWarmer(endpoint, athensWarmCoordinator)

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
	case i.ghFetchCoord.Ch() <- mav:
		log.Debug("enqueued %s onto the in-memory crawler", mav)
		return nil
	case <-ctx.Done():
		return errors.WithStack(fmt.Errorf(
			"Failed to start crawling module %s due to context timeout",
			mav.Module,
		))
	}
}

func (i *inMemoryCrawler) Wait(context.Context) error {
	select {
	case <-i.athensWarmCoord.DoneCh():
		log.Debug(
			"The Athens warmer stopped, cleaning up tickers/contexts and error-ing",
		)
		i.stopTickers()
		i.stopContexts()
		return errors.WithStack(fmt.Errorf("Athens fetcher stopped"))
	case <-i.ghFetchCoord.DoneCh():
		log.Debug(
			"The GitHub fetcher stopped, cleaning up tickers/contexts and error-ing",
		)
		i.stopTickers()
		i.stopContexts()
		return errors.WithStack(fmt.Errorf("Github fetcher stopped"))
	}
}

func (i *inMemoryCrawler) stopTickers() {
	i.ghFetchCoord.Ticker().Stop()
	i.athensWarmCoord.Ticker().Stop()
}

func (i *inMemoryCrawler) stopContexts() {
	i.ghFetchCoord.StopCtx()
	i.athensWarmCoord.StopCtx()
}
