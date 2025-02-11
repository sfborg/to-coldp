package tocoldp

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/sfborg/to-coldp/internal/ent"
	"github.com/sfborg/to-coldp/pkg/config"
)

type tocoldp struct {
	cfg config.Config
	clf ent.ColDPFiles
}

func New(
	cfg config.Config,
	clf ent.ColDPFiles,
) ToCoLDP {
	res := tocoldp{
		cfg: cfg,
		clf: clf,
	}
	return &res
}

func (t *tocoldp) Export(outputPath string) error {
	base := t.cfg.CacheColdpDir

	slog.Info("Exporting Metadata file")
	path := filepath.Join(base, "Metadata.json")
	err := t.clf.Meta(path)
	if err != nil {
		slog.Error("Cannot create Metadata.json file")
		return err
	}

	slog.Info("Exporting Reference file")
	path = filepath.Join(base, "Reference.tsv")
	err = t.clf.Reference(path)
	if err != nil {
		slog.Error("Cannot create Reference.tsv file", "error", err)
		return err
	}

	if t.cfg.WithNameUsage {
	} else {
		slog.Info("Exporting Name file")
		path = filepath.Join(base, "Name.tsv")
		err = t.clf.Name(path)
		if err != nil {
			slog.Error("Cannot create Name.tsv file", "error", err)
			return err
		}

		slog.Info("Exporting Taxon file")
		path = filepath.Join(base, "Taxon.tsv")
		err = t.clf.Taxon(path)
		if err != nil {
			slog.Error("Cannot create Taxon.tsv file", "error", err)
			return err
		}

		slog.Info("Exporting Synonym file")
		path = filepath.Join(base, "Synonym.tsv")
		err = t.clf.Synonym(path)
		if err != nil {
			slog.Error("Cannot create Synonym.tsv file", "error", err)
			return err
		}
	}

	slog.Info("Exporting Vernacular file")
	path = filepath.Join(base, "Vernacular.tsv")
	err = t.clf.Vernacular(path)
	if err != nil {
		slog.Error("Cannot create Vernacular.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting NameRelation file")
	path = filepath.Join(base, "NameRelation.tsv")
	err = t.clf.NameRelation(path)
	if err != nil {
		slog.Error("Cannot create NameRelation.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting TypeMaterial file")
	path = filepath.Join(base, "TypeMaterial.tsv")
	err = t.clf.TypeMaterial(path)
	if err != nil {
		slog.Error("Cannot create TypeMaterial.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting Distribution file")
	path = filepath.Join(base, "Distribution.tsv")
	err = t.clf.Distribution(path)
	if err != nil {
		slog.Error("Cannot create Distribution.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting Media file")
	path = filepath.Join(base, "Media.tsv")
	err = t.clf.Media(path)
	if err != nil {
		slog.Error("Cannot create Media.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting Treatment file")
	path = filepath.Join(base, "Treatment.tsv")
	err = t.clf.Treatment(path)
	if err != nil {
		slog.Error("Cannot create Treatment.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting SpeciesEstimate file")
	path = filepath.Join(base, "SpeciesEstimate.tsv")
	err = t.clf.SpeciesEstimate(path)
	if err != nil {
		slog.Error("Cannot create SpeciesEstimate.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting TaxonProperty file")
	path = filepath.Join(base, "TaxonProperty.tsv")
	err = t.clf.TaxonProperty(path)
	if err != nil {
		slog.Error("Cannot create TaxonProperty.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting SpeciesInteraction file")
	path = filepath.Join(base, "SpeciesInteraction.tsv")
	err = t.clf.SpeciesInteraction(path)
	if err != nil {
		slog.Error("Cannot create SpeciesInteraction.tsv file", "error", err)
		return err
	}

	slog.Info("Exporting TaxonConceptRelation file")
	path = filepath.Join(base, "TaxonConceptRelation.tsv")
	err = t.clf.TaxonConceptRelation(path)
	if err != nil {
		slog.Error("Cannot create TaxonConceptRelation.tsv file", "error", err)
		return err
	}

	if !strings.HasSuffix(outputPath, ".zip") {
		outputPath += ".zip"
	}
	slog.Info("Creating Zip file", "path", outputPath)
	err = t.clf.CreateZip(outputPath)
	if err != nil {
		slog.Error("Cannot create zipped archive", "error", err)
		return err
	}

	slog.Info("Creation of CoLDP archive finished successfully")
	return nil
}
