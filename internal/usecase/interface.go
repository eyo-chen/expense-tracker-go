package usecase

import "github.com/OYE0303/expense-tracker-go/internal/domain"

type UserModel interface {
	// Create inserts a new user into the database.
	Create(name, email, passwordHash string, countryID int) error

	// FindByEmail returns a user by email.
	FindByEmail(email string) (*domain.User, error)
}
