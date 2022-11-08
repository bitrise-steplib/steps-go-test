package step

import (
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-go-test/filesystem"
)

type CodeCoverageCollector interface {
	PrepareAndReturnCurrentPackageCoverageOutputPath(string) (string, error)
	CollectCoverageResultsAndReset() error
	FinishCollectionAndReturnPathToCollectedResults() string
}

func CreateDefaultCodeCoverageCollector(logger log.Logger) CodeCoverageCollector {
	return &filesystembasedCollector{logger, "", ""}
}

type filesystembasedCollector struct {
	logger                               log.Logger
	amalgamatedCodeCoveragePath          string
	currentPackageCodeCoverageOutputPath string
}

func (c *filesystembasedCollector) PrepareAndReturnCurrentPackageCoverageOutputPath(outputDir string) (string, error) {
	packageCodeCoveragePath, err := filesystem.CreatePackageCodeCoverageFile()
	if err != nil {
		return "", err
	}

	codeCoveragePath, err := filesystem.CodeCoveragePath(outputDir)
	if err != nil {
		return "", err
	}

	c.amalgamatedCodeCoveragePath = codeCoveragePath
	c.currentPackageCodeCoverageOutputPath = packageCodeCoveragePath
	return codeCoveragePath, nil
}

func (c *filesystembasedCollector) CollectCoverageResultsAndReset() error {
	if err := filesystem.AppendPackageCoverageAndRecreate(c.currentPackageCodeCoverageOutputPath, c.amalgamatedCodeCoveragePath, c.logger); err != nil {
		return err
	}
	return nil
}

func (c *filesystembasedCollector) FinishCollectionAndReturnPathToCollectedResults() string {
	return c.amalgamatedCodeCoveragePath
}
