package maincateg

import (
	"net/http"
	"slices"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/handler/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/ctxutil"
	"github.com/OYE0303/expense-tracker-go/pkg/errutil"
	"github.com/OYE0303/expense-tracker-go/pkg/jsonutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/validator"
)

type MainCategHandler struct {
	MainCateg interfaces.MainCategUC
}

func NewMainCategHandler(m interfaces.MainCategUC) *MainCategHandler {
	return &MainCategHandler{MainCateg: m}
}

func (m *MainCategHandler) CreateMainCateg(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		IconID int64  `json:"icon_id"`
	}
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	categ := domain.MainCateg{
		Name: input.Name,
		Type: domain.CvtToMainCategType(input.Type),
		Icon: domain.Icon{
			ID: input.IconID,
		},
	}

	v := validator.New()
	if !v.CreateMainCateg(&categ) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	if err := m.MainCateg.Create(&categ, user.ID); err != nil {
		errors := []error{
			domain.ErrIconNotFound,
			domain.ErrUniqueNameUserType,
			domain.ErrUniqueIconUser,
		}

		if slices.Contains(errors, err) {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusCreated, nil, nil); err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (m *MainCategHandler) GetAllMainCateg(w http.ResponseWriter, r *http.Request) {
	qType := r.URL.Query().Get("type")
	categType := domain.CvtToMainCategType(qType)
	user := ctxutil.GetUser(r)

	categs, err := m.MainCateg.GetAll(user.ID, categType)
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	resp := cvtToGetAllMainCategResp(categs)
	respData := map[string]interface{}{
		"categories": resp.Categories,
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (m *MainCategHandler) UpdateMainCateg(w http.ResponseWriter, r *http.Request) {
	id, err := jsonutil.ReadID(r)
	if err != nil {
		logger.Error("jsonutil.ReadID failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	var input struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		IconID int64  `json:"icon_id"`
	}
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	categ := domain.MainCateg{
		ID:   id,
		Name: input.Name,
		Type: domain.CvtToMainCategType(input.Type),
		Icon: domain.Icon{
			ID: input.IconID,
		},
	}

	v := validator.New()
	if !v.UpdateMainCateg(&categ) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	if err := m.MainCateg.Update(&categ, user.ID); err != nil {
		errors := []error{
			domain.ErrUniqueNameUserType,
			domain.ErrUniqueIconUser,
			domain.ErrMainCategNotFound,
			domain.ErrIconNotFound,
		}
		if slices.Contains(errors, err) {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, nil, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (m *MainCategHandler) DeleteMainCateg(w http.ResponseWriter, r *http.Request) {
	id, err := jsonutil.ReadID(r)
	if err != nil {
		logger.Error("jsonutil.ReadID failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	if err := m.MainCateg.Delete(id); err != nil {
		logger.Error("m.MainCateg.Delete failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, nil, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
