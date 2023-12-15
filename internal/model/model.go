package model

import "database/sql"

type Model struct {
	User      UserModel
	MainCateg MainCategModel
}

func New(db *sql.DB) *Model {
	return &Model{
		User:      *newUserModel(db),
		MainCateg: *newMainCategModel(db),
	}
}
