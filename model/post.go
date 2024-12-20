package model

import (
	"net/http"
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

func PostFromFormData(r http.Request) Post {
	return Post{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Url:         r.FormValue("url"),
	}
}

func (p Post) Validate() ValidationErrors {
	var errs ValidationErrors

	if p.Title == "" {
		errs = append(errs, ValidationError{Field: "title", Message: "El titulo es requerido"})
	}
	if p.Description == "" {
		errs = append(errs, ValidationError{Field: "description", Message: "La descripción es requerida"})
	}
	if p.Url == "" {
		errs = append(errs, ValidationError{Field: "url", Message: "La URL es requerida"})
	}

	// Comentado por mientras pq me da hueva hacer el crud de autores, atte: El Sebas
	/* if p.Author.ID == uuid.Nil {
		errs = append(errs, validator.ValidationError{Field: "author", Message: "El autor es requerido"})
	} */

	return errs
}
