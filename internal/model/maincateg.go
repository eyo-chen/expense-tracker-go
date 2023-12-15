package model

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
)

type MainCategModel struct {
	DB *sql.DB
}

func newMainCategModel(db *sql.DB) *MainCategModel {
	return &MainCategModel{DB: db}
}

type MainCateg struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (m *MainCategModel) Create(categ *domain.MainCateg, userID int64, iconID int64) error {
	stmt := `INSERT INTO main_categories (name, type, user_id, icon_id) VALUES (?, ?, ?, ?)`

	categType := "1"
	if categ.Type == "expense" {
		categType = "2"
	}

	if _, err := m.DB.Exec(stmt, categ.Name, categType, userID, iconID); err != nil {
		return err
	}

	return nil
}

func (m *MainCategModel) GetOneByUserID(userID int64, name string) (*domain.MainCateg, error) {
	stmt := `SELECT id, name, type FROM main_categories WHERE user_id = ? AND name = ?`

	var categ MainCateg
	if err := m.DB.QueryRow(stmt, userID, name).Scan(&categ.ID, &categ.Name, &categ.Type); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}

		return nil, err
	}

	return cvtToDomainMainCateg(&categ), nil
}

func cvtToDomainMainCateg(c *MainCateg) *domain.MainCateg {
	return &domain.MainCateg{
		ID:   c.ID,
		Name: c.Name,
		Type: c.Type,
	}
}
