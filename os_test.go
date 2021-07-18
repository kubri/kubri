package appcast_test

import (
	"testing"

	"github.com/abemedia/appcast"
)

func TestDetectOS(t *testing.T) {
	tests := []struct {
		in   string
		want appcast.OS
	}{
		{"test", appcast.Unknown},
		{"test.dmg", appcast.MacOS},
		{"test.pkg", appcast.MacOS},
		{"test.exe", appcast.Windows},
		{"test.msi", appcast.Windows},
		{"test.msix", appcast.Windows},
		{"test_32bit.exe", appcast.Windows32},
		{"test_x86.msi", appcast.Windows32},
		{"test_i386.msix", appcast.Windows32},
		{"test_ia32.msix", appcast.Windows32},
		{"test_64-bit.exe", appcast.Windows64},
		{"test_x86_64.msi", appcast.Windows64},
		{"test_x64.msi", appcast.Windows64},
		{"test_amd64.msix", appcast.Windows64},
	}

	for _, test := range tests {
		if got := appcast.DetectOS(test.in); got != test.want {
			t.Errorf("%s: want '%s' got '%s'", test.in, test.want, got)
		}
	}
}
