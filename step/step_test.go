package step

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-go-test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGoTestRunner_Run_WhenTestSucceedItWritesCodeCoverageToFile(t *testing.T) {
	// Create run options
	packages := []string{"./..."}
	outputDir := "output"
	opts := RunOpts{
		Packages:  packages,
		OutputDir: outputDir,
	}

	// Expected result
	wantCoveragePth := filepath.Join(outputDir, "go_code_coverage.txt")
	wantRunResult := &RunResult{
		CodeCoveragePth: wantCoveragePth,
	}

	// Create mocks
	mockCmdFactory := new(mocks.Factory)
	mockCmd := new(mocks.Command)
	mockFileManager := new(mocks.FileManager)
	mockPathProvider := new(mocks.PathProvider)

	// It runs go test command with code coverage enabled
	tmpDir := "tmp_dir"
	packageCoveragePth := filepath.Join(tmpDir, "profile.out")
	mockCmdFactory.On("Create", "go", []string{"test", "-v", "-race", "-coverprofile=" + packageCoveragePth, "-covermode=atomic", packages[0]}, mock.Anything).Return(mockCmd)
	mockCmd.On("Run").Return(nil)
	mockCmd.On("PrintableCommandArgs").Return("")

	// It recreates package coverage file for every packages' go test command run
	mockPathProvider.On("CreateTempDir", mock.Anything).Return(tmpDir, nil)
	mockFileManager.On("Open", packageCoveragePth).Return(strings.NewReader(""), nil)
	mockFileManager.On("Create", packageCoveragePth).Return(nil, nil)
	mockFileManager.On("MkdirAll", outputDir, mock.Anything).Return(nil)
	mockFileManager.On("RemoveAll", packageCoveragePth).Return(nil)

	// It writes code coverage to file
	mockFileManager.On("Open", wantCoveragePth).Return(strings.NewReader(""), nil)
	mockFileManager.On("Write", wantCoveragePth, mock.Anything, mock.Anything).Return(nil)

	s := GoTestRunner{
		logger:         log.NewLogger(),
		inputParser:    nil,
		envRepo:        nil,
		cmdFactory:     mockCmdFactory,
		outputExporter: nil,
		pathProvider:   mockPathProvider,
		fileManager:    mockFileManager,
	}
	runResult, err := s.Run(opts)
	require.NoError(t, err)
	require.Equal(t, wantRunResult, runResult)
	mockCmdFactory.AssertExpectations(t)
	mockCmd.AssertExpectations(t)
	mockFileManager.AssertExpectations(t)
	mockPathProvider.AssertExpectations(t)
}

func TestGoTestRunner_Run_WhenTestFailsItReturnsAnError(t *testing.T) {
	// Create run options
	packages := []string{"./..."}
	outputDir := "output"
	opts := RunOpts{
		Packages:  packages,
		OutputDir: outputDir,
	}

	// Create mocks
	mockCmdFactory := new(mocks.Factory)
	mockCmd := new(mocks.Command)
	mockFileManager := new(mocks.FileManager)
	mockPathProvider := new(mocks.PathProvider)

	// It runs go test command with code coverage enabled
	tmpDir := "tmp_dir"
	packageCoveragePth := filepath.Join(tmpDir, "profile.out")
	mockCmdFactory.On("Create", "go", []string{"test", "-v", "-race", "-coverprofile=" + packageCoveragePth, "-covermode=atomic", packages[0]}, mock.Anything).Return(mockCmd)
	mockCmd.On("Run").Return(fmt.Errorf("exit status 1"))
	mockCmd.On("PrintableCommandArgs").Return("")

	// It recreates package coverage file for every packages' go test command run
	mockPathProvider.On("CreateTempDir", mock.Anything).Return(tmpDir, nil)
	mockFileManager.On("Create", packageCoveragePth).Return(nil, nil)
	mockFileManager.On("MkdirAll", outputDir, mock.Anything).Return(nil)

	s := GoTestRunner{
		logger:         log.NewLogger(),
		inputParser:    nil,
		envRepo:        nil,
		cmdFactory:     mockCmdFactory,
		outputExporter: nil,
		pathProvider:   mockPathProvider,
		fileManager:    mockFileManager,
	}
	gotResult, err := s.Run(opts)
	require.EqualError(t, err, "go test failed: exit status 1")
	require.Nil(t, gotResult)
	mockCmdFactory.AssertExpectations(t)
	mockCmd.AssertExpectations(t)
	mockFileManager.AssertExpectations(t)
	mockPathProvider.AssertExpectations(t)
}
