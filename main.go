package main

import (
	"context"
	"os"
	"time"

	"github.com/arschles/crathens/pkg/config"
	"github.com/arschles/crathens/pkg/log"
	"github.com/arschles/crathens/pkg/queue"
	"github.com/arschles/crathens/pkg/resp"
	"github.com/google/go-github/github"
	"github.com/kelseyhightower/envconfig"
	"github.com/parnurzeal/gorequest"
)

func main() {
	cfg := new(config.Config)
	if err := envconfig.Process("crathens", cfg); err != nil {
		log.Err("Processing configuration (%s)", err)
		os.Exit(1)
	}
	log.SetDebug(cfg.Debug)
	log.Debug("Configuration:\n%s", cfg)

	ctx := context.Background()
	cl := gorequest.New()

	tport := &github.UnauthenticatedRateLimitedTransport{
		ClientID:     cfg.GHClientID,
		ClientSecret: cfg.GHClientSecret,
	}
	ghCl := github.NewClient(tport.Client())

	res := new(resp.Catalog)
	log.Debug("getting catalog endpoint")
	resp, _, err := cl.Get(cfg.Endpoint + "/catalog").EndStruct(res)
	if err != nil {
		log.Err("getting the /catalog endpoint (%s)", err)
		os.Exit(1)
	} else if resp.StatusCode != 200 {
		log.Err("/catalog status code was %d", resp.StatusCode)
		os.Exit(1)
	}

	crawler := queue.InMemory(
		ctx,
		cfg.Endpoint,
		ghCl,
		cfg.GHTickDur(),
		cfg.AthensTickDur(),
	)
	for _, modAndVer := range res.ModsAndVersions {
		log.Debug("enqueueing %s", modAndVer)
		if err := crawler.Enqueue(ctx, modAndVer); err != nil {
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
