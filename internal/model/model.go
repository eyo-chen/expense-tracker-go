package model

import "database/sql"

type Model struct {
	User UserModel
}

func New(db *sql.DB) *Model {
	return &Model{
		User: *newUserModel(db),
	}
}
