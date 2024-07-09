package postgres

import (
	"context"
	"database/sql"

	"github.com/m4dison/my-telegram-bot/internal/models"
	"github.com/m4dison/my-telegram-bot/internal/repository/memory"
)

var _ memory.UserRepository = &UserStore{}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(dataSourceName string) (*UserStore, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &UserStore{db: db}, nil
}

func (s *UserStore) AddUser(ctx context.Context, user models.User) error {
	query := `INSERT INTO users (name, birthday) VALUES ($1, $2)`
	_, err := s.db.ExecContext(ctx, query, user.Name, user.Birthday)
	return err
}

func (s *UserStore) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT name, birthday FROM users`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Name, &user.Birthday); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
