package interfaces

import (
	"context"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
)

// UserRepo is the interface that wraps the basic methods for user repository.
type UserRepo interface {
	// Create inserts a new user into the database.
	Create(name, email, passwordHash string) error

	// FindByEmail returns a user by email.
	FindByEmail(email string) (domain.User, error)

	// GetInfo returns a user by id.
	GetInfo(userID int64) (domain.User, error)

	// Update updates a user.
	Update(ctx context.Context, userID int64, opt domain.UpdateUserOpt) error
}

// MainCategRepo is the interface that wraps the basic methods for main category repository.
type MainCategRepo interface {
	// Create inserts a new main category into the database.
	Create(ctx context.Context, categ domain.MainCateg, userID int64) error

	// GetAll returns all main categories by user id.
	GetAll(ctx context.Context, userID int64, transType domain.TransactionType) ([]domain.MainCateg, error)

	// Update updates a main category.
	Update(ctx context.Context, categ domain.MainCateg) error

	// Delete deletes a main category.
	Delete(id int64) error

	// GetByID returns a main category by id and user id.
	GetByID(id, userID int64) (*domain.MainCateg, error)

	// BatchCreate inserts multiple main categories into the database.
	BatchCreate(ctx context.Context, categs []domain.MainCateg, userID int64) error
}

// SubCategRepo is the interface that wraps the basic methods for sub category repository.
type SubCategRepo interface {
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

// IconRepo is the interface that wraps the basic methods for icon repository.
type IconRepo interface {
	// GetByID returns an icon by id.
	GetByID(ctx context.Context, id int64) (domain.DefaultIcon, error)

	// List returns all icons.
	List() ([]domain.DefaultIcon, error)

	// GetByIDs returns icons by ids.
	GetByIDs(ids []int64) (map[int64]domain.DefaultIcon, error)
}

// UserIconRepo is the interface that wraps the basic methods for user icon repository.
type UserIconRepo interface {
	// Create inserts a new user icon into the database.
	Create(ctx context.Context, userIcon domain.UserIcon) error

	// GetByUserID returns user icons by user id.
	GetByUserID(ctx context.Context, userID int64) ([]domain.UserIcon, error)

	// GetByID returns a user icon by id and user id.
	GetByID(ctx context.Context, id, userID int64) (domain.UserIcon, error)
}

// TransactionRepo is the interface that wraps the basic methods for transaction repository.
type TransactionRepo interface {
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

	// GetMonthlyAggregatedData returns monthly aggregated data.
	GetMonthlyAggregatedData(ctx context.Context, date time.Time) ([]domain.MonthlyAggregatedData, error)
}

// MonthlyTransRepo is the interface that wraps the basic methods for monthly transaction repository.
type MonthlyTransRepo interface {
	// Create inserts new monthly transactions into the database.
	Create(ctx context.Context, date time.Time, trans []domain.MonthlyAggregatedData) error

	// GetByUserIDAndMonthDate returns monthly aggregated data by user id and month date.
	GetByUserIDAndMonthDate(ctx context.Context, userID int64, monthDate time.Time) (domain.AccInfo, error)
}

// RedisService is the interface that wraps the basic methods for redis service.
type RedisService interface {
	// GetByFunc returns a value by key. If the value is not found, it will call the function to get the value and cache it.
	GetByFunc(ctx context.Context, key string, ttl time.Duration, f func() (string, error)) (string, error)

	// GetDel returns a value by key and delete it.
	GetDel(ctx context.Context, key string) (string, error)

	// Set sets a value by key.
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

// S3Service is the interface that wraps the basic methods for s3 service.
type S3Service interface {
	// PutObjectUrl returns a pre-signed URL to upload an object to S3.
	PutObjectUrl(ctx context.Context, objectKey string, lifetimeSecs int64) (string, error)

	// GetObjectUrl returns a pre-signed URL to get an object from S3.
	GetObjectUrl(ctx context.Context, objectKey string, lifetimeSecs int64) (string, error)

	// DeleteObject deletes an object from S3.
	DeleteObject(ctx context.Context, objectKey string) error
}

// StockService is the interface that wraps the basic methods for stock service.
type StockService interface {
	// Create creates a new stock.
	Create(ctx context.Context, stock domain.CreateStock) (string, error)

	// GetPortfolioInfo returns the portfolio info.
	GetPortfolioInfo(ctx context.Context, userID int32) (domain.Portfolio, error)

	// GetStockInfo returns the stock info.
	GetStockInfo(ctx context.Context, userID int32) (domain.AllStockInfo, error)
}

// HistoricalPortfolioService is the interface that wraps the basic methods for historical portfolio service.
type HistoricalPortfolioService interface {
	// Create creates a new historical portfolio record.
	Create(ctx context.Context, userID int32, date time.Time) error

	// GetPortfolioValue returns portfolio values over time.
	GetPortfolioValue(ctx context.Context, userID int32, dateOption string) ([]string, []float64, error)
}
