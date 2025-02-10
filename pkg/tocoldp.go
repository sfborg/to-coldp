package tocoldp

import (
	"path/filepath"

	"github.com/gnames/coldp/ent/coldp"
	"github.com/sfborg/to-coldp/internal/ent"
	"github.com/sfborg/to-coldp/pkg/config"
)

type tocoldp struct {
	cfg  config.Config
	clf  ent.ColDPFiles
	cldp coldp.Builder
}

func New(
	cfg config.Config,
	clf ent.ColDPFiles,
	cldp coldp.Builder,
) ToCoLDP {
	res := tocoldp{
		cfg:  cfg,
		clf:  clf,
		cldp: cldp,
	}
	return &res
}

func (t *tocoldp) Export(path string) error {
	clCfg := t.cldp.Config()
	pathMeta := filepath.Join(clCfg.BuilderDir, "metadata.json")
	err := t.clf.Meta(pathMeta)
	if err != nil {
		return err
	}
	return nil
}
