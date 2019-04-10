package queue

import (
	"context"
	"log"
	"sync"
	"time"

	gh "github.com/arschles/crathens/pkg/github"
	"github.com/google/go-github/github"
)

type TagFetcher struct {
	ghCl  *github.Client
	adder func(string, string)
	wg    sync.WaitGroup
}

func NewTagFetcher(
	ghCl *github.Client,
	adder func(string, string),
) *TagFetcher {
	return &TagFetcher{
		ghCl:  ghCl,
		adder: adder,
	}
}

func (tf *TagFetcher) Fetch(ctx context.Context, mod string) error {
	owner, repo, err := gh.SplitModule(mod)
	if err != nil {
		return err
	}
	tf.wg.Add(1)
	go func() {
		defer tf.wg.Done()
		time.Sleep(500 * time.Millisecond)
		tags, _, err := tf.ghCl.Repositories.ListTags(ctx, owner, repo, nil)
		if err != nil {
			log.Fatal(err)
		}
		for _, tag := range tags {
			tf.adder(mod, *tag.Name)
		}
	}()
	return nil
}

func (tf *TagFetcher) Wait() {
	tf.wg.Wait()
}
