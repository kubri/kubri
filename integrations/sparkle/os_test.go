package sparkle_test

import (
	"errors"
	"testing"

	"github.com/kubri/kubri/integrations/sparkle"
)

func TestOS_MarshalText(t *testing.T) {
	tests := []struct {
		in   sparkle.OS
		want string
		err  error
	}{
		{sparkle.MacOS, "macos", nil},
		{sparkle.Windows, "windows", nil},
		{sparkle.Windows32, "windows-x86", nil},
		{sparkle.Windows64, "windows-x64", nil},
		{sparkle.WindowsARM64, "windows-arm64", nil},
		{sparkle.Unknown, "", nil},
		{255, "", sparkle.ErrUnknownOS},
	}

	for _, test := range tests {
		got, err := test.in.MarshalText()
		if !errors.Is(err, test.err) {
			t.Errorf("%T(%d): want error '%v' got '%v'", test.in, test.in, test.err, err)
		} else if string(got) != test.want {
			t.Errorf("%T(%d): want '%s' got '%s'", test.in, test.in, test.want, got)
		}
	}
}

func TestOS_UnmarshalText(t *testing.T) {
	tests := []struct {
		in   string
		want sparkle.OS
		err  error
	}{
		{"macos", sparkle.MacOS, nil},
		{"windows", sparkle.Windows, nil},
		{"windows-x86", sparkle.Windows32, nil},
		{"windows-x64", sparkle.Windows64, nil},
		{"windows-arm64", sparkle.WindowsARM64, nil},
		{"", sparkle.Unknown, nil},
		{"foo", sparkle.Unknown, sparkle.ErrUnknownOS},
	}

	for _, test := range tests {
		var got sparkle.OS
		err := got.UnmarshalText([]byte(test.in))
		if !errors.Is(err, test.err) {
			t.Errorf("%s: want error '%v' got '%v'", test.in, test.err, err)
		} else if got != test.want {
			t.Errorf("%s: want '%s' got '%s'", test.in, test.want, got)
		}
	}
}

func TestIsOS(t *testing.T) {
	tests := []struct {
		a, b sparkle.OS
		want bool
	}{
		{sparkle.MacOS, sparkle.MacOS, true},
		{sparkle.MacOS, sparkle.Windows, false},
		{sparkle.Windows32, sparkle.Windows, true},
		{sparkle.Windows64, sparkle.Windows, true},
		{sparkle.WindowsARM64, sparkle.Windows, true},
		{sparkle.Unknown, sparkle.MacOS, false},
		{sparkle.MacOS, sparkle.Unknown, true},
		{255, sparkle.Windows, false},
	}

	for _, test := range tests {
		if got := sparkle.IsOS(test.a, test.b); got != test.want {
			t.Errorf("%s == %s: want '%t' got '%t'", test.a, test.b, test.want, got)
		}
	}
}

func TestDetectOS(t *testing.T) {
	tests := []struct {
		in   string
		want sparkle.OS
	}{
		{"test", sparkle.Unknown},
		{"test.dmg", sparkle.MacOS},
		{"test.pkg", sparkle.MacOS},
		{"test.exe", sparkle.Windows},
		{"test.msi", sparkle.Windows},
		{"test_i386_amd64.msi", sparkle.Windows}, // Ambiguous
		{"test_32bit.exe", sparkle.Windows32},
		{"test_x86.msi", sparkle.Windows32},
		{"test_i386.msi", sparkle.Windows32},
		{"test_i686.msi", sparkle.Windows32},
		{"test_ia32.msi", sparkle.Windows32},
		{"test_64-bit.exe", sparkle.Windows64},
		{"test_x86_64.msi", sparkle.Windows64},
		{"test_x64.msi", sparkle.Windows64},
		{"test_amd64.msi", sparkle.Windows64},
		{"test_arm64.msi", sparkle.WindowsARM64},
		{"test_aarch64.msi", sparkle.WindowsARM64},
	}

	for _, test := range tests {
		if got := sparkle.DetectOS(test.in); got != test.want {
			t.Errorf("%s: want '%s' got '%s'", test.in, test.want, got)
		}
	}
}
