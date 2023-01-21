// Package memory is an in-memory simulator of a blob source for use in tests.
package memory

import (
	"github.com/abemedia/appcast/target"
	"github.com/abemedia/appcast/target/blob/internal/blob"
	_ "gocloud.dev/blob/memblob" // blob driver
)

type Config struct{}

func New(c Config) (target.Target, error) {
	return blob.New("mem://", "", "mem:/")
}
