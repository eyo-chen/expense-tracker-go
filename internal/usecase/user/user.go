package user

import (
	"context"
	"errors"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/auth"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	packageName = "usecase/user"
)

type UC struct {
	user  interfaces.UserRepo
	redis interfaces.RedisService
}

func New(u interfaces.UserRepo, r interfaces.RedisService) *UC {
	return &UC{user: u, redis: r}
}

func (u *UC) Signup(ctx context.Context, user domain.User) (domain.Token, error) {
	_, err := u.user.FindByEmail(user.Email)
	if err != nil && err != domain.ErrEmailNotFound {
		return domain.Token{}, err
	}
	if err == nil {
		return domain.Token{}, domain.ErrEmailAlreadyExists
	}

	passwordHash, err := auth.GenerateHashPassword(user.Password)
	if err != nil {
		logger.Error("auth.GenerateHashPassword failed", "package", packageName, "err", err)
		return domain.Token{}, err
	}

	if err := u.user.Create(user.Name, user.Email, passwordHash); err != nil {
		return domain.Token{}, err
	}

	userWithID, err := u.user.FindByEmail(user.Email)
	if err != nil {
		return domain.Token{}, err
	}

	accessToken, err := genJWTToken(userWithID)
	if err != nil {
		return domain.Token{}, err
	}
	refreshToken, err := genRefreshToken()
	if err != nil {
		return domain.Token{}, err
	}

	// it's ok to fail
	_ = u.redis.Set(ctx, hashToken(refreshToken), userWithID.Email, 7*24*time.Hour)

	return domain.Token{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (u *UC) Login(ctx context.Context, user domain.User) (domain.Token, error) {
	userByEmail, err := u.user.FindByEmail(user.Email)
	if err != nil {
		if errors.Is(err, domain.ErrEmailNotFound) {
			return domain.Token{}, domain.ErrAuthentication
		}

		return domain.Token{}, err
	}

	if !auth.CompareHashPassword(user.Password, userByEmail.Password_hash) {
		return domain.Token{}, domain.ErrAuthentication
	}

	token, err := genJWTToken(userByEmail)
	if err != nil {
		return domain.Token{}, err
	}
	refreshToken, err := genRefreshToken()
	if err != nil {
		return domain.Token{}, err
	}

	// it's ok to fail
	_ = u.redis.Set(ctx, hashToken(refreshToken), userByEmail.Email, 7*24*time.Hour)

	return domain.Token{
		Access:  token,
		Refresh: refreshToken,
	}, nil
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
