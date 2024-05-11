package user

import (
	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/auth"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

type UserUC struct {
	User interfaces.UserModel
}

func NewUserUC(u interfaces.UserModel) *UserUC {
	return &UserUC{User: u}
}

type Claims struct {
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	jwt.RegisteredClaims
}

func (u *UserUC) Signup(user *domain.User) (string, error) {
	userByEmail, err := u.User.FindByEmail(user.Email)
	if err != nil {
		logger.Error("u.User.FindByEmail failed", "package", "usecase", "err", err)
		return "", err
	}

	if userByEmail != nil {
		return "", domain.ErrDataAlreadyExists
	}

	passwordHash, err := auth.GenerateHashPassword(user.Password)
	if err != nil {
		logger.Error("auth.GenerateHashPassword failed", "package", "usecase", "err", err)
		return "", err
	}

	if err := u.User.Create(user.Name, user.Email, passwordHash); err != nil {
		logger.Error("u.User.Insert failed", "package", "usecase", "err", err)
		return "", err
	}

	userWithID, err := u.User.FindByEmail(user.Email)
	if err != nil {
		logger.Error("u.User.FindByEmail failed", "package", "usecase", "err", err)
		return "", err
	}

	token, err := genJWTToken(*userWithID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUC) Login(user *domain.User) (string, error) {
	userByEmail, err := u.User.FindByEmail(user.Email)
	if err != nil {
		logger.Error("u.User.FindByEmail failed", "package", "usecase", "err", err)
		return "", err
	}

	if userByEmail == nil {
		return "", domain.ErrAuthentication
	}

	if !auth.CompareHashPassword(user.Password, userByEmail.Password_hash) {
		return "", domain.ErrAuthentication
	}

	token, err := genJWTToken(*userByEmail)
	if err != nil {
		return "", err
	}

	return token, nil
}
