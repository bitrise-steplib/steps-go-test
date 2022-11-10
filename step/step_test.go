package step

import (
	"reflect"
	"testing"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

func TestStep_Run(t *testing.T) {
	type fields struct {
		collector   CodeCoverageCollector
		env         env.Repository
		inputParser stepconf.InputParser
		logger      log.Logger
		testRunner  TestRunner
	}
	type args struct {
		config *Config
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       *Result
		wantErr    bool
		wasTestRun bool
	}{
		{
			name: "With one package expect Result to contain collected file path and test to run",
			fields: fields{
				collector:   testCollector{"out", "collected_file_path"},
				env:         testEnv{},
				inputParser: testParser{},
				logger:      log.NewLogger(),
				testRunner:  &fakeTestRunner{},
			},
			args: args{&Config{
				OutputDir: "test_outputdir",
				Packages:  []string{"package1"},
			}},
			want:       &Result{"collected_file_path"},
			wantErr:    false,
			wasTestRun: true,
		},
		{
			name: "When packages are empty expect success with no tests run",
			fields: fields{
				collector:   testCollector{"out", "collected_file_path"},
				env:         testEnv{},
				inputParser: testParser{},
				logger:      log.NewLogger(),
				testRunner:  &fakeTestRunner{},
			},
			args: args{&Config{
				OutputDir: "test_outputdir",
				Packages:  []string{},
			}},
			want:       &Result{"collected_file_path"},
			wantErr:    false,
			wasTestRun: false,
		},
		{
			name: "When packages are nil expect error with no tests run",
			fields: fields{
				collector:   testCollector{"out", "collected_file_path"},
				env:         testEnv{},
				inputParser: testParser{},
				logger:      log.NewLogger(),
				testRunner:  &fakeTestRunner{},
			},
			args: args{&Config{
				OutputDir: "test_outputdir",
				Packages:  nil,
			}},
			want:       &Result{"collected_file_path"},
			wantErr:    false,
			wasTestRun: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Step{
				collector:   tt.fields.collector,
				env:         tt.fields.env,
				inputParser: tt.fields.inputParser,
				logger:      tt.fields.logger,
				testRunner:  tt.fields.testRunner,
			}
			got, err := s.Run(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Step.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Step.Run() = %v, want %v", got, tt.want)
			}

			ftr, ok := tt.fields.testRunner.(*fakeTestRunner)
			if !ok {
				panic("Test failed due to inability to cast interface to expected instance")
			}

			if ftr.wasRun != tt.wasTestRun {
				t.Errorf("Test run status did not match expectation (%v == %v)", ftr.wasRun, tt.wasTestRun)
			}
		})
	}
}

func TestStep_ProcessConfig(t *testing.T) {
	collector := testCollector{}
	logger := log.NewLogger()
	testRunner := fakeTestRunner{}

	tests := []struct {
		name    string
		env     env.Repository
		want    *Config
		wantErr bool
	}{
		{
			"Empty env vars",
			testEnv{},
			nil,
			true,
		},
		{
			"Only deploy dir env var set",
			testEnv{map[string]string{"BITRISE_DEPLOY_DIR": "dir"}},
			nil,
			true,
		},
		{
			"Only packages env var set",
			testEnv{map[string]string{"packages": "pkg"}},
			nil,
			true,
		},
		{
			"Valid config available",
			testEnv{map[string]string{"BITRISE_DEPLOY_DIR": "dir", "packages": "pkg"}},
			&Config{
				OutputDir: "dir",
				Packages:  []string{"pkg"},
			},
			false,
		},
		{
			"Valid config available with multiple packages",
			testEnv{map[string]string{"BITRISE_DEPLOY_DIR": "dir", "packages": "pkg\npkg2"}},
			&Config{
				OutputDir: "dir",
				Packages:  []string{"pkg", "pkg2"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Step{
				collector:   collector,
				env:         tt.env,
				inputParser: stepconf.NewInputParser(tt.env),
				logger:      logger,
				testRunner:  &testRunner,
			}
			got, err := s.ProcessConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("Step.ProcessConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Step.ProcessConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test Collector

type testCollector struct {
	outputPathComponent  string
	collectedResultsPath string
}

func (c testCollector) PrepareAndReturnCurrentPackageCoverageOutputPath(prepend string) (string, error) {
	return prepend + c.outputPathComponent, nil
}

func (c testCollector) CollectCoverageResultsAndReset() error {
	return nil
}

func (c testCollector) FinishCollectionAndReturnPathToCollectedResults() string {
	return c.collectedResultsPath
}

// Test Environment

type testEnv struct {
	val map[string]string
}

func (e testEnv) List() []string {
	panic("not implemented") // TODO: Implement
}

func (e testEnv) Unset(key string) error {
	panic("not implemented") // TODO: Implement
}

func (e testEnv) Get(key string) string {
	v := e.val[key]
	return v
}

func (e testEnv) Set(key string, value string) error {
	panic("not implemented") // TODO: Implement
}

// Test Parser

type testParser struct{}

func (p testParser) Parse(input interface{}) error {
	panic("not implemented") // TODO: Implement
}

// Test TestRunner

type fakeTestRunner struct {
	wasRun bool
}

func (r *fakeTestRunner) RunTest(config TestConfig) error {
	r.wasRun = true
	return nil
}
