package azureblob_test

import (
	"path"
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	_ "github.com/abemedia/appcast/source/blob/azureblob"
	"github.com/abemedia/appcast/source/blob/internal/test"
)

func TestAzureblob(t *testing.T) {
	host := emulator.AzureBlob(t, "bucket")
	repo := "bucket/downloads/test"

	test.Run(t, "azblob://"+repo, func(version, asset string) string {
		return "http://" + host + "/devstoreaccount1/" + path.Join(repo, version, asset)
	})
}
