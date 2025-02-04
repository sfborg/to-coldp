package tocoldp

import (
	"database/sql"

	"github.com/gnames/coldp/ent/coldp"
	"github.com/sfborg/sflib/ent/sfga"
	"github.com/sfborg/to-coldp/pkg/config"
)

type tocoldp struct {
	cfg  config.Config
	db   *sql.DB
	cldp coldp.Builder
}

func New(
	cfg config.Config,
	sfdb sfga.DB,
	cldp coldp.Builder,
) (ToCoLDP, error) {
	res := tocoldp{
		cfg:  cfg,
		cldp: cldp,
	}

	db, err := sfdb.Connect()
	if err != nil {
		return nil, err
	}

	res.db = db
	return &res, nil
}

func (t *tocoldp) Export(path string) error {
	return nil
}
