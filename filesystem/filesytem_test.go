package filesystem

import (
	"os"
	"path"
	"testing"
)

func TestCreatePackageCodeCoverageFile(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Integration verification test", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreatePackageCodeCoverageFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePackageCodeCoverageFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("CreatePackageCodeCoverageFile() = '%v'", got)
			}

			os.Remove(got)
		})
	}
}

func TestCodeCoveragePath(t *testing.T) {
	os.Setenv("BITRISE_DEPLOY_DIR", "test_dir")

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Integration CodeCoveragePath test", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CodeCoveragePath()
			if (err != nil) != tt.wantErr {
				t.Errorf("CodeCoveragePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			dir := path.Dir(got)
			defer func() { os.Remove(dir) }()

			if _, err := os.Stat(dir); err != nil {
				t.Errorf("CodeCoveragePath() = %v, err %v", got, err)
			}
		})
	}
}
