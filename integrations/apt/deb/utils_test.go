package deb_test

import "time"

type stringer struct {
	S string
}

func (s stringer) String() string {
	return s.S
}

type marshaler struct {
	S string
}

func (s *marshaler) MarshalText() ([]byte, error) {
	return []byte(s.S), nil
}

func (s *marshaler) UnmarshalText(text []byte) error {
	s.S = string(text)
	return nil
}

type errMarshaler struct {
	E error
}

func (s errMarshaler) MarshalText() ([]byte, error) {
	return nil, s.E
}

func (s *errMarshaler) UnmarshalText([]byte) error {
	return s.E
}

type record struct {
	String    string
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
	Stringer  stringer
	Marshaler *marshaler
	Date      time.Time
}
