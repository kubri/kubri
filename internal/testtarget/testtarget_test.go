package testtarget_test

import (
	"testing"

	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/internal/testtarget"
)

func TestFile(t *testing.T) {
	tgt := testtarget.New()
	test.Target(t, tgt, func(asset string) string {
		return "https://example.com/" + asset
	})
}
