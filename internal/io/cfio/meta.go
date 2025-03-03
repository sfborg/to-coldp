package cfio

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gnames/coldp/ent/coldp"
	"github.com/gnames/gnlib"
)

func (c *cfio) Meta(path string) error {
	res, err := c.dbMeta()
	if err != nil {
		return err
	}

	bs, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bs, 0644)
	return nil
}

func (c *cfio) dbMeta() (*coldp.Meta, error) {
	q := `
SELECT
	col__doi, col__title, col__alias, col__description, col__issued,
	col__version, col__keywords, col__geographic_scope, col__taxonomic_scope,
	col__temporal_scope, col__confidence, col__completeness, col__license,
	col__url, col__logo, col__label, col__citation, col__private
  FROM metadata LIMIT 1
`
	res := coldp.Meta{}
	row := c.db.QueryRow(q)
	var keywords string
	err := row.Scan(
		&res.DOI, &res.Title, &res.Alias, &res.Description,
		&res.Issued, &res.Version, &keywords,
		&res.GeographicScope, &res.TaxonomicScope, &res.TemporalScope,
		&res.Confidence, &res.Completeness, &res.License, &res.URL,
		&res.Logo, &res.Label, &res.Citation, &res.Private,
	)
	if err != nil {
		return nil, err
	}
	if keywords != "" {
		words := strings.Split(keywords, ",")
		words = gnlib.Map(words, func(word string) string {
			return strings.TrimSpace(word)
		})
		res.Keywords = words
	}
	res.Contact, err = c.getActor("contact")
	if err != nil {
		return nil, err
	}
	res.Publisher, err = c.getActor("publisher")
	if err != nil {
		return nil, err
	}
	res.Editors, err = c.getActors("editor")
	if err != nil {
		return nil, err
	}
	res.Creators, err = c.getActors("creator")
	if err != nil {
		return nil, err
	}
	res.Contributors, err = c.getActors("contributor")
	if err != nil {
		return nil, err
	}
	res.Sources, err = c.getSources()
	return &res, nil
}

func (c *cfio) getActor(table string) (*coldp.Actor, error) {
	q := fmt.Sprintf(`
SELECT
		col__orcid, col__given, col__family, col__rorid, col__organisation,
	  col__email, col__url, col__note
	FROM %s
	LIMIT 1
`, table)
	row := c.db.QueryRow(q)
	var res coldp.Actor

	err := row.Scan(
		&res.Orcid, &res.Given, &res.Family, &res.RorID, &res.Organization,
		&res.Email, &res.URL, &res.Note,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *cfio) getActors(table string) ([]coldp.Actor, error) {
	var res []coldp.Actor
	q := fmt.Sprintf(`
SELECT
		col__orcid, col__given, col__family, col__rorid, col__organisation,
		col__email, col__url, col__note
	FROM %s
`, table)
	rows, err := c.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var act coldp.Actor
		err := rows.Scan(
			&act.Orcid, &act.Given, &act.Family, &act.RorID, &act.Organization,
			&act.Email, &act.URL, &act.Note,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, act)
	}

	return res, nil
}

func (c *cfio) getSources() ([]coldp.Source, error) {
	var res []coldp.Source
	q := `
SELECT
		col__type, col__title, col__authors, col__issued, col__:wisbn
	FROM source
`
	rows, err := c.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var src coldp.Source
		var authors string
		err = rows.Scan(&src.Type, &src.Title, &authors, &src.Issued, &src.ISBN)
		if err != nil {
			return nil, err
		}
		es := strings.Split(authors, ",|")
		auAry := gnlib.Map(es, func(s string) any {
			s = strings.TrimSpace(s)
			var res any = s
			return res
		})
		src.Authors = auAry
		res = append(res, src)
	}
	return res, nil
}
