package model

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type UserModel struct {
	DB *sql.DB
}

func newUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

type User struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password_hash string `json:"password_hash"`
}

// Create inserts a new user into the database.
func (m *UserModel) Create(name, email, passwordHash string) error {
	stmt := `INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`

	if _, err := m.DB.Exec(stmt, name, email, passwordHash); err != nil {
		logger.Error("users INSERT m.DB.Exec", "err", err)
		return err
	}

	return nil
}

// FindByEmail returns a user by email.
func (m *UserModel) FindByEmail(email string) (*domain.User, error) {
	stmt := `SELECT id, name, email, password_hash FROM users WHERE email = ?`

	var user User
	if err := m.DB.QueryRow(stmt, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password_hash); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		logger.Error("users SELECT m.DB.QueryRow", "err", err)
		return nil, err
	}

	return cvtToDomainUser(&user), nil
}

func cvtToDomainUser(u *User) *domain.User {
	return &domain.User{
		ID:            u.ID,
		Name:          u.Name,
		Email:         u.Email,
		Password_hash: u.Password_hash,
	}
}
