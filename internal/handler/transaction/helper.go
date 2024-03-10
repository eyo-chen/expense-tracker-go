package transaction

import (
	"net/http"
	"strconv"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

func genGetAllQuery(r *http.Request) (domain.GetQuery, error) {
	rawStartDate := r.URL.Query().Get("start_date")
	rawEndDate := r.URL.Query().Get("end_date")
	rawMainCategID := r.URL.Query().Get("main_category_id")
	rawSubCategID := r.URL.Query().Get("sub_category_id")

	var query domain.GetQuery

	if rawStartDate != "" {
		query.StartDate = &rawStartDate
	}

	if rawEndDate != "" {
		query.EndDate = &rawEndDate
	}

	if rawMainCategID != "" {
		id, err := strconv.ParseInt(rawMainCategID, 10, 64)
		if err != nil {
			logger.Error("strconv.ParseInt failed", "package", packageName, "err", err)
			return domain.GetQuery{}, err
		}

		query.MainCategID = &id
	}

	if rawSubCategID != "" {
		id, err := strconv.ParseInt(rawSubCategID, 10, 64)
		if err != nil {
			logger.Error("strconv.ParseInt failed", "package", packageName, "err", err)
			return domain.GetQuery{}, err
		}

		query.SubCategID = &id
	}

	return query, nil
}
