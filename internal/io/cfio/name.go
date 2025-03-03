package cfio

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) Name(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	col__id, col__alternative_id, col__source_id, col__scientific_name,
	col__authorship, col__rank_id, col__uninomial, col__genus,
	col__infrageneric_epithet, col__specific_epithet,
	col__infraspecific_epithet, col__cultivar_epithet, col__notho_id,
	col__original_spelling, col__combination_authorship,
	col__combination_authorship_id, col__combination_ex_authorship,
	col__combination_ex_authorship_id, col__combination_authorship_year,
	col__basionym_authorship, col__basionym_authorship_id,
	col__basionym_ex_authorship, col__basionym_ex_authorship_id,
	col__basionym_authorship_year, col__code_id, col__status_id,
	col__reference_id, col__published_in_year, col__published_in_page,
	col__published_in_page_link, col__gender_id, col__gender_agreement,
	col__etymology, col__link, col__remarks, col__modified, col__modified_by
FROM name
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var n coldp.Name
		if count == 0 {
			err := writer.Write(n.Headers())
			if err != nil {
				return err
			}
		}
		count++

		var rank, notho, code, status, gender string
		err = rows.Scan(
			&n.ID, &n.AlternativeID, &n.SourceID, &n.ScientificName,
			&n.Authorship, &rank, &n.Uninomial, &n.Genus, &n.InfragenericEpithet,
			&n.SpecificEpithet, &n.InfraspecificEpithet, &n.CultivarEpithet,
			&notho, &n.OriginalSpelling, &n.CombinationAuthorship,
			&n.CombinationAuthorshipID, &n.CombinationExAuthorship,
			&n.CombinationExAuthorshipID, &n.CombinationAuthorshipYear,
			&n.BasionymAuthorship, &n.BasionymAuthorshipID,
			&n.BasionymExAuthorship, &n.BasionymExAuthorshipID,
			&n.BasionymAuthorshipYear, &code, &status, &n.ReferenceID,
			&n.PublishedInYear, &n.PublishedInPage, &n.PublishedInPageLink,
			&gender, &n.GenderAgreement, &n.Etymology, &n.Link,
			&n.Remarks, &n.Modified, &n.ModifiedBy,
		)
		if err != nil {
			return err
		}

		n.Rank = coldp.NewRank(rank)
		n.Notho = coldp.NewNamePart(notho)
		n.Code = coldp.NewNomCode(code)
		n.Status = coldp.NewNomStatus(status)
		n.Gender = coldp.NewGender(gender)

		var origSpelling, gndrAgr string
		if n.OriginalSpelling.Valid {
			origSpelling = strconv.FormatBool(n.OriginalSpelling.Bool)
		}
		if n.GenderAgreement.Valid {
			gndrAgr = strconv.FormatBool(n.GenderAgreement.Bool)
		}

		row := []string{
			n.ID, n.AlternativeID, n.SourceID, n.BasionymID, n.ScientificName,
			n.Authorship, n.Rank.String(), n.Uninomial, n.Genus,
			n.InfragenericEpithet, n.SpecificEpithet, n.InfraspecificEpithet,
			n.CultivarEpithet, n.Notho.String(), origSpelling,
			n.CombinationAuthorship, n.CombinationAuthorshipID,
			n.CombinationExAuthorship, n.CombinationExAuthorshipID,
			n.CombinationAuthorshipYear, n.BasionymAuthorship,
			n.BasionymAuthorshipID, n.BasionymExAuthorship, n.BasionymExAuthorshipID,
			n.BasionymAuthorshipYear, n.Code.String(), n.Status.String(),
			n.ReferenceID, n.PublishedInYear, n.PublishedInPage,
			n.PublishedInPageLink, n.Gender.String(), gndrAgr, n.Etymology, n.Link,
			n.Remarks, n.Modified, n.ModifiedBy,
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
