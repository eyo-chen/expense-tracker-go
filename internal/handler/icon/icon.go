package icon

import (
	"net/http"

	"github.com/eyo-chen/expense-tracker-go/internal/handler/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/ctxutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/errutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/jsonutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	PackageName = "handler/icon"
)

type Hlr struct {
	Icon interfaces.IconUC
}

func New(i interfaces.IconUC) *Hlr {
	return &Hlr{
		Icon: i,
	}
}

type icon struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

func (h *Hlr) List(w http.ResponseWriter, r *http.Request) {
	icons, err := h.Icon.List()
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	respData := map[string]interface{}{
		"icons": icons,
	}
	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", PackageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *Hlr) ListByUserID(w http.ResponseWriter, r *http.Request) {
	user := ctxutil.GetUser(r)
	icons, err := h.Icon.ListByUserID(r.Context(), user.ID)
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	respData := map[string]interface{}{
		"icons": cvtToIcon(icons),
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", PackageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
