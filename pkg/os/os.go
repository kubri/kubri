package os

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dlclark/regexp2"
)

type OS uint8

const (
	Unknown OS = iota
	MacOS
	Windows
	Windows64
	Windows32
)

func (os OS) String() string {
	switch os {
	default:
		return ""
	case MacOS:
		return "macos"
	case Windows:
		return "windows"
	case Windows64:
		return "windows-x64"
	case Windows32:
		return "windows-x86"
	}
}

func (os OS) MarshalText() ([]byte, error) {
	return []byte(os.String()), nil
}

func (os *OS) UnmarshalText(text []byte) error {
	s := string(text)
	for i := 0; i < 5; i++ {
		if OS(i).String() == s {
			*os = OS(i)
			return nil
		}
	}
	return fmt.Errorf("unknown os: %s", text)
}

func Is(constraint, os OS) bool {
	if constraint == os {
		return true
	}

	if os == Windows {
		return constraint > Windows
	}

	return false
}

//nolint:gochecknoglobals
var (
	reWin64 = regexp2.MustCompile(`amd64|x64|x86[\W_]?64|64[\W_]?bit`, regexp2.IgnoreCase)
	reWin32 = regexp2.MustCompile(`386|x86(?![\W_]?64)|ia32|32[\W_]?bit`, regexp2.IgnoreCase)
)

func Detect(name string) OS {
	s := strings.ToLower(name)
	ext := filepath.Ext(s)
	switch ext {
	case "":
	case ".dmg", ".pkg", ".mpkg":
		return MacOS
	case ".exe", ".msi", ".msix", ".msixbundle", ".appx", ".appxbundle", ".appinstaller":
		is64, _ := reWin64.MatchString(s)
		is32, _ := reWin32.MatchString(s)
		switch {
		case is64 == is32:
		case is64:
			return Windows64
		case is32:
			return Windows32
		}
		return Windows
	}

	return Unknown
}
