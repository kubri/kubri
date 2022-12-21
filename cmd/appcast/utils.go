package main

import (
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"

	"github.com/abemedia/appcast"
	"github.com/abemedia/appcast/source"
	"gopkg.in/yaml.v3"
)

func getDir() string {
	if dir := os.Getenv("APPCAST_CONFIG_PATH"); dir != "" {
		return dir
	}
	dir, _ := os.UserConfigDir()
	dir = filepath.Join(dir, "appcast")
	_ = os.MkdirAll(dir, os.ModePerm)
	return dir
}

func readConfig(path string) (*appcast.Config, error) {
	c := &appcast.Config{}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}

func parseSource(src string) (*source.Source, error) {
	if !strings.Contains(src, "://") {
		src = "local://" + src
	}
	s := &source.Source{}
	if err := s.UnmarshalText([]byte(src)); err != nil {
		return nil, err
	}
	return s, nil
}

func readKey[T any](path string, unmarshaler func([]byte) (T, error)) (T, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		var zero T
		return zero, err
	}
	block, _ := pem.Decode(b)
	return unmarshaler(block.Bytes)
}
