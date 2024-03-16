package interfaces

import (
	"context"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

// UserUC is the interface that wraps the basic methods for user usecase.
type UserUC interface {
	// Signup registers a user.
	Signup(user *domain.User) error

	// Login logs in a user.
	Login(user *domain.User) (string, error)
}

// MainCategUC is the interface that wraps the basic methods for main category usecase.
type MainCategUC interface {
	// Create creates a main category.
	Create(categ *domain.MainCateg, userID int64) error

	// GetAll returns all main categories by user id.
	GetAll(userID int64, transType domain.TransactionType) ([]domain.MainCateg, error)

	// Update updates a main category.
	Update(categ *domain.MainCateg, userID int64) error

	// Delete deletes a main category.
	Delete(id int64) error
}

// SubCategUC is the interface that wraps the basic methods for sub category usecase.
type SubCategUC interface {
	// Create creates a sub category.
	Create(categ *domain.SubCateg, userID int64) error

	// GetAll returns all sub categories by user id.
	GetAll(userID int64) ([]*domain.SubCateg, error)

	// GetByMainCategID returns all sub categories by user id and main category id.
	GetByMainCategID(userID, mainCategID int64) ([]*domain.SubCateg, error)

	// Update updates a sub category.
	Update(categ *domain.SubCateg, userID int64) error

	// Delete deletes a sub category.
	Delete(id int64) error
}

// TransactionUC is the interface that wraps the basic methods for transaction usecase.
type TransactionUC interface {
	// Create creates a transaction.
	Create(ctx context.Context, trans domain.CreateTransactionInput) error

	// GetAll returns all transactions by query and user id.
	GetAll(ctx context.Context, query domain.GetQuery, user domain.User) ([]domain.Transaction, error)

	// GetAccInfo returns the accunulated information by user id.
	GetAccInfo(ctx context.Context, query domain.GetAccInfoQuery, user domain.User) (domain.AccInfo, error)

	// Delete deletes a transaction by id.
	Delete(ctx context.Context, id int64, user domain.User) error
}

// IconUC is the interface that wraps the basic methods for icon usecase.
type IconUC interface {
	// List returns all icons.
	List() ([]domain.Icon, error)
}
