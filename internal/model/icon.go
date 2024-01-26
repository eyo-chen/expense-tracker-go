package model

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

type IconModel struct {
	DB *sql.DB
}

func newIconModel(db *sql.DB) *IconModel {
	return &IconModel{DB: db}
}

type Icon struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

func (m *IconModel) GetByID(id int64) (*domain.Icon, error) {
	stmt := `SELECT id, url FROM icons WHERE id = ?`

	var icon Icon
	if err := m.DB.QueryRow(stmt, id).Scan(&icon.ID, &icon.URL); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrIconNotFound
		}

		return nil, err
	}

	return cvtToDomainIcons(&icon), nil
}

func cvtToDomainIcons(i *Icon) *domain.Icon {
	return &domain.Icon{
		ID:  i.ID,
		URL: i.URL,
	}
}
