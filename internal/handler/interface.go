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
}
