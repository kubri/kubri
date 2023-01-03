package azureblob_test

import (
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	_ "github.com/abemedia/appcast/target/blob/azureblob"
	"github.com/abemedia/appcast/target/internal/test"
)

func TestAzureblob(t *testing.T) {
	emulator.AzureBlob(t, "bucket")
	test.Run(t, "azblob://bucket/folder")
}
