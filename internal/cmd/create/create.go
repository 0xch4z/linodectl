package create

import (
	"fmt"

	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/spf13/cobra"
)

const (
	createExample = `  # Create a linode instance
  linodectl create linode my-linode --region us-east --image ubuntu --type g6-standard-1

  # Create a linode based on a preset
  linodectl create linode my-linode --preset prod-machine

  # Create an LKE Cluster
  linodectl create lkecluster lke123 -v1.21 --region ap-west --pool g6-standard-3:3`
)

func NameFromCommandArgs(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("exactly one NAME is required, got %d", len(args))
	}
	return args[0], nil
}

func NewCmdCreate(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	cmds := &cobra.Command{
		Use:     "create [resource]",
		Short:   "Create a Linode resource",
		Example: createExample,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmds.AddCommand(NewCmdCreateInstance(f, ioStreams))
	cmds.AddCommand(NewCmdCreateLKECluister(f, ioStreams))
	return cmds
}
