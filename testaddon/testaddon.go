package testaddon

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-steplib/steps-go-test/filemanager"
)

const (
	resultDescriptorFileName = "test-info.json"
)

// TestInfo ...
type TestInfo struct {
	Name string `json:"test-name"`
}

type Exporter interface {
	PrepareTestResultExport(testName string) (*os.File, error)
}

type exporter struct {
	envRepo     env.Repository
	fileManager filemanager.FileManager
}

func NewExporter(envRepo env.Repository, fileManager filemanager.FileManager) Exporter {
	return exporter{
		envRepo:     envRepo,
		fileManager: fileManager,
	}
}

func (e exporter) PrepareTestResultExport(testName string) (*os.File, error) {
	testInfo := &TestInfo{
		Name: testName,
	}

	stepTestResultDir := e.envRepo.Get("BITRISE_TEST_RESULT_DIR")
	testResultFileName := replaceUnsupportedFilenameCharacters(testName)
	testResultDir := filepath.Join(stepTestResultDir, testResultFileName)

	if err := e.fileManager.MkdirAll(testResultDir, os.ModePerm); err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(testInfo)
	if err != nil {
		return nil, err
	}

	testInfoPath := filepath.Join(testResultDir, resultDescriptorFileName)
	if err := e.fileManager.Write(testInfoPath, string(bytes), os.ModePerm); err != nil {
		return nil, err
	}

	testReportFilePth := filepath.Join(testResultDir, "test_report.xml")
	testReportFile, err := e.fileManager.Create(testReportFilePth)
	if err != nil {
		return nil, err
	}

	return testReportFile, nil
}

func replaceUnsupportedFilenameCharacters(s string) string {
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ":", "-")
	return s
}
