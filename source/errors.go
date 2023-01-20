package source

import "errors"

var (
	ErrMissingSource  = errors.New("missing source")
	ErrNoReleaseFound = errors.New("no release found")
	ErrAssetNotFound  = errors.New("asset not found")
)
