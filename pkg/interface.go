package tocoldp

type ToCoLDP interface {
	// Exports SFGA archive to CoLDP archive
	// where path is the filepath where CoLDP file is being
	// created.
	Export(path string) error
}
