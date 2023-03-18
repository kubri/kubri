package github

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/abemedia/appcast/target"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Config struct {
	Owner  string
	Repo   string
	Branch string
	Folder string
}

type githubTarget struct {
	client *github.RepositoriesService
	owner  string
	repo   string
	branch string
	path   string
}

func New(c Config) (target.Target, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	client := github.NewClient(oauth2.NewClient(ctx, ts)).Repositories

	// Ensure config is valid.
	if c.Branch == "" {
		repo, _, err := client.Get(ctx, c.Owner, c.Repo)
		if err != nil {
			return nil, err
		}
		c.Branch = *repo.DefaultBranch
	} else {
		_, _, err := client.GetBranch(ctx, c.Owner, c.Repo, c.Branch)
		if err != nil {
			return nil, err
		}
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
	file, _, _, err := t.client.GetContents(ctx, t.owner, t.repo, path.Join(t.path, filename), opt)
	if err != nil {
		return nil, err
	}

	content, err := file.GetContent()
	if err != nil {
		return nil, err
	}

	return io.NopCloser(strings.NewReader(content)), nil
}

func (t *githubTarget) Sub(dir string) target.Target {
	sub := *t
	sub.path = filepath.Join(t.path, dir)
	return &sub
}

func (t *githubTarget) URL(_ context.Context, filename string) (string, error) {
	return "https://raw.githubusercontent.com/" + path.Join(t.owner, t.repo, t.branch, t.path, filename), nil
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
	if err != nil && (res == nil || res.StatusCode != 404) {
		return err
	}

	opt := &github.RepositoryContentFileOptions{Content: w.Bytes()}
	if w.t.branch != "" {
		opt.Branch = &w.t.branch
	}

	if res.StatusCode == http.StatusNotFound {
		opt.Message = github.String("Create " + w.path)
		_, _, err = w.t.client.CreateFile(w.ctx, w.t.owner, w.t.repo, w.path, opt)
	} else {
		opt.Message = github.String("Update " + w.path)
		opt.SHA = file.SHA
		_, _, err = w.t.client.UpdateFile(w.ctx, w.t.owner, w.t.repo, w.path, opt)
	}

	return err
}
