package main

import (
	"os"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-go-test/step"
)

func main() {
	logger := log.NewLogger()
	status, err := run(logger)
	if err != nil {
		logger.Errorf(err.Error())
	}
	os.Exit(int(status))
}

type RunStatus int

const (
	Success RunStatus = iota
	Failure
)

func run(logger log.Logger) (RunStatus, error) {
	envRepository := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepository)
	testRunner := step.CreateDefaultTestRunner(envRepository, logger)

	step := step.CreateStep(envRepository, inputParser, logger, testRunner)
	config, err := step.ProcessConfig()
	if err != nil {
		return Failure, err
	}
	result, err := step.Run(config)
	if err != nil {
		return Failure, err
	}
	if err := step.ExportOutputs(result); err != nil {
		return Failure, err
	}
	return Success, nil
}
