package pkginfo_test

import "time"

type marshaler struct {
	S string
}

func (s *marshaler) UnmarshalText(text []byte) error {
	s.S = string(text)
	return nil
}

type errMarshaler struct {
	E error
}

func (s *errMarshaler) UnmarshalText([]byte) error {
	return s.E
}

type record struct {
	String    string
	Pointer   *string
	Struct    recordStruct
	Slice     []string
	Hex       [4]byte
	Int       int
	Int8      int8
	Int16     int16
	Int32     int32
	Int64     int64
	Uint      uint
	Uint8     uint8
	Uint16    uint16
	Uint32    uint32
	Uint64    uint64
	Float32   float32
	Float64   float64
	Bool      bool
	Marshaler *marshaler
	Date      time.Time
}

type recordStruct struct {
	unexported string //nolint:unused
	Ignored    string `pkginfo:"-"`
	String     string
	Int        int
	Bool       bool
}
