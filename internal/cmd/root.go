package cmd

import (
	"io"

	"github.com/0xch4z/linodectl/internal/cli/genericoptions"
	"github.com/0xch4z/linodectl/internal/cmd/config"
	"github.com/0xch4z/linodectl/internal/cmd/create"
	"github.com/0xch4z/linodectl/internal/cmd/delete"
	"github.com/0xch4z/linodectl/internal/cmd/edit"
	"github.com/0xch4z/linodectl/internal/cmd/get"
	cmdutil "github.com/0xch4z/linodectl/internal/cmd/util"
	"github.com/0xch4z/linodectl/internal/cmd/whoami"
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
	cmds.AddCommand(edit.NewCmdEdit(f, ioStreams))
	cmds.AddCommand(delete.NewCmdDelete(f, ioStreams))
	cmds.AddCommand(config.NewCmdConfig(f, ioStreams))
	cmds.AddCommand(whoami.NewCmdWhoami(f, ioStreams))
	profileFlags.AddFlags(cmds)
	return cmds
}
