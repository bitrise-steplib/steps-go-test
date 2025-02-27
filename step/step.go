package step

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/bitrise-steplib/steps-go-test/filemanager"
	"github.com/bitrise-steplib/steps-go-test/testaddon"
	"github.com/jstemmer/go-junit-report/v2/junit"
	"github.com/jstemmer/go-junit-report/v2/parser/gotest"
)

type Inputs struct {
	Packages  string `env:"packages,required"`
	OutputDir string `env:"output_dir,required"`
}

type Config struct {
	Packages  []string
	OutputDir string
}

type GoTestRunner struct {
	logger            log.Logger
	inputParser       stepconf.InputParser
	envRepo           env.Repository
	cmdFactory        command.Factory
	outputExporter    OutputExporter
	pathProvider      pathutil.PathProvider
	fileManager       filemanager.FileManager
	testaddonExporter testaddon.Exporter
}

func NewGoTestRunner(
	logger log.Logger,
	inputParser stepconf.InputParser,
	envRepo env.Repository,
	cmdFactory command.Factory,
	outputExporter OutputExporter,
	pathProvider pathutil.PathProvider,
	fileManager filemanager.FileManager,
	testaddonExporter testaddon.Exporter,
) GoTestRunner {
	return GoTestRunner{
		logger:            logger,
		inputParser:       inputParser,
		envRepo:           envRepo,
		cmdFactory:        cmdFactory,
		outputExporter:    outputExporter,
		pathProvider:      pathProvider,
		fileManager:       fileManager,
		testaddonExporter: testaddonExporter,
	}
}

func (s GoTestRunner) ProcessInputs() (Config, error) {
	var inputs Inputs
	if err := s.inputParser.Parse(&inputs); err != nil {
		return Config{}, fmt.Errorf("issue with input: %w", err)
	}

	stepconf.Print(inputs)
	s.logger.Println()

	var packages []string
	packagesSplit := strings.Split(inputs.Packages, "\n")
	for _, p := range packagesSplit {
		p = strings.TrimSpace(p)
		if p != "" {
			packages = append(packages, p)
		}
	}

	return Config{
		Packages:  packages,
		OutputDir: inputs.OutputDir,
	}, nil
}

type RunOpts struct {
	Packages  []string
	OutputDir string
}

type RunResult struct {
	CodeCoveragePth     string
	TestRunLogFilePaths map[string]string
}

func (s GoTestRunner) Run(opts RunOpts) (*RunResult, error) {
	packageCodeCoveragePth, err := s.createPackageCodeCoverageFile()
	if err != nil {
		return nil, fmt.Errorf("failed to create package code coverage file: %w", err)
	}

	codeCoveragePth, err := s.codeCoveragePath(opts.OutputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get code coverage path: %w", err)
	}

	testRunLogFilePaths := map[string]string{}

	for _, p := range opts.Packages {
		logFile, err := s.testRunLogTmpFile()
		if err != nil {
			return nil, fmt.Errorf("failed to create tmp file for test run logs: %w", err)
		}
		defer func() {
			if err := logFile.Close(); err != nil {
				s.logger.Warnf("Failed to close test run log file: %s", err)
			}
		}()

		outWriter := io.MultiWriter(os.Stdout, logFile)
		errWriter := io.MultiWriter(os.Stderr, logFile)

		cmd := s.cmdFactory.Create("go", []string{"test", "-v", "-race", "-coverprofile=" + packageCodeCoveragePth, "-covermode=atomic", p}, &command.Opts{
			Stdout: outWriter,
			Stderr: errWriter,
		})
		s.logger.Printf("$ %s", cmd.PrintableCommandArgs())
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("go test failed: %w", err)
		}

		if err := s.appendPackageCoverageAndRecreate(packageCodeCoveragePth, codeCoveragePth); err != nil {
			return nil, fmt.Errorf("failed to append package coverage: %w", err)
		}

		testRunLogFilePaths[p] = logFile.Name()
	}

	return &RunResult{
		CodeCoveragePth:     codeCoveragePth,
		TestRunLogFilePaths: testRunLogFilePaths,
	}, nil
}

type ExportOpts struct {
	CodeCoveragePth     string
	TestRunLogFilePaths map[string]string
}

func (s GoTestRunner) ExportOutput(opts ExportOpts) error {
	if err := s.outputExporter.ExportOutput("GO_CODE_COVERAGE_REPORT_PATH", opts.CodeCoveragePth); err != nil {
		return fmt.Errorf("failed to export GO_CODE_COVERAGE_REPORT_PATH=%s: %w", opts.CodeCoveragePth, err)
	}

	s.logger.Donef("\ncode coverage is available at: GO_CODE_COVERAGE_REPORT_PATH=%s", opts.CodeCoveragePth)

	idx := 0
	for pkg, testRunLogFilePth := range opts.TestRunLogFilePaths {
		idx++

		testName := fmt.Sprintf("Test run #%d (%s)", idx, pkg)
		testName = s.testaddonExporter.ReplaceUnsupportedFilenameCharacters(testName)

		testRunLogFile, err := s.fileManager.Open(testRunLogFilePth)
		if err != nil {
			return fmt.Errorf("failed to open test run log file: %w", err)
		}

		report, err := gotest.NewParser().Parse(testRunLogFile)
		if err != nil {
			return fmt.Errorf("failed to parse go test report: %w", err)
		}
		reportFile, err := s.testaddonExporter.PrepareTestResultExport(testName)
		if err != nil {
			return fmt.Errorf("failed to prepare test result export: %w", err)
		}

		testsuits := junit.CreateFromReport(report, "")
		if err := testsuits.WriteXML(reportFile); err != nil {
			return fmt.Errorf("failed to write test report: %w", err)
		}
		s.logger.Printf("\nTest report is available at: %s", reportFile.Name())
	}

	return nil
}

func (s GoTestRunner) testRunLogTmpFile() (*os.File, error) {
	tmpDir, err := s.pathProvider.CreateTempDir("go-test")
	if err != nil {
		return nil, fmt.Errorf("failed to create tmp dir for test run logs: %w", err)
	}
	pth := filepath.Join(tmpDir, "test_run.log")

	f, err := s.fileManager.Create(pth)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (s GoTestRunner) createPackageCodeCoverageFile() (string, error) {
	tmpDir, err := s.pathProvider.CreateTempDir("go-test")
	if err != nil {
		return "", fmt.Errorf("failed to create tmp dir for code coverage reports: %w", err)
	}
	pth := filepath.Join(tmpDir, "profile.out")

	if _, err := s.fileManager.Create(pth); err != nil {
		return "", err
	}

	return pth, nil
}

func (s GoTestRunner) codeCoveragePath(outputDir string) (string, error) {
	if err := s.fileManager.MkdirAll(outputDir, 0777); err != nil {
		return "", fmt.Errorf("failed to create BITRISE_DEPLOY_DIR: %w", err)
	}

	return filepath.Join(outputDir, "go_code_coverage.txt"), nil
}

func (s GoTestRunner) appendPackageCoverageAndRecreate(packageCoveragePth, coveragePth string) error {
	// Read current package coverage report file
	packageCoverageFile, err := s.fileManager.Open(packageCoveragePth)
	if err != nil {
		return fmt.Errorf("failed to open package code coverage report file: %w", err)
	}
	packageCoverageFileContent, err := io.ReadAll(packageCoverageFile)
	if err != nil {
		return fmt.Errorf("failed to read package code coverage report file: %w", err)
	}

	// Append package coverage report to the main coverage report file
	var coverageFileContent []byte
	coverageFile, err := s.fileManager.Open(coveragePth)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to open code coverage report file: %w", err)
		}
	} else {
		coverageFileContent, err = io.ReadAll(coverageFile)
		if err != nil {
			return fmt.Errorf("failed to read code coverage report file: %w", err)
		}
	}

	coverageFileContent = append(coverageFileContent, packageCoverageFileContent...)

	if err := s.fileManager.Write(coveragePth, string(coverageFileContent), 0777); err != nil {
		return fmt.Errorf("failed to write code coverage report file: %w", err)
	}

	// Recreate package coverage report file
	if err := s.fileManager.RemoveAll(packageCoveragePth); err != nil {
		return fmt.Errorf("failed to remove package code coverage report file: %w", err)
	}

	if _, err := s.fileManager.Create(packageCoveragePth); err != nil {
		return fmt.Errorf("failed to create package code coverage report file: %w", err)
	}

	return nil
}
