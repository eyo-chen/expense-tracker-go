package user

import (
	"net/http"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/handler/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/errutil"
	"github.com/OYE0303/expense-tracker-go/pkg/jsonutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/validator"
)

type UserHandler struct {
	User interfaces.UserUC
}

func NewUserHandler(user interfaces.UserUC) *UserHandler {
	return &UserHandler{User: user}
}

func (u UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJson failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if !v.Signup(input.Email, input.Password, input.Name) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}
	if err := u.User.Signup(&user); err != nil {
		if err == domain.ErrDataAlreadyExists {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		logger.Error("u.User.Signup failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusCreated, nil, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (u UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJson failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if !v.Login(input.Email, input.Password) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := domain.User{
		Email:    input.Email,
		Password: input.Password,
	}

	token, err := u.User.Login(&user)
	if err != nil {
		if err == domain.ErrAuthentication {
			errutil.AuthenticationErrorResponse(w, r, err)
			return
		}

		logger.Error("u.User.Login failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	resp := map[string]interface{}{
		"token": token,
	}
	if err := jsonutil.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
