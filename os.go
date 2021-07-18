package appcast

import (
	"fmt"
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

func matchOS(constraint, os OS) bool {
	if constraint == os {
		return true
	}

	if os == Windows {
		return constraint > Windows
	}

	return false
}
