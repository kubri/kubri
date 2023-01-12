package azureblob_test

import (
	"net/url"
	"path"
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	_ "github.com/abemedia/appcast/source/blob/azureblob"
	"github.com/abemedia/appcast/source/blob/internal/test"
)

func TestAzureblob(t *testing.T) {
	host := emulator.AzureBlob(t, "bucket")
	repo := "downloads/test"

	test.Run(t, "azblob://bucket/"+repo, func(version, asset string) string {
		return "http://" + host + "/devstoreaccount1/bucket/" + url.PathEscape(path.Join(repo, version, asset))
	})
}
