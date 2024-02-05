package test

import (
	"bytes"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
)

func YAML(s string) []byte {
	return []byte(strings.ReplaceAll(heredoc.Doc(s), "\t", "  "))
}

func JoinYAML(s ...string) []byte {
	b := make([][]byte, 0, len(s))
	for _, v := range s {
		b = append(b, YAML(v))
	}
	return bytes.Join(b, []byte("\n"))
}
