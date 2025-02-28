package main

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/errorutil"
	"github.com/bitrise-io/go-utils/v2/exitcode"
	"github.com/bitrise-io/go-utils/v2/fileutil"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/bitrise-steplib/steps-go-test/filemanager"
	"github.com/bitrise-steplib/steps-go-test/step"
	"github.com/bitrise-steplib/steps-go-test/testaddon"
)

func main() {
	os.Exit(int(run()))
}

func run() exitcode.ExitCode {
	logger := log.NewLogger()

	goTestRunner := createGoTestRunner(logger)
	config, err := goTestRunner.ProcessInputs()
	if err != nil {
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to process Step inputs: %w", err)))
		return exitcode.Failure
	}

	runOpts := step.RunOpts{
		Package:     config.Package,
		Covermode:   config.Covermode,
		TestOptions: config.TestOptions,
		OutputDir:   config.OutputDir,
	}
	runResult, runErr := goTestRunner.Run(runOpts)
	if runResult == nil {
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to execute Step main logic: %w", runErr)))
		return exitcode.Failure
	}

	testReportName := config.TestReportName
	if testReportName == "" {
		testReportName = config.Package
	}
	exportOpts := step.ExportOpts{
		TestRunLogPth:   runResult.TestRunLogPath,
		CodeCoveragePth: runResult.CodeCoveragePth,
		TestReportName:  testReportName,
	}
	exportErr := goTestRunner.ExportOutput(exportOpts)

	if runErr != nil && exportErr != nil {
		logger.Warnf(errorutil.FormattedError(fmt.Errorf("Failed to export Step outputs: %w", exportErr)))
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to execute Step main logic: %w", runErr)))
		return exitcode.Failure
	}
	if runErr != nil {
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to execute Step main logic: %w", runErr)))
		return exitcode.Failure
	}
	if exportErr != nil {
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to export Step outputs: %w", exportErr)))
		return exitcode.Failure
	}

	return exitcode.Success
}

func createGoTestRunner(logger log.Logger) step.GoTestRunner {
	envRepo := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepo)
	cmdFactory := command.NewFactory(envRepo)
	exporter := export.NewExporter(cmdFactory)
	pathProvider := pathutil.NewPathProvider()
	fileManager := filemanager.New(fileutil.NewFileManager())
	testAddonExporter := testaddon.NewExporter(envRepo, fileManager)

	return step.NewGoTestRunner(logger, inputParser, envRepo, cmdFactory, &exporter, pathProvider, fileManager, testAddonExporter)
}
