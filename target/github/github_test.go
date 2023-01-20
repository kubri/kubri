package github_test

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/abemedia/appcast/target/github"
	"github.com/abemedia/appcast/target/internal/test"
	gh "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func TestGithub(t *testing.T) {
	owner := "abemedia"
	repo := "appcast-test"

	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		t.Skip("Missing environment variable: GITHUB_TOKEN")
	}

	t.Cleanup(func() {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client := gh.NewClient(oauth2.NewClient(ctx, ts))

		file, _, _, err := client.Repositories.GetContents(ctx, owner, repo, "path/to/file", nil)
		if err != nil {
			t.Fatal(err)
		}

		opt := &gh.RepositoryContentFileOptions{Message: gh.String("Delete path/to/file"), SHA: file.SHA}
		_, _, err = client.Repositories.DeleteFile(ctx, owner, repo, "path/to/file", opt)
		if err != nil {
			t.Fatal(err)
		}
	})

	tgt, err := github.New(github.Config{Owner: owner, Repo: repo})
	if err != nil {
		t.Fatal(err)
	}

	test.Run(t, tgt, func(asset string) string {
		return "https://raw.githubusercontent.com/" + path.Join(owner, repo, "master", asset)
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
