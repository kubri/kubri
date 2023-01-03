// Package memory is an in-memory simulator of a blob source for use in tests.
package memory

import (
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/internal/blob"
	_ "gocloud.dev/blob/memblob" // blob driver
)

func New(c source.Config) (*source.Source, error) {
	return blob.New("mem://", "", "mem:/")
}

//nolint:gochecknoinits
func init() { source.Register("mem", New) }
