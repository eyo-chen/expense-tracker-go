package user

import (
	"context"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/model/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/auth"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	packageName = "usecase/user"
)

type UC struct {
	user  interfaces.UserModel
	redis interfaces.RedisService
}

func New(u interfaces.UserModel, r interfaces.RedisService) *UC {
	return &UC{user: u, redis: r}
}

func (u *UC) Signup(user domain.User) (string, error) {
	_, err := u.user.FindByEmail(user.Email)
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

	if err := u.user.Create(user.Name, user.Email, passwordHash); err != nil {
		return "", err
	}

	userWithID, err := u.user.FindByEmail(user.Email)
	if err != nil {
		return "", err
	}

	token, err := genJWTToken(userWithID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UC) Login(user domain.User) (string, error) {
	userByEmail, err := u.user.FindByEmail(user.Email)
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

func (u *UC) Token(ctx context.Context, refreshToken string) (domain.Token, error) {
	hashedToken := hashToken(refreshToken)
	userEmail, err := u.redis.GetDel(ctx, hashedToken)
	if err != nil {
		return domain.Token{}, err
	}

	user, err := u.user.FindByEmail(userEmail)
	if err != nil {
		return domain.Token{}, err
	}

	accessToken, err := genJWTToken(user)
	if err != nil {
		return domain.Token{}, err
	}

	newRefreshToken, err := genRefreshToken()
	if err != nil {
		return domain.Token{}, err
	}

	// it's ok to fail
	_ = u.redis.Set(ctx, hashToken(newRefreshToken), userEmail, 7*24*time.Hour)

	return domain.Token{
		Access:  accessToken,
		Refresh: newRefreshToken,
	}, nil
}

func (u *UC) GetInfo(userID int64) (domain.User, error) {
	return u.user.GetInfo(userID)
}
