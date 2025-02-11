package ent

import clcfg "github.com/gnames/coldp/config"

type ColDPFiles interface {
	Config() clcfg.Config
	CreateZip(path string) error

	Meta(path string) error
	Author(path string) error
	Reference(path string) error
	Name(path string) error
	Taxon(path string) error
	Synonym(path string) error
	Vernacular(path string) error
	NameRelation(path string) error
	TypeMaterial(path string) error
	Distribution(path string) error
	Media(path string) error
	Treatment(path string) error
	SpeciesEstimate(path string) error
	TaxonProperty(path string) error
	SpeciesInteraction(path string) error
	TaxonConceptRelation(path string) error
}
