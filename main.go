package main

import (
	"context"
	"os"
	"time"

	"github.com/arschles/crathens/pkg/log"
	"github.com/arschles/crathens/pkg/queue"
	"github.com/arschles/crathens/pkg/resp"
	"github.com/google/go-github/github"
	"github.com/parnurzeal/gorequest"
)

const endpoint = "https://athens.azurefd.net"

func main() {
	ctx := context.Background()
	cl := gorequest.New()

	tport := &github.UnauthenticatedRateLimitedTransport{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}
	ghCl := github.NewClient(tport.Client())

	res := new(resp.Catalog)
	resp, _, err := cl.Get(endpoint + "/catalog").EndStruct(res)
	if err != nil {
		log.Err("getting the /catalog endpoint (%s)", err)
		os.Exit(1)
	} else if resp.StatusCode != 200 {
		log.Err("/catalog status code was %d", resp.StatusCode)
		os.Exit(1)
	}

	var crawler queue.Crawler
	for _, modAndVer := range res.ModsAndVersions {
		toCtx, done := context.WithTimeout(ctx, 500*time.Millisecond)
		defer done()
		if err := crawler.Enqueue(toCtx, modAndVer); err != nil {
			log.Warn("crawling %s (%s)", modAndVer, err)
		}
	}

	waitCtx, done := context.WithTimeout(ctx, 1*time.Second)
	defer done()
	if err := crawler.Wait(waitCtx); err != nil {
		log.Err("Waiting for crawler failed (%s)", err)
		os.Exit(1)
	}
}
