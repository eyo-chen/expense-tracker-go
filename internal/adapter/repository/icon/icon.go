package icon

import (
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
)

const (
	PackageName = "model/icon"
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

func (r *Repo) List() ([]domain.Icon, error) {
	stmt := `SELECT id, url FROM icons`

	rows, err := r.DB.Query(stmt)
	if err != nil {
		logger.Error("m.DB.Query failed", "package", PackageName, "err", err)
		return nil, err
	}
	defer rows.Close()

	var icons []domain.Icon
	for rows.Next() {
		var icon Icon
		if err := rows.Scan(&icon.ID, &icon.URL); err != nil {
			logger.Error("rows.Scan failed", "package", PackageName, "err", err)
			return nil, err
		}

		icons = append(icons, cvtToDomainIcon(icon))
	}
	defer rows.Close()

	return icons, nil
}

func (r *Repo) GetByID(id int64) (domain.Icon, error) {
	stmt := `SELECT id, url FROM icons WHERE id = ?`

	var icon Icon
	if err := r.DB.QueryRow(stmt, id).Scan(&icon.ID, &icon.URL); err != nil {
		if err == sql.ErrNoRows {
			return domain.Icon{}, domain.ErrIconNotFound
		}

		return domain.Icon{}, err
	}

	return cvtToDomainIcon(icon), nil
}

func (r *Repo) GetByIDs(ids []int64) (map[int64]domain.Icon, error) {
	stmt := `SELECT id, url FROM icons WHERE id IN (`
	for i := range ids {
		if i == 0 {
			stmt += "?"
		} else {
			stmt += ", ?"
		}

		if i == len(ids)-1 {
			stmt += ")"
		}
	}

	idsInterface := make([]interface{}, len(ids))
	for i, id := range ids {
		idsInterface[i] = id
	}

	var icons []Icon
	rows, err := r.DB.Query(stmt, idsInterface...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var icon Icon
		if err := rows.Scan(&icon.ID, &icon.URL); err != nil {
			return nil, err
		}

		icons = append(icons, icon)
	}
	defer rows.Close()

	if len(icons) == 0 {
		return nil, domain.ErrIconNotFound
	}

	return cvtToIDToDomainIcon(icons), nil
}
