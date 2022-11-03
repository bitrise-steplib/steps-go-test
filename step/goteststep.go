package step

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-go-test/filesystem"
	"github.com/bitrise-tools/go-steputils/tools"
)

type Step struct {
}

type Config struct {
}

type Result struct {
}

func CreateStep() Step {
	return Step{}
}

func (s Step) ProcessConfig() (*Config, error) {
	return &Config{}, nil
}

func (s Step) Run(config *Config) (*Result, error) {
	packages := os.Getenv("packages")

	log.Infof("Configs:")
	log.Printf("- packages: %s", packages)

	if packages == "" {
		return nil, errors.New("Required input not defined: packages")
	}

	log.Infof("\nRunning go test...")

	packageCodeCoveragePth, err := filesystem.CreatePackageCodeCoverageFile()
	if err != nil {
		return nil, err
	}

	codeCoveragePth, err := filesystem.CodeCoveragePath()
	if err != nil {
		return nil, err
	}

	for _, p := range strings.Split(packages, "\n") {
		cmd := command.NewWithStandardOuts("go", "test", "-v", "-race", "-coverprofile="+packageCodeCoveragePth, "-covermode=atomic", p)

		log.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("go test failed: %w", err)
		}

		if err := filesystem.AppendPackageCoverageAndRecreate(packageCodeCoveragePth, codeCoveragePth); err != nil {
			return nil, err
		}
	}

	if err := tools.ExportEnvironmentWithEnvman("GO_CODE_COVERAGE_REPORT_PATH", codeCoveragePth); err != nil {
		return nil, fmt.Errorf("Failed to export GO_CODE_COVERAGE_REPORT_PATH=%s", codeCoveragePth)
	}

	log.Donef("\ncode coverage is available at: GO_CODE_COVERAGE_REPORT_PATH=%s", codeCoveragePth)

	return &Result{}, nil
}

func (s Step) ExportOutputs(result *Result) error {
	return nil
}
