package delete

import (
	"context"
	"fmt"

	"github.com/0xch4z/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/0xch4z/linodectl/internal/cmd/util"
	"github.com/0xch4z/linodectl/internal/resource/resourceref"
	"github.com/0xch4z/linodectl/internal/resource/stackscript"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type DeleteStackScriptOptions struct {
	refs resourceref.List

	genericoptions.PaginationFlags
	genericoptions.ProfileFlags
	genericoptions.PrinterFlags
	stackscript.FilterFlags
	cmdutil.IOStreams
}

func NewDeleteStackScriptOptions(ioStreams cmdutil.IOStreams) *DeleteStackScriptOptions {
	return &DeleteStackScriptOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdDeleteStackScript(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewDeleteStackScriptOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "stackscript [NAME] [args...]",
		Aliases: []string{"stackscripts", "ss"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(f, ioStreams, args); err != nil {
				return err
			}
			return o.Run(f, cmd)
		},
	}

	o.PaginationFlags.AddFlags(cmd)
	o.ProfileFlags.AddFlags(cmd)
	o.PrinterFlags.AddFlags(cmd)
	o.FilterFlags.AddFlags(cmd)
	return cmd
}

func (o *DeleteStackScriptOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	if o.refs, err = resourceref.ListFromArgs(args); err != nil {
		return err
	}
	return nil
}

func (o *DeleteStackScriptOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	filter := o.Filter(o.refs.Label())

	if len(o.refs) == 0 {
		return fmt.Errorf("need filter")
	}

	filterBytes, err := filter.MarshalJSON()
	if err != nil {
		return err
	}

	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}
	ctx := context.Background()

	stackScripts, err := client.ListStackscripts(ctx, &linodego.ListOptions{
		PageOptions: o.PageOptions(),
		PageSize:    o.PageOptions().Results,
		Filter:      string(filterBytes),
	})
	if err != nil {
		return err
	}

	toDelete := stackscript.FilterByRefs(stackScripts, o.refs)
	for _, ss := range toDelete {
		if err := client.DeleteStackscript(context.Background(), ss.ID); err != nil {
			return err
		}
		fmt.Fprintf(o.Out, "StackScript %q (%d) deleted...\n", ss.Label, ss.ID)
	}
	return nil
}
