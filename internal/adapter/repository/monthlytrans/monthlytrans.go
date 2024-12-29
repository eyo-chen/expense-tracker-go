package monthlytrans

import (
	"context"
	"database/sql"
	"errors"
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

func (r *Repo) GetByUserIDAndMonthDate(ctx context.Context, userID int64, monthDate time.Time) (domain.AccInfo, error) {
	query := `SELECT * FROM monthly_transactions WHERE user_id = ? AND month_date = ?`

	var mt MonthlyTrans
	row := r.DB.QueryRowContext(ctx, query, userID, monthDate)
	if err := row.Scan(&mt.ID, &mt.UserID, &mt.MonthDate, &mt.TotalExpense, &mt.TotalIncome, &mt.CreatedAt, &mt.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.AccInfo{}, domain.ErrDataNotFound
		}

		logger.Error("r.DB.QueryRowContext failed", "package", packageName, "err", err)
		return domain.AccInfo{}, err
	}

	return domain.AccInfo{
		TotalExpense: mt.TotalExpense,
		TotalIncome:  mt.TotalIncome,
		TotalBalance: mt.TotalIncome - mt.TotalExpense,
	}, nil
}

func (r *Repo) Update(ctx context.Context, userID int64, monthDate time.Time, transType domain.TransactionType, amount float64) error {
	if !transType.IsValid() {
		logger.Error("invalid transaction type", "package", packageName, "err", domain.ErrInvalidTransType)
		return domain.ErrInvalidTransType
	}

	var query string
	if transType == domain.TransactionTypeIncome {
		query = `UPDATE monthly_transactions SET total_income = total_income + ? WHERE user_id = ? AND month_date = ?`
	} else {
		query = `UPDATE monthly_transactions SET total_expense = total_expense + ? WHERE user_id = ? AND month_date = ?`
	}

	if _, err := r.DB.ExecContext(ctx, query, amount, userID, monthDate); err != nil {
		logger.Error("r.DB.ExecContext failed", "package", packageName, "err", err)
		return err
	}

	return nil
}
