package azureblob_test

import (
	"net/url"
	"path"
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	"github.com/abemedia/appcast/source/blob/azureblob"
	"github.com/abemedia/appcast/source/blob/internal/test"
)

func TestAzureblob(t *testing.T) {
	host := emulator.AzureBlob(t, "bucket")

	s, err := azureblob.New(azureblob.Config{Bucket: "bucket", Folder: "folder"})
	if err != nil {
		t.Fatal(err)
	}

	test.Run(t, s, func(version, asset string) string {
		return "http://" + host + "/devstoreaccount1/bucket/" + url.PathEscape(path.Join("folder", version, asset))
	})
}
