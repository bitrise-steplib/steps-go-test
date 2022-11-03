package step

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-go-test/filesystem"
)

const CodeCoverageReportPathExportName = "GO_CODE_COVERAGE_REPORT_PATH"

type Step struct {
	env    env.Repository
	logger log.Logger
}

type Config struct {
}

type Result struct {
	codeCoveragePath string
}

func CreateStep(env env.Repository, logger log.Logger) Step {
	return Step{
		env:    env,
		logger: logger,
	}
}

func (s Step) ProcessConfig() (*Config, error) {
	return &Config{}, nil
}

func (s Step) Run(config *Config) (*Result, error) {
	packages := os.Getenv("packages")

	s.logger.Infof("Configs:")
	s.logger.Printf("- packages: %s", packages)

	if packages == "" {
		return nil, errors.New("Required input not defined: packages")
	}

	s.logger.Infof("\nRunning go test...")

	packageCodeCoveragePth, err := filesystem.CreatePackageCodeCoverageFile()
	if err != nil {
		return nil, err
	}

	codeCoveragePth, err := filesystem.CodeCoveragePath()
	if err != nil {
		return nil, err
	}

	for _, p := range strings.Split(packages, "\n") {
		args := []string{
			"test",
			"-v",
			"-race",
			"-coverageprofile=" + packageCodeCoveragePth,
			"-covermode=atomic",
			p,
		}
		cmd := command.NewFactory(s.env).Create("go", args, nil)

		s.logger.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("go test failed: %w", err)
		}

		if err := filesystem.AppendPackageCoverageAndRecreate(packageCodeCoveragePth, codeCoveragePth); err != nil {
			return nil, err
		}
	}

	return &Result{
		codeCoveragePth,
	}, nil
}

func (s Step) ExportOutputs(result *Result) error {
	if err := s.env.Set(CodeCoverageReportPathExportName, result.codeCoveragePath); err != nil {
		return fmt.Errorf("Failed to export %s=%s", CodeCoverageReportPathExportName, result.codeCoveragePath)
	}

	s.logger.Donef("\ncode coverage is available at: %s=%s", CodeCoverageReportPathExportName, result.codeCoveragePath)
	return nil
}
