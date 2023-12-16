package model

import "database/sql"

type Model struct {
	User      UserModel
	MainCateg MainCategModel
	Icon      IconModel
}

func New(db *sql.DB) *Model {
	return &Model{
		User:      *newUserModel(db),
		MainCateg: *newMainCategModel(db),
		Icon:      *newIconModel(db),
	}
}
