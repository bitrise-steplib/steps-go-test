package main

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/errorutil"
	"github.com/bitrise-io/go-utils/v2/fileutil"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/bitrise-steplib/steps-go-test/step"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.NewLogger()

	goTestRunner := createGoTestRunner(logger)
	config, err := goTestRunner.ProcessInputs()
	if err != nil {
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to process Step inputs: %w", err)))
		return 1
	}

	runOpts := step.RunOpts{Packages: config.Packages}
	runResult, err := goTestRunner.Run(runOpts)
	if err != nil {
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to execute Step: %w", err)))
		return 1
	}

	exportOpts := step.ExportOpts{CodeCoveragePth: runResult.CodeCoveragePth}
	if err := goTestRunner.ExportOutput(exportOpts); err != nil {
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to export Step outputs: %w", err)))
		return 1
	}

	return 0
}

func createGoTestRunner(logger log.Logger) step.GoTestRunner {
	envRepo := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepo)
	cmdFactory := command.NewFactory(envRepo)
	exporter := export.NewExporter(cmdFactory)
	pathProvider := pathutil.NewPathProvider()
	fileManager := fileutil.NewFileManager()

	return step.NewGoTestRunner(logger, inputParser, envRepo, cmdFactory, exporter, pathProvider, fileManager)
}
