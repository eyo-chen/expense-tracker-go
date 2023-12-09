package usecase

import (
	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/auth"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type userUC struct {
	User UserModel
}

func newUserUC(UserModel UserModel) *userUC {
	return &userUC{User: UserModel}
}

type Claims struct {
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

func (u *userUC) Signup(user *domain.User) error {
	userByEmail, err := u.User.FindByEmail(user.Email)
	if err != nil {
		logger.Error("u.User.FindByEmail failed", "package", "usecase", "err", err)
		return err
	}

	if userByEmail != nil {
		return domain.ErrDataAlreadyExists
	}

	passwordHash, err := auth.GenerateHashPassword(user.Password)
	if err != nil {
		logger.Error("auth.GenerateHashPassword failed", "package", "usecase", "err", err)
		return err
	}

	if err := u.User.Create(user.Name, user.Email, passwordHash); err != nil {
		logger.Error("u.User.Insert failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}
