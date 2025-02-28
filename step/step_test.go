package step

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-go-test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGoTestRunner_Run_WhenTestSucceeds(t *testing.T) {
	// Create run options
	pkg := "./..."
	outputDir := t.TempDir()
	opts := RunOpts{
		Package:   pkg,
		Covermode: "atomic",
		OutputDir: outputDir,
	}

	// Expected result
	testRunLogFile := createTmpFile(outputDir, "go_test_run.log", t)
	codeCoverageFile := createTmpFile(outputDir, "go_code_coverage.out", t)
	wantRunResult := &RunResult{
		CodeCoveragePth: codeCoverageFile.Name(),
		TestRunLogPath:  testRunLogFile.Name(),
	}

	// Create mocks
	mockCmdFactory := mocks.NewFactory(t)
	mockCmd := mocks.NewCommand(t)
	mockFileManager := mocks.NewFileManager(t)
	mockPathProvider := mocks.NewPathProvider(t)

	// It runs go test command with code coverage enabled
	mockCmdFactory.On("Create", "go", []string{"test", "-v", "-covermode=atomic", "-coverprofile=" + codeCoverageFile.Name(), pkg}, mock.Anything).Return(mockCmd)
	mockCmd.On("Run").Return(nil)
	mockCmd.On("PrintableCommandArgs").Return("")

	// It creates package coverage file and test run log file
	mockFileManager.On("MkdirAll", outputDir, mock.Anything).Return(nil)
	mockFileManager.On("Create", codeCoverageFile.Name()).Return(codeCoverageFile, nil)
	mockFileManager.On("Create", testRunLogFile.Name()).Return(testRunLogFile, nil)

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
	require.NoError(t, err)
	require.Equal(t, wantRunResult, gotResult)
}

func TestGoTestRunner_Run_WhenTestFails(t *testing.T) {
	// Create run options
	pkg := "./..."
	outputDir := t.TempDir()
	opts := RunOpts{
		Package:   pkg,
		OutputDir: outputDir,
	}

	// Expected result
	testRunLogFile := createTmpFile(outputDir, "go_test_run.log", t)
	wantRunResult := &RunResult{
		TestRunLogPath: testRunLogFile.Name(),
	}

	// Create mocks
	mockCmdFactory := mocks.NewFactory(t)
	mockCmd := mocks.NewCommand(t)
	mockFileManager := mocks.NewFileManager(t)
	mockPathProvider := mocks.NewPathProvider(t)

	// It runs go test command with code coverage enabled
	mockCmdFactory.On("Create", "go", []string{"test", "-v", pkg}, mock.Anything).Return(mockCmd)
	mockCmd.On("Run").Return(fmt.Errorf("exit status 1"))
	mockCmd.On("PrintableCommandArgs").Return("")

	// It creates package coverage file and test run log file
	mockFileManager.On("MkdirAll", outputDir, mock.Anything).Return(nil)
	mockFileManager.On("Create", testRunLogFile.Name()).Return(testRunLogFile, nil)

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
	require.EqualError(t, err, "exit status 1")
	require.Equal(t, wantRunResult, gotResult)
}

func createTmpFile(tmpDir string, name string, t *testing.T) *os.File {
	tmpFilePth := filepath.Join(tmpDir, name)
	f, err := os.Create(tmpFilePth)
	require.NoError(t, err)
	return f
}
