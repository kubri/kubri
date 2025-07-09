package github_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	gh "github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/source/github"
)

func TestGithub(t *testing.T) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		t.Skip("Missing environment variable: GITHUB_TOKEN")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client := gh.NewClient(oauth2.NewClient(t.Context(), ts))
	user, _, err := client.Users.Get(t.Context(), "")
	if err != nil {
		t.Fatal(err)
	}

	owner := user.GetLogin()
	repo := fmt.Sprintf("test_%d", time.Now().UnixNano())

	_, _, err = client.Repositories.Create(t.Context(), "", &gh.Repository{Name: &repo})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, err := client.Repositories.Delete(context.Background(), owner, repo)
		if err != nil {
			t.Fatal(err)
		}
	})

	_, _, err = client.Repositories.CreateFile(t.Context(), owner, repo, "test", &gh.RepositoryContentFileOptions{
		Message: gh.String("test"),
		Content: []byte("test"),
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range test.SourceWant() {
		_, _, err = client.Repositories.CreateRelease(t.Context(), owner, repo, &gh.RepositoryRelease{
			TagName: gh.String(r.Version),
			Body:    gh.String(r.Description),
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
