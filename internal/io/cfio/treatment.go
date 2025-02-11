package cfio

import (
	"encoding/csv"
	"os"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) Treatment(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	taxon_id, source_id, document, format, modified, modified_by
FROM treatment
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var tr coldp.Treatment
		if count == 0 {
			err := writer.Write(tr.Headers())
			if err != nil {
				return err
			}
		}
		count++

		err = rows.Scan(
			&tr.TaxonID, &tr.SourceID, &tr.Document, &tr.Format,
			&tr.Modified, &tr.ModifiedBy,
		)
		if err != nil {
			return err
		}

		row := []string{
			tr.TaxonID, tr.SourceID, tr.Document, tr.Format,
			tr.Modified, tr.ModifiedBy,
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
