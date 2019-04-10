package queue

import "github.com/arschles/crathens/pkg/log"

func athensWarmer(coord *coordinator) {
	for range coord.ticker.C {
		select {
		case <-coord.ctx.Done():
			return
		case mod := <-coord.ch:
			log.Info("warming athens with %s", mod)
		}
	}
}
