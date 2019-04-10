package queue

import (
	"fmt"

	"github.com/arschles/crathens/pkg/log"
	"github.com/parnurzeal/gorequest"
	"github.com/souz9/errlist"
)

func modInfoPath(endpoint, mod, ver string) string {
	return fmt.Sprintf("%s/%s/@v/%s.info", endpoint, mod, ver)
}

func athensWarmer(endpoint string, coord *coordinator) {
	for range coord.ticker.C {
		select {
		case <-coord.ctx.Done():
			return
		case mod := <-coord.ch:
			resp, _, errs := gorequest.
				New().
				Get(modInfoPath(endpoint, mod.Module, mod.Version)).
				End()
			if len(errs) > 0 {
				// TODO: send an error back on a chan
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
