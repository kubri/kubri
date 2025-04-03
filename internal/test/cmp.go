package test

import (
	"crypto/rsa"
	"log"
	"path"
	"reflect"
	"testing/fstest"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// CompareErrorMessages compares errors by message.
func CompareErrorMessages() cmp.Option {
	return cmp.FilterValues(
		func(x, y any) bool {
			_, ok1 := x.(error)
			_, ok2 := y.(error)
			return ok1 && ok2
		},
		cmp.Comparer(func(a, b any) bool {
			if a == nil || b == nil {
				return a == b
			}
			return a.(error).Error() == b.(error).Error() //nolint:forcetypeassert
		}),
	)
}

// CompareLoggers compares instances of [log.Logger].
func CompareLoggers() cmp.Option {
	return cmp.Comparer(func(a, b *log.Logger) bool {
		return a.Prefix() == b.Prefix() && a.Flags() == b.Flags() && a.Writer() == b.Writer()
	})
}

// ComparePGPKeys compares PGP keys' fingerprints.
func ComparePGPKeys() cmp.Option {
	return cmp.Comparer(func(a, b *crypto.Key) bool {
		if a == nil || b == nil {
			return a == b
		}
		return a.GetFingerprint() == b.GetFingerprint()
	})
}

// CompareRSAPrivateKeys compares RSA private keys.
func CompareRSAPrivateKeys() cmp.Option {
	return cmp.Comparer(func(a, b *rsa.PrivateKey) bool {
		if a == nil || b == nil {
			return a == b
		}
		return a.Equal(b)
	})
}

// ExportAll exports all unexported fields.
func ExportAll() cmp.Option {
	return cmp.Exporter(func(reflect.Type) bool {
		return true
	})
}

// IgnoreFunctions ignores all functions.
func IgnoreFunctions() cmp.Option {
	return cmp.FilterPath(func(p cmp.Path) bool {
		return p.Last().Type().Kind() == reflect.Func
	}, cmp.Ignore())
}

// IgnoreKeys ignores keys in a map that match the given patterns.
func IgnoreKeys(patterns ...string) cmp.Option {
	return cmpopts.IgnoreMapEntries(func(k string, _ any) bool {
		for _, pattern := range patterns {
			if ok, _ := path.Match(pattern, k); ok {
				return true
			}
		}
		return false
	})
}

// IgnoreFSMeta ignores the metadata of a file in a [fstest.MapFS].
func IgnoreFSMeta() cmp.Option {
	return cmpopts.IgnoreFields(fstest.MapFile{}, "ModTime", "Sys")
}

// CompareFSStrict compares file systems by transforming them into [fstest.MapFS].
func CompareFSStrict() cmp.Option {
	return cmp.Transformer("fs.FS", ReadFS)
}

// CompareFS compares file systems by transforming them into [fstest.MapFS].
// It ignores the modification time and underlying data source.
func CompareFS() cmp.Option {
	return cmp.Options{CompareFSStrict(), IgnoreFSMeta()}
}
