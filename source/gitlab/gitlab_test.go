package gitlab_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	gl "github.com/xanzy/go-gitlab"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/source/gitlab"
)

func TestGitlab(t *testing.T) {
	token, ok := os.LookupEnv("GITLAB_TOKEN")
	if !ok {
		t.Skip("Missing environment variable: GITLAB_TOKEN")
	}

	client, err := gl.NewClient(token)
	if err != nil {
		t.Fatal(err)
	}

	user, _, err := client.Users.CurrentUser()
	if err != nil {
		t.Fatal(err)
	}

	owner := user.Username
	repo := fmt.Sprintf("test_%d", time.Now().UnixNano())
	pid := owner + "/" + repo

	_, _, err = client.Projects.CreateProject(&gl.CreateProjectOptions{
		Name:       &repo,
		Visibility: gl.Ptr(gl.PublicVisibility),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, err = client.Projects.DeleteProject(pid)
		if err != nil {
			t.Fatal(err)
		}
	})

	_, _, err = client.RepositoryFiles.CreateFile(pid, "test", &gl.CreateFileOptions{
		Branch:        gl.Ptr("main"),
		CommitMessage: gl.Ptr("test"),
		Content:       gl.Ptr("test"),
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range test.SourceWant() {
		_, _, err = client.Releases.CreateRelease(pid, &gl.CreateReleaseOptions{
			Description: gl.Ptr(r.Description),
			Ref:         gl.Ptr("main"),
			TagName:     gl.Ptr(r.Version),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	s, err := gitlab.New(gitlab.Config{Owner: owner, Repo: repo})
	if err != nil {
		t.Fatal(err)
	}

	test.Source(t, s, func(version, asset string) string {
		links, _, _ := client.ReleaseLinks.ListReleaseLinks(pid, version, nil)
		for _, link := range links {
			if link.Name == asset {
				return link.URL
			}
		}
		return ""
	})
}
