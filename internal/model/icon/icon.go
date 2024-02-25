package icon

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

const (
	PackageName = "model/icon"
)

type IconModel struct {
	DB *sql.DB
}

func NewIconModel(db *sql.DB) *IconModel {
	return &IconModel{DB: db}
}

type Icon struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

func (m *IconModel) List() ([]domain.Icon, error) {
	stmt := `SELECT id, url FROM icons`

	rows, err := m.DB.Query(stmt)
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

	return icons, nil
}

func (m *IconModel) GetByID(id int64) (domain.Icon, error) {
	stmt := `SELECT id, url FROM icons WHERE id = ?`

	var icon Icon
	if err := m.DB.QueryRow(stmt, id).Scan(&icon.ID, &icon.URL); err != nil {
		if err == sql.ErrNoRows {
			return domain.Icon{}, domain.ErrIconNotFound
		}

		return domain.Icon{}, err
	}

	return cvtToDomainIcon(icon), nil
}
