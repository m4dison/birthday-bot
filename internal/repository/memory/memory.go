package memory

import (
	"sync"

	"github.com/m4dison/my-telegram-bot/internal/models"
)

type UserStore struct {
	mu    sync.Mutex
	users map[string]models.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]models.User),
	}
}

func (s *UserStore) AddUser(user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.Name] = user
	return nil
}

func (s *UserStore) GetAllUsers() []models.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	users := make([]models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}
