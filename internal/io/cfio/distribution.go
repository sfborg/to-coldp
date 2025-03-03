package cfio

import (
	"encoding/csv"
	"os"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) Distribution(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	col__taxon_id, col__source_id, col__area, col__area_id, col__gazetteer_id,
	col__status_id, col__reference_id, col__remarks, col__modified,
	col__modified_by
FROM distribution
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var dist coldp.Distribution
		if count == 0 {
			err := writer.Write(dist.Headers())
			if err != nil {
				return err
			}
		}
		count++

		var gztr, status string
		err = rows.Scan(
			&dist.TaxonID, &dist.SourceID, &dist.Area, &dist.AreaID,
			&gztr, &status, &dist.ReferenceID, &dist.Remarks,
			&dist.Modified, &dist.ModifiedBy,
		)
		if err != nil {
			return err
		}

		dist.Gazetteer = coldp.NewGazetteerEnt(gztr)
		dist.Status = coldp.NewDistrStatus(status)

		row := []string{
			dist.TaxonID, dist.SourceID, dist.Area, dist.AreaID,
			dist.Gazetteer.String(), dist.Status.String(), dist.ReferenceID,
			dist.Remarks, dist.Modified, dist.ModifiedBy,
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
