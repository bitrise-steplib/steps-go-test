package main

import (
	"os"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-go-test/step"
)

func main() {
	status, err := run()
	if err != nil {
		log.Errorf(err.Error())
	}
	os.Exit(int(status))
}

type RunStatus int

const (
	Success RunStatus = iota
	Failure
)

func run() (RunStatus, error) {
	step := step.CreateStep()
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
