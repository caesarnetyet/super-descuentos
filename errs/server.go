package errs

import "errors"

var (
	ErrInvalidID          = errors.New("id inválido")
	ErrInvalidJSON        = errors.New("JSON inválido")
	ErrPostNotFound       = errors.New("post no encontrado")
	ErrInvalidQueryParams = errors.New("query params inválidos")
	ErrInvalidPost        = errors.New("post inválido")
	ErrPostAlreadyExists  = errors.New("el post ya existe")
	ErrInternalError      = errors.New("error interno del servidor")
)
