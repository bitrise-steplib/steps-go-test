package step

// OutputExporter ...
// TODO: export.NewExporter should return an interface
type OutputExporter interface {
	ExportOutput(key, value string) error
}
