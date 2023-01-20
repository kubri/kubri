package azureblob_test

import (
	"net/url"
	"path"
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	"github.com/abemedia/appcast/target/blob/azureblob"
	"github.com/abemedia/appcast/target/internal/test"
)

func TestAzureblob(t *testing.T) {
	host := emulator.AzureBlob(t, "bucket")

	tgt, err := azureblob.New(azureblob.Config{Bucket: "bucket", Folder: "folder"})
	if err != nil {
		t.Fatal(err)
	}

	test.Run(t, tgt, func(asset string) string {
		return "http://" + host + "/devstoreaccount1/bucket/" + url.PathEscape(path.Join("folder", asset))
	})
}
