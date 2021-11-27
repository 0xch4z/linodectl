package get

import (
	"context"

	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/printer"
	"github.com/Charliekenney23/linodectl/internal/resource/stackscript"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type GetStackScriptOptions struct {
	// Label (optional) is the name of an instance to fetch
	Label string

	genericoptions.PaginationFlags
	genericoptions.ProfileFlags
	genericoptions.PrinterFlags
	stackscript.FilterFlags
	cmdutil.IOStreams
}

func NewGetStackScriptOptions(ioStreams cmdutil.IOStreams) *GetStackScriptOptions {
	return &GetStackScriptOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdGetStackScript(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewGetStackScriptOptions(ioStreams)

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

func (o *GetStackScriptOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	if len(args) == 1 {
		o.Label = args[0]
	}

	return nil
}

func (o *GetStackScriptOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	filter := o.Filter(o.Label)

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

	resourceList := stackscript.NewList(stackScripts)
	p := printer.New(o.Out)
	return p.PrintResources(context.Background(), resourceList, printer.ResourcePrintOptions{
		Columns:         o.Fields(),
		SortBy:          o.SortBy(),
		OmitHeader:      o.NoHeader(),
		DescendingOrder: o.Descending(),
	})
}
