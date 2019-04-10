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
	log.Info("Configuration:\n%s", cfg)

	ctx := context.Background()
	cl := gorequest.New()

	tport := &github.UnauthenticatedRateLimitedTransport{
		ClientID:     cfg.GHClientID,
		ClientSecret: cfg.GHClientSecret,
	}
	ghCl := github.NewClient(tport.Client())

	res := new(resp.Catalog)
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
		time.Duration(cfg.GHTickDurSec)*time.Second,
		time.Duration(cfg.AthensTickDurSec)*time.Second,
	)
	for _, modAndVer := range res.ModsAndVersions {
		// TODO: collate all the versions for a single module, so that
		// we don't have tons of redundant GH requests
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
