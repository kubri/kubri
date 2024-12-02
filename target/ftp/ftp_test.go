package ftp_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"net"
	"net/textproto"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/target/ftp"
)

func TestFTP(t *testing.T) {
	minPort := 20000 + rand.IntN(999)*10
	maxPort := minPort + 10
	portRange := fmt.Sprintf("%[1]d-%[2]d:%[1]d-%[2]d", minPort, maxPort)

	ctx := context.Background()

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context: "testdata",
				Repo:    "vsftpd-alpine",
			},
			ExposedPorts: []string{"21", portRange},
			Env: map[string]string{
				"FTP_USER":      "user",
				"FTP_PASS":      "password",
				"PASV_MIN_PORT": strconv.Itoa(minPort),
				"PASV_MAX_PORT": strconv.Itoa(maxPort),
				"PASV_ADDRESS":  "127.0.0.1",
			},
			WaitingFor: wait.ForListeningPort("21"),
		},
		Started: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Terminate(ctx)

	address, err := c.PortEndpoint(ctx, "21", "")
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("FTP_USER", "user")
	t.Setenv("FTP_PASSWORD", "password")

	tgt, err := ftp.New(ftp.Config{
		Address: address,
		Folder:  "folder",
		URL:     "http://dl.example.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	test.Target(t, tgt, func(asset string) string {
		return "http://dl.example.com/" + asset
	})

	t.Run("Error", func(t *testing.T) {
		tests := []struct {
			name     string
			config   ftp.Config
			user     string
			password string
			err      error
		}{
			{
				name:   "DialError",
				config: ftp.Config{Address: ""},
				err:    &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("missing address")},
			},
			{
				name:     "LoginError",
				config:   ftp.Config{Address: address},
				user:     "wrong",
				password: "wrong",
				err:      &textproto.Error{Code: 530, Msg: "Login incorrect."},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Setenv("FTP_USER", tt.user)
				t.Setenv("FTP_PASSWORD", tt.password)
				_, err := ftp.New(tt.config)
				if diff := cmp.Diff(tt.err, err, test.ExportAll()); diff != "" {
					t.Error(diff)
				}
			})
		}
	})
}
