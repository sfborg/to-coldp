package cfio

import (
	"encoding/csv"
	"os"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) Media(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	col__taxon_id, col__source_id, col__url, col__type, col__format,
	col__title, col__created, col__creator, col__license, col__link,
	col__remarks, col__modified, col__modified_by
FROM media
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var md coldp.Media
		if count == 0 {
			err := writer.Write(md.Headers())
			if err != nil {
				return err
			}
		}
		count++

		err = rows.Scan(
			&md.TaxonID, &md.SourceID, &md.URL, &md.Type, &md.Format,
			&md.Title, &md.Created, &md.Creator, &md.License, &md.Link,
			&md.Remarks, &md.Modified, &md.ModifiedBy,
		)
		if err != nil {
			return err
		}

		row := []string{
			md.TaxonID, md.SourceID, md.URL, md.Type, md.Format,
			md.Title, md.Created, md.Creator, md.License, md.Link,
			md.Remarks, md.Modified, md.ModifiedBy,
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
