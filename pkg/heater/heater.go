package heater

import (
	"fmt"
	"net/http"

	"github.com/arschles/crathens/pkg/queue"
	"github.com/parnurzeal/gorequest"
	"github.com/souz9/errlist"
)

func Heat(proxyURL string, mat queue.ModuleAndTag) error {
	cl := gorequest.New()
	url := fmt.Sprintf("%s/%s", proxyURL, proxyPath(mat))
	resp, _, err := cl.Get(url).End()
	if err != nil {
		return errlist.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s failed with code %d", url, resp.StatusCode)
	}
	return nil
}

func proxyPath(mat queue.ModuleAndTag) string {
	return fmt.Sprintf("%s/%s.info", mat.Module, mat.Name)
}
