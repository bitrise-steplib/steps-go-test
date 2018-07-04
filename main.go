package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

func failf(format string, args ...interface{}) {
	log.Errorf(format, args...)
	os.Exit(1)
}

func main() {
	packages := os.Getenv("packages")
	codeCoveragePath := os.Getenv("code_coverage_path")

	log.Infof("Configs:")
	log.Printf("- packages: %s", packages)
	log.Printf("- code_coverage_path: %s", codeCoveragePath)

	if packages == "" {
		failf("Required input not defined: packages")
	}

	log.Infof("\nRunning go test...")

	args := []string{"test", "-v"}

	var packageCoverProfilePth string
	if codeCoveragePath != "" {
		if err := os.MkdirAll(filepath.Dir(codeCoveragePath), 0777); err != nil {
			failf("Failed to make dir: %s", err)
		}

		tmpDir, err := pathutil.NormalizedOSTempDirPath("go-test")
		if err != nil {
			failf("Failed to create tmp dir: %s", err)
		}
		packageCoverProfilePth = filepath.Join(tmpDir, "profile.out")

		args = append(args, "-race", "-coverprofile="+packageCoverProfilePth, "-covermode=atomic")
	}

	for _, p := range strings.Split(packages, "\n") {
		cmd := command.NewWithStandardOuts("go", append(args, p)...)

		log.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			failf("go test failed: %s", err)
		}

		if codeCoveragePath != "" {
			coverProfileContent, err := fileutil.ReadStringFromFile(packageCoverProfilePth)
			if err != nil {
				failf("Failed to read package cover profile: %s", err)
			}

			if err := fileutil.AppendStringToFile(codeCoveragePath, coverProfileContent); err != nil {
				failf("Failed to append package cover profile: %s", err)
			}

			if err := os.RemoveAll(packageCoverProfilePth); err != nil {
				failf("Failed to append package cover profile: %s", err)
			}
		}
	}
}
