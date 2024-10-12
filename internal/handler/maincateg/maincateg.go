package maincateg

import (
	"net/http"
	"slices"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/handler/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/ctxutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/errutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/jsonutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/validator"
)

type Hlr struct {
	MainCateg interfaces.MainCategUC
}

func New(m interfaces.MainCategUC) *Hlr {
	return &Hlr{MainCateg: m}
}

type createInput struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	IconType string `json:"icon_type"`
	IconID   int64  `json:"icon_id"`
}

func (h *Hlr) Create(w http.ResponseWriter, r *http.Request) {
	var input createInput
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	categ := domain.CreateMainCategInput{
		Name:     input.Name,
		Type:     domain.CvtToTransactionType(input.Type),
		IconType: domain.CvtToIconType(input.IconType),
		IconID:   input.IconID,
	}

	v := validator.New()
	if !v.CreateMainCateg(categ) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	if err := h.MainCateg.Create(r.Context(), categ, user.ID); err != nil {
		errors := []error{
			domain.ErrIconNotFound,
			domain.ErrUserIconNotFound,
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

func (h *Hlr) GetAll(w http.ResponseWriter, r *http.Request) {
	qType := r.URL.Query().Get("type")
	categType := domain.CvtToTransactionType(qType)
	user := ctxutil.GetUser(r)
	ctx := r.Context()

	categs, err := h.MainCateg.GetAll(ctx, user.ID, categType)
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

var updateInput struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	IconType string `json:"icon_type"`
	IconID   int64  `json:"icon_id"`
}

func (h *Hlr) Update(w http.ResponseWriter, r *http.Request) {
	id, err := jsonutil.ReadID(r)
	if err != nil {
		logger.Error("jsonutil.ReadID failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	if err := jsonutil.ReadJson(w, r, &updateInput); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	categ := domain.UpdateMainCategInput{
		ID:       id,
		Name:     updateInput.Name,
		Type:     domain.CvtToTransactionType(updateInput.Type),
		IconType: domain.CvtToIconType(updateInput.IconType),
		IconID:   updateInput.IconID,
	}

	v := validator.New()
	if !v.UpdateMainCateg(categ) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	if err := h.MainCateg.Update(r.Context(), categ, user.ID); err != nil {
		errors := []error{
			domain.ErrUniqueNameUserType,
			domain.ErrUniqueIconUser,
			domain.ErrUserIconNotFound,
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

func (h *Hlr) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := jsonutil.ReadID(r)
	if err != nil {
		logger.Error("jsonutil.ReadID failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	if err := h.MainCateg.Delete(id); err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, nil, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
