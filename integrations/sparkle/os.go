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
	Windows64
	Windows32
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
	case Windows64:
		return "windows-x64"
	case Windows32:
		return "windows-x86"
	}
}

func (os OS) MarshalText() ([]byte, error) {
	if os > Windows32 {
		return nil, ErrUnknownOS
	}
	return []byte(os.String()), nil
}

func (os *OS) UnmarshalText(text []byte) error {
	s := string(text)
	for i := Unknown; i <= Windows32; i++ {
		if i.String() == s {
			*os = i
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrUnknownOS, text)
}

func IsOS(os, target OS) bool {
	switch target {
	case os, Unknown:
		return true
	case Windows:
		return os > Windows
	}
	return false
}

//nolint:gochecknoglobals
var (
	reWin64 = regexp2.MustCompile(`amd64|x64|x86[\W_]?64|64[\W_]?bit`, regexp2.IgnoreCase)
	reWin32 = regexp2.MustCompile(`386|x86(?![\W_]?64)|ia32|32[\W_]?bit`, regexp2.IgnoreCase)
)

func DetectOS(name string) OS {
	ext := path.Ext(name)
	switch ext {
	case "":
	case ".dmg", ".pkg", ".mpkg":
		return MacOS
	case ".exe", ".msi":
		is64, _ := reWin64.MatchString(name)
		is32, _ := reWin32.MatchString(name)
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
