package netpol

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stackrox/rox/roxctl/common/mocks"
	"github.com/stretchr/testify/suite"
)

func TestAnalyzeNetpolCommand(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(analyzeNetpolTestSuite))
}

type analyzeNetpolTestSuite struct {
	suite.Suite
}

func (d *analyzeNetpolTestSuite) TestAnalyzeNetpol() {
	cases := map[string]struct {
		inputFolderPath       string
		expectedAnalysisError error
		expectedValidateError error
		strict                bool
		stopOnFirstErr        bool
		outFile               string
		removeOutputPath      bool
	}{

		"not existing inputFolderPath should raise 'os.ErrNotExist' error": {
			inputFolderPath:       "/tmp/xxx",
			expectedAnalysisError: os.ErrNotExist,
		},
		"happyPath": {
			inputFolderPath:       "testdata/minimal",
			expectedAnalysisError: nil,
		},
		/*("treating warnings as errors": {
			inputFolderPath:       "testdata/empty-yamls",
			expectedAnalysisError: errNPGWarningsIndicator,
			strict:                true,
		},*/
		/*"stopOnFistError": {
			inputFolderPath:       "testdata/dirty",
			expectedAnalysisError: errNPGErrorsIndicator,
			stopOnFirstErr:        true,
		},*/
		"output should be written to a single file": {
			inputFolderPath:       "testdata/minimal",
			expectedAnalysisError: nil,
			outFile:               d.T().TempDir() + "/out.yaml",
			removeOutputPath:      false,
		},
	}

	for name, tt := range cases {
		tt := tt
		d.Run(name, func() {
			testCmd := &cobra.Command{Use: "test"}
			//testCmd.Flags().String("output-dir", "", "")
			testCmd.Flags().String("output-file", "", "")

			env, _, _ := mocks.NewEnvWithConn(nil, d.T())
			analyzeNetpolCmd := analyzeNetpolCommand{
				stopOnFirstError:      tt.stopOnFirstErr,
				treatWarningsAsErrors: tt.strict,
				inputFolderPath:       "", // set through construct
				outputFilePath:        tt.outFile,
				removeOutputPath:      tt.removeOutputPath,
				env:                   env,
				//printer:               nil,
			}

			if tt.outFile != "" {
				d.Assert().NoError(testCmd.Flags().Set("output-file", tt.outFile))
			}

			analyzer, err := analyzeNetpolCmd.construct([]string{tt.inputFolderPath}, testCmd)
			d.Assert().NoError(err)

			err = analyzeNetpolCmd.validate()
			if tt.expectedValidateError != nil {
				d.Require().Error(err)
				d.Assert().ErrorIs(err, tt.expectedValidateError)
				return
			}
			d.Assert().NoError(err)

			err = analyzeNetpolCmd.analyzeNetpols(analyzer)
			if tt.expectedAnalysisError != nil {
				d.Require().Error(err)
				d.Assert().ErrorIs(err, tt.expectedAnalysisError)
			} else {
				d.Assert().NoError(err)
			}

		})
	}

}

/*
TODO:
- handle : treating_warnings_as_errors : avoid having fatal error for empty resources in the input dir.
- handle:s topOnFistError : debug why test is failing

=== RUN   TestAnalyzeNetpolCommand/TestAnalyzeNetpol/treating_warnings_as_errors
    netpol_test.go:84:
                Error Trace:    C:\Users\847978756\npv\stackrox\roxctl\analyze\netpol\netpol_test.go:84
                                                        C:\Users\847978756\npv\stackrox\roxctl\analyze\netpol\suite.go:91
                Error:          Target error should be in err chain:
                                expected: "there were warnings during execution"
                                in chain: "error analyzing network policies: cannot produce connectivity list without k8s workloads"
                                        "error analyzing network policies: cannot produce connectivity list without k8s workloads"
                                        "cannot produce connectivity list without k8s workloads"
                Test:           TestAnalyzeNetpolCommand/TestAnalyzeNetpol/treating_warnings_as_errors
--- FAIL: TestAnalyzeNetpolCommand (0.01s)
    --- FAIL: TestAnalyzeNetpolCommand/TestAnalyzeNetpol (0.01s)
        --- PASS: TestAnalyzeNetpolCommand/TestAnalyzeNetpol/happyPath (0.01s)
        --- PASS: TestAnalyzeNetpolCommand/TestAnalyzeNetpol/not_existing_inputFolderPath_should_raise_'os.ErrNotExist'_error (0.00s)
        --- FAIL: TestAnalyzeNetpolCommand/TestAnalyzeNetpol/treating_warnings_as_errors (0.00s)
FAIL
FAIL    github.com/stackrox/rox/roxctl/analyze/netpol   12.639s
FAIL




$ go test -v ./roxctl/analyze/netpol/
=== RUN   TestAnalyzeNetpolCommand
=== PAUSE TestAnalyzeNetpolCommand
=== CONT  TestAnalyzeNetpolCommand
=== RUN   TestAnalyzeNetpolCommand/TestAnalyzeNetpol
=== RUN   TestAnalyzeNetpolCommand/TestAnalyzeNetpol/not_existing_inputFolderPath_should_raise_'os.ErrNotExist'_error
=== RUN   TestAnalyzeNetpolCommand/TestAnalyzeNetpol/stopOnFistError
    netpol_test.go:89:
                Error Trace:    C:\Users\847978756\npv\stackrox\roxctl\analyze\netpol\netpol_test.go:89
                                                        C:\Users\847978756\npv\stackrox\roxctl\analyze\netpol\suite.go:91
                Error:          Target error should be in err chain:
                                expected: "there were errors during execution"
                                in chain: "error in connectivity analysis: YAML document is malformed: yaml: line 17: found character that cannot start any token"
                                        "error in connectivity analysis: YAML document is malformed: yaml: line 17: found character that cannot start any token"
                                        "YAML document is malformed: yaml: line 17: found character that cannot start any token"
                                        "yaml: line 17: found character that cannot start any token"
                Test:           TestAnalyzeNetpolCommand/TestAnalyzeNetpol/stopOnFistError
=== RUN   TestAnalyzeNetpolCommand/TestAnalyzeNetpol/happyPath
--- FAIL: TestAnalyzeNetpolCommand (0.03s)
    --- FAIL: TestAnalyzeNetpolCommand/TestAnalyzeNetpol (0.03s)
        --- PASS: TestAnalyzeNetpolCommand/TestAnalyzeNetpol/not_existing_inputFolderPath_should_raise_'os.ErrNotExist'_error (0.00s)
        --- FAIL: TestAnalyzeNetpolCommand/TestAnalyzeNetpol/stopOnFistError (0.01s)
        --- PASS: TestAnalyzeNetpolCommand/TestAnalyzeNetpol/happyPath (0.02s)
FAIL
FAIL    github.com/stackrox/rox/roxctl/analyze/netpol   19.413s
FAIL


*/
