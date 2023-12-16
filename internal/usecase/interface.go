package usecase

import "github.com/OYE0303/expense-tracker-go/internal/domain"

// UserModel is the interface that wraps the basic methods for user model.
type UserModel interface {
	// Create inserts a new user into the database.
	Create(name, email, passwordHash string) error

	// FindByEmail returns a user by email.
	FindByEmail(email string) (*domain.User, error)
}

// MainCategModel is the interface that wraps the basic methods for main category model.
type MainCategModel interface {
	// Create inserts a new main category into the database.
	Create(categ *domain.MainCateg, userID int64) error

	// Update updates a main category.
	Update(categ *domain.MainCateg) error

	// Delete deletes a main category.
	Delete(id int64) error

	// GetByID returns a main category by id.
	GetByID(id int64) (*domain.MainCateg, error)

	// GetOneByUserID returns a main category by user id and name.
	GetOneByUserID(userID int64, name string) (*domain.MainCateg, error)
}

// SubCategModel is the interface that wraps the basic methods for sub category model.
type SubCategModel interface {
	// Create inserts a new sub category into the database.
	Create(categ *domain.SubCateg, userID int64) error

	// Update updates a sub category.
	Update(categ *domain.SubCateg) error

	// GetByID returns a sub category by id.
	GetByID(id int64) (*domain.SubCateg, error)

	// GetOneByUserID returns a sub category by user id and name.
	GetOneByUserID(userID int64, name string) (*domain.SubCateg, error)
}

// IconModel is the interface that wraps the basic methods for icon model.
type IconModel interface {
	// GetByID returns an icon by id.
	GetByID(id int64) (*domain.Icon, error)
}
