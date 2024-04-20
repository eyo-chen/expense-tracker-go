package transaction

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

func genGetTransOpt(r *http.Request) (domain.GetTransOpt, error) {
	var opt domain.GetTransOpt

	rawStartDate := r.URL.Query().Get("start_date")
	if rawStartDate != "" {
		date, err := time.Parse(time.DateOnly, rawStartDate)
		if err != nil {
			return domain.GetTransOpt{}, err
		}
		opt.Filter.StartDate = &date
	}

	rawEndDate := r.URL.Query().Get("end_date")
	if rawEndDate != "" {
		date, err := time.Parse(time.DateOnly, rawEndDate)
		if err != nil {
			return domain.GetTransOpt{}, err
		}
		opt.Filter.EndDate = &date
	}

	mainCategIDs, err := genMainCategIDs(r)
	if err != nil {
		return domain.GetTransOpt{}, err
	}
	opt.Filter.MainCategIDs = mainCategIDs

	subCategIDs, err := genSubCategIDs(r)
	if err != nil {
		return domain.GetTransOpt{}, err
	}
	opt.Filter.SubCategIDs = subCategIDs

	opt.Cursor.NextKey = r.URL.Query().Get("next_key")

	rawSize := r.URL.Query().Get("size")
	if rawSize != "" {
		size, err := strconv.Atoi(rawSize)
		if err != nil {
			return domain.GetTransOpt{}, err
		}
		opt.Cursor.Size = size
	}

	return opt, nil
}

func genGetAccInfoQuery(r *http.Request) domain.GetAccInfoQuery {
	rawStartDate := r.URL.Query().Get("start_date")
	rawEndDate := r.URL.Query().Get("end_date")

	var query domain.GetAccInfoQuery

	if rawStartDate != "" {
		query.StartDate = &rawStartDate
	}

	if rawEndDate != "" {
		query.EndDate = &rawEndDate
	}

	return query
}

func genGetMonthlyDataRange(r *http.Request) (time.Time, time.Time, error) {
	rawStartDate := r.URL.Query().Get("start_date")
	rawEndDate := r.URL.Query().Get("end_date")

	startDate, err := time.Parse(time.DateOnly, rawStartDate)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("start date must be in YYYY-MM-DD format")
	}

	endDate, err := time.Parse(time.DateOnly, rawEndDate)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("end date must be in YYYY-MM-DD format")
	}

	return startDate, endDate, nil
}

func genChartDateRange(r *http.Request) (domain.ChartDateRange, error) {
	rawStartDate := r.URL.Query().Get("start_date")
	rawEndDate := r.URL.Query().Get("end_date")

	start, err := time.Parse(time.DateOnly, rawStartDate)
	if err != nil {
		return domain.ChartDateRange{}, errors.New("start date must be in YYYY-MM-DD format")
	}

	end, err := time.Parse(time.DateOnly, rawEndDate)
	if err != nil {
		return domain.ChartDateRange{}, errors.New("end date must be in YYYY-MM-DD format")
	}

	return domain.ChartDateRange{
		Start: start,
		End:   end,
	}, nil
}

func genMainCategIDs(r *http.Request) ([]int64, error) {
	rawMainCategIDs := r.URL.Query().Get("main_category_ids")
	if rawMainCategIDs == "" {
		return nil, nil
	}

	strSlice := strings.Split(rawMainCategIDs, ",")
	intSlice := make([]int64, len(strSlice))

	for i, str := range strSlice {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		intSlice[i] = int64(num)
	}

	return intSlice, nil
}

func genSubCategIDs(r *http.Request) ([]int64, error) {
	rawSubCategIDs := r.URL.Query().Get("sub_category_ids")
	if rawSubCategIDs == "" {
		return nil, nil
	}

	strSlice := strings.Split(rawSubCategIDs, ",")
	intSlice := make([]int64, len(strSlice))

	for i, str := range strSlice {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		intSlice[i] = int64(num)
	}

	return intSlice, nil
}
