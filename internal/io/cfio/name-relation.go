package cfio

import (
	"encoding/csv"
	"os"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) NameRelation(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	name_id, related_name_id, source_id, type_id, reference_id, remarks,
	modified, modified_by
FROM name_relation
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var nr coldp.NameRelation
		if count == 0 {
			err := writer.Write(nr.Headers())
			if err != nil {
				return err
			}
		}
		count++

		var typ string
		err = rows.Scan(
			&nr.NameID, &nr.RelatedNameID, &nr.SourceID, &typ, &nr.RelatedNameID,
			&nr.Remarks, &nr.Modified, &nr.ModifiedBy,
		)
		if err != nil {
			return err
		}

		nr.Type = coldp.NewNomRelType(typ)

		row := []string{
			nr.NameID, nr.RelatedNameID, nr.SourceID, nr.Type.String(),
			nr.RelatedNameID, nr.Remarks, nr.Modified, nr.ModifiedBy,
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
