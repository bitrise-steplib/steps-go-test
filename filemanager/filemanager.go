package filemanager

import (
	"os"

	"github.com/bitrise-io/go-utils/v2/fileutil"
)

// TODO: Check what functions are used

// FileManager ...
// TODO: fileutil.FileManager interface should be extended with more methods
type FileManager interface {
	Create(name string) (*os.File, error)
	MkdirAll(path string, perm os.FileMode) error
	Open(path string) (*os.File, error)
	Write(path string, value string, perm os.FileMode) error
}

type fileManager struct {
	fileManager fileutil.FileManager
}

func New(manager fileutil.FileManager) FileManager {
	return fileManager{fileManager: manager}
}

func (f fileManager) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (f fileManager) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (f fileManager) Open(path string) (*os.File, error) {
	return f.fileManager.Open(path)
}

func (f fileManager) Write(path string, value string, perm os.FileMode) error {
	return f.fileManager.Write(path, value, perm)
}
