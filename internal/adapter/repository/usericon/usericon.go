package usericon

import "database/sql"

type Repo struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

type UserIcon struct {
	ID        int64
	UserID    int64 `gofacto:"foreignKey,struct:User"`
	ObjectKey string
}
