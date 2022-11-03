package filesystem

import (
	"os"
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
