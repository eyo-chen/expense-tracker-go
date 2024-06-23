package user

import (
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/model/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/auth"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

const (
	packageName = "usecase/user"
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

func (u *UserUC) Signup(user domain.User) (string, error) {
	_, err := u.User.FindByEmail(user.Email)
	if err != nil && err != domain.ErrEmailNotFound {
		return "", err
	}
	if err == nil {
		return "", domain.ErrEmailAlreadyExists
	}

	passwordHash, err := auth.GenerateHashPassword(user.Password)
	if err != nil {
		logger.Error("auth.GenerateHashPassword failed", "package", packageName, "err", err)
		return "", err
	}

	if err := u.User.Create(user.Name, user.Email, passwordHash); err != nil {
		return "", err
	}

	userWithID, err := u.User.FindByEmail(user.Email)
	if err != nil {
		return "", err
	}

	token, err := genJWTToken(userWithID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUC) Login(user domain.User) (string, error) {
	userByEmail, err := u.User.FindByEmail(user.Email)
	if err != nil {
		if err == domain.ErrEmailNotFound {
			return "", domain.ErrAuthentication
		}

		return "", err
	}

	if !auth.CompareHashPassword(user.Password, userByEmail.Password_hash) {
		return "", domain.ErrAuthentication
	}

	token, err := genJWTToken(userByEmail)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUC) GetInfo(userID int64) (domain.User, error) {
	return u.User.GetInfo(userID)
}
