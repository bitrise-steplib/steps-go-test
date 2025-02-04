package step

import (
	"io"
	"os"

	"github.com/bitrise-io/go-utils/v2/fileutil"
)

// OutputExporter ...
// TODO: export.NewExporter should return an interface
type OutputExporter interface {
	ExportOutput(key, value string) error
}

// FileManager ...
// TODO: fileutil.FileManager interface should be extended with more methods
type FileManager interface {
	Create(name string) (*os.File, error)
	MkdirAll(path string, perm os.FileMode) error
	Open(path string) (io.Reader, error)
	Write(path string, value string, perm os.FileMode) error
	RemoveAll(path string) error
}

type fileManager struct {
	fileManager fileutil.FileManager
}

func NewFileManager(manager fileutil.FileManager) FileManager {
	return fileManager{fileManager: manager}
}

func (f fileManager) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (f fileManager) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (f fileManager) Open(path string) (io.Reader, error) {
	return f.fileManager.Open(path)
}

func (f fileManager) Write(path string, value string, perm os.FileMode) error {
	return f.fileManager.Write(path, value, perm)
}

func (f fileManager) RemoveAll(path string) error {
	return f.fileManager.RemoveAll(path)
}
