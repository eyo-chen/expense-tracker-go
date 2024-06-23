package user

import (
	"github.com/eyo-chen/expense-tracker-go/internal/domain"
)

func genUpdateStmtAndVal(opt domain.UpdateUserOpt, userID int64) (string, []interface{}) {
	var (
		stmt string
		val  []interface{}
	)

	stmt += `UPDATE users SET `

	if opt.IsSetInitCategory != nil {
		stmt += `is_set_init_category = ? `
		val = append(val, *opt.IsSetInitCategory)
	}

	stmt += `WHERE id = ?`
	val = append(val, userID)
	return stmt, val
}
