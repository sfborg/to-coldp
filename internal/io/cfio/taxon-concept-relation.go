package cfio

import (
	"encoding/csv"
	"os"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) TaxonConceptRelation(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	taxon_id, related_taxon_id, source_id, type, reference_id,
	remarks, modified, modified_by
FROM taxon_concept_relation
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var tcr coldp.TaxonConceptRelation
		if count == 0 {
			err := writer.Write(tcr.Headers())
			if err != nil {
				return err
			}
		}
		count++

		var typ string
		err = rows.Scan(
			&tcr.TaxonID, &tcr.RelatedTaxonID, &tcr.SourceID, &typ,
			&tcr.ReferenceID, &tcr.Remarks, &tcr.Modified, &tcr.ModifiedBy,
		)
		if err != nil {
			return err
		}

		tcr.Type = coldp.NewTaxonConceptRelType(typ)

		row := []string{
			tcr.TaxonID, tcr.RelatedTaxonID, tcr.SourceID, tcr.Type.String(),
			tcr.ReferenceID, tcr.Remarks, tcr.Modified, tcr.ModifiedBy,
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
