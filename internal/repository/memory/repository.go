package memory

import (
	"context"

	"github.com/m4dison/my-telegram-bot/internal/models"
)

// user, repo: add updating ...
// user, repo: remove some feature
// linus torvalds unix github
type UserRepository interface {
	AddUser(ctx context.Context, user models.User) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
}
