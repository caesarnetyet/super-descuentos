package relational

import (
	"database/sql"
	"super-descuentos/relational/repository"
)

type SQLStore struct {
	Queries *repository.Queries
	DB      *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		Queries: repository.New(db),
		DB:      db,
	}
}
