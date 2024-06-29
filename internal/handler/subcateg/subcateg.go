package subcateg

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

type SubCategHandler struct {
	SubCateg interfaces.SubCategUC
}

func NewSubCategHandler(s interfaces.SubCategUC) *SubCategHandler {
	return &SubCategHandler{
		SubCateg: s,
	}
}

func (s *SubCategHandler) CreateSubCateg(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		MainCategID int64  `json:"main_category_id"`
	}
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	categ := domain.SubCateg{
		Name:        input.Name,
		MainCategID: input.MainCategID,
	}

	v := validator.New()
	if !v.CreateSubCateg(&categ) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	if err := s.SubCateg.Create(&categ, user.ID); err != nil {
		if err == domain.ErrUniqueNameUserMainCateg {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		logger.Error("s.SubCateg.Create failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (s *SubCategHandler) GetByMainCategID(w http.ResponseWriter, r *http.Request) {
	id, err := jsonutil.ReadID(r)
	if err != nil {
		logger.Error("jsonutil.ReadID failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	user := ctxutil.GetUser(r)
	categs, err := s.SubCateg.GetByMainCategID(user.ID, id)
	if err != nil {
		logger.Error("s.SubCateg.GetByMainCategID failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	respData := map[string]interface{}{
		"categories": categs,
	}
	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (s *SubCategHandler) UpdateSubCateg(w http.ResponseWriter, r *http.Request) {
	id, err := jsonutil.ReadID(r)
	if err != nil {
		logger.Error("jsonutil.ReadID failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	var input struct {
		Name        string `json:"name"`
		MainCategID int64  `json:"main_category_id"`
	}
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	categ := domain.SubCateg{
		ID:          id,
		Name:        input.Name,
		MainCategID: input.MainCategID,
	}

	v := validator.New()
	if !v.UpdateSubCateg(&categ) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	if err := s.SubCateg.Update(&categ, user.ID); err != nil {
		if err == domain.ErrUniqueNameUserMainCateg {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		logger.Error("s.SubCateg.Update failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (s *SubCategHandler) DeleteSubCateg(w http.ResponseWriter, r *http.Request) {
	id, err := jsonutil.ReadID(r)
	if err != nil {
		logger.Error("jsonutil.ReadID failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	if err := s.SubCateg.Delete(id); err != nil {
		logger.Error("s.SubCateg.Delete failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, nil, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
