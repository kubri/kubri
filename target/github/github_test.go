package github_test

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	gh "github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/target/github"
)

func TestGithub(t *testing.T) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		t.Skip("Missing environment variable: GITHUB_TOKEN")
	}

	ctx := t.Context()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client := gh.NewClient(oauth2.NewClient(ctx, ts))
	user, _, err := client.Users.Get(ctx, "")
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
			repo := fmt.Sprintf("test_%d", time.Now().UnixNano())

			r, _, err := client.Repositories.Create(ctx, "", &gh.Repository{Name: &repo})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				_, err := client.Repositories.Delete(ctx, owner, repo)
				if err != nil {
					t.Fatal(err)
				}
			})

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
			}, test.WithDelay(time.Second))

			t.Run("Error", func(t *testing.T) {
				_, err = github.New(github.Config{Owner: "owner", Repo: "repo", Branch: tc.branch})
				if err == nil {
					t.Fatal("should fail for invalid repo")
				}
			})
		})
	}
}
