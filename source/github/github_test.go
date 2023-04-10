package github_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/source/github"
	gh "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func TestGithub(t *testing.T) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		t.Skip("Missing environment variable: GITHUB_TOKEN")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client := gh.NewClient(oauth2.NewClient(ctx, ts))
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		t.Fatal(err)
	}

	owner := user.GetLogin()
	repo := fmt.Sprintf("test_%d", time.Now().UnixNano())

	_, _, err = client.Repositories.Create(ctx, "", &gh.Repository{Name: &repo})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, err := client.Repositories.Delete(ctx, owner, repo)
		if err != nil {
			t.Fatal(err)
		}
	})

	opt := &gh.RepositoryContentFileOptions{
		Message: gh.String("test"),
		Content: []byte("test"),
	}
	_, _, err = client.Repositories.CreateFile(ctx, owner, repo, "test", opt)
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range test.SourceWant() {
		opt := &gh.RepositoryRelease{TagName: &r.Version, Body: &r.Description}
		_, _, err = client.Repositories.CreateRelease(ctx, owner, repo, opt)
		if err != nil {
			t.Fatal(err)
		}
	}

	s, err := github.New(github.Config{Owner: owner, Repo: repo})
	if err != nil {
		t.Fatal(err)
	}

	test.Source(t, s, func(version, asset string) string {
		return "https://github.com/" + owner + "/" + repo + "/releases/download/" + version + "/" + asset
	})
}
