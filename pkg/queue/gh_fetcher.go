package queue

import (
	gh "github.com/arschles/crathens/pkg/github"
	"github.com/arschles/crathens/pkg/log"
	"github.com/arschles/crathens/pkg/resp"
	"github.com/google/go-github/github"
)

func ghFetcher(
	coord *coordinator,
	ghCl *github.Client,
	nextCh chan<- resp.ModuleAndVersion,
) {
	for range coord.ticker.C {
		select {
		case <-coord.ctx.Done():
			log.Debug("GitHub fetcher exiting because the context is done")
			return
		case mod := <-coord.ch:
			tags, err := gh.FetchTags(coord.ctx, ghCl, mod.Module)
			if err != nil {
				log.Warn("fetching GH tags for %s (%s)", mod, err)
			}
			for _, tag := range tags {
				newMod := mod
				mod.Version = tag
				select {
				case <-coord.ctx.Done():
					return
				case nextCh <- newMod:
				}
			}
		}
	}
}
