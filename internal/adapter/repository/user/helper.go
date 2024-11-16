package user

import (
	"strings"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
)

func genUpdateStmtAndVal(opt domain.UpdateUserOpt, userID int64) (string, []interface{}) {
	val := []interface{}{}

	var sb strings.Builder
	sb.WriteString(`UPDATE users SET `)

	if opt.IsSetInitCategory != nil {
		sb.WriteString(`is_set_init_category = ? `)
		val = append(val, *opt.IsSetInitCategory)
	}

	sb.WriteString(`WHERE id = ?`)
	val = append(val, userID)
	return sb.String(), val
}
