package errs

import "errors"

var (
	ErrAuthorNotFound = errors.New("el autor no existe en el sistema")
	ErrPostNotFound   = errors.New("post no encontrado")
	ErrInvalidID      = errors.New("id inválido")
	ErrPostExists     = errors.New("el post ya existe")
)
