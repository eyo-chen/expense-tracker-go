package transaction

import "github.com/OYE0303/expense-tracker-go/internal/domain"

func getAllQStmt(query domain.GetQuery) string {
	qStmt := `SELECT t.id, t.user_id, t.type, t.price, t.note, t.date, mc.id, mc.name, mc.type, sc.id, sc.name, i.id, i.url
						FROM transactions AS t
						INNER JOIN main_categories AS mc 
						ON t.main_category_id = mc.id
						INNER JOIN sub_categories AS sc 
						ON t.sub_category_id = sc.id
						INNER JOIN icons AS i
						ON mc.icon_id = i.id
						WHERE t.user_id = ?`

	if query.StartDate != nil && query.EndDate != nil {
		qStmt += " AND date BETWEEN ? AND ?"
	}

	if query.StartDate != nil {
		qStmt += " AND date >= ?"
	}

	if query.EndDate != nil {
		qStmt += " AND date <= ?"
	}

	if query.MainCategID != nil {
		qStmt += " AND mc.id = ?"
	}

	if query.SubCategID != nil {
		qStmt += " AND sc.id = ?"
	}

	return qStmt
}

func getAllArgs(query domain.GetQuery, userID int64) []interface{} {
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

	if query.MainCategID != nil {
		args = append(args, *query.MainCategID)
	}

	if query.SubCategID != nil {
		args = append(args, *query.SubCategID)
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
