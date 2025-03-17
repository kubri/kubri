package sparkle

import (
	"errors"
	"fmt"
	"path"

	"github.com/dlclark/regexp2"
)

type OS uint8

const (
	Unknown OS = iota
	MacOS
	Windows
	Windows32
	Windows64
	WindowsARM64
)

var ErrUnknownOS = errors.New("unknown os")

func (os OS) String() string {
	switch os {
	default:
		return ""
	case MacOS:
		return "macos"
	case Windows:
		return "windows"
	case Windows32:
		return "windows-x86"
	case Windows64:
		return "windows-x64"
	case WindowsARM64:
		return "windows-arm64"
	}
}

func (os OS) MarshalText() ([]byte, error) {
	if os > WindowsARM64 {
		return nil, ErrUnknownOS
	}
	return []byte(os.String()), nil
}

func (os *OS) UnmarshalText(text []byte) error {
	s := string(text)
	for i := Unknown; i <= WindowsARM64; i++ {
		if i.String() == s {
			*os = i
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrUnknownOS, text)
}

func isOS(os, target OS) bool {
	switch target {
	case os, Unknown:
		return true
	case Windows:
		return os >= Windows && os <= WindowsARM64
	}
	return false
}

//nolint:gochecknoglobals
var (
	reWin32    = regexp2.MustCompile(`386|686|x86(?![\W_]?64)|ia32|32[\W_]?bit`, regexp2.IgnoreCase)
	reWin64    = regexp2.MustCompile(`amd64|x64|x86[\W_]?64|64[\W_]?bit`, regexp2.IgnoreCase)
	reWinARM64 = regexp2.MustCompile(`arm64|aarch64|a64`, regexp2.IgnoreCase)
)

func detectOS(name string) OS {
	ext := path.Ext(name)
	switch ext {
	case "":
	case ".dmg", ".pkg", ".mpkg":
		return MacOS
	case ".exe", ".msi":
		is32, _ := reWin32.MatchString(name)
		is64, _ := reWin64.MatchString(name)
		isARM64, _ := reWinARM64.MatchString(name)

		var matched int
		if is32 {
			matched++
		}
		if is64 {
			matched++
		}
		if isARM64 {
			matched++
		}

		switch {
		case matched > 1: // Ambiguous case.
		case is32:
			return Windows32
		case is64:
			return Windows64
		case isARM64:
			return WindowsARM64
		}
		return Windows
	}

	return Unknown
}
