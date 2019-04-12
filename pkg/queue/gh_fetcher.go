package queue

import (
	"github.com/arschles/crathens/pkg/ctx"
	gh "github.com/arschles/crathens/pkg/github"
	"github.com/arschles/crathens/pkg/log"
	"github.com/arschles/crathens/pkg/resp"
	"github.com/google/go-github/github"
)

func ghFetcher(
	coord ctx.Coordinator,
	ghCl *github.Client,
	nextCh chan<- resp.ModuleAndVersion,
) {
	for range coord.Ticker().C {
		select {
		case <-coord.Done():
			log.Debug("GitHub fetcher exiting because the context is done")
			return
		case mod := <-coord.Ch():
			tags, err := gh.FetchTags(coord, ghCl, mod.Module)
			if err != nil {
				log.Warn("fetching GH tags for %s (%s)", mod, err)
			}
			for _, tag := range tags {
				newMod := mod
				mod.Version = tag
				select {
				case <-coord.Done():
					return
				case nextCh <- newMod:
				}
			}
		}
	}
}
