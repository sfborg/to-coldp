package ent

type ColDPFiles interface {
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
}
