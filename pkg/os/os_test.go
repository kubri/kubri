package os_test

import (
	"errors"
	"testing"

	"github.com/abemedia/appcast/pkg/os"
)

func TestMarshalText(t *testing.T) {
	tests := []struct {
		in   os.OS
		want string
		err  error
	}{
		{os.MacOS, "macos", nil},
		{os.Windows, "windows", nil},
		{os.Windows64, "windows-x64", nil},
		{os.Windows32, "windows-x86", nil},
		{os.Unknown, "", nil},
		{99, "", os.ErrUnknownOS},
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

func TestUnmarshalText(t *testing.T) {
	tests := []struct {
		in   string
		want os.OS
		err  error
	}{
		{"macos", os.MacOS, nil},
		{"windows", os.Windows, nil},
		{"windows-x64", os.Windows64, nil},
		{"windows-x86", os.Windows32, nil},
		{"", os.Unknown, nil},
		{"foo", os.Unknown, os.ErrUnknownOS},
	}

	for _, test := range tests {
		var got os.OS
		err := got.UnmarshalText([]byte(test.in))
		if !errors.Is(err, test.err) {
			t.Errorf("%s: want error '%v' got '%v'", test.in, test.err, err)
		} else if got != test.want {
			t.Errorf("%s: want '%s' got '%s'", test.in, test.want, got)
		}
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		a, b os.OS
		want bool
	}{
		{os.MacOS, os.MacOS, true},
		{os.MacOS, os.Windows, false},
		{os.Windows32, os.Windows, true},
		{os.Windows64, os.Windows, true},
		{os.Unknown, os.MacOS, false},
	}

	for _, test := range tests {
		if got := os.Is(test.a, test.b); got != test.want {
			t.Errorf("%s == %s: want '%t' got '%t'", test.a, test.b, test.want, got)
		}
	}
}

func TestDetect(t *testing.T) {
	tests := []struct {
		in   string
		want os.OS
	}{
		{"test", os.Unknown},
		{"test.dmg", os.MacOS},
		{"test.pkg", os.MacOS},
		{"test.exe", os.Windows},
		{"test.msi", os.Windows},
		{"test.msix", os.Windows},
		{"test_32bit.exe", os.Windows32},
		{"test_x86.msi", os.Windows32},
		{"test_i386.msix", os.Windows32},
		{"test_ia32.msix", os.Windows32},
		{"test_64-bit.exe", os.Windows64},
		{"test_x86_64.msi", os.Windows64},
		{"test_x64.msi", os.Windows64},
		{"test_amd64.msix", os.Windows64},
	}

	for _, test := range tests {
		if got := os.Detect(test.in); got != test.want {
			t.Errorf("%s: want '%s' got '%s'", test.in, test.want, got)
		}
	}
}
