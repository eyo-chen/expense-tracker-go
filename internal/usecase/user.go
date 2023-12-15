package usecase

import (
	"os"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/auth"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

type userUC struct {
	User UserModel
}

func newUserUC(u UserModel) *userUC {
	return &userUC{User: u}
}

type Claims struct {
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	jwt.RegisteredClaims
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

func (u *userUC) Login(user *domain.User) (string, error) {
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

	key := []byte(os.Getenv("JWT_SECRET_KEY"))
	claims := Claims{
		UserID:    userByEmail.ID,
		UserName:  userByEmail.Name,
		UserEmail: userByEmail.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(key)
	if err != nil {
		logger.Error("token.SignedString failed", "package", "usecase", "err", err)
		return "", err
	}

	return s, nil
}
