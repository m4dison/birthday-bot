package memory

import "github.com/m4dison/my-telegram-bot/internal/models"

type UserRepository interface {
	AddUser(user models.User) error
	GetAllUsers() []models.User
	// Другие методы, если нужно
}
