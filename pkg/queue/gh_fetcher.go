package queue

import (
	"context"
	"time"

	gh "github.com/arschles/crathens/pkg/github"
	"github.com/arschles/crathens/pkg/log"
	"github.com/google/go-github/github"
)

func ghFetcher(
	ctx context.Context,
	ghCl *github.Client,
	modCh <-chan string,
	nextCh chan<- github.RepositoryTag,
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
				if err := send(ctx, nextCh, tag.Name); err != nil {
					log.Warn("queueing tag %s (%s)", tag.Name, err)
				}
			}
		}
	}
}
