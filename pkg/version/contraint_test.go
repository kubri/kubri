package version_test

import (
	"testing"

	"github.com/kubri/kubri/pkg/version"
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

		// TODO: Decide if leading/trailing commas should be a failure or ignored.
		{",!=0.1.0,", "0.1.0", false},
		{", ", "0.1.0", true},
	}

	for _, test := range tests {
		if _, err := version.NewConstraint(test.constraint); err != nil {
			t.Errorf("%q failed to parse: %s", test.constraint, err)
		}
		if test.want != version.Check(test.constraint, test.version) {
			t.Errorf("%q should return %t for %q", test.constraint, test.want, test.version)
		}
	}
}

func TestConstraintError(t *testing.T) {
	tests := []string{
		"a",
		"!1",
		"va",
		"1a",
		"=a",
		"v1.0.",
		"v1.0.a",
		"*1",
	}

	for _, v := range tests {
		_, err := version.NewConstraint(v)
		if err == nil {
			t.Errorf("%q should return error", v)
		}
		if version.Check(v, "1.0.0") {
			t.Errorf("%q should return false", v)
		}
	}
}
