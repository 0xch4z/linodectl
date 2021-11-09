package get

import (
	"fmt"

	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/spf13/cobra"
)

func NameFromCommandArgs(args []string) (string, error) {
	fmt.Println(args)
	if len(args) != 1 {
		return "", fmt.Errorf("exactly one NAME is required, got %d", len(args))
	}
	return args[0], nil
}

func NewCmdGet(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "get [resource]",
		Short: "Get a resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmds.AddCommand(NewCmdGetInstance(f, ioStreams))
	return cmds
}
