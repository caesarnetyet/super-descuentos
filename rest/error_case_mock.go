package rest

import (
	"net/http"

	"github.com/google/uuid"
)

var ErrCaseMock = []struct {
	Name           string
	Method         string
	Path           string
	Body           interface{}
	ExpectedStatus int
	ExpectedBody   string
}{
	{
		Name:           "Get Non-existent Post",
		Method:         "GET",
		Path:           "/posts/" + uuid.New().String(),
		ExpectedStatus: http.StatusNotFound,
		ExpectedBody:   "{\"message\":\"post no encontrado\"}\n",
	},
	{
		Name:           "Invalid UUID",
		Method:         "GET",
		Path:           "/posts/invalid-uuid",
		ExpectedStatus: http.StatusBadRequest,
		ExpectedBody:   "{\"message\":\"id inválido\"}\n",
	},
	{
		Name:           "Invalid JSON",
		Method:         "POST",
		Path:           "/posts",
		Body:           []byte(`{"invalid json"`),
		ExpectedStatus: http.StatusBadRequest,
		ExpectedBody:   "{\"message\":\"JSON inválido\"}\n",
	},
}
