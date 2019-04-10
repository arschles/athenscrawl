package queue

import (
	"context"
	"time"

	gh "github.com/arschles/crathens/pkg/github"
	"github.com/arschles/crathens/pkg/log"
	"github.com/arschles/crathens/pkg/resp"
	"github.com/google/go-github/github"
)

func ghFetcher(
	ctx context.Context,
	ghCl *github.Client,
	modCh <-chan string,
	nextCh chan<- resp.ModuleAndVersion,
	ticker *time.Ticker,
) {
	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		case mod := <-modCh:
			tags, err := gh.FetchTags(ctx, ghCl, mod)
			if err != nil {
				log.Warn("fetching GH tags for %s (%s)", mod, err)
			}
			for _, tag := range tags {
				modAndVer := resp.ModuleAndVersion{
					Module:  mod,
					Version: tag,
				}
				select {
				case <-ctx.Done():
					log.Warn(
						"failed sending tag %s due to context cancel",
						tag,
					)
				case nextCh <- modAndVer:
				}
			}
		}
	}
}
