package tocoldp

import (
	"database/sql"

	"github.com/gnames/coldp/ent/coldp"
	"github.com/sfborg/to-coldp/pkg/config"
)

type tocoldp struct {
	cfg  config.Config
	db   *sql.DB
	cldp coldp.Builder
}

func New(
	cfg config.Config,
	sfdb *sql.DB,
	cldp coldp.Builder,
) ToCoLDP {
	res := tocoldp{
		cfg:  cfg,
		cldp: cldp,
		db:   sfdb,
	}
	return &res
}

func (t *tocoldp) Export(path string) error {
	return nil
}
