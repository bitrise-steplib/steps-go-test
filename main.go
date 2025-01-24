package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/fileutil"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
)

func failf(logger log.Logger, format string, args ...interface{}) {
	logger.Errorf(format, args...)
	os.Exit(1)
}

func createPackageCodeCoverageFile(pathProvider pathutil.PathProvider) (string, error) {
	tmpDir, err := pathProvider.CreateTempDir("go-test")
	if err != nil {
		return "", fmt.Errorf("failed to create tmp dir for code coverage reports: %w", err)
	}
	pth := filepath.Join(tmpDir, "profile.out")

	// TODO: don't call os.Create directly
	if _, err := os.Create(pth); err != nil {
		return "", err
	}

	return pth, nil
}

func codeCoveragePath(envRepo env.Repository) (string, error) {
	deployDir := envRepo.Get("BITRISE_DEPLOY_DIR")
	if deployDir == "" {
		return "", fmt.Errorf("BITRISE_DEPLOY_DIR env not set")
	}

	// TODO: don't call os.MkdirAll directly
	if err := os.MkdirAll(deployDir, 0777); err != nil {
		return "", fmt.Errorf("failed to create BITRISE_DEPLOY_DIR: %w", err)
	}

	return filepath.Join(deployDir, "go_code_coverage.txt"), nil
}

func appendPackageCoverageAndRecreate(fileManager fileutil.FileManager, packageCoveragePth, coveragePth string) error {
	// Read current package coverage report file
	packageCoverageFile, err := fileManager.Open(packageCoveragePth)
	if err != nil {
		return fmt.Errorf("failed to open package code coverage report file: %w", err)
	}
	packageCoverageFileContent, err := io.ReadAll(packageCoverageFile)
	if err != nil {
		return fmt.Errorf("failed to read package code coverage report file: %w", err)
	}

	// Append package coverage report to the main coverage report file
	coverageFile, err := fileManager.Open(coveragePth)
	if err != nil {
		return fmt.Errorf("failed to open package code coverage report file: %w", err)
	}
	coverageFileContent, err := io.ReadAll(coverageFile)
	if err != nil {
		return fmt.Errorf("failed to read package code coverage report file: %w", err)
	}

	coverageFileContent = append(coverageFileContent, packageCoverageFileContent...)

	if err := fileManager.Write(coveragePth, string(coverageFileContent), 0777); err != nil {
		return fmt.Errorf("failed to write package code coverage report file: %w", err)
	}

	// Recreate package coverage report file
	if err := fileManager.RemoveAll(packageCoveragePth); err != nil {
		return fmt.Errorf("failed to remove package code coverage report file: %w", err)
	}

	// TODO: don't call os.Create directly
	if _, err := os.Create(packageCoveragePth); err != nil {
		return fmt.Errorf("failed to create package code coverage report file: %w", err)
	}

	return nil
}

func main() {
	envRepo := env.NewRepository()
	cmdFactory := command.NewFactory(envRepo)
	exporter := export.NewExporter(cmdFactory)
	logger := log.NewLogger()
	pathProvider := pathutil.NewPathProvider()
	fileManager := fileutil.NewFileManager()

	packages := envRepo.Get("packages")

	logger.Infof("Configs:")
	logger.Printf("- packages: %s", packages)

	if packages == "" {
		failf(logger, "Required input not defined: packages")
	}

	logger.Infof("\nRunning go test...")

	packageCodeCoveragePth, err := createPackageCodeCoverageFile(pathProvider)
	if err != nil {
		failf(logger, err.Error())
	}

	codeCoveragePth, err := codeCoveragePath(envRepo)
	if err != nil {
		failf(logger, err.Error())
	}

	for _, p := range strings.Split(packages, "\n") {
		cmd := cmdFactory.Create("go", []string{"test", "-v", "-race", "-coverprofile=" + packageCodeCoveragePth, "-covermode=atomic", p}, &command.Opts{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		})
		logger.Printf("$ %s", cmd.PrintableCommandArgs())
		if err := cmd.Run(); err != nil {
			failf(logger, "go test failed: %s", err)
		}

		if err := appendPackageCoverageAndRecreate(fileManager, packageCodeCoveragePth, codeCoveragePth); err != nil {
			failf(logger, err.Error())
		}
	}

	if err := exporter.ExportOutput("GO_CODE_COVERAGE_REPORT_PATH", codeCoveragePth); err != nil {
		failf(logger, "Failed to export GO_CODE_COVERAGE_REPORT_PATH=%s", codeCoveragePth)
	}

	logger.Donef("\ncode coverage is available at: GO_CODE_COVERAGE_REPORT_PATH=%s", codeCoveragePth)
}
