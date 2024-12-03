package model

import (
	"net/http"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func UserFromFormData(r http.Request) User {
	return User{
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
	}
}

func (user User) Validate() ValidationErrors {
	var errs ValidationErrors

	if user.Name == "" {
		errs = append(errs, ValidationError{
			Field:   "name",
			Message: "el nombre es requerido",
		})
	}

	if user.Email == "" {
		errs = append(errs, ValidationError{
			Field:   "email",
			Message: "el correo electrónico es requerido",
		})
	}

	return errs
}
