package main

import (
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/tools"

	"github.com/bitrise-steplib/steps-go-test/filesystem"
)

func failf(format string, args ...interface{}) {
	log.Errorf(format, args...)
	os.Exit(1)
}

func main() {
	status, _ := run()
	// if err != nil {
	// log.Errorf(err)
	// }
	os.Exit(int(status))
}

type RunStatus int

const (
	Success RunStatus = iota
	Failure
)

func run() (RunStatus, error) {
	packages := os.Getenv("packages")

	log.Infof("Configs:")
	log.Printf("- packages: %s", packages)

	if packages == "" {
		failf("Required input not defined: packages")
	}

	log.Infof("\nRunning go test...")

	packageCodeCoveragePth, err := filesystem.CreatePackageCodeCoverageFile()
	if err != nil {
		return Failure, err
	}

	codeCoveragePth, err := filesystem.CodeCoveragePath()
	if err != nil {
		failf(err.Error())
	}

	for _, p := range strings.Split(packages, "\n") {
		cmd := command.NewWithStandardOuts("go", "test", "-v", "-race", "-coverprofile="+packageCodeCoveragePth, "-covermode=atomic", p)

		log.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			failf("go test failed: %s", err)
		}

		if err := filesystem.AppendPackageCoverageAndRecreate(packageCodeCoveragePth, codeCoveragePth); err != nil {
			failf(err.Error())
		}
	}

	if err := tools.ExportEnvironmentWithEnvman("GO_CODE_COVERAGE_REPORT_PATH", codeCoveragePth); err != nil {
		failf("Failed to export GO_CODE_COVERAGE_REPORT_PATH=%s", codeCoveragePth)
	}

	log.Donef("\ncode coverage is available at: GO_CODE_COVERAGE_REPORT_PATH=%s", codeCoveragePth)
	return Success, nil
}
