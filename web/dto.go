package web

import "super-descuentos/model"

type HandleCreatePostFormRequest struct {
	Title       string `json:"title" form:"title"`
	Content     string `json:"content" form:"content"`
	AuthorEmail string `json:"author_email" form:"author_email"`
	Url         string `json:"url" form:"url"`
}

func (r HandleCreatePostFormRequest) Validate() model.ValidationErrors {
	var errs model.ValidationErrors

	if r.Title == "" {
		errs = append(errs, model.ValidationError{
			Field:   "title",
			Message: "el título es requerido",
		})
	}

	if r.Content == "" {
		errs = append(errs, model.ValidationError{
			Field:   "content",
			Message: "el contenido es requerido",
		})
	}

	if r.AuthorEmail == "" {
		errs = append(errs, model.ValidationError{
			Field:   "author_email",
			Message: "el correo electrónico del autor es requerido",
		})
	}

	if r.Url == "" {
		errs = append(errs, model.ValidationError{
			Field:   "url",
			Message: "la URL es requerida",
		})
	}

	return errs
}

type HandleCreateAuthorFormRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r HandleCreateAuthorFormRequest) Validate() model.ValidationErrors {
	var problems model.ValidationErrors

	if r.Name == "" {
		problems = append(problems, model.ValidationError{
			Field:   "name",
			Message: "el nombre es requerido",
		})
	}

	if r.Email == "" {
		problems = append(problems, model.ValidationError{
			Field:   "email",
			Message: "el correo electrónico es requerido",
		})
	}

	return problems
}
