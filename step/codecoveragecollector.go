package step

import (
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-go-test/filesystem"
)

type CodeCoverageCollector interface {
	PrepareAndReturnCoverageOutputPath(string) (string, error)
	CollectCoverageResultsAndReset() error
}

func CreateDefaultCodeCoverageCollector(logger log.Logger) CodeCoverageCollector {
	return &filesystembasedCollector{logger, "", ""}
}

type filesystembasedCollector struct {
	logger                      log.Logger
	amalgamatedCodeCoveragePath string
	codeCoverageOutputPath      string
}

func (c *filesystembasedCollector) PrepareAndReturnCoverageOutputPath(outputDir string) (string, error) {
	packageCodeCoveragePath, err := filesystem.CreatePackageCodeCoverageFile()
	if err != nil {
		return "", err
	}

	codeCoveragePath, err := filesystem.CodeCoveragePath(outputDir)
	if err != nil {
		return "", err
	}

	c.amalgamatedCodeCoveragePath = packageCodeCoveragePath
	c.codeCoverageOutputPath = codeCoveragePath
	return codeCoveragePath, nil
}

func (c *filesystembasedCollector) CollectCoverageResultsAndReset() error {
	if err := filesystem.AppendPackageCoverageAndRecreate(c.amalgamatedCodeCoveragePath, c.codeCoverageOutputPath, c.logger); err != nil {
		return err
	}
	return nil
}
