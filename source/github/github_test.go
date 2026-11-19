package github_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	gh "github.com/google/go-github/v83/github"

	"github.com/kubri/kubri/internal/githubclient"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/source/github"
)

func TestGithub(t *testing.T) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		t.Skip("Missing environment variable: GITHUB_TOKEN")
	}

	client := githubclient.New(token)
	user, _, err := client.Users.Get(t.Context(), "")
	if err != nil {
		t.Fatal(err)
	}

	owner := user.GetLogin()
	repo := fmt.Sprintf("test_%d", time.Now().UnixNano())

	_, _, err = client.Repositories.Create(t.Context(), "", &gh.Repository{Name: &repo, AutoInit: new(true)})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, err := client.Repositories.Delete(context.Background(), owner, repo)
		if err != nil {
			t.Fatal(err)
		}
	})

	// Wait for the repository to become available.
	for range 10 {
		_, resp, err := client.Repositories.Get(t.Context(), owner, repo)
		if err == nil {
			break
		}
		if resp.StatusCode == http.StatusNotFound {
			time.Sleep(time.Second)
			continue
		}
		t.Fatalf("unexpected error waiting for repo to be available: %v", err)
	}

	for _, r := range test.SourceWant() {
		_, _, err = client.Repositories.CreateRelease(t.Context(), owner, repo, &gh.RepositoryRelease{
			TagName: new(r.Version),
			Body:    new(r.Description),
		})
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
