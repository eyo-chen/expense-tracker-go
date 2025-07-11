package icon

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	packageName = "adapter/repository/icon"
)

type Repo struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

type Icon struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

func (r *Repo) GetByID(ctx context.Context, id int64) (domain.DefaultIcon, error) {
	stmt := `SELECT id, url FROM icons WHERE id = ?`

	var icon Icon
	if err := r.DB.QueryRowContext(ctx, stmt, id).Scan(&icon.ID, &icon.URL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.DefaultIcon{}, domain.ErrIconNotFound
		}

		logger.Error("get icon by id r.DB.QueryRowContext", "err", err, "package", packageName)
		return domain.DefaultIcon{}, err
	}

	return cvtToDomainDefaultIcon(icon), nil
}

func (r *Repo) List() ([]domain.DefaultIcon, error) {
	stmt := `SELECT id, url FROM icons`

	rows, err := r.DB.Query(stmt)
	if err != nil {
		logger.Error("r.DB.Query failed", "package", packageName, "err", err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Unable to close rows", "package", packageName, "err", err)
		}
	}()

	var icons []domain.DefaultIcon
	for rows.Next() {
		var icon Icon
		if err := rows.Scan(&icon.ID, &icon.URL); err != nil {
			logger.Error("rows.Scan failed", "package", packageName, "err", err)
			return nil, err
		}

		icons = append(icons, cvtToDomainDefaultIcon(icon))
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Unable to close rows", "package", packageName, "err", err)
		}
	}()

	return icons, nil
}

func (r *Repo) GetByIDs(ids []int64) (map[int64]domain.DefaultIcon, error) {
	var sb strings.Builder
	sb.WriteString(`SELECT id, url FROM icons WHERE id IN (`)

	for i := range ids {
		if i == 0 {
			sb.WriteString("?")
		} else {
			sb.WriteString(", ?")
		}

		if i == len(ids)-1 {
			sb.WriteString(")")
		}
	}

	idsInterface := make([]interface{}, len(ids))
	for i, id := range ids {
		idsInterface[i] = id
	}

	var icons []Icon
	rows, err := r.DB.Query(sb.String(), idsInterface...)
	if err != nil {
		logger.Error("r.DB.Query failed", "package", packageName, "err", err)
		return nil, err
	}

	for rows.Next() {
		var icon Icon
		if err := rows.Scan(&icon.ID, &icon.URL); err != nil {
			logger.Error("rows.Scan failed", "package", packageName, "err", err)
			return nil, err
		}

		icons = append(icons, icon)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("Unable to close rows", "package", packageName, "err", err)
		}
	}()

	if len(icons) == 0 {
		return nil, domain.ErrIconNotFound
	}

	return cvtToIDToDomainDefaultIcon(icons), nil
}
