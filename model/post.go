package model

import (
	"super-descuentos/validator"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Url          string    `json:"url"`
	Author       User      `json:"author"`
	Likes        int       `json:"likes"`
	ExpireTime   time.Time `json:"expire_time"`
	CreationTime time.Time `json:"creation_time"`
}

func (p Post) Validate() validator.ValidationErrors {
	var errs validator.ValidationErrors

	if p.Title == "" {
		errs = append(errs, validator.ValidationError{Field: "title", Message: "El titulo es requerido"})
	}
	if p.Description == "" {
		errs = append(errs, validator.ValidationError{Field: "description", Message: "La descripción es requerida"})
	}
	if p.Url == "" {
		errs = append(errs, validator.ValidationError{Field: "url", Message: "La URL es requerida"})
	}
	if p.Author.ID == uuid.Nil {
		errs = append(errs, validator.ValidationError{Field: "author", Message: "El autor es requerido"})
	}

	return errs
}
