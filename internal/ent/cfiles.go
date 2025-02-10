package ent

type ColDPFiles interface {
	Meta(path string) error
	Author(path string) error
	Reference(path string) error
	Name(path string) error
}
