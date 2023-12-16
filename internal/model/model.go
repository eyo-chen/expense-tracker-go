package model

import "database/sql"

type Model struct {
	User      UserModel
	MainCateg MainCategModel
	SubCateg  SubCategModel
	Icon      IconModel
}

func New(db *sql.DB) *Model {
	return &Model{
		User:      *newUserModel(db),
		MainCateg: *newMainCategModel(db),
		SubCateg:  *newSubCategModel(db),
		Icon:      *newIconModel(db),
	}
}
