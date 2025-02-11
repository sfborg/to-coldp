package cfio

import (
	"encoding/csv"
	"os"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) SpeciesInteraction(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	taxon_id, related_taxon_id, source_id, related_taxon_scientific_name,
	type, reference_id, remarks, modified, modified_by
FROM species_interaction
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var si coldp.SpeciesInteraction
		if count == 0 {
			err := writer.Write(si.Headers())
			if err != nil {
				return err
			}
		}
		count++

		var typ string
		err = rows.Scan(
			&si.TaxonID, &si.RelatedTaxonID, &si.SourceID,
			&si.RelatedTaxonScientificName, &typ, &si.ReferenceID, &si.Remarks,
			&si.Modified, &si.ModifiedBy,
		)
		if err != nil {
			return err
		}

		si.Type = coldp.NewSpInteractionType(typ)

		row := []string{
			si.TaxonID, si.RelatedTaxonID, si.SourceID,
			si.RelatedTaxonScientificName, si.Type.String(), si.ReferenceID,
			si.Remarks, si.Modified, si.ModifiedBy,
		}
		err := writer.Write(row)
		if err != nil {
			return err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	// remove the file if it is empty
	if count == 0 {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}
