// Package testtarget is an in-memory simulator of a target for use in tests.
package testtarget

import (
	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/target"
	_ "gocloud.dev/blob/memblob" // blob driver
)

func New() target.Target {
	t, _ := blob.NewTarget("mem://", "", "https://example.com/")
	return t
}
