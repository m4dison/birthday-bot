package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/m4dison/my-telegram-bot/internal/models"
	"github.com/m4dison/my-telegram-bot/internal/service"
)

var _ service.UserRepository = &UserStore{}

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
	// Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}

	// Вставляем данные в таблицу users
	query := `INSERT INTO users (name, birthday) VALUES ($1, $2)`
	_, err = tx.ExecContext(ctx, query, user.Name, user.Birthday)
	if err != nil {
		tx.Rollback() // Если произошла ошибка, откатываем транзакцию
		return fmt.Errorf("could not insert user: %v", err)
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
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
