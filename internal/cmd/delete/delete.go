package delete

import (
	cmdutil "github.com/0xch4z/linodectl/internal/cmd/util"
	"github.com/spf13/cobra"
)

func NewCmdDelete(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "delete [resource]",
		Short: "Delete a resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmds.AddCommand(NewCmdDeleteInstance(f, ioStreams))
	cmds.AddCommand(NewCmdDeleteStackScript(f, ioStreams))
	cmds.AddCommand(NewCmdDeleteBucket(f, ioStreams))
	cmds.AddCommand(NewCmdDeleteLKECluster(f, ioStreams))
	return cmds
}
