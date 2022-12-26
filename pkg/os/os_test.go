package os_test

import (
	"testing"

	"github.com/abemedia/appcast/pkg/os"
)

func TestDetectOS(t *testing.T) {
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
