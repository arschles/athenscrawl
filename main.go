package main

import (
	"context"
	"log"
	"os"
	"strings"

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

	res := new(catalogRes)
	resp, _, err := cl.Get(endpoint + "/catalog").EndStruct(res)
	if err != nil {
		log.Fatal(err)
	} else if resp.StatusCode != 200 {
		log.Fatalf("status code was %d", resp.StatusCode)
	}
	for _, mav := range res.ModsAndVersions {
		log.Printf("module %s", mav.Module)
		if strings.HasPrefix(mav.Module, "github.com") {
			log.Printf("----> using GH API to get versions")
			owner, repo, err := splitGHModule(mav.Module)
			if err != nil {
				log.Fatal(err)
			}
			tags, _, err := ghCl.Repositories.ListTags(ctx, owner, repo, nil)
			if err != nil {
				log.Fatal(err)
			}
			for _, tag := range tags {
				log.Print(tag.Name)
			}
		}
	}
}
