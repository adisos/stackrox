package analyze

import (
	"github.com/spf13/cobra"
	"github.com/stackrox/rox/roxctl/analyze/netpol"
	"github.com/stackrox/rox/roxctl/common/environment"
)

// Command defines the generate command tree
func Command(cliEnvironment environment.Environment) *cobra.Command {
	c := &cobra.Command{
		Use: "analyze",
	}

	c.AddCommand(
		netpol.Command(cliEnvironment),
	)
	return c
}
