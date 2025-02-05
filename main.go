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
	"github.com/bitrise-steplib/steps-go-test/step"
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

	runOpts := step.RunOpts(config)
	runResult, err := goTestRunner.Run(runOpts)
	if err != nil {
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to execute Step: %w", err)))
		return exitcode.Failure
	}

	exportOpts := step.ExportOpts{CodeCoveragePth: runResult.CodeCoveragePth}
	if err := goTestRunner.ExportOutput(exportOpts); err != nil {
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to export Step outputs: %w", err)))
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
	fileManager := step.NewFileManager(fileutil.NewFileManager())

	return step.NewGoTestRunner(logger, inputParser, envRepo, cmdFactory, &exporter, pathProvider, fileManager)
}
