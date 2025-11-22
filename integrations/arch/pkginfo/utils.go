package pkginfo

import (
	"reflect"
	"strings"
	"time"
	"unsafe"
)

//nolint:gochecknoglobals
var (
	sep   = []byte(" = ")
	equal = []byte("=")

	dateType = reflect.TypeFor[time.Time]()
)

func getFieldName(field reflect.StructField) string {
	if !field.IsExported() {
		return ""
	}

	name, _, _ := strings.Cut(field.Tag.Get("pkginfo"), ",")
	if name == "" {
		return strings.ToLower(field.Name)
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
