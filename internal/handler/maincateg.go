package handler

import (
	"net/http"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/ctxutil"
	"github.com/OYE0303/expense-tracker-go/pkg/errutil"
	"github.com/OYE0303/expense-tracker-go/pkg/jsonutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/validator"
)

type mainCategHandler struct {
	MainCateg MainCategUC
}

func newMainCategHandler(m MainCategUC) *mainCategHandler {
	return &mainCategHandler{MainCateg: m}
}

func (m *mainCategHandler) CreateMainCateg(w http.ResponseWriter, r *http.Request) {
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

	mainCategType := domain.CvtToMainCategType(input.Type)
	categ := domain.MainCateg{
		Name: input.Name,
		Type: mainCategType,
		Icon: &domain.Icon{
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
		if err == domain.ErrDataAlreadyExists || err == domain.ErrDataNotFound {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		logger.Error("m.MainCateg.Add failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusCreated, nil, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (m *mainCategHandler) GetAllMainCateg(w http.ResponseWriter, r *http.Request) {
	user := ctxutil.GetUser(r)
	categs, err := m.MainCateg.GetAll(user.ID)
	if err != nil {
		logger.Error("m.MainCateg.GetAll failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	// TODO: refactor the following logic
	type icon struct {
		ID  int64  `json:"id"`
		URL string `json:"url"`
	}
	type resp struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		Icon icon   `json:"icon"`
	}

	var respCategs []*resp
	for _, categ := range categs {
		respCategs = append(respCategs, &resp{
			ID:   categ.ID,
			Name: categ.Name,
			Type: categ.Type.String(),
		})
	}

	respData := map[string]interface{}{
		"categories": respCategs,
	}
	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (m *mainCategHandler) UpdateMainCateg(w http.ResponseWriter, r *http.Request) {
	id, err := jsonutil.ReadID(r)
	if err != nil {
		logger.Error("jsonutil.ReadID failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	var input struct {
		Name   string `json:"name"`
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
		Icon: &domain.Icon{
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
		if err == domain.ErrDataAlreadyExists || err == domain.ErrDataNotFound {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		logger.Error("m.MainCateg.Update failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, nil, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (m *mainCategHandler) DeleteMainCateg(w http.ResponseWriter, r *http.Request) {
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
