package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"super-descuentos/errs"
)

// ValidationError representa un error de validación con detalles
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors es una colección de errores de validación
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var errMsgs []string
	for _, err := range e {
		errMsgs = append(errMsgs, err.Error())
	}
	return strings.Join(errMsgs, "; ")
}

// Validator es una interfaz que deben implementar las estructuras que requieren validación
type Validator interface {
	Validate() ValidationErrors
}

// DecodeAndValidate es una función genérica que decodifica y valida una petición HTTP
// T debe ser un tipo que implemente Validator
func DecodeAndValidate[T Validator](r *http.Request) (T, error) {
	var payload T

	// Verificar Content-Type
	// contentType := r.Header.Get("Content-Type")
	// if contentType != "application/json" {
	// 	return payload, ValidationError{
	// 		Field:   "Content-Type",
	// 		Message: "el tipo de contenido debe ser application/json",
	// 	}
	// }

	// Decodificar el cuerpo de la petición
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return payload, ValidationError{
			Field:   "body",
			Message: errs.ErrInvalidJSON.Error(),
		}
	}
	defer r.Body.Close()

	// Validar la estructura
	if errs := payload.Validate(); len(errs) > 0 {
		return payload, errs
	}

	return payload, nil
}
