// Package github provides a target implementation for GitHub.
package github

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/kubri/kubri/target"
)

// Config represents the configuration for a GitHub target.
type Config struct {
	Owner  string
	Repo   string
	Branch string
	Folder string
}

// New returns a new GitHub target.
func New(c Config) (target.Target, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	client := github.NewClient(oauth2.NewClient(ctx, ts)).Repositories

	// Ensure config is valid.
	repo, _, err := client.Get(ctx, c.Owner, c.Repo)
	if err != nil {
		return nil, err
	}

	if c.Branch == "" {
		c.Branch = *repo.DefaultBranch
	}

	t := &githubTarget{
		client: client,
		owner:  c.Owner,
		repo:   c.Repo,
		branch: c.Branch,
		path:   c.Folder,
	}

	return t, nil
}

type githubTarget struct {
	client *github.RepositoriesService
	owner  string
	repo   string
	branch string
	path   string
}

func (t *githubTarget) NewWriter(ctx context.Context, filename string) (io.WriteCloser, error) {
	w := &fileWriter{
		t:    t,
		ctx:  ctx,
		path: path.Join(t.path, filename),
	}
	return w, nil
}

func (t *githubTarget) NewReader(ctx context.Context, filename string) (io.ReadCloser, error) {
	opt := &github.RepositoryContentGetOptions{Ref: t.branch}
	file, _, r, err := t.client.GetContents(ctx, t.owner, t.repo, path.Join(t.path, filename), opt)
	if err != nil {
		if r.StatusCode == http.StatusNotFound {
			return nil, &fs.PathError{Op: "read", Path: filename, Err: fs.ErrNotExist}
		}
		return nil, err
	}

	content, err := file.GetContent()
	if err != nil {
		return nil, err
	}

	return io.NopCloser(strings.NewReader(content)), nil
}

func (t *githubTarget) Remove(ctx context.Context, filename string) error {
	path := path.Join(t.path, filename)
	getOpt := &github.RepositoryContentGetOptions{Ref: t.branch}
	file, _, r, err := t.client.GetContents(ctx, t.owner, t.repo, path, getOpt)
	if err != nil {
		if r != nil && r.StatusCode == http.StatusNotFound {
			return &fs.PathError{Op: "remove", Path: filename, Err: fs.ErrNotExist}
		}
		return err
	}
	_, _, err = t.client.DeleteFile(ctx, t.owner, t.repo, path, &github.RepositoryContentFileOptions{
		Message: github.String("Delete " + path),
		Branch:  &t.branch,
		SHA:     file.SHA,
	})
	return err
}

func (t *githubTarget) Sub(dir string) target.Target {
	sub := *t
	sub.path = path.Join(t.path, dir)
	return &sub
}

func (t *githubTarget) URL(_ context.Context, filename string) (string, error) {
	return "https://raw.githubusercontent.com/" + path.Join(t.owner, t.repo, t.branch, t.path, filename), nil
}

func (g *githubTarget) ReadDir(ctx context.Context, p string) ([]fs.DirEntry, error) {
	opt := &github.RepositoryContentGetOptions{Ref: g.branch}
	fileContent, dirContents, _, err := g.client.GetContents(ctx, g.owner, g.repo, p, opt)
	if err != nil {
		return nil, err
	}
	if fileContent != nil {
		return nil, fs.ErrInvalid
	}

	modified := make([]func() time.Time, len(dirContents))
	for i, content := range dirContents {
		modified[i] = func() time.Time {
			opts := &github.CommitsListOptions{
				SHA:         g.branch,
				Path:        content.GetPath(),
				ListOptions: github.ListOptions{PerPage: 1},
			}

			commits, _, err := g.client.ListCommits(context.Background(), g.owner, g.repo, opts)
			if err != nil {
				panic(err)
			}
			if len(commits) == 0 {
				panic("no commits found for " + content.GetPath())
			}

			return commits[0].GetCommit().GetCommitter().GetDate()
		}
	}

	entries := make([]fs.DirEntry, len(dirContents))
	for i, content := range dirContents {
		entries[i] = &githubDirEntry{content, modified[i]}
	}
	return entries, nil
}

type fileWriter struct {
	bytes.Buffer

	t    *githubTarget
	ctx  context.Context //nolint:containedctx
	path string
}

func (w *fileWriter) Close() error {
	getOpt := &github.RepositoryContentGetOptions{Ref: w.t.branch}
	file, _, res, err := w.t.client.GetContents(w.ctx, w.t.owner, w.t.repo, w.path, getOpt)
	if err != nil && (res == nil || res.StatusCode != http.StatusNotFound) {
		return err
	}

	opt := &github.RepositoryContentFileOptions{Content: w.Bytes()}
	if w.t.branch != "" {
		opt.Branch = &w.t.branch
	}

	if res.StatusCode == http.StatusNotFound {
		opt.Message = github.String("Create " + w.path)
		_, _, err = w.t.client.CreateFile(w.ctx, w.t.owner, w.t.repo, w.path, opt)

		// Retry if writing failed due to race condition.
		// This can occur when creating a file and updating it right away in which case it might still return 404.
		// if e, ok := err.(*github.ErrorResponse); ok &&
		// 	e.Response.StatusCode == http.StatusUnprocessableEntity &&
		// 	e.Message == "Invalid request.\n\n\"sha\" wasn't supplied." {
		// 	return w.Close()
		// }
	} else {
		opt.Message = github.String("Update " + w.path)
		opt.SHA = file.SHA
		_, _, err = w.t.client.UpdateFile(w.ctx, w.t.owner, w.t.repo, w.path, opt)
	}

	return err
}
