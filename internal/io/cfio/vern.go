package cfio

import (
	"encoding/csv"
	"os"
	"strconv"

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

		var sex string
		err = rows.Scan(
			&vrn.TaxonID, &vrn.SourceID, &vrn.Name, &vrn.Transliteration,
			&vrn.Language, &vrn.Preferred, &vrn.Country, &vrn.Area,
			&sex, &vrn.ReferenceID, &vrn.Remarks, &vrn.Modified, &vrn.ModifiedBy,
		)
		if err != nil {
			return err
		}
		vrn.Sex = coldp.NewSex(sex)
		var pref string
		if vrn.Preferred.Valid {
			pref = strconv.FormatBool(vrn.Preferred.Bool)
		}

		row := []string{
			vrn.TaxonID, vrn.SourceID, vrn.Name, vrn.Transliteration, vrn.Language,
			pref, vrn.Country, vrn.Area, vrn.Sex.String(), vrn.ReferenceID,
			vrn.Remarks, vrn.Modified, vrn.ModifiedBy,
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
