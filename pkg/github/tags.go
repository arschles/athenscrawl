package github

import (
	"context"

	"github.com/google/go-github/github"
)

// FetchTags returns all the tags for the given module string
func FetchTags(
	ctx context.Context,
	ghCl *github.Client,
	mod string,
) ([]string, error) {
	owner, repo, err := SplitModule(mod)
	if err != nil {
		return nil, err
	}
	tags, _, err := ghCl.Repositories.ListTags(ctx, owner, repo, nil)
	if err != nil {
		return nil, err
	}
	vers := make([]string, len(tags))
	for i, tag := range tags {
		vers[i] = *tag.Name
	}
	return vers, nil
}
