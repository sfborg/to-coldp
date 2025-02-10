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
	id, alternative_id, source_id, scientific_name, authorship,
	rank_id, uninomial, genus, infrageneric_epithet, specific_epithet,
	infraspecific_epithet, cultivar_epithet, notho_id, original_spelling,
	combination_authorship, combination_authorship_id, 
	combination_ex_authorship, combination_ex_authorship_id,
	combination_authorship_year, basionym_authorship, basionym_authorship_id,
	basionym_ex_authorship, basionym_ex_authorship_id,
	basionym_authorship_year, code_id, status_id, reference_id, published_in_year,
	published_in_page, published_in_page_link, gender_id, gender_agreement,
	etymology, link, remarks, modified, modified_by
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

		origSpelling := strconv.FormatBool(n.OriginalSpelling.Bool)
		gndrAgr := strconv.FormatBool(n.GenderAgreement.Bool)

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

	// remove the file if it is empty
	if count == 0 {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}
