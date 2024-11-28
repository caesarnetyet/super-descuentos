package errs

import "errors"

var (
	ErrInvalidJSON        = errors.New("JSON inválido")
	ErrInvalidQueryParams = errors.New("query params inválidos")
)
