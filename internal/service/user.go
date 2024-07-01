package service

import (
	"log"
	"sync"
	"time"

	"github.com/m4dison/my-telegram-bot/internal/models"
	"github.com/m4dison/my-telegram-bot/internal/repository/memory"
)

// Add user in repo

// Get user by name
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

func (s *UserService) AddUser(user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.repo.AddUser(user)
}

func (s *UserService) CheckBirthdays() []models.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	users := s.repo.GetAllUsers()
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
	return usersWithBirthday
}
