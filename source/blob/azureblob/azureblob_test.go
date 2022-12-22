package azureblob_test

import (
	"context"
	"log"
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/azureblob"
	"github.com/abemedia/appcast/source/blob/internal/testutils"
	azblob "gocloud.dev/blob/azureblob"
)

func TestBlobAzure(t *testing.T) {
	host := testutils.TestContainer(t, testutils.Container{
		Image:   "mcr.microsoft.com/azure-storage/azurite:latest",
		Port:    10000,
		Command: []string{"azurite-blob", "--blobHost", "0.0.0.0"},
	})

	t.Setenv("AZURE_STORAGE_ACCOUNT", "devstoreaccount1")
	t.Setenv("AZURE_STORAGE_KEY", "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==")
	t.Setenv("AZURE_STORAGE_DOMAIN", host)
	t.Setenv("AZURE_STORAGE_PROTOCOL", "http")

	dir := "downloads/test"
	repo := "bucket/" + dir

	client, err := azblob.NewDefaultServiceClient(azblob.ServiceURL("http://" + host + "/devstoreaccount1"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.CreateContainer(context.Background(), "bucket", nil)
	if err != nil {
		log.Fatal(err)
	}

	s, err := azureblob.New(source.Config{Repo: repo})
	if err != nil {
		t.Fatal(err)
	}

	makeURL := func(version, asset string) string {
		return "http://" + host + "/devstoreaccount1/" + filepath.Join(repo, version, asset)
	}

	testutils.TestBlob(t, s, makeURL)
}
