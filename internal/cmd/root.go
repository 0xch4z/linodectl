package cmd

import (
	"io"

	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	"github.com/Charliekenney23/linodectl/internal/cmd/config"
	"github.com/Charliekenney23/linodectl/internal/cmd/create"
	"github.com/Charliekenney23/linodectl/internal/cmd/delete"
	"github.com/Charliekenney23/linodectl/internal/cmd/get"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/cmd/whoami"
	"github.com/spf13/cobra"
)

func NewRootCommand(f cmdutil.Factory, in io.Reader, out, err io.Writer) *cobra.Command {
	var profileFlags genericoptions.ProfileFlags

	cmds := &cobra.Command{
		Use:   "linodectl",
		Short: "linodectl manages Linode resources",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	ioStreams := cmdutil.IOStreams{In: in, Out: out, Err: err}
	cmds.AddCommand(get.NewCmdGet(f, ioStreams))
	cmds.AddCommand(create.NewCmdCreate(f, ioStreams))
	cmds.AddCommand(delete.NewCmdDelete(f, ioStreams))
	cmds.AddCommand(config.NewCmdConfig(f, ioStreams))
	cmds.AddCommand(whoami.NewCmdWhoami(f, ioStreams))
	profileFlags.AddFlags(cmds)
	return cmds
}
