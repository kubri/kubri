package source

import "errors"

var (
	ErrReleaseNotFound = errors.New("release not found")
	ErrAssetNotFound   = errors.New("asset not found")
)
