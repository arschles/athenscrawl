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

func main() {
	endpoint := "https://athens.azurefd.net"
	envEndpoint := os.Getenv("GOPROXY")
	if envEndpoint != "" {
		endpoint = envEndpoint
	}

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

	crawler := queue.InMemory(
		ctx,
		ghCl,
		100*time.Millisecond, // TODO: make configurable
		100*time.Millisecond, // TODO: make configurable
	)
	for _, modAndVer := range res.ModsAndVersions {
		toCtx, done := context.WithTimeout(ctx, 500*time.Millisecond)
		defer done()
		// TODO: collate all the versions for a single module, so that
		// we don't have tons of redundant GH requests
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
