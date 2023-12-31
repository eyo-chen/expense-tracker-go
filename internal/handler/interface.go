package handler

import (
	"context"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

// UserUC is the interface that wraps the basic methods for user usecase.
type UserUC interface {
	// Signup is a function that registers a user.
	Signup(user *domain.User) error

	// Login is a function that logs in a user.
	Login(user *domain.User) (string, error)
}

// MainCategUC is the interface that wraps the basic methods for main category usecase.
type MainCategUC interface {
	// Create is a function that creates a main category.
	Create(categ *domain.MainCateg, userID int64) error

	// GetAll is a function that returns all main categories by user id.
	GetAll(userID int64) ([]*domain.MainCateg, error)

	// Update is a function that updates a main category.
	Update(categ *domain.MainCateg, userID int64) error

	// Delete is a function that deletes a main category.
	Delete(id int64) error
}

// SubCategUC is the interface that wraps the basic methods for sub category usecase.
type SubCategUC interface {
	// Create is a function that creates a sub category.
	Create(categ *domain.SubCateg, userID int64) error

	// GetAll is a function that returns all sub categories by user id.
	GetAll(userID int64) ([]*domain.SubCateg, error)

	// GetByMainCategID is a function that returns all sub categories by user id and main category id.
	GetByMainCategID(userID, mainCategID int64) ([]*domain.SubCateg, error)

	// Update is a function that updates a sub category.
	Update(categ *domain.SubCateg, userID int64) error

	// Delete is a function that deletes a sub category.
	Delete(id int64) error
}

// TransactionUC is the interface that wraps the basic methods for transaction usecase.
type TransactionUC interface {
	// Create is a function that creates a transaction.
	Create(ctx context.Context, user *domain.User, transaction *domain.Transaction) error
}
