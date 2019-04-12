package queue

import (
	"fmt"

	"github.com/arschles/crathens/pkg/ctx"
	"github.com/arschles/crathens/pkg/log"
	"github.com/parnurzeal/gorequest"
	"github.com/souz9/errlist"
)

func modInfoPath(endpoint, mod, ver string) string {
	return fmt.Sprintf("%s/%s/@v/%s.info", endpoint, mod, ver)
}

func athensWarmer(endpoint string, coord ctx.Coordinator) {
	for range coord.Ticker().C {
		select {
		case <-coord.Done():
			log.Debug("Athens warmer exiting because the context is done")
			return
		case mod := <-coord.Ch():
			resp, _, errs := gorequest.
				New().
				Get(modInfoPath(endpoint, mod.Module, mod.Version)).
				End()
			if len(errs) > 0 {
				log.Warn(
					"fetching GH tags for %s (%s)",
					mod,
					errlist.Error(errs),
				)
			} else if resp.StatusCode != 200 {
				log.Warn(
					"fetching GH tags for %s returned status code %d",
					mod,
					resp.StatusCode,
				)
			} else {
				log.Debug("Warmed Athens with %s", mod)
			}
		}
	}
}
