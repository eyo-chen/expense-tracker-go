package usecase

import "github.com/OYE0303/expense-tracker-go/internal/domain"

// UserModel is the interface that wraps the basic methods for user model.
type UserModel interface {
	// Create inserts a new user into the database.
	Create(name, email, passwordHash string) error

	// FindByEmail returns a user by email.
	FindByEmail(email string) (*domain.User, error)
}

// MainCategModel is the interface that wraps the basic methods for main category model.
type MainCategModel interface {
	// Create inserts a new main category into the database.
	Create(categ *domain.MainCateg, userID int64, iconID int64) error

	// GetOneByUserID returns a main category by user id and name.
	GetOneByUserID(userID int64, name string) (*domain.MainCateg, error)
}
