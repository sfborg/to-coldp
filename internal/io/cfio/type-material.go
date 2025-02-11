package cfio

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) TypeMaterial(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	id, source_id, name_id, citation, status_id, institution_code, catalog_number,
	reference_id, locality, country, latitude, longitude, altitude, host,
	sex_id, date, collector, associated_sequences, link, remarks, modified,
	modified_by
FROM type_material
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var tm coldp.TypeMaterial
		if count == 0 {
			err := writer.Write(tm.Headers())
			if err != nil {
				return err
			}
		}
		count++

		var sex, status string
		err = rows.Scan(
			&tm.ID, &tm.SourceID, &tm.NameID, &tm.Citation, &status,
			&tm.InstitutionCode, &tm.CatalogNumber, &tm.ReferenceID, &tm.Locality,
			&tm.Country, &tm.Latitude, &tm.Longitude, &tm.Altitude, &tm.Host,
			&sex, &tm.Date, &tm.Collector, &tm.AssociatedSequences, &tm.Link,
			&tm.Remarks, &tm.Modified, &tm.ModifiedBy,
		)
		if err != nil {
			return err
		}

		tm.Sex = coldp.NewSex(sex)
		tm.Status = coldp.NewTypeStatus(status)

		var lat, long, alt string
		if tm.Latitude.Valid {
			lat = strconv.FormatFloat(tm.Latitude.Float64, 'f', 6, 64)
		}
		if tm.Longitude.Valid {
			long = strconv.FormatFloat(tm.Longitude.Float64, 'f', 6, 64)
		}
		if tm.Altitude.Valid {
			alt = strconv.Itoa(int(tm.Altitude.Int64))
		}

		row := []string{
			tm.ID, tm.SourceID, tm.NameID, tm.Citation, tm.Status.String(),
			tm.InstitutionCode, tm.CatalogNumber, tm.ReferenceID, tm.Locality,
			tm.Country, lat, long, alt, tm.Host, tm.Sex.String(), tm.Date,
			tm.Collector, tm.AssociatedSequences, tm.Link, tm.Remarks, tm.Modified,
			tm.ModifiedBy,
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
