package emulator

import (
	"archive/tar"
	"bytes"
	"context"
	"strings"
	"testing"
	"unicode"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/exec"
)

type Container struct{ testcontainers.Container }

func (c *Container) Exec(t *testing.T, script string) string {
	t.Helper()

	var buf bytes.Buffer
	opt := exec.ProcessOptionFunc(func(opts *exec.ProcessOptions) {
		if opts.Reader != nil {
			_, _ = stdcopy.StdCopy(&buf, &buf, opts.Reader)
		}
	})

	code, _, err := c.Container.Exec(context.Background(), []string{"sh", "-c", script}, opt)
	if err != nil {
		t.Fatal(err)
	}

	s := strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(buf.String())
	t.Logf("%s\n%s", script, s)

	if code != 0 {
		t.FailNow()
	}

	return strings.TrimRightFunc(s, unicode.IsSpace)
}

func Image(t *testing.T, name string) Container {
	t.Helper()
	return runContainer(t, testcontainers.ContainerRequest{
		Image:              name,
		HostConfigModifier: func(hc *container.HostConfig) { hc.NetworkMode = "host" },
		Entrypoint:         []string{"tail", "-f", "/dev/null"},
	})
}

func Build(t *testing.T, dockerfile string) Container {
	t.Helper()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	hdr := &tar.Header{
		Name: "Dockerfile",
		Mode: 0o600,
		Size: int64(len(dockerfile)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write([]byte(dockerfile)); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}

	return runContainer(t, testcontainers.ContainerRequest{
		FromDockerfile:     testcontainers.FromDockerfile{ContextArchive: bytes.NewReader(buf.Bytes())},
		HostConfigModifier: func(hc *container.HostConfig) { hc.NetworkMode = "host" },
	})
}

func runContainer(t *testing.T, cr testcontainers.ContainerRequest) Container {
	t.Helper()
	ctx := context.Background()
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: cr,
		Started:          true,
		Logger:           nopLogger{},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Terminate(ctx) })
	return Container{c}
}

type nopLogger struct{}

func (nopLogger) Printf(string, ...any) {}
func (nopLogger) Print(...any)          {}
