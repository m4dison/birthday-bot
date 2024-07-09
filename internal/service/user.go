package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/m4dison/my-telegram-bot/internal/models"
	"github.com/m4dison/my-telegram-bot/internal/repository/memory"
)

// gomock uber
// go:generate
// go generate ./...
type UserService struct {
	repo memory.UserRepository
	mu   *sync.Mutex
}

func NewUserService(repo memory.UserRepository, mu *sync.Mutex) *UserService {
	return &UserService{
		repo: repo,
		mu:   mu,
	}
}

func (s *UserService) AddUser(ctx context.Context, user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.repo.AddUser(ctx, user)
}

func (s *UserService) CheckBirthdays(ctx context.Context) ([]models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	today := time.Now()
	log.Printf("Todays time looks like this %v", today)

	var usersWithBirthday []models.User
	for _, user := range users {
		log.Printf("Checking birthday for user: %s, Birthday: %s", user.Name, user.Birthday.Format("2006-01-02"))
		if user.Birthday.Month() == today.Month() && user.Birthday.Day() == today.Day() {
			log.Printf("User %s has a birthday today!", user.Name)
			usersWithBirthday = append(usersWithBirthday, user)
		}
	}
	return usersWithBirthday, nil
}
