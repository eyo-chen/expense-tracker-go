package usericon

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/go-sql-driver/mysql"
)

var (
	packageName = "adapter/repository/usericon"
)

type Repo struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

type userIcon struct {
	ID        int64
	UserID    int64 `gofacto:"foreignKey,struct:User"`
	ObjectKey string
}

func (r *Repo) Create(ctx context.Context, userIcon domain.UserIcon) error {
	stmt := `INSERT INTO user_icons (user_id, object_key) VALUES (?, ?)`
	_, err := r.DB.ExecContext(ctx, stmt, userIcon.UserID, userIcon.ObjectKey)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == domain.ErrMySQLForeignKeyConstraintViolation {
			logger.Error("user_icons INSERT foreign key constraint violation", "err", err, "package", packageName)
			return domain.ErrUserNotFound
		}
		logger.Error("user_icons INSERT r.DB.ExecContext", "err", err, "package", packageName)
		return err
	}

	return nil
}

func (r *Repo) GetByUserID(ctx context.Context, userID int64) ([]domain.UserIcon, error) {
	stmt := `SELECT id, user_id, object_key FROM user_icons WHERE user_id = ?`
	rows, err := r.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		logger.Error("user_icons SELECT r.DB.Query", "err", err, "package", packageName)
		return nil, err
	}

	var userIcons []domain.UserIcon
	for rows.Next() {
		var userIcon userIcon
		if err := rows.Scan(&userIcon.ID, &userIcon.UserID, &userIcon.ObjectKey); err != nil {
			logger.Error("user_icons SELECT rows.Scan", "err", err, "package", packageName)
			return nil, err
		}
		userIcons = append(userIcons, cvtToDomainUserIcon(userIcon))
	}

	return userIcons, nil
}

func (r *Repo) GetByObjectKeyAndUserID(ctx context.Context, objectKey string, userID int64) (domain.UserIcon, error) {
	stmt := `SELECT id, user_id, object_key FROM user_icons WHERE user_id = ? AND object_key = ?`

	var userIcon userIcon
	if err := r.DB.QueryRowContext(ctx, stmt, userID, objectKey).Scan(&userIcon.ID, &userIcon.UserID, &userIcon.ObjectKey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.UserIcon{}, domain.ErrUserIconNotFound
		}

		logger.Error("get user icon by object key and user id r.DB.QueryRowContext", "err", err, "package", packageName)
		return domain.UserIcon{}, err
	}

	return cvtToDomainUserIcon(userIcon), nil
}
