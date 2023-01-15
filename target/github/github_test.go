package github_test

import (
	"context"
	"os"
	"testing"

	_ "github.com/abemedia/appcast/target/github"
	"github.com/abemedia/appcast/target/internal/test"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func TestGithub(t *testing.T) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		t.Skip("Missing environment variable: GITHUB_TOKEN")
	}

	t.Cleanup(func() {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client := github.NewClient(oauth2.NewClient(ctx, ts))

		file, _, _, err := client.Repositories.GetContents(ctx, "abemedia", "appcast-test", "folder/file", nil)
		if err != nil {
			t.Fatal(err)
		}

		_, _, err = client.Repositories.DeleteFile(ctx, "abemedia", "appcast-test", "folder/file", &github.RepositoryContentFileOptions{
			Message: github.String("Delete folder/file"),
			SHA:     file.SHA,
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	test.Run(t, "github://abemedia/appcast-test")
}
