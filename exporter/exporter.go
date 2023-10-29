package exporter

type Exporter interface {
	Output() (string, error)
}
