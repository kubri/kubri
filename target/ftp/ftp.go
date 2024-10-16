// Package ftp provides a target implementation for the File Transfer Protocol.
package ftp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/jlaffaye/ftp"

	"github.com/kubri/kubri/target"
)

type Config struct {
	Address string
	Folder  string
	URL     string
}

func New(c Config) (target.Target, error) {
	conn, err := ftp.Dial(c.Address)
	if err != nil {
		return nil, err
	}
	user := os.Getenv("FTP_USER")
	pass := os.Getenv("FTP_PASSWORD")
	if user != "" {
		if err := conn.Login(user, pass); err != nil {
			_ = conn.Quit()
			return nil, err
		}
	}
	return &ftpTarget{conn, c.Folder, c.URL}, nil
}

type ftpTarget struct {
	conn *ftp.ServerConn
	path string
	url  string
}

func (t *ftpTarget) NewWriter(_ context.Context, filename string) (io.WriteCloser, error) {
	return &fileWriter{t: t, path: path.Join(t.path, filename)}, nil
}

func (t *ftpTarget) NewReader(_ context.Context, filename string) (io.ReadCloser, error) {
	rd, err := t.conn.Retr(path.Join(t.path, filename))
	return rd, mapError("read", filename, err)
}

func (t *ftpTarget) Remove(_ context.Context, filename string) error {
	return mapError("delete", filename, t.conn.Delete(path.Join(t.path, filename)))
}

func (t *ftpTarget) Sub(dir string) target.Target {
	u, _ := url.JoinPath(t.url, dir)
	return &ftpTarget{conn: t.conn, path: path.Join(t.path, dir), url: u}
}

func (t *ftpTarget) URL(_ context.Context, filename string) (string, error) {
	return url.JoinPath(t.url, filename)
}

type fileWriter struct {
	bytes.Buffer

	t    *ftpTarget
	path string
}

func (w *fileWriter) Close() error {
	var base string
	dirs := strings.Split(path.Dir(w.path), "/")
	for _, dir := range dirs {
		base = path.Join(base, dir)
		_ = w.t.conn.MakeDir(base)
	}
	return mapError("write", w.path, w.t.conn.Stor(w.path, w))
}

func mapError(op, name string, err error) error {
	var tErr *textproto.Error
	if !errors.As(err, &tErr) {
		return err
	}

	switch tErr.Code {
	case ftp.StatusFileUnavailable:
		return &fs.PathError{Op: op, Path: name, Err: fs.ErrNotExist}
	default:
		return err
	}
}
