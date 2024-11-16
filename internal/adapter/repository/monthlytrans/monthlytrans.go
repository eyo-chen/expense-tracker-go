package monthlytrans

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/errorutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	packageName    = "adapter/repository/monthlytrans"
	uniqueUserDate = "monthly_transactions.unique_user_month_date"
)

type Repo struct {
	DB *sql.DB
}

type MonthlyTrans struct {
	ID           int64
	UserID       int64 `gofacto:"foreignKey,struct:User"`
	MonthDate    time.Time
	TotalExpense float64
	TotalIncome  float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func New(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

func (r *Repo) Create(ctx context.Context, date time.Time, trans []domain.MonthlyAggregatedData) error {
	var sb strings.Builder
	sb.WriteString("INSERT INTO monthly_transactions (user_id, month_date, total_expense, total_income) VALUES ")

	args := make([]interface{}, 0, len(trans)*4)
	for i, t := range trans {
		sb.WriteString("(?, ?, ?, ?)")
		if i < len(trans)-1 {
			sb.WriteString(", ")
		}

		args = append(args, t.UserID, date, t.TotalExpense, t.TotalIncome)
	}

	stmt := sb.String()
	if _, err := r.DB.ExecContext(ctx, stmt, args...); err != nil {
		if errorutil.ParseError(err, uniqueUserDate) {
			return domain.ErrUniqueUserDate
		}

		logger.Error("r.DB.Exec failed", "package", packageName, "err", err)
		return err
	}

	return nil
}
