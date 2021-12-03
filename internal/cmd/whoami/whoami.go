package whoami

import (
	"context"
	"fmt"

	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/spf13/cobra"
)

type WhoamiOptions struct {
	genericoptions.ProfileFlags
	cmdutil.IOStreams
}

func NewWhoamiOptions(ioStreams cmdutil.IOStreams) *WhoamiOptions {
	return &WhoamiOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdWhoami(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewWhoamiOptions(ioStreams)

	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Introspect",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(f, cmd)
		},
	}

	o.ProfileFlags.AddFlags(cmd)

	return cmd
}

func (o *WhoamiOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}

	profile, err := client.GetProfile(context.Background())
	if err != nil {
		return err
	}

	fmt.Fprintln(o.Out, profile.Username)
	return nil
}
