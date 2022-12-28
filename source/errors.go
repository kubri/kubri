package source

import "errors"

var (
	ErrMissingSource   = errors.New("missing source")
	ErrReleaseNotFound = errors.New("release not found")
	ErrAssetNotFound   = errors.New("asset not found")
)
