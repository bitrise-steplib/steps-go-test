package step

import (
	"fmt"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

type TestConfig struct {
	PackageName        string
	CoverageReportPath string
}

type TestRunner interface {
	RunTest(config TestConfig) error
}

func CreateDefaultTestRunner(env env.Repository, logger log.Logger) TestRunner {
	return &gocliTestRunner{env, logger}
}

type gocliTestRunner struct {
	env    env.Repository
	logger log.Logger
}

func (r gocliTestRunner) RunTest(config TestConfig) error {
	args := []string{
		"test",
		"-v",
		"-race",
		"-coverprofile=" + config.CoverageReportPath,
		"-covermode=atomic",
		config.PackageName,
	}
	cmd := command.NewFactory(r.env).Create("go", args, nil)

	r.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go test failed: %w", err)
	}
	return nil
}
