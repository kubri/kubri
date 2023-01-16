package azureblob_test

import (
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	"github.com/abemedia/appcast/target/blob/azureblob"
	"github.com/abemedia/appcast/target/internal/test"
)

func TestAzureblob(t *testing.T) {
	emulator.AzureBlob(t, "bucket")

	tgt, err := azureblob.New(azureblob.Config{Bucket: "bucket", Folder: "folder"})
	if err != nil {
		t.Fatal(err)
	}

	test.Run(t, tgt)
}
