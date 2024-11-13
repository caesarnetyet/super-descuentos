package store

import (
	"super-descuentos/errs"
	"super-descuentos/model"
	"sync"

	"github.com/google/uuid"
)

type InMemoryStore struct {
	sync.RWMutex
	posts map[uuid.UUID]model.Post
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		posts: make(map[uuid.UUID]model.Post),
	}
}

func (s *InMemoryStore) CreatePost(post model.Post) error {
	s.Lock()
	defer s.Unlock()

	if post.ID == uuid.Nil {
		post.ID = uuid.New()
	}
	s.posts[post.ID] = post
	return nil
}

func (s *InMemoryStore) DeletePost(id uuid.UUID) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.posts[id]; !exists {
		return errs.ErrPostNotFound
	}
	delete(s.posts, id)
	return nil
}

func (s *InMemoryStore) UpdatePost(id uuid.UUID, post model.Post) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.posts[id]; !exists {
		return errs.ErrPostNotFound
	}
	post.ID = id // asegurarse de que el ID del post sea el mismo que el ID proporcionado
	s.posts[id] = post
	return nil
}

func (s *InMemoryStore) GetPost(id uuid.UUID) (model.Post, error) {
	s.RLock()
	defer s.RUnlock()

	post, exists := s.posts[id]
	if !exists {
		return model.Post{}, errs.ErrPostNotFound
	}
	return post, nil
}

func (s *InMemoryStore) GetPosts() ([]model.Post, error) {
	s.RLock()
	defer s.RUnlock()

	posts := make([]model.Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}
	return posts, nil
}
