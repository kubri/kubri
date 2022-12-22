package testutils

import (
	"context"
	"log"
	"net"
	"strconv"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Container struct {
	Image      string
	Port       int
	Env        map[string]string
	Command    []string
	Entrypoint []string
	Wait       wait.Strategy
}

type Service struct {
	Host string

	container testcontainers.Container
}

func (s *Service) Close() {
	terminate(s.container)
}

func RunContainer(c Container) (*Service, error) {
	ctx := context.Background()
	port := strconv.Itoa(c.Port)

	w := c.Wait
	if w == nil {
		w = wait.ForListeningPort(nat.Port(port))
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        c.Image,
			ExposedPorts: []string{port},
			Env:          c.Env,
			Cmd:          c.Command,
			Entrypoint:   c.Entrypoint,
			WaitingFor:   w,
		},
		Started: true,
		Logger:  NopLogger{},
	})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			terminate(container)
		}
	}()

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	return &Service{Host: net.JoinHostPort(hostIP, mappedPort.Port()), container: container}, nil
}

func TestContainer(t *testing.T, c Container) string {
	t.Helper()
	s, err := RunContainer(c)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Close)
	return s.Host
}

func terminate(c testcontainers.Container) {
	if err := c.Terminate(context.Background()); err != nil {
		log.Printf("Error shutting down container: %s", err)
	}
}

type NopLogger struct{}

func (NopLogger) Printf(string, ...any) {}
func (NopLogger) Print(...any)          {}
