package deb

import (
	"bytes"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"
)

//nolint:gochecknoglobals
var (
	bufPool = sync.Pool{New: func() any { return bytes.NewBuffer(make([]byte, 1024)) }}

	colon = []byte(":")
	space = []byte(" ")
	nl    = []byte("\n")

	dateType = reflect.TypeOf((*time.Time)(nil)).Elem()
)

func getFieldName(field reflect.StructField) string {
	if !field.IsExported() {
		return ""
	}

	name, _, _ := strings.Cut(field.Tag.Get("deb"), ",")
	if name == "" {
		return field.Name
	}
	if name[0] == '-' {
		return ""
	}

	return name
}

func trim(s []byte) []byte {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	n := len(s)
	for n > i && (s[n-1] == ' ' || s[n-1] == '\t' || s[n-1] == '\n' || s[n-1] == '\r') {
		n--
	}
	return s[i:n]
}

func btoa(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
