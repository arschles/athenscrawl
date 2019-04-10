package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/arschles/crathens/pkg/queue"
	"github.com/arschles/crathens/pkg/resp"
	"github.com/google/go-github/github"
	"github.com/parnurzeal/gorequest"
)

const endpoint = "https://athens.azurefd.net"

func main() {
	ctx := context.Background()
	cl := gorequest.New()

	// https://github.com/settings/applications/1045516
	tport := &github.UnauthenticatedRateLimitedTransport{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}
	ghCl := github.NewClient(tport.Client())

	res := new(resp.Catalog)
	resp, _, err := cl.Get(endpoint + "/catalog").EndStruct(res)
	if err != nil {
		log.Fatal(err)
	} else if resp.StatusCode != 200 {
		log.Fatalf("status code was %d", resp.StatusCode)
	}

	var crawler queue.Crawler
	for _, modAndVer := range res.ModsAndVersions {
		toCtx, done := context.WithTimeout(ctx, 500 * time.Millisecond)
		defer done()
		if err := crawler.Enqueue(toCtx, modAndVer); err != nil {
			log.Printf("Error crawling %s (%s)", modAndVer, err)
		}
	}

	waitCtx, done := context.WithTimeout(ctx, 1 * time.Second)
	defer done()
	if err := crawler.Wait(waitCtx); err != nil {
		log.Fatalf("Waiting for crawler failed (%s)", err)
	}

// 	tagCrawler := queue.NewTagCrawler()
// 	fetcher := queue.NewTagFetcher(ghCl, tagCrawler.Add)

// 	go func() {
// 		for tag := range tagCrawler.Ch() {
// 			log.Printf("running go get for %s@%s", tag.Module, tag.Name)
// 		}
// 	}()

// 	for _, mav := range res.ModsAndVersions {
// 		log.Printf("module %s", mav.Module)
// 		if strings.HasPrefix(mav.Module, "github.com") {
// 			log.Printf("----> using GH API to get versions for %s", mav.Module)
// 			if err := fetcher.Fetch(ctx, mav.Module); err != nil {
// 				log.Fatal(err)
// 			}
// 		}
// 	}
// 	fetcher.Wait()
// }
