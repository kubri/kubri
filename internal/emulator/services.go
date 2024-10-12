package emulator

import (
	"context"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/fullstorydev/emulators/storage/gcsemu"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gocloud.dev/blob/azureblob"
)

func AzureBlob(t *testing.T, bucket string) string {
	t.Helper()

	c := runContainer(t, testcontainers.ContainerRequest{
		Image:        "mcr.microsoft.com/azure-storage/azurite:latest",
		ExposedPorts: []string{"10000"},
		Cmd:          []string{"azurite-blob", "--blobHost", "0.0.0.0"},
		WaitingFor:   wait.ForListeningPort(nat.Port("10000")),
	})
	host, err := c.PortEndpoint(context.Background(), "10000", "")
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("AZURE_STORAGE_ACCOUNT", "devstoreaccount1")
	t.Setenv("AZURE_STORAGE_KEY",
		"Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==")
	t.Setenv("AZURE_STORAGE_DOMAIN", host)
	t.Setenv("AZURE_STORAGE_PROTOCOL", "http")

	u := azureblob.ServiceURL("http://" + host + "/devstoreaccount1")
	client, err := azureblob.NewDefaultClient(u, azureblob.ContainerName(bucket))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = client.Create(context.Background(), nil); err != nil {
		t.Fatal(err)
	}

	return host
}

func GCS(t *testing.T, bucket string) string {
	t.Helper()

	emu, err := gcsemu.NewServer(":0", gcsemu.Options{})
	if err != nil {
		t.Fatal(err)
	}

	if err = emu.InitBucket(bucket); err != nil {
		t.Fatal(err)
	}

	t.Setenv("STORAGE_EMULATOR_HOST", emu.Addr)

	return emu.Addr
}

func S3(t *testing.T, bucket string) string {
	t.Helper()

	c := runContainer(t, testcontainers.ContainerRequest{
		Image:        "adobe/s3mock:latest",
		ExposedPorts: []string{"9090"},
		Env:          map[string]string{"initialBuckets": bucket},
		WaitingFor:   wait.ForHTTP("/").WithPort("9090").WithStatusCodeMatcher(nil),
	})
	host, err := c.PortEndpoint(context.Background(), "9090", "http")
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	return host
}
