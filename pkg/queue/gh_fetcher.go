package queue

import (
	"github.com/arschles/crathens/pkg/log"
)
func ghFetcher(
	ctx context.Context,
	ghCl *github.Client,
	modCh <-chan string,
	nextCh chan<- github.RepositoryTag,
	ticker *time.Ticker,
) {
	for range <-ticker.C {
		switch {
		case <-ctx.Done():
			return
		case mod := <-modCh:
			tags, err := gh.FetchTags(ctx, ghCl, mod)
			if err != nil {
				log.Err("fetching GH tags for %s (%s)", mod, err)
			}
			for _, tag := range tags {
				send(ctx, nextCh, tag.Name)
			}
		}
	}
}
