package netpol

import (
	"path/filepath"

	npguard "github.com/np-guard/netpol-analyzer/pkg/netpol/connlist"
	"github.com/spf13/cobra"
	"github.com/stackrox/rox/roxctl/common/environment"
	"github.com/stackrox/rox/roxctl/common/printer"
)

type analyzeNetpolCommand struct {
	// Properties that are bound to cobra flags.
	inputFolderPath string

	// injected or constructed values
	env     environment.Environment
	printer printer.ObjectPrinter
}

// Command defines the netpol command tree
func Command(cliEnvironment environment.Environment) *cobra.Command {
	analyzeNetpolCmd := &analyzeNetpolCommand{env: cliEnvironment}
	c := &cobra.Command{
		Use:  "netpol <folder-path>",
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			analyzeNetpolCmd.env.Logger().WarnfLn("This is a Technology Preview feature. Red Hat does not recommend using Technology Preview features in production.")
			res, err := analyzeNetpolCmd.construct(args, c)
			analyzeNetpolCmd.printNetpolAnalysis(res)
			/*synth, err := generateNetpolCmd.construct(args, c)
			if err != nil {
				return err
			}
			if err := generateNetpolCmd.validate(); err != nil {
				return err
			}
			return generateNetpolCmd.generateNetpol(synth)*/
			return err
		},
	}
	/*c.Flags().BoolVar(&generateNetpolCmd.treatWarningsAsErrors, "strict", false, "treat warnings as errors")
	c.Flags().BoolVar(&generateNetpolCmd.stopOnFirstError, "fail", false, "fail on the first encountered error")
	c.Flags().BoolVar(&generateNetpolCmd.removeOutputPath, "remove", false, "remove the output path if it already exists")
	c.Flags().StringVarP(&generateNetpolCmd.outputFolderPath, "output-dir", "d", "", "save generated policies into target folder - one file per policy")
	c.Flags().StringVarP(&generateNetpolCmd.outputFilePath, "output-file", "f", "", "save and merge generated policies into a single yaml file")*/
	return c
}
func (cmd *analyzeNetpolCommand) construct(args []string, c *cobra.Command) (string, error) {
	cmd.inputFolderPath = args[0]
	//return npguard.NewPoliciesSynthesizer(opts...), nil
	conns, err := npguard.FromDir(cmd.inputFolderPath, filepath.WalkDir)
	if err != nil {
		return "", err
	}
	res := npguard.ConnectionsListToString(conns)
	return res, err

}
