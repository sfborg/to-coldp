package ent

type ColDPFiles interface {
	Meta(path string) error
	Reference(path string) error
}
