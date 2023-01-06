package blob

import (
	// Import all common blob providers to register them.
	_ "github.com/abemedia/appcast/target/blob/azureblob"
	_ "github.com/abemedia/appcast/target/blob/file"
	_ "github.com/abemedia/appcast/target/blob/gcs"
	_ "github.com/abemedia/appcast/target/blob/s3"
)
