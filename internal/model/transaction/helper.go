package transaction

import (
	"bytes"
	"fmt"
	"reflect"
	"unicode"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

func getAllQStmt(opt domain.GetTransOpt, decodedNextKey domain.DecodedNextKey, t Transaction) string {
	qStmt := `SELECT t.id, t.user_id, t.type, t.price, t.note, t.date, mc.id, mc.name, mc.type, sc.id, sc.name, i.id, i.url
						FROM transactions AS t
						INNER JOIN main_categories AS mc 
						ON t.main_category_id = mc.id
						INNER JOIN sub_categories AS sc 
						ON t.sub_category_id = sc.id
						INNER JOIN icons AS i
						ON mc.icon_id = i.id
						WHERE t.user_id = ?`

	if opt.Filter.StartDate != nil && opt.Filter.EndDate != nil {
		qStmt += " AND date BETWEEN ? AND ?"
	}

	if opt.Filter.StartDate != nil {
		qStmt += " AND date >= ?"
	}

	if opt.Filter.EndDate != nil {
		qStmt += " AND date <= ?"
	}

	if opt.Filter.MinPrice != nil {
		qStmt += " AND price >= ?"
	}

	if opt.Filter.MaxPrice != nil {
		qStmt += " AND price <= ?"
	}

	if opt.Filter.MainCategIDs != nil {
		qStmt += " AND mc.id IN (?"
		for i := 1; i < len(opt.Filter.MainCategIDs); i++ {
			qStmt += ", ?"
		}
		qStmt += ")"
	}

	if opt.Filter.SubCategIDs != nil {
		qStmt += " AND sc.id IN (?"
		for i := 1; i < len(opt.Filter.SubCategIDs); i++ {
			qStmt += ", ?"
		}
		qStmt += ")"
	}

	if len(decodedNextKey) != 0 {
		for key := range decodedNextKey {
			s := fmt.Sprintf(" AND t.%s < ?", genDBFieldNames(key, t))
			qStmt += s
		}
	}

	if opt.Cursor.Size != 0 {
		qStmt += " ORDER BY t.id DESC LIMIT ?"
	}

	return qStmt
}

// genDBFieldNames generates db field names from struct field names
// e.g. "UserID" -> "user_id", "MainCategID" -> "main_category_id"
func genDBFieldNames(key string, t Transaction) string {
	val := reflect.ValueOf(t)

	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		if fieldName != key {
			continue
		}

		t := val.Type().Field(i).Tag.Get("esql")
		fmt.Println("t: ", t)
		if t == "" {
			return camelToSnake(key)
		}

		return t
	}

	return ""
}

func camelToSnake(input string) string {
	var buf bytes.Buffer

	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 && unicode.IsLower(rune(input[i-1])) {
				buf.WriteRune('_')
			}
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
		}
	}

	return buf.String()
}

func getAllArgs(opt domain.GetTransOpt, decodedNextKey domain.DecodedNextKey, userID int64) []interface{} {
	var args []interface{}
	args = append(args, userID)

	if opt.Filter.StartDate != nil && opt.Filter.EndDate != nil {
		args = append(args, *opt.Filter.StartDate, *opt.Filter.EndDate)
	}

	if opt.Filter.StartDate != nil {
		args = append(args, *opt.Filter.StartDate)
	}

	if opt.Filter.EndDate != nil {
		args = append(args, *opt.Filter.EndDate)
	}

	if opt.Filter.MinPrice != nil {
		args = append(args, *opt.Filter.MinPrice)
	}

	if opt.Filter.MaxPrice != nil {
		args = append(args, *opt.Filter.MaxPrice)
	}

	if opt.Filter.MainCategIDs != nil {
		for _, id := range opt.Filter.MainCategIDs {
			args = append(args, id)
		}
	}

	if opt.Filter.SubCategIDs != nil {
		for _, id := range opt.Filter.SubCategIDs {
			args = append(args, id)
		}
	}

	if len(decodedNextKey) != 0 {
		for _, v := range decodedNextKey {
			args = append(args, v)
		}
	}

	if opt.Cursor.Size != 0 {
		args = append(args, opt.Cursor.Size)
	}

	return args
}

func getAccInfoQStmt(query domain.GetAccInfoQuery) string {
	qStmt := `SELECT
						SUM(CASE WHEN type = '1' THEN price ELSE 0 END) AS total_income,
						SUM(CASE WHEN type = '2' THEN price ELSE 0 END) AS total_expense,
						SUM(CASE WHEN type = '1' THEN price ELSE -price END) AS total_balance
	          FROM transactions
						WHERE user_id = ?
						`

	if query.StartDate != nil && query.EndDate != nil {
		qStmt += " AND date BETWEEN ? AND ?"
	}

	if query.StartDate != nil {
		qStmt += " AND date >= ?"
	}

	if query.EndDate != nil {
		qStmt += " AND date <= ?"
	}

	qStmt += " GROUP BY user_id"

	return qStmt
}

func getAccInfoArgs(query domain.GetAccInfoQuery, userID int64) []interface{} {
	var args []interface{}
	args = append(args, userID)

	if query.StartDate != nil && query.EndDate != nil {
		args = append(args, *query.StartDate, *query.EndDate)
	}

	if query.StartDate != nil {
		args = append(args, *query.StartDate)
	}

	if query.EndDate != nil {
		args = append(args, *query.EndDate)
	}

	return args
}

func getGetDailyBarChartDataQuery(mainCategIDs []int64) string {
	qStmt := `
	  SELECT DATE_FORMAT(date, '%Y-%m-%d') AS date,
		       SUM(price)
		FROM transactions
		WHERE user_id = ?
		AND type = ?
		AND date BETWEEN ? AND ?
	`

	if mainCategIDs != nil {
		qStmt += "AND main_category_id IN (?"
		for i := 1; i < len(mainCategIDs); i++ {
			qStmt += ", ?"
		}
		qStmt += ")"
	}

	qStmt += `GROUP BY date
						ORDER BY date`

	return qStmt
}

func genGetDailyBarChartDataArgs(userID int64, transactionType domain.TransactionType, dateRange domain.ChartDateRange, mainCategIDs []int64) []interface{} {
	l := 4
	if mainCategIDs != nil {
		l += len(mainCategIDs)
	}

	args := make([]interface{}, 0, l)
	args = append(args, userID, transactionType.ToModelValue(), dateRange.Start, dateRange.End)

	for _, id := range mainCategIDs {
		args = append(args, id)
	}

	return args
}

func getGetMonthlyBarChartDataQuery(mainCategIDs []int64) string {
	qStmt := `
		SELECT YEAR(date),
					 LPAD(MONTH(date), 2, '0') AS month,
					 SUM(price)
		FROM transactions
		WHERE user_id = ?
		AND type = ?
		AND date BETWEEN ? AND ?
		`

	if mainCategIDs != nil {
		qStmt += "AND main_category_id IN (?"
		for i := 1; i < len(mainCategIDs); i++ {
			qStmt += ", ?"
		}
		qStmt += ")"
	}

	qStmt += `GROUP BY YEAR(date), LPAD(MONTH(date), 2, '0')
						ORDER BY YEAR(date), LPAD(MONTH(date), 2, '0')`

	return qStmt
}

func getGetMonthlyBarChartDataArgs(userID int64, transactionType domain.TransactionType, dateRange domain.ChartDateRange, mainCategIDs []int64) []interface{} {
	l := 4
	if mainCategIDs != nil {
		l += len(mainCategIDs)
	}

	args := make([]interface{}, 0, l)
	args = append(args, userID, transactionType.ToModelValue(), dateRange.Start, dateRange.End)

	for _, id := range mainCategIDs {
		args = append(args, id)
	}

	return args
}
