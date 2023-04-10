package github_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/target/github"
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

	tgt, err := github.New(github.Config{Owner: owner, Repo: repo})
	if err != nil {
		t.Fatal(err)
	}

	test.Target(t, tgt, func(asset string) string {
		return "https://raw.githubusercontent.com/" + path.Join(owner, repo, "main", asset)
	})

	_, err = github.New(github.Config{Owner: owner, Repo: repo, Branch: "foo"})
	if err == nil {
		t.Fatal("should fail for invalid branch")
	}

	_, err = github.New(github.Config{Owner: owner, Repo: "foo"})
	if err == nil {
		t.Fatal("should fail for invalid repo")
	}
}
