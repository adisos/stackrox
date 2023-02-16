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
		"treating warnings as errors": {
			inputFolderPath:       "testdata/empty-yamls",
			expectedAnalysisError: errNPGWarningsIndicator,
			strict:                true,
		},
		"stopOnFistError": {
			inputFolderPath:       "testdata/dirty",
			expectedAnalysisError: errNPGErrorsIndicator,
			stopOnFirstErr:        true,
		},
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
