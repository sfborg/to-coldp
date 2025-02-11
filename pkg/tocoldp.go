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

func (t *tocoldp) Export(outputPath string) error {
	clCfg := t.cldp.Config()

	slog.Info("Exporting Metadata file")
	path := filepath.Join(clCfg.BuilderDir, "Metadata.json")
	err := t.clf.Meta(path)
	if err != nil {
		slog.Error("Cannot create Metadata.json file")
		return err
	}

	slog.Info("Exporting Reference file")
	path = filepath.Join(clCfg.BuilderDir, "Reference.tsv")
	err = t.clf.Reference(path)
	if err != nil {
		slog.Error("Cannot create Reference.tsv file", "error", err)
		return err
	}

	if t.cfg.WithNameUsage {
	} else {
		slog.Info("Exporting Name file")
		path = filepath.Join(clCfg.BuilderDir, "Name.tsv")
		err = t.clf.Name(path)
		if err != nil {
			slog.Error("Cannot create Name.tsv file", "error", err)
			return err
		}

		slog.Info("Exporting Taxon file")
		path = filepath.Join(clCfg.BuilderDir, "Taxon.tsv")
		err = t.clf.Taxon(path)
		if err != nil {
			slog.Error("Cannot create Taxon.tsv file", "error", err)
			return err
		}

		slog.Info("Exporting Synonym file")
		path = filepath.Join(clCfg.BuilderDir, "Synonym.tsv")
		err = t.clf.Synonym(path)
		if err != nil {
			slog.Error("Cannot create Synonym.tsv file", "error", err)
			return err
		}
	}

	slog.Info("Exporting Vernacular file")
	path = filepath.Join(clCfg.BuilderDir, "Vernacular.tsv")
	err = t.clf.Vernacular(path)
	if err != nil {
		slog.Error("Cannot create Vernacular.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting NameRelation file")
	path = filepath.Join(clCfg.BuilderDir, "NameRelation.tsv")
	err = t.clf.NameRelation(path)
	if err != nil {
		slog.Error("Cannot create NameRelation.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting TypeMaterial file")
	path = filepath.Join(clCfg.BuilderDir, "TypeMaterial.tsv")
	err = t.clf.TypeMaterial(path)
	if err != nil {
		slog.Error("Cannot create TypeMaterial.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting Distribution file")
	path = filepath.Join(clCfg.BuilderDir, "Distribution.tsv")
	err = t.clf.Distribution(path)
	if err != nil {
		slog.Error("Cannot create Distribution.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting Media file")
	path = filepath.Join(clCfg.BuilderDir, "Media.tsv")
	err = t.clf.Media(path)
	if err != nil {
		slog.Error("Cannot create Media.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting Treatment file")
	path = filepath.Join(clCfg.BuilderDir, "Treatment.tsv")
	err = t.clf.Treatment(path)
	if err != nil {
		slog.Error("Cannot create Treatment.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting SpeciesEstimate file")
	path = filepath.Join(clCfg.BuilderDir, "SpeciesEstimate.tsv")
	err = t.clf.SpeciesEstimate(path)
	if err != nil {
		slog.Error("Cannot create SpeciesEstimate.tsv file", "error", err)
		return err
	}

	slog.Info("Creation of CoLDP archive finished successfully")
	return nil
}
