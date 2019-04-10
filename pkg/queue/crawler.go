package queue

import (
	"context"

	"github.com/arschles/crathens/pkg/resp"
)

// Crawler accepts module/version pairs, fetches version lists for
// the module, and then submits the new versions to Athens
type Crawler interface {
	Enqueue(context.Context, resp.ModuleAndVersion) error
	Wait(context.Context) error
}
