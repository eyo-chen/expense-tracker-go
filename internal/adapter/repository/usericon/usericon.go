package usericon

import (
	"context"
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
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
	_, err := r.DB.Exec(stmt, userIcon.UserID, userIcon.ObjectKey)
	if err != nil {
		logger.Error("user_icons INSERT r.DB.Exec", "err", err, "package", packageName)
		return err
	}

	return nil
}

func (r *Repo) GetByUserID(ctx context.Context, userID int64) ([]domain.UserIcon, error) {
	stmt := `SELECT id, user_id, object_key FROM user_icons WHERE user_id = ?`
	rows, err := r.DB.Query(stmt, userID)
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
