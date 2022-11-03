package filesystem

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/v2/pathutil"
)

func CreatePackageCodeCoverageFile() (string, error) {
	tmpDir, err := pathutil.NewPathProvider().CreateTempDir("go-test")
	if err != nil {
		return "", fmt.Errorf("Failed to create tmp dir for code coverage reports: %s", err)
	}
	pth := filepath.Join(tmpDir, "profile.out")
	if _, err := os.Create(pth); err != nil {
		return "", err
	}
	return pth, nil
}

func CodeCoveragePath() (string, error) {
	deployDir := os.Getenv("BITRISE_DEPLOY_DIR")
	if deployDir == "" {
		return "", fmt.Errorf("BITRISE_DEPLOY_DIR env not set")
	}
	if err := os.MkdirAll(deployDir, 0777); err != nil {
		return "", fmt.Errorf("Failed to create BITRISE_DEPLOY_DIR: %s", err)
	}
	return filepath.Join(deployDir, "go_code_coverage.txt"), nil
}

func AppendPackageCoverageAndRecreate(packageCoveragePth, coveragePth string) error {
	bytes, err := os.ReadFile(packageCoveragePth)
	if err != nil {
		return fmt.Errorf("Failed to read package code coverage report file: %s", err)
	}

	if err := appendBytesToFile(coveragePth, bytes); err != nil {
		return fmt.Errorf("Failed to append package code coverage report: %s", err)
	}

	if err := os.RemoveAll(packageCoveragePth); err != nil {
		return fmt.Errorf("Failed to remove package code coverage report file: %s", err)
	}
	if _, err := os.Create(packageCoveragePth); err != nil {
		return fmt.Errorf("Failed to create package code coverage report file: %s", err)
	}
	return nil
}

// Stolen from v1 of fileutil
// https://github.com/bitrise-io/go-utils/blob/9e20aaef213f7fe1e50290b4a5f78edb1e518713/fileutil/fileutil.go
func appendBytesToFile(pth string, fileCont []byte) error {
	if pth == "" {
		return errors.New("No path provided")
	}

	var file *os.File
	filePerm, err := os.Lstat(pth)
	if err != nil {
		// create the file
		file, err = os.Create(pth)
	} else {
		// open for append
		file, err = os.OpenFile(pth, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePerm.Mode())
	}
	if err != nil {
		// failed to create or open-for-append the file
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(" [!] Failed to close file:", err)
		}
	}()

	if _, err := file.Write(fileCont); err != nil {
		return err
	}

	return nil
}
