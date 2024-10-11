package test

import (
	"crypto/rsa"
	"log"
	"reflect"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/google/go-cmp/cmp"
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
