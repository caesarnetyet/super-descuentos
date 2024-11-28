package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"super-descuentos/model"
	"super-descuentos/rest"
	"super-descuentos/store"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetPosts(t *testing.T) {
	store := store.NewInMemoryStore()
	server := rest.NewServer(store)

	// Crear post de prueba
	post1 := model.Post{
		ID:           uuid.New(),
		Title:        "Test Post 1",
		Description:  "Description 1",
		Url:          "https://example.com/1",
		Author:       model.User{ID: uuid.New(), Name: "Author 1", Email: "author1@example.com"},
		CreationTime: time.Now(),
	}
	post2 := model.Post{
		ID:           uuid.New(),
		Title:        "Test Post 2",
		Description:  "Description 2",
		Url:          "https://example.com/2",
		Author:       model.User{ID: uuid.New(), Name: "Author 2", Email: "author2@example.com"},
		CreationTime: time.Now(),
	}

	ctx := context.Background()

	store.CreatePost(ctx, post1)
	store.CreatePost(ctx, post2)

	req := httptest.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)

	}

	var posts []model.Post
	json.NewDecoder(w.Body).Decode(&posts)

	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}
}

func TestCRUDOperations(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		setupFunc      func(*store.InMemoryStore) uuid.UUID
	}{
		{
			name:   "Create Post",
			method: "POST",
			path:   "/posts",
			body: model.Post{
				Title:       "New Post",
				Description: "Description",
				Url:         "wompwomp",
				Author:      model.User{ID: uuid.New()},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Get Post",
			method:         "GET",
			path:           "/posts/", // ID will be appended
			expectedStatus: http.StatusOK,
			setupFunc: func(store *store.InMemoryStore) uuid.UUID {
				post := model.Post{ID: uuid.New(), Title: "Test Post"}
				store.CreatePost(ctx, post)
				return post.ID
			},
		},
		{
			name:   "Update Post",
			method: "PUT",
			path:   "/posts/", // ID will be appended
			body: model.Post{
				Title:       "Updated Post",
				Description: "Updated Description",
				Url:         "wompwomp",
				Author:      model.User{ID: uuid.New()},
			},
			expectedStatus: http.StatusOK,
			setupFunc: func(store *store.InMemoryStore) uuid.UUID {
				post := model.Post{ID: uuid.New(), Title: "Original Post"}
				store.CreatePost(ctx, post)
				return post.ID
			},
		},
		{
			name:           "Delete Post",
			method:         "DELETE",
			path:           "/posts/", // ID will be appended
			expectedStatus: http.StatusNoContent,
			setupFunc: func(store *store.InMemoryStore) uuid.UUID {
				post := model.Post{ID: uuid.New(), Title: "To Delete"}
				store.CreatePost(ctx, post)
				return post.ID
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := store.NewInMemoryStore()
			server := rest.NewServer(store)

			var path string
			if test.setupFunc != nil {
				id := test.setupFunc(store)
				path = test.path + id.String()
			} else {
				path = test.path
			}

			var body bytes.Buffer
			if test.body != nil {
				json.NewEncoder(&body).Encode(test.body)
			}

			req := httptest.NewRequest(test.method, path, &body)
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			if w.Code != test.expectedStatus {
				var body interface{}
				json.NewDecoder(w.Body).Decode(&body)
				t.Errorf("Expected status %d, got %d. details: %q", test.expectedStatus, w.Code, body)
			}
		})
	}
}

func TestErrorCases(t *testing.T) {
	tests := []struct {
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

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store := store.NewInMemoryStore()
			server := rest.NewServer(store)

			var body bytes.Buffer
			if test.Body != nil {
				if jsonBytes, ok := test.Body.([]byte); ok {
					body = *bytes.NewBuffer(jsonBytes)
				} else {
					json.NewEncoder(&body).Encode(test.Body)
				}
			}

			req := httptest.NewRequest(test.Method, test.Path, &body)
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			if w.Code != test.ExpectedStatus {
				body := w.Body.String()
				t.Errorf("Expected status %d, got %d. details: %q", test.ExpectedStatus, w.Code, body)
			}

			if w.Body.String() != test.ExpectedBody {
				t.Errorf("Expected body %s, got %s", test.ExpectedBody, w.Body.String())
			}
		})
	}
}
