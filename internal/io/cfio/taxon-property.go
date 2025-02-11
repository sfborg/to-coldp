package cfio

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) TaxonProperty(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	taxon_id, source_id, property, value, reference_id, page, ordinal,
	remarks, modified, modified_by
FROM taxon_property
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var tp coldp.TaxonProperty
		if count == 0 {
			err := writer.Write(tp.Headers())
			if err != nil {
				return err
			}
		}
		count++

		err = rows.Scan(
			&tp.TaxonID, &tp.SourceID, &tp.Property, &tp.Value, &tp.ReferenceID,
			&tp.Page, &tp.Ordinal, &tp.Remarks, &tp.Modified, &tp.ModifiedBy,
		)
		if err != nil {
			return err
		}

		var ord string
		if tp.Ordinal.Valid {
			ord = strconv.Itoa(int(tp.Ordinal.Int64))
		}

		row := []string{
			tp.TaxonID, tp.SourceID, tp.Property, tp.Value, tp.ReferenceID,
			tp.Page, ord, tp.Remarks, tp.Modified, tp.ModifiedBy,
		}
		err := writer.Write(row)
		if err != nil {
			return err
		}
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
