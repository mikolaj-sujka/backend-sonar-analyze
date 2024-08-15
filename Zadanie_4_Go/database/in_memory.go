package database

import (
	"go_app/models"
	"sync"
)

type UserStore struct {
	mu    sync.Mutex
	users map[string]*models.User 
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]*models.User),
	}
}

func (s *UserStore) AddUser(user *models.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
}

func (s *UserStore) GetUser(id string) (*models.User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, exists := s.users[id]
	return user, exists
}

func (s *UserStore) GetUserByEmail(email string) (*models.User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range s.users {
		if user.Email == email {
			return user, true
		}
	}
	return nil, false
}
