package apt

import "time"

func SetTime(t time.Time) {
	timeNow = func() time.Time { return t }
}

//nolint:gochecknoglobals
var (
	Compress              = compress
	Decompress            = decompress
	CompressionExtensions = compressionExtensions
)
