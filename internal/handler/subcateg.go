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

type subCategHandler struct {
	SubCateg SubCategUC
}

func newSubCategHandler(s SubCategUC) *subCategHandler {
	return &subCategHandler{
		SubCateg: s,
	}
}

func (s *subCategHandler) CreateSubCateg(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		MainCategID int64  `json:"main_category_id"`
	}
	if err := jsutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsutil.ReadJSON failed", "package", "handler", "err", err)
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
		if err == domain.ErrDataAlreadyExists {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		logger.Error("s.SubCateg.Create failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
