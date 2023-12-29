package model

import "database/sql"

type Model struct {
	User      UserModel
	MainCateg MainCategModel
	SubCateg  SubCategModel
	Icon      IconModel
}

func New(mysqlDB *sql.DB) *Model {
	return &Model{
		User:      *newUserModel(mysqlDB),
		MainCateg: *newMainCategModel(mysqlDB),
		SubCateg:  *newSubCategModel(mysqlDB),
		Icon:      *newIconModel(mysqlDB),
	}
}
