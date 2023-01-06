// Package memory is an in-memory simulator of a blob source for use in tests.
package memory

import (
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
	"github.com/abemedia/appcast/target/blob/internal/blob"
	_ "gocloud.dev/blob/memblob" // blob driver
)

func New(c source.Config) (target.Target, error) {
	return blob.New("mem://", "")
}

//nolint:gochecknoinits
func init() { target.Register("mem", New) }
