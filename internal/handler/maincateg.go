package handler

import (
	"net/http"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/ctxutil"
	"github.com/OYE0303/expense-tracker-go/pkg/errutil"
	"github.com/OYE0303/expense-tracker-go/pkg/jsutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/validator"
)

type mainCategHandler struct {
	MainCateg MainCategUC
}

func newMainCategHandler(m MainCategUC) *mainCategHandler {
	return &mainCategHandler{MainCateg: m}
}

func (m *mainCategHandler) AddMainCateg(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		IconID int64  `json:"icon_id"`
	}
	if err := jsutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsutil.ReadJSON failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	categ := domain.MainCateg{
		Name:   input.Name,
		Type:   input.Type,
		IconID: input.IconID,
	}

	v := validator.New()
	if !v.AddMainCateg(&categ) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	if err := m.MainCateg.Add(&categ, user.ID); err != nil {
		if err == domain.ErrDataAlreadyExists || err == domain.ErrDataNotFound {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		logger.Error("m.MainCateg.Add failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsutil.WriteJSON(w, http.StatusCreated, nil, nil); err != nil {
		logger.Error("jsutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
