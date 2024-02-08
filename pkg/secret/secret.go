package secret

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExists   = errors.New("key already exists")
	ErrEnvironment = errors.New("key set via environment variable")
)

func Get(key string) ([]byte, error) {
	if data := getEnv(key); data != "" {
		return []byte(data), nil
	}

	if path := getPathEnv(key); path != "" {
		return os.ReadFile(path)
	}

	path := filepath.Join(dir(), key)
	if _, err := os.Stat(path); err == nil {
		return os.ReadFile(path)
	}

	return nil, ErrKeyNotFound
}

func Put(key string, data []byte) error {
	if data := getEnv(key); data != "" {
		return ErrEnvironment
	}

	if path := getPathEnv(key); path != "" {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return ErrKeyExists
		}
		return os.WriteFile(path, data, 0o600)
	}

	path := filepath.Join(dir(), key)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return ErrKeyExists
	}
	return os.WriteFile(path, data, 0o600)
}

func Delete(key string) error {
	if data := getEnv(key); data != "" {
		return ErrEnvironment
	}

	if path := getPathEnv(key); path != "" {
		return os.Remove(path)
	}

	path := filepath.Join(dir(), key)
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}

	return ErrKeyNotFound
}

func getEnv(key string) string {
	return os.Getenv("KUBRI_" + strings.ToUpper(key))
}

func getPathEnv(key string) string {
	return os.Getenv("KUBRI_" + strings.ToUpper(key) + "_PATH")
}

func dir() string {
	if dir := os.Getenv("KUBRI_PATH"); dir != "" {
		return dir
	}
	dir, _ := os.UserConfigDir()
	dir = filepath.Join(dir, "kubri")
	_ = os.MkdirAll(dir, os.ModePerm)
	return dir
}
