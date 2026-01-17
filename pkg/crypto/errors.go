package crypto //nolint:revive // var-naming

import "errors"

var (
	ErrWrongKeyType = errors.New("wrong key type")
	ErrInvalidKey   = errors.New("invalid key")
)
