package github

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type githubTarget struct {
	client *github.Client
	owner  string
	repo   string
	path   string
}

func New(c source.Config) (target.Target, error) {
	owner, repo, ok := strings.Cut(c.Repo, "/")
	if !ok {
		return nil, fmt.Errorf("invalid repo: %s", c.Repo)
	}

	var client *http.Client
	if c.Token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.Token})
		client = oauth2.NewClient(context.Background(), ts)
	}

	s := &githubTarget{
		client: github.NewClient(client),
		owner:  owner,
		repo:   repo,
	}
	return s, nil
}

func (t *githubTarget) NewWriter(ctx context.Context, filename string) (io.WriteCloser, error) {
	return &fileWriter{t: t, ctx: ctx, path: path.Join(t.path, filename)}, nil
}

func (t *githubTarget) NewReader(ctx context.Context, filename string) (io.ReadCloser, error) {
	file, _, _, err := t.client.Repositories.GetContents(ctx, t.owner, t.repo, path.Join(t.path, filename), nil)
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

type fileWriter struct {
	bytes.Buffer

	t    *githubTarget
	ctx  context.Context //nolint:containedctx
	path string
}

func (w *fileWriter) Close() error {
	file, _, res, err := w.t.client.Repositories.GetContents(w.ctx, w.t.owner, w.t.repo, w.path, nil)
	if err != nil && (res == nil || res.StatusCode != 404) {
		return err
	}

	opt := &github.RepositoryContentFileOptions{Content: w.Bytes()}

	if res.StatusCode == http.StatusNotFound {
		opt.Message = github.String("Create " + w.path)
		_, _, err = w.t.client.Repositories.CreateFile(w.ctx, w.t.owner, w.t.repo, w.path, opt)
	} else {
		opt.Message = github.String("Update " + w.path)
		opt.SHA = file.SHA
		_, _, err = w.t.client.Repositories.UpdateFile(w.ctx, w.t.owner, w.t.repo, w.path, opt)
	}

	return err
}

//nolint:gochecknoinits
func init() { target.Register("github", New) }
