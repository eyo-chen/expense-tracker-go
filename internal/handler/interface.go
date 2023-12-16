package handler

import "github.com/OYE0303/expense-tracker-go/internal/domain"

// UserUC is the interface that wraps the basic methods for user usecase.
type UserUC interface {
	// Signup is a function that registers a user.
	Signup(user *domain.User) error

	// Login is a function that logs in a user.
	Login(user *domain.User) (string, error)
}

// MainCategUC is the interface that wraps the basic methods for main category usecase.
type MainCategUC interface {
	// Add is a function that adds a main category.
	Add(categ *domain.MainCateg, userID int64) error

	// Update is a function that updates a main category.
	Update(categ *domain.MainCateg, userID int64) error

	// Delete is a function that deletes a main category.
	Delete(id int64) error
}

// SubCategUC is the interface that wraps the basic methods for sub category usecase.
type SubCategUC interface {
	// Create is a function that creates a sub category.
	Create(categ *domain.SubCateg, userID int64) error
}
