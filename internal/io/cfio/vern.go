package cfio

import (
	"encoding/csv"
	"os"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) Vernacular(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	taxon_id, source_id, name, transliteration, language, preferred,
	country, area, sex_id, reference_id, remarks, modified, modified_by
FROM vernacular
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var vrn coldp.Vernacular
		if count == 0 {
			err := writer.Write(vrn.Headers())
			if err != nil {
				return err
			}
		}
		count++

		err = rows.Scan()
		if err != nil {
			return err
		}

		row := []string{}
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
