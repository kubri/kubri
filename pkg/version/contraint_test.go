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
		{"", "1.0.0", true},
		{"*", "1.0.0", true},
		{"1", "1.0.0", true},
		{"v1", "1.0.0", true},
		{"1", "v1.0.0", true},
		{"1.0.0", "1.0.0", true},
		{"=1", "1.0.0", true},
		{" = 1 ", " 1.0.0 ", true},
		{"=1", "1.1.0", false},
		{"!=1", "1.0.0", false},
		{"!=1", "2.0.0", true},
		{">1", "1.0.0", false},
		{">1", "1.1.0-pre", true},
		{">1", "1.0.0-pre", false},
		{">1", "2.0.0", true},
		{"<2", "1.0.0", true},
		{"<2", "2.0.0", false},
		{">=1", "1.0.0", true},
		{">=1", "0.1.0", false},
		{">=1", "2.0.0", true},
		{"<=2", "2.0.0", true},
		{"<=2", "2.1.0", false},
		{"<=2", "1.0.0", true},
		{"~1", "1.0.0", true},
		{"~1", "1.1.0", true},
		{"~1", "2", false},
		{"~0.1.1", "0.2.0", false},
		{"~1.1.1", "1.2.0", false},
		{"~1.1", "1.0.0", false},
		{"~1.1", "1.1.0", true},
		{"~1.1", "1.1.5", true},
		{"~1.1", "1.2", false},
		{"~>1.1", "1.0.0", false},
		{"~>1.1", "1.1.0", true},
		{"~>1.1", "1.1.5", true},
		{"~>1.1", "1.2", false},
		{"~=1.1", "1.0.0", false},
		{"~=1.1", "1.1.0", true},
		{"~=1.1", "1.1.5", true},
		{"~=1.1", "1.2", false},
		{"^1.1", "1.0.0", false},
		{"^1.0", "1.1.0", true},
		{"^1.0", "2.0.0", false},
		{"^0.1", "0.2.0", false},
		{"1.*", "1.0.0", true},
		{"1.*", "1.1.0", true},
		{"1.*", "2.0.0", false},
		{">=1, <2", "1.0.0", true},
		{">=1, <2", "1.1.0", true},
		{">=1, <2", "2.0.0", false},
		{">=1, <2", "0.1.0", false},
	}

	for _, test := range tests {
		if test.want != version.Check(test.constraint, test.version) {
			t.Errorf("%q should return %t for %q", test.constraint, test.want, test.version)
		}
	}
}
