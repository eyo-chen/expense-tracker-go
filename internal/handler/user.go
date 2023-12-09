package handler

import (
	"net/http"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/errutil"
	"github.com/OYE0303/expense-tracker-go/pkg/jsutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type userHandler struct {
	User UserUC
}

func newUserHandler(user UserUC) *userHandler {
	return &userHandler{User: user}
}

func (u userHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := jsutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsutil.ReadJson failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
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

	if err := jsutil.WriteJSON(w, http.StatusCreated, nil, nil); err != nil {
		logger.Error("jsutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
