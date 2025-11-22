package desc

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

	percent = []byte("%")
	nl      = []byte("\n")

	dateType = reflect.TypeFor[time.Time]()
)

func getFieldName(field reflect.StructField) string {
	if !field.IsExported() {
		return ""
	}

	name, _, _ := strings.Cut(field.Tag.Get("desc"), ",")
	if name == "" {
		return strings.ToUpper(field.Name)
	}
	if name[0] == '-' {
		return ""
	}

	return name
}

func atob(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func btoa(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
