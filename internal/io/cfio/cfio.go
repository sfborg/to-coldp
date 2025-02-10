package cfio

import (
	"database/sql"

	cfiles "github.com/sfborg/to-coldp/internal/ent"
)

type cfio struct {
	db *sql.DB
}

func New(db *sql.DB) cfiles.ColDPFiles {
	res := cfio{db: db}
	return &res
}
