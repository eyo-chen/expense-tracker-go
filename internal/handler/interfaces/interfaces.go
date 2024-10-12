package interfaces

import (
	"context"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
)

// UserUC is the interface that wraps the basic methods for user usecase.
type UserUC interface {
	// Signup registers a user.
	Signup(ctx context.Context, user domain.User) (domain.Token, error)

	// Login logs in a user.
	Login(ctx context.Context, user domain.User) (domain.Token, error)

	// GetInfo returns the user information by user id.
	GetInfo(userID int64) (domain.User, error)

	// Token returns the access token and refresh token by refresh token.
	Token(ctx context.Context, refreshToken string) (domain.Token, error)
}

// MainCategUC is the interface that wraps the basic methods for main category usecase.
type MainCategUC interface {
	// Create creates a main category.
	Create(ctx context.Context, categ domain.CreateMainCategInput, userID int64) error

	// GetAll returns all main categories by user id.
	GetAll(ctx context.Context, userID int64, transType domain.TransactionType) ([]domain.MainCateg, error)

	// Update updates a main category.
	Update(ctx context.Context, categ domain.UpdateMainCategInput, userID int64) error

	// Delete deletes a main category.
	Delete(id int64) error
}

// SubCategUC is the interface that wraps the basic methods for sub category usecase.
type SubCategUC interface {
	// Create creates a sub category.
	Create(categ *domain.SubCateg, userID int64) error

	// GetByMainCategID returns all sub categories by user id and main category id.
	GetByMainCategID(userID, mainCategID int64) ([]*domain.SubCateg, error)

	// Update updates a sub category.
	Update(categ *domain.SubCateg, userID int64) error

	// Delete deletes a sub category.
	Delete(id int64) error
}

// TransactionUC is the interface that wraps the basic methods for transaction usecase.
type TransactionUC interface {
	// Create creates a transaction.
	Create(ctx context.Context, trans domain.CreateTransactionInput) error

	// GetAll returns all transactions by query option and user id.
	GetAll(ctx context.Context, opt domain.GetTransOpt, user domain.User) ([]domain.Transaction, domain.Cursor, error)

	// Update updates a transaction.
	Update(ctx context.Context, trans domain.UpdateTransactionInput, user domain.User) error

	// Delete deletes a transaction by id.
	Delete(ctx context.Context, id int64, user domain.User) error

	// GetAccInfo returns the accunulated information by user id.
	GetAccInfo(ctx context.Context, query domain.GetAccInfoQuery, user domain.User) (domain.AccInfo, error)

	// GetBarChartData returns bar chart data.
	GetBarChartData(ctx context.Context, chartDateRange domain.ChartDateRange, timeRangeType domain.TimeRangeType, transactionType domain.TransactionType, mainCategIDs []int64, user domain.User) (domain.ChartData, error)

	// GetPieChartData returns pie chart data.
	GetPieChartData(ctx context.Context, dataRange domain.ChartDateRange, transactionType domain.TransactionType, user domain.User) (domain.ChartData, error)

	// GetLineChartData returns line chart data.
	GetLineChartData(ctx context.Context, chartDateRange domain.ChartDateRange, timeRangeType domain.TimeRangeType, user domain.User) (domain.ChartData, error)

	// GetMonthlyData returns monthly data.
	GetMonthlyData(ctx context.Context, dateRange domain.GetMonthlyDateRange, user domain.User) ([]domain.TransactionType, error)
}

// IconUC is the interface that wraps the basic methods for icon usecase.
type IconUC interface {
	// List returns all icons.
	List() ([]domain.DefaultIcon, error)

	// ListByUserID returns all icons by user id.
	ListByUserID(ctx context.Context, userID int64) ([]domain.Icon, error)
}

// UserIconUC is the interface that wraps the basic methods for user icon usecase.
type UserIconUC interface {
	// GetPutObjectURL returns the put object url.
	GetPutObjectURL(ctx context.Context, fileName string, userID int64) (string, error)

	// Create creates a user icon.
	Create(ctx context.Context, fileName string, userID int64) error
}

// InitDataUC is the interface that wraps the basic methods for init data usecase.
type InitDataUC interface {
	// List returns the initial data.
	List() (domain.InitData, error)

	// Create creates the initial data.
	Create(ctx context.Context, data domain.InitData, userID int64) error
}
