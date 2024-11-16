package transaction

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
)

func getAllQStmt(opt domain.GetTransOpt, decodedNextKeys domain.DecodedNextKeys, t Transaction) string {
	var sb strings.Builder

	sb.WriteString(`SELECT t.id, t.user_id, t.type, t.price, t.note, t.date, mc.id, mc.name, mc.type, mc.icon_type, mc.icon_data, sc.id, sc.name
									FROM transactions AS t
									LEFT JOIN main_categories AS mc 
									ON t.main_category_id = mc.id
									LEFT JOIN sub_categories AS sc 
									ON t.sub_category_id = sc.id
									WHERE t.user_id = ?
									`)

	if opt.Search.Keyword != nil {
		sb.WriteString(" AND MATCH (note) AGAINST (? IN NATURAL LANGUAGE MODE)")
	}

	if opt.Filter.StartDate != nil && opt.Filter.EndDate != nil {
		sb.WriteString(" AND date BETWEEN ? AND ?")
	}

	if opt.Filter.StartDate != nil {
		sb.WriteString(" AND date >= ?")
	}

	if opt.Filter.EndDate != nil {
		sb.WriteString(" AND date <= ?")
	}

	if opt.Filter.MinPrice != nil {
		sb.WriteString(" AND price >= ?")
	}

	if opt.Filter.MaxPrice != nil {
		sb.WriteString(" AND price <= ?")
	}

	if opt.Filter.MainCategIDs != nil {
		sb.WriteString(" AND mc.id IN (?")
		for i := 1; i < len(opt.Filter.MainCategIDs); i++ {
			sb.WriteString(", ?")
		}
		sb.WriteString(")")
	}

	if opt.Filter.SubCategIDs != nil {
		sb.WriteString(" AND sc.id IN (?")
		for i := 1; i < len(opt.Filter.SubCategIDs); i++ {
			sb.WriteString(", ?")
		}
		sb.WriteString(")")
	}

	// construct the next key query statement
	// now, we only support 1 or 2 next keys
	// when it's 1, it means there's no sorting(sort by id)
	// when it's 2, it means there's sorting(sort by id and other field)
	if len(decodedNextKeys) != 0 {
		if len(decodedNextKeys) == 1 {
			sb.WriteString(fmt.Sprintf(" AND t.%s %s ?", genDBFieldNames(decodedNextKeys[0].Field, t), domain.GetOperandFromSort(opt.Sort)))
		}

		if len(decodedNextKeys) == 2 {
			// AND col_1 < or > val_1
			// OR (col_1 = val_1 AND col_2 < or > val_2)
			// the shorted version is: AND (col_1, col_2) < or > (val_1, val_2)
			sb.WriteString(fmt.Sprintf(" AND t.%s %s ?", genDBFieldNames(decodedNextKeys[0].Field, t), domain.GetOperandFromSort(opt.Sort)))
			sb.WriteString(fmt.Sprintf(" OR (t.%s = ? AND t.%s %s ?)", genDBFieldNames(decodedNextKeys[0].Field, t), genDBFieldNames(decodedNextKeys[1].Field, t), domain.GetOperandFromSort(opt.Sort)))
		}
	}

	if opt.Sort != nil {
		sb.WriteString(fmt.Sprintf(" ORDER BY t.%s %s, t.id %s", opt.Sort.By.String(), opt.Sort.Dir.String(), opt.Sort.Dir.String()))
	}

	if opt.Cursor.Size != 0 {
		sb.WriteString(" LIMIT ?")
	}

	return sb.String()
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

		t := val.Type().Field(i).Tag.Get("mysqlf")
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

func getAllArgs(opt domain.GetTransOpt, decodedNextKeys domain.DecodedNextKeys, userID int64) []interface{} {
	var args []interface{}
	args = append(args, userID)

	if opt.Search.Keyword != nil {
		args = append(args, *opt.Search.Keyword)
	}

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

	if len(decodedNextKeys) != 0 {
		if len(decodedNextKeys) == 1 {
			args = append(args, decodedNextKeys[0].Value)
		}

		if len(decodedNextKeys) == 2 {
			args = append(args, decodedNextKeys[0].Value, decodedNextKeys[0].Value, decodedNextKeys[1].Value)
		}
	}

	if opt.Cursor.Size != 0 {
		args = append(args, opt.Cursor.Size)
	}

	return args
}

func getAccInfoQStmt(query domain.GetAccInfoQuery) string {
	var sb strings.Builder

	sb.WriteString(`SELECT
									SUM(CASE WHEN type = '1' THEN price ELSE 0 END) AS total_income,
									SUM(CASE WHEN type = '2' THEN price ELSE 0 END) AS total_expense,
									SUM(CASE WHEN type = '1' THEN price ELSE -price END) AS total_balance
									FROM transactions
									WHERE user_id = ?
									`)

	if query.StartDate != nil && query.EndDate != nil {
		sb.WriteString(" AND date BETWEEN ? AND ?")
	}

	if query.StartDate != nil {
		sb.WriteString(" AND date >= ?")
	}

	if query.EndDate != nil {
		sb.WriteString(" AND date <= ?")
	}

	sb.WriteString(" GROUP BY user_id")

	return sb.String()
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
	var sb strings.Builder

	sb.WriteString(`SELECT 
									DATE_FORMAT(date, '%Y-%m-%d') AS date,
									SUM(price)
									FROM transactions
									WHERE user_id = ?
									AND type = ?
									AND date BETWEEN ? AND ?
									`)

	if mainCategIDs != nil {
		sb.WriteString("AND main_category_id IN (?")
		for i := 1; i < len(mainCategIDs); i++ {
			sb.WriteString(", ?")
		}
		sb.WriteString(")")
	}

	sb.WriteString(`GROUP BY date
						      ORDER BY date`)

	return sb.String()
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
	var sb strings.Builder

	sb.WriteString(`SELECT
									YEAR(date),
									LPAD(MONTH(date), 2, '0') AS month,
									SUM(price)
									FROM transactions
									WHERE user_id = ?
									AND type = ?
									AND date BETWEEN ? AND ?
									`)

	if mainCategIDs != nil {
		sb.WriteString("AND main_category_id IN (?")
		for i := 1; i < len(mainCategIDs); i++ {
			sb.WriteString(", ?")
		}
		sb.WriteString(")")
	}

	sb.WriteString(`GROUP BY YEAR(date), LPAD(MONTH(date), 2, '0')
						      ORDER BY YEAR(date), LPAD(MONTH(date), 2, '0')`)

	return sb.String()
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
