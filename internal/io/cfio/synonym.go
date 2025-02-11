package cfio

import (
	"encoding/csv"
	"os"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) Synonym(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	id, taxon_id, source_id, name_id, name_phrase, according_to_id,
	status_id, reference_id, link, remarks, modified, modified_by
FROM synonym
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var syn coldp.Synonym
		if count == 0 {
			err := writer.Write(syn.Headers())
			if err != nil {
				return err
			}
		}
		count++

		var status string
		err = rows.Scan(
			&syn.ID, &syn.TaxonID, &syn.SourceID, &syn.NameID, &syn.NamePhrase,
			&syn.AccordingToID, &status, &syn.ReferenceID, &syn.Link,
			&syn.Remarks, &syn.Modified, &syn.ModifiedBy,
		)
		if err != nil {
			return err
		}

		syn.Status = coldp.NewTaxonomicStatus(status)

		row := []string{
			syn.ID, syn.TaxonID, syn.SourceID, syn.NameID, syn.NamePhrase,
			syn.AccordingToID, syn.Status.String(), syn.ReferenceID, syn.Link,
			syn.Remarks, syn.Modified, syn.ModifiedBy,
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
