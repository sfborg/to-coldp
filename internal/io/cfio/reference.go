package cfio

import (
	"encoding/csv"
	"os"

	"github.com/gnames/coldp/ent/coldp"
)

func (c *cfio) Reference(path string) error {
	f, err := os.Create(path) // Create/open the file
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	q := `
SELECT
	col__id, col__alternative_id, col__source_id, col__citation, col__type_id,
	col__author, col__author_id, col__editor, col__editor_id, col__title,
	col__title_short, col__container_author, col__container_title,
	col__container_title_short, col__issued, col__accessed,
	col__collection_title, col__collection_editor, col__volume, col__issue,
	col__edition, col__page, col__publisher, col__publisher_place,
	col__version, col__isbn, col__issn, col__doi, col__link, col__remarks,
	col__modified, col__modified_by
FROM reference
`
	rows, err := c.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var ref coldp.Reference
		if count == 0 {
			err := writer.Write(ref.Headers())
			if err != nil {
				return err
			}
		}
		count++

		var tp string
		err = rows.Scan(
			&ref.ID, &ref.AlternativeID, &ref.SourceID, &ref.Citation,
			&tp, &ref.Author, &ref.AuthorID, &ref.Editor,
			&ref.EditorID, &ref.Title, &ref.TitleShort, &ref.ContainerAuthor,
			&ref.ContainerTitle, &ref.ContainerTitleShort,
			&ref.Issued, &ref.Accessed, &ref.CollectionTitle,
			&ref.CollectionEditor, &ref.Volume, &ref.Issue, &ref.Edition,
			&ref.Page, &ref.Publisher, &ref.PublisherPlace, &ref.Version,
			&ref.ISBN, &ref.ISSN, &ref.DOI, &ref.Link, &ref.Remarks,
			&ref.Modified, &ref.ModifiedBy,
		)
		if err != nil {
			return err
		}

		ref.Type = coldp.NewReferenceType(tp)

		row := []string{
			ref.ID, ref.AlternativeID, ref.SourceID, ref.Citation,
			ref.Type.String(), ref.Author, ref.AuthorID, ref.Editor,
			ref.EditorID, ref.Title, ref.TitleShort, ref.ContainerAuthor,
			ref.ContainerTitle, ref.ContainerTitleShort,
			ref.Issued, ref.Accessed, ref.CollectionTitle,
			ref.CollectionEditor, ref.Volume, ref.Issue, ref.Edition,
			ref.Page, ref.Publisher, ref.PublisherPlace, ref.Version,
			ref.ISBN, ref.ISSN, ref.DOI, ref.Link, ref.Remarks,
			ref.Modified, ref.ModifiedBy,
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
