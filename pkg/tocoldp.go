package tocoldp

import (
	"log/slog"
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

	slog.Info("Exporting Metadata file")
	pathMeta := filepath.Join(clCfg.BuilderDir, "Metadata.json")
	err := t.clf.Meta(pathMeta)
	if err != nil {
		slog.Error("Cannot create Metadata.json file")
		return err
	}

	slog.Info("Exporting Reference file")
	pathRef := filepath.Join(clCfg.BuilderDir, "Reference.tsv")
	err = t.clf.Reference(pathRef)
	if err != nil {
		slog.Error("Cannot create Reference.tsv file", "error", err)
		return err
	}

	slog.Info("Creation of CoLDP archive finished successfully")
	return nil
}
