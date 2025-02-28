package step

import (
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
	shellquote "github.com/kballard/go-shellquote"
)

type Inputs struct {
	Package        string `env:"package,required"`
	TestOptions    string `env:"test_options"`
	TestReportName string `env:"test_report_name"`
	OutputDir      string `env:"output_dir,required"`
}

type Config struct {
	Package        string
	TestOptions    []string
	TestReportName string
	OutputDir      string
}

type GoTestRunner struct {
	logger            log.Logger
	inputParser       stepconf.InputParser
	envRepo           env.Repository
	cmdFactory        command.Factory
	outputExporter    OutputExporter
	pathProvider      pathutil.PathProvider
	fileManager       filemanager.FileManager
	testAddonExporter testaddon.Exporter
}

func NewGoTestRunner(
	logger log.Logger,
	inputParser stepconf.InputParser,
	envRepo env.Repository,
	cmdFactory command.Factory,
	outputExporter OutputExporter,
	pathProvider pathutil.PathProvider,
	fileManager filemanager.FileManager,
	testAddonExporter testaddon.Exporter,
) GoTestRunner {
	return GoTestRunner{
		logger:            logger,
		inputParser:       inputParser,
		envRepo:           envRepo,
		cmdFactory:        cmdFactory,
		outputExporter:    outputExporter,
		pathProvider:      pathProvider,
		fileManager:       fileManager,
		testAddonExporter: testAddonExporter,
	}
}

func (s GoTestRunner) ProcessInputs() (*Config, error) {
	var inputs Inputs
	if err := s.inputParser.Parse(&inputs); err != nil {
		return nil, fmt.Errorf("issue with input: %w", err)
	}

	stepconf.Print(inputs)
	s.logger.Println()

	pkg := strings.TrimSpace(inputs.Package)
	testReportName := strings.TrimSpace(inputs.TestReportName)
	testOptions, err := shellquote.Split(inputs.TestOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test options: %w", err)
	}

	return &Config{
		Package:        pkg,
		TestOptions:    testOptions,
		TestReportName: testReportName,
		OutputDir:      inputs.OutputDir,
	}, nil
}

type RunOpts struct {
	Package     string
	TestOptions []string
	OutputDir   string
}

type RunResult struct {
	CodeCoveragePth string
	TestRunLogPath  string
}

func (s GoTestRunner) Run(opts RunOpts) (*RunResult, error) {
	codeCoverageFile, err := s.createCodeCoverageFile(opts.OutputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create package code coverage file: %w", err)
	}
	defer func() {
		if err := codeCoverageFile.Close(); err != nil {
			s.logger.Warnf("Failed to close code coverage file: %s", err)
		}
	}()

	testRunLogFile, err := s.testRunLogFile(opts.OutputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create tmp file for test run logs: %w", err)
	}
	defer func() {
		if err := testRunLogFile.Close(); err != nil {
			s.logger.Warnf("Failed to close test run log file: %s", err)
		}
	}()

	outWriter := io.MultiWriter(os.Stdout, testRunLogFile)
	errWriter := io.MultiWriter(os.Stderr, testRunLogFile)

	args := []string{"test", "-v", "-coverprofile=" + codeCoverageFile.Name(), "-covermode=atomic"}
	if len(opts.TestOptions) > 0 {
		args = append(args, opts.TestOptions...)
	}
	args = append(args, opts.Package)

	cmd := s.cmdFactory.Create("go", args, &command.Opts{
		Stdout: outWriter,
		Stderr: errWriter,
	})
	s.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return &RunResult{
		CodeCoveragePth: codeCoverageFile.Name(),
		TestRunLogPath:  testRunLogFile.Name(),
	}, cmd.Run()
}

type ExportOpts struct {
	TestRunLogPth   string
	CodeCoveragePth string
	TestReportName  string
}

func (s GoTestRunner) ExportOutput(opts ExportOpts) error {
	if err := s.outputExporter.ExportOutput("GO_TEST_RUN_LOG_PATH", opts.TestRunLogPth); err != nil {
		return fmt.Errorf("failed to export GO_TEST_RUN_LOG_PATH=%s: %w", opts.TestRunLogPth, err)
	}
	s.logger.Donef("\ngo test run log file is available at: GO_TEST_RUN_LOG_PATH=%s", opts.TestRunLogPth)

	if err := s.outputExporter.ExportOutput("GO_CODE_COVERAGE_REPORT_PATH", opts.CodeCoveragePth); err != nil {
		return fmt.Errorf("failed to export GO_CODE_COVERAGE_REPORT_PATH=%s: %w", opts.CodeCoveragePth, err)
	}
	s.logger.Donef("code coverage file is available at: GO_CODE_COVERAGE_REPORT_PATH=%s", opts.CodeCoveragePth)

	testRunLogFile, err := s.fileManager.Open(opts.TestRunLogPth)
	if err != nil {
		return fmt.Errorf("failed to open test run log file: %w", err)
	}
	defer func() {
		if err := testRunLogFile.Close(); err != nil {
			s.logger.Warnf("Failed to close test run log file: %s", err)
		}
	}()

	report, err := gotest.NewParser().Parse(testRunLogFile)
	if err != nil {
		return fmt.Errorf("failed to parse go test report: %w", err)
	}

	reportFile, err := s.testAddonExporter.PrepareTestResultExport(opts.TestReportName)
	if err != nil {
		return fmt.Errorf("failed to prepare test result export: %w", err)
	}

	testSuits := junit.CreateFromReport(report, "")
	if err := testSuits.WriteXML(reportFile); err != nil {
		return fmt.Errorf("failed to write test report: %w", err)
	}

	s.logger.Donef("test report exported, report file is available at: %s", reportFile.Name())

	return nil
}

func (s GoTestRunner) testRunLogFile(outputDir string) (*os.File, error) {
	// TODO: do this earlier
	if err := s.fileManager.MkdirAll(outputDir, 0777); err != nil {
		return nil, fmt.Errorf("failed to create BITRISE_DEPLOY_DIR: %w", err)
	}

	pth := filepath.Join(outputDir, "go_test_run.log")
	f, err := s.fileManager.Create(pth)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (s GoTestRunner) createCodeCoverageFile(outputDir string) (*os.File, error) {
	// TODO: do this earlier
	if err := s.fileManager.MkdirAll(outputDir, 0777); err != nil {
		return nil, fmt.Errorf("failed to create BITRISE_DEPLOY_DIR: %w", err)
	}

	pth := filepath.Join(outputDir, "go_code_coverage.out")
	f, err := s.fileManager.Create(pth)
	if err != nil {
		return nil, err
	}

	return f, nil
}
