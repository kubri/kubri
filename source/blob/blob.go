package blob

import (
	// Import all common blob providers to register them.
	_ "github.com/abemedia/appcast/source/blob/azureblob"
	_ "github.com/abemedia/appcast/source/blob/file"
	_ "github.com/abemedia/appcast/source/blob/gcs"
	_ "github.com/abemedia/appcast/source/blob/s3"
)
