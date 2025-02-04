package coldpio

import cfiles "github.com/sfborg/to-coldp/internal/ent"

type coldpio struct {
}

func New() cfiles.ColDPFiles {
	res := coldpio{}
	return &res
}
