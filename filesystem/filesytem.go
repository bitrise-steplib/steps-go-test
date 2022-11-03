package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

func CreatePackageCodeCoverageFile() (string, error) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("go-test")
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
	content, err := fileutil.ReadStringFromFile(packageCoveragePth)
	if err != nil {
		return fmt.Errorf("Failed to read package code coverage report file: %s", err)
	}

	if err := fileutil.AppendStringToFile(coveragePth, content); err != nil {
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
