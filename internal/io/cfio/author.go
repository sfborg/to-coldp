package cfio

import (
	"encoding/csv"
	"os"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) Author(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	col__id, col__source_id, col__alternative_id, col__given, col__family,
	col__suffix, col__abbreviation_botany, col__alternative_names, col__sex_id,
	col__country, col__birth, col__birth_place, col__death, col__affiliation,
	col__interest, col__reference_id, col__link, col__remarks, col__modified,
	col__modified_by
FROM author
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var au coldp.Author
		if count == 0 {
			err := writer.Write(au.Headers())
			if err != nil {
				return err
			}
		}
		count++

		var sex string
		err = rows.Scan(
			&au.ID, &au.SourceID, &au.AlternativeID, &au.Given, &au.Family, &au.Suffix,
			&au.AbbreviationBotany, &au.AlternativeNames, &sex, &au.Country, &au.Birth,
			&au.BirthPlace, &au.Death, &au.Affiliation, &au.Interest, &au.ReferenceID,
			&au.Link, &au.Remarks, &au.Modified, &au.ModifiedBy,
		)

		if err != nil {
			return err
		}

		au.Sex = coldp.NewSex(sex)

		row := []string{
			au.ID, au.SourceID, au.AlternativeID, au.Given, au.Family,
			au.Suffix, au.AbbreviationBotany, au.AlternativeNames, au.Sex.String(),
			au.Country, au.Birth, au.BirthPlace, au.Death, au.Affiliation,
			au.Interest, au.ReferenceID, au.Link, au.Remarks,
			au.Modified, au.ModifiedBy,
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
