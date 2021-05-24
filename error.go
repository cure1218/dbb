package dbb

import (
	"errors"
)

var (
	ErrVerHandlerExists      error = errors.New("Version handler exists.")
	ErrNothingChanged        error = errors.New("Nothing changed.")
	ErrHandlerFuncAssertFail error = errors.New("Handler function assertion fail.")
)
