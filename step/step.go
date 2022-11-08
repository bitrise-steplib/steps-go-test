package step

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

const CodeCoverageReportPathExportName = "GO_CODE_COVERAGE_REPORT_PATH"

type Step struct {
	collector   CodeCoverageCollector
	env         env.Repository
	inputParser stepconf.InputParser
	logger      log.Logger
	testRunner  TestRunner
}

type Config struct {
	Packages []string
}

type Result struct {
	codeCoveragePath string
}

type envvars struct {
	packages string `env:"packages"`
}

func CreateStep(
	collector CodeCoverageCollector,
	env env.Repository,
	inputParser stepconf.InputParser,
	logger log.Logger,
	testRunner TestRunner) Step {
	return Step{
		collector:   collector,
		env:         env,
		inputParser: inputParser,
		logger:      logger,
		testRunner:  testRunner,
	}
}

func (s Step) ProcessConfig() (*Config, error) {
	var envvars envvars
	if err := s.inputParser.Parse(&envvars); err != nil {
		return nil, err
	}

	return &Config{
		Packages: strings.Split(envvars.packages, "\n"),
	}, nil
}

func (s Step) Run(config *Config) (*Result, error) {
	packages := config.Packages

	s.logger.Infof("Configs:")
	s.logger.Printf("- packages: %s", packages)

	if len(packages) == 0 {
		return nil, errors.New("Required input not defined: packages")
	}

	s.logger.Infof("\nRunning go test...")

	codeCoveragePath, err := s.collector.PrepareAndReturnCoverageOutputPath()
	if err != nil {
		return nil, err
	}

	for _, p := range config.Packages {
		testConfig := TestConfig{
			PackageName:        p,
			CoverageReportPath: codeCoveragePath,
		}

		if err := s.testRunner.RunTest(testConfig); err != nil {
			return nil, err
		}

		if err := s.collector.CollectCoverageResultsAndReset(); err != nil {
			return nil, err
		}
	}

	return &Result{
		codeCoveragePath,
	}, nil
}

func (s Step) ExportOutputs(result *Result) error {
	if err := s.env.Set(CodeCoverageReportPathExportName, result.codeCoveragePath); err != nil {
		return fmt.Errorf("Failed to export %s=%s", CodeCoverageReportPathExportName, result.codeCoveragePath)
	}

	s.logger.Donef("\ncode coverage is available at: %s=%s", CodeCoverageReportPathExportName, result.codeCoveragePath)
	return nil
}
