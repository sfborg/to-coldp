package cfio

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) Taxon(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	id, alternative_id, source_id, parent_id, ordinal, branch_length,
	name_id, name_phrase, according_to_id, according_to_page,
	according_to_page_link, scrutinizer, scrutinizer_id, status_id,
	extinct, temporal_range_start_id, temporal_range_end_id, environment_id,
	species, section, subgenus, genus, subtribe, tribe, subfamily, family,
	superfamily, suborder, "order", subclass, class, subphylum, phylum,
	kingdom, reference_id, link, remarks, modified, modified_by
FROM taxon
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var tx coldp.Taxon
		if count == 0 {
			err := writer.Write(tx.Headers())
			if err != nil {
				return err
			}
		}
		count++
		var status, start, end, env string

		err = rows.Scan(
			&tx.ID, &tx.AlternativeID, &tx.SourceID, &tx.ParentID, &tx.Ordinal,
			&tx.BranchLength, &tx.NameID, &tx.NamePhrase, &tx.AccordingToID,
			&tx.AccordingToPage, &tx.AccordingToPageLink, &tx.Scrutinizer,
			&tx.ScrutinizerID, &status, &tx.Extinct, &start, &end, &env, &tx.Species,
			&tx.Section, &tx.Subgenus, &tx.Genus, &tx.Subtribe, &tx.Tribe,
			&tx.Subfamily, &tx.Family, &tx.Superfamily, &tx.Suborder, &tx.Order,
			&tx.Subclass, &tx.Class, &tx.Subphylum, &tx.Phylum, &tx.Kingdom,
			&tx.ReferenceID, &tx.Link, &tx.Remarks, &tx.Modified, &tx.ModifiedBy,
		)
		if err != nil {
			return err
		}

		tx.TemporalRangeStart = coldp.NewGeoTime(start)
		tx.TemporalRangeEnd = coldp.NewGeoTime(end)
		tx.Environment = coldp.NewEnvironment(env)
		taxStatus := coldp.NewTaxonomicStatus(status)

		var ordinal, brLen, prov, extinct string
		if tx.Ordinal.Valid {
			ordinal = strconv.Itoa(int(tx.Ordinal.Int64))
		}
		if tx.BranchLength.Valid {
			brLen = strconv.Itoa(int(tx.BranchLength.Int64))
		}
		if taxStatus == coldp.ProvisionallyAcceptedTS {
			prov = "true"
		}
		if tx.Extinct.Valid {
			extinct = strconv.FormatBool(tx.Extinct.Bool)
		}

		row := []string{
			tx.ID, tx.AlternativeID, tx.SourceID, tx.ParentID, ordinal, brLen,
			tx.NameID, tx.NamePhrase, tx.AccordingToID, tx.AccordingToPage,
			tx.AccordingToPageLink, tx.Scrutinizer, tx.ScrutinizerID, prov, extinct,
			tx.TemporalRangeStart.String(), tx.TemporalRangeEnd.String(),
			tx.Environment.String(), tx.Species, tx.Section, tx.Subgenus, tx.Genus,
			tx.Subtribe, tx.Tribe, tx.Subfamily, tx.Family, tx.Superfamily,
			tx.Suborder, tx.Order, tx.Subclass, tx.Class, tx.Subphylum, tx.Phylum,
			tx.Kingdom, tx.ReferenceID, tx.Link, tx.Remarks, tx.Modified,
			tx.ModifiedBy,
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
