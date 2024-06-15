package interfaces

import (
	"context"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

// UserModel is the interface that wraps the basic methods for user model.
type UserModel interface {
	// Create inserts a new user into the database.
	Create(name, email, passwordHash string) error

	// FindByEmail returns a user by email.
	FindByEmail(email string) (domain.User, error)

	// GetInfo returns a user by id.
	GetInfo(userID int64) (domain.User, error)

	// Update updates a user.
	Update(ctx context.Context, userID int64, opt domain.UpdateUserOpt) error
}

// MainCategModel is the interface that wraps the basic methods for main category model.
type MainCategModel interface {
	// Create inserts a new main category into the database.
	Create(categ *domain.MainCateg, userID int64) error

	// GetAll returns all main categories by user id.
	GetAll(ctx context.Context, userID int64, transType domain.TransactionType) ([]domain.MainCateg, error)

	// Update updates a main category.
	Update(categ *domain.MainCateg) error

	// Delete deletes a main category.
	Delete(id int64) error

	// GetByID returns a main category by id and user id.
	GetByID(id, userID int64) (*domain.MainCateg, error)

	// BatchCreate inserts multiple main categories into the database.
	BatchCreate(ctx context.Context, categs []domain.MainCateg, userID int64) error
}

// SubCategModel is the interface that wraps the basic methods for sub category model.
type SubCategModel interface {
	// Create inserts a new sub category into the database.
	Create(categ *domain.SubCateg, userID int64) error

	// Update updates a sub category.
	Update(categ *domain.SubCateg) error

	// GetByMainCategID returns all sub categories by user id and main category id.
	GetByMainCategID(userID, mainCategID int64) ([]*domain.SubCateg, error)

	// Delete deletes a sub category.
	Delete(id int64) error

	// GetByID returns a sub category by id and user id.
	GetByID(id, userID int64) (*domain.SubCateg, error)

	// BatchCreate inserts multiple sub categories into the database.
	BatchCreate(ctx context.Context, categs []domain.SubCateg, userID int64) error
}

// IconModel is the interface that wraps the basic methods for icon model.
type IconModel interface {
	// GetByID returns an icon by id.
	GetByID(id int64) (domain.Icon, error)

	// List returns all icons.
	List() ([]domain.Icon, error)

	// GetByIDs returns icons by ids.
	GetByIDs(ids []int64) (map[int64]domain.Icon, error)
}

// TransactionModel is the interface that wraps the basic methods for transaction model.
type TransactionModel interface {
	// Create inserts a new transaction into the database.
	Create(ctx context.Context, trans domain.CreateTransactionInput) error

	// GetAll returns all transactions by user id and query option.
	GetAll(ctx context.Context, query domain.GetTransOpt, userID int64) ([]domain.Transaction, domain.DecodedNextKeys, error)

	// Update updates a transaction.
	Update(ctx context.Context, trans domain.UpdateTransactionInput) error

	// Delete deletes a transaction by id.
	Delete(ctx context.Context, id int64) error

	// GetAccInfo returns accumulated information by user id and query.
	GetAccInfo(ctx context.Context, query domain.GetAccInfoQuery, userID int64) (domain.AccInfo, error)

	// GetByIDAndUserID returns a transaction by id and user id. Note that returned transaction is not included main category, sub category, and icon.
	GetByIDAndUserID(ctx context.Context, id, userID int64) (domain.Transaction, error)

	// GetDailyBarChartData returns bar chart data grouped by date.
	GetDailyBarChartData(ctx context.Context, dateRange domain.ChartDateRange, transactionType domain.TransactionType, mainCategIDs []int64, userID int64) (domain.DateToChartData, error)

	// GetMonthlyBarChartData returns bar chart data grouped by month.
	GetMonthlyBarChartData(ctx context.Context, dateRange domain.ChartDateRange, transactionType domain.TransactionType, mainCategIDs []int64, userID int64) (domain.DateToChartData, error)

	// GetPieChartData returns pie chart data.
	GetPieChartData(ctx context.Context, dataRange domain.ChartDateRange, transactionType domain.TransactionType, userID int64) (domain.ChartData, error)

	// GetDailyLineChartData returns line chart data grouped by date.
	GetDailyLineChartData(ctx context.Context, dateRange domain.ChartDateRange, userID int64) (domain.DateToChartData, error)

	// GetMonthlyLineChartData returns line chart data grouped by month.
	GetMonthlyLineChartData(ctx context.Context, dateRange domain.ChartDateRange, userID int64) (domain.DateToChartData, error)

	// GetMonthlyData returns monthly data.
	GetMonthlyData(ctx context.Context, dateRange domain.GetMonthlyDateRange, userID int64) (domain.MonthDayToTransactionType, error)
}
