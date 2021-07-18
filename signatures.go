package appcast

import (
	"bytes"
	"encoding"
	"encoding/csv"
	"io"
)

type signatures map[[2]string]string

func (s signatures) Get(filename, algorithm string) string {
	return s[[2]string{filename, algorithm}]
}

func (s signatures) Set(filename, algorithm, signature string) {
	s[[2]string{filename, algorithm}] = signature
}

func (s signatures) UnmarshalText(b []byte) error {
	rd := csv.NewReader(bytes.NewReader(b))
	rd.Comma = '\t'
	rd.ReuseRecord = true
	rd.FieldsPerRecord = 3

	for {
		record, err := rd.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		s[[2]string{record[0], record[1]}] = record[2]
	}
}

func (s signatures) MarshalText() ([]byte, error) {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	w.Comma = '\t'

	for k, v := range s {
		err := w.Write([]string{k[0], k[1], v})
		if err != nil {
			return nil, err
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var (
	_ encoding.TextMarshaler   = (signatures)(nil)
	_ encoding.TextUnmarshaler = (signatures)(nil)
)
