package filesystem

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/log"
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

			panicIfErr(os.Remove(got))
		})
	}
}

func TestCodeCoveragePath(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Integration CodeCoveragePath test", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CodeCoveragePath("test_dir")
			if (err != nil) != tt.wantErr {
				t.Errorf("CodeCoveragePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			dir := path.Dir(got)
			defer func() { panicIfErr(os.Remove(dir)) }()

			if _, err := os.Stat(dir); err != nil {
				t.Errorf("CodeCoveragePath() = %v, err %v", got, err)
			}
		})
	}
}

func TestAppendPackageCoverageAndRecreate(t *testing.T) {
	logger := log.NewLogger()

	type args struct {
		packageCoveragePth string
		coveragePth        string
	}
	type contents struct {
		pkg string
		cov string
	}
	tests := []struct {
		name     string
		args     args
		contents contents
		wantErr  bool
	}{
		{
			"Integration AppendPackageCoverageAndRecreate",
			args{"pkg.txt", "cov.txt"},
			contents{"foo", "bar"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panicIfErr(os.WriteFile(tt.args.packageCoveragePth, []byte(tt.contents.pkg), 0644))
			panicIfErr(os.WriteFile(tt.args.coveragePth, []byte(tt.contents.cov), 0644))

			defer func() {
				panicIfErr(os.Remove(tt.args.packageCoveragePth))
				if _, err := os.Stat(tt.args.coveragePth); err == nil {
					panicIfErr(os.Remove(tt.args.coveragePth))
				}
			}()

			if err := AppendPackageCoverageAndRecreate(tt.args.packageCoveragePth, tt.args.coveragePth, logger); (err != nil) != tt.wantErr {
				t.Errorf("AppendPackageCoverageAndRecreate() error = %v, wantErr %v", err, tt.wantErr)
			}

			bytes, err := os.ReadFile(tt.args.coveragePth)
			if err != nil {
				t.Errorf("Error reading file %v; error = %v", tt.args.packageCoveragePth, err)
			}

			contents := string(bytes)
			if !strings.HasPrefix(contents, tt.contents.cov) || !strings.HasSuffix(contents, tt.contents.pkg) {
				t.Errorf("Error with order of contents: '%v'", contents)
			}
		})
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
