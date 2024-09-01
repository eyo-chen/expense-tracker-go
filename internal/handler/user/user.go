package user

import (
	"net/http"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/ctxutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/errutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/jsonutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/validator"
)

type Hlr struct {
	User interfaces.UserUC
}

func New(user interfaces.UserUC) *Hlr {
	return &Hlr{User: user}
}

func (h *Hlr) Signup(w http.ResponseWriter, r *http.Request) {
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
	token, err := h.User.Signup(r.Context(), user)
	if err != nil {
		if err == domain.ErrEmailAlreadyExists {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		errutil.ServerErrorResponse(w, r, err)
		return
	}

	resp := map[string]interface{}{
		"token": token,
	}
	if err := jsonutil.WriteJSON(w, http.StatusCreated, resp, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *Hlr) Login(w http.ResponseWriter, r *http.Request) {
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

	token, err := h.User.Login(r.Context(), user)
	if err != nil {
		if err == domain.ErrAuthentication {
			errutil.AuthenticationErrorResponse(w, r, err)
			return
		}

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

func (h *Hlr) Token(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.URL.Query().Get("refresh_token")
	v := validator.New()
	if !v.Token(refreshToken) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	token, err := h.User.Token(r.Context(), refreshToken)
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	resp := map[string]interface{}{
		"access_token":  token.Access,
		"refresh_token": token.Refresh,
	}
	if err := jsonutil.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *Hlr) GetInfo(w http.ResponseWriter, r *http.Request) {
	userCtx := ctxutil.GetUser(r)
	user, err := h.User.GetInfo(userCtx.ID)
	if err != nil {
		if err == domain.ErrUserIDNotFound {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		errutil.ServerErrorResponse(w, r, err)
		return
	}

	resp := map[string]interface{}{
		"id":                   user.ID,
		"name":                 user.Name,
		"email":                user.Email,
		"is_set_init_category": user.IsSetInitCategory,
	}
	if err := jsonutil.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
