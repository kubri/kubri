package target_test

import (
	"context"
	"io"
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
	"github.com/google/go-cmp/cmp"
)

type fakeTarget struct {
	source.Config
}

func (*fakeTarget) NewWriter(ctx context.Context, path string) (io.WriteCloser, error) {
	return nil, nil
}

func (*fakeTarget) NewReader(ctx context.Context, path string) (io.ReadCloser, error) {
	return nil, nil
}

func (*fakeTarget) Sub(dir string) target.Target {
	return nil
}

func TestOpen(t *testing.T) {
	target.Register("fake", func(c source.Config) (target.Target, error) {
		return &fakeTarget{c}, nil
	})

	t.Setenv("FAKE_TOKEN", "fake")
	s, err := target.Open("fake://user/repo")
	if err != nil {
		t.Fatal(err)
	}

	want := &fakeTarget{source.Config{Repo: "user/repo", Token: "fake"}}
	if diff := cmp.Diff(want, s); diff != "" {
		t.Error(diff)
	}
}
