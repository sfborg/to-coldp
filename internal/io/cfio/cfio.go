package cfio

import (
	"database/sql"

	"github.com/gnames/gnsys"
	"github.com/sfborg/to-coldp/internal/ent"
)

type cfio struct {
	db  *sql.DB
	dir string
}

func New(db *sql.DB, dir string) (ent.ColDPFiles, error) {
	res := cfio{db: db, dir: dir}
	err := gnsys.MakeDir(dir)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
