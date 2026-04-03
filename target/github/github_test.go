package github_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	gh "github.com/google/go-github/v83/github"
	"golang.org/x/oauth2"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/target/github"
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

	tests := []struct {
		name   string
		branch string
	}{
		{"DefaultBranch", ""},
		{"WithBranch", "foo"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := fmt.Sprintf("test_%d", time.Now().UnixNano())

			r, _, err := client.Repositories.Create(t.Context(), "", &gh.Repository{Name: &repo})
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

			tgt, err := github.New(github.Config{Owner: user.GetLogin(), Repo: repo, Branch: tc.branch})
			if err != nil {
				t.Fatal(err)
			}

			branch := tc.branch
			if branch == "" {
				branch = r.GetDefaultBranch()
			}

			test.Target(t, tgt, func(asset string) string {
				return "https://raw.githubusercontent.com/" + path.Join(owner, repo, branch, asset)
			}, test.WithDelay(3*time.Second))

			t.Run("Error", func(t *testing.T) {
				_, err = github.New(github.Config{Owner: "owner", Repo: "repo", Branch: tc.branch})
				if err == nil {
					t.Fatal("should fail for invalid repo")
				}
			})
		})
	}
}
