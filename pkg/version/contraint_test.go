package version_test

import (
	"testing"

	"github.com/abemedia/appcast/pkg/version"
)

func TestConstraint(t *testing.T) {
	tests := []struct {
		constraint string
		version    string
		want       bool
	}{
		{"v1", "v1.0.0", true},
		{"1.0.0", "1.0.0", true},
		{"= v1", "v1.0.0", true},
		{"= v1", "v1.1.0", false},
		{"!= v1", "v1.0.0", false},
		{"!= v1", "v2.0.0", true},
		{"> v1", "v1.0.0", false},
		{"> v1", "v1.1.0-pre", true},
		{"> v1", "v1.0.0-pre", false},
		{"> v1", "v2.0.0", true},
		{"< v2", "v1.0.0", true},
		{"< v2", "v2.0.0", false},
		{">= v1", "v1.0.0", true},
		{">= v1", "v0.1.0", false},
		{">= v1", "v2.0.0", true},
		{"<= v2", "v2.0.0", true},
		{"<= v2", "v2.1.0", false},
		{"<= v2", "v1.0.0", true},
		{">= v1, < v2", "v1.0.0", true},
		{">= v1, < v2", "v1.1.0", true},
		{">= v1, < v2", "v2.0.0", false},
		{">= v1, < v2", "v0.1.0", false},
	}

	for _, test := range tests {
		if test.want != version.Check(test.constraint, test.version) {
			t.Errorf("%q should return %t for %q", test.constraint, test.want, test.version)
		}
	}
}
