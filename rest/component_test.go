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
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Get Non-existent Post",
			method:         "GET",
			path:           "/posts/" + uuid.New().String(),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "{\"message\":\"post no encontrado\"}\n",
		},
		{
			name:           "Invalid UUID",
			method:         "GET",
			path:           "/posts/invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"message\":\"id inválido\"}\n",
		},
		{
			name:           "Invalid JSON",
			method:         "POST",
			path:           "/posts",
			body:           []byte(`{"invalid json"`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"message\":\"JSON inválido\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := store.NewInMemoryStore()
			server := rest.NewServer(store)

			var body bytes.Buffer
			if tt.body != nil {
				if jsonBytes, ok := tt.body.([]byte); ok {
					body = *bytes.NewBuffer(jsonBytes)
				} else {
					json.NewEncoder(&body).Encode(tt.body)
				}
			}

			req := httptest.NewRequest(tt.method, tt.path, &body)
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				body := w.Body.String()
				t.Errorf("Expected status %d, got %d. details: %q", tt.expectedStatus, w.Code, body)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %s, got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}
