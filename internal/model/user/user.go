package user

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type UserModel struct {
	DB *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

type User struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password_hash string `json:"password_hash"`
}

func (m *UserModel) Create(name, email, passwordHash string) error {
	stmt := `INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`

	if _, err := m.DB.Exec(stmt, name, email, passwordHash); err != nil {
		logger.Error("users INSERT m.DB.Exec", "err", err)
		return err
	}

	return nil
}

func (m *UserModel) FindByEmail(email string) (domain.User, error) {
	stmt := `SELECT id, name, email, password_hash FROM users WHERE email = ?`

	var user User
	if err := m.DB.QueryRow(stmt, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password_hash); err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrEmailNotFound
		}

		logger.Error("users SELECT m.DB.QueryRow", "err", err)
		return domain.User{}, err
	}

	return cvtToDomainUser(user), nil
}

func (m *UserModel) GetInfo(userID int64) (domain.User, error) {
	stmt := `SELECT id, name, email FROM users WHERE id = ?`

	var user User
	if err := m.DB.QueryRow(stmt, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Password_hash); err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrUserIDNotFound
		}

		logger.Error("users SELECT m.DB.QueryRow", "err", err)
		return domain.User{}, err
	}

	return cvtToDomainUser(user), nil
}

func cvtToDomainUser(u User) domain.User {
	return domain.User{
		ID:            u.ID,
		Name:          u.Name,
		Email:         u.Email,
		Password_hash: u.Password_hash,
	}
}
