package cfio

import (
	"database/sql"

	clcfg "github.com/gnames/coldp/config"
	"github.com/sfborg/to-coldp/internal/ent"
)

type cfio struct {
	db    *sql.DB
	clCfg clcfg.Config
}

func New(db *sql.DB, clCfg clcfg.Config) ent.ColDPFiles {
	res := cfio{db: db, clCfg: clCfg}
	return &res
}

func (c *cfio) Config() clcfg.Config {
	return c.clCfg
}
