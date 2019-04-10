package queue

import (
	"context"

	"github.com/arschles/crathens/pkg/resp"
)

type Crawler interface {
	Crawl(context.Context, resp.ModuleAndVersion) error
	Wait(context.Context) error
}
