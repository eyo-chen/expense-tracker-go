package handler

import "github.com/OYE0303/expense-tracker-go/internal/domain"

type UserUC interface {
	// Signup is a function that registers a user.
	Signup(user *domain.User) error
}
