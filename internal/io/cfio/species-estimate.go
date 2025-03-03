package cfio

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) SpeciesEstimate(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	col__taxon_id, col__source_id, col__estimate, col__type_id,
	col__reference_id, col__remarks, col__modified, col__modified_by
FROM species_estimate
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var se coldp.SpeciesEstimate
		if count == 0 {
			err := writer.Write(se.Headers())
			if err != nil {
				return err
			}
		}
		count++

		var typ string
		err = rows.Scan(
			&se.TaxonID, &se.SourceID, &se.Estimate, &typ, &se.ReferenceID,
			&se.Remarks, &se.Modified, &se.ModifiedBy,
		)
		if err != nil {
			return err
		}

		se.Type = coldp.NewEstimateType(typ)
		var est string
		if se.Estimate.Valid {
			est = strconv.Itoa(int(se.Estimate.Int64))
		}

		row := []string{
			se.TaxonID, se.SourceID, est, se.Type.String(),
			se.ReferenceID, se.Remarks, se.Modified, se.ModifiedBy,
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
