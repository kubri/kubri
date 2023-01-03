package emulator

import (
	"context"
	"log"
	"testing"

	"github.com/fullstorydev/emulators/storage/gcsemu"
	"github.com/testcontainers/testcontainers-go/wait"
	azblob "gocloud.dev/blob/azureblob"
)

func AzureBlob(t *testing.T, bucket string) string {
	t.Helper()

	host := TestContainer(t, Container{
		Image:   "mcr.microsoft.com/azure-storage/azurite:latest",
		Port:    10000,
		Command: []string{"azurite-blob", "--blobHost", "0.0.0.0"},
	})

	t.Setenv("AZURE_STORAGE_ACCOUNT", "devstoreaccount1")
	t.Setenv("AZURE_STORAGE_KEY", "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==") //nolint:lll
	t.Setenv("AZURE_STORAGE_DOMAIN", host)
	t.Setenv("AZURE_STORAGE_PROTOCOL", "http")

	client, err := azblob.NewDefaultServiceClient(azblob.ServiceURL("http://" + host + "/devstoreaccount1"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.CreateContainer(context.Background(), bucket, nil)
	if err != nil {
		log.Fatal(err)
	}

	return host
}

func GCS(t *testing.T, bucket string) string {
	t.Helper()

	emu, err := gcsemu.NewServer(":0", gcsemu.Options{})
	if err != nil {
		log.Fatal(err)
	}

	if err = emu.InitBucket(bucket); err != nil {
		log.Fatal(err)
	}

	t.Setenv("STORAGE_EMULATOR_HOST", emu.Addr)

	return emu.Addr
}

func S3(t *testing.T, bucket string) string {
	t.Helper()

	host := TestContainer(t, Container{
		Image: "adobe/s3mock:latest",
		Port:  9090,
		Env:   map[string]string{"initialBuckets": bucket},
		Wait:  wait.ForHTTP("/").WithPort("9090").WithStatusCodeMatcher(nil),
	})

	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_REGION", "us-east-1")

	return host
}
