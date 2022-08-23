package get

import (
	"context"

	"github.com/0xch4z/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/0xch4z/linodectl/internal/cmd/util"
	"github.com/0xch4z/linodectl/internal/printer"
	"github.com/0xch4z/linodectl/internal/resource/instance"
	"github.com/0xch4z/linodectl/internal/resource/resourceref"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type GetInstanceOptions struct {
	refs resourceref.List

	genericoptions.PaginationFlags
	genericoptions.ProfileFlags
	genericoptions.PrinterFlags
	instance.FilterFlags
	cmdutil.IOStreams
}

func NewGetInstanceOptions(ioStreams cmdutil.IOStreams) *GetInstanceOptions {
	return &GetInstanceOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdGetInstance(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewGetInstanceOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "instance [NAME] [args...]",
		Aliases: []string{"instances", "linode", "linodes"},
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

func (o *GetInstanceOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	if o.refs, err = resourceref.ListFromArgs(args); err != nil {
		return err
	}
	return nil
}

func (o *GetInstanceOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	filter := o.Filter(o.refs.Label())

	filterBytes, err := filter.MarshalJSON()
	if err != nil {
		return err
	}

	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}
	ctx := context.Background()

	instances, err := client.ListInstances(ctx, &linodego.ListOptions{
		PageOptions: o.PageOptions(),
		PageSize:    o.PageOptions().Results,
		Filter:      string(filterBytes),
	})
	if err != nil {
		return err
	}

	if o.LKECluster() != 0 {
		if instances, err = instance.FilterLKECluster(ctx, client, o.LKECluster(), instances); err != nil {
			return err
		}
	}

	if len(o.refs) > 0 {
		instances = instance.FilterByRefs(instances, o.refs)
	}

	resourceList := instance.NewList(instances)
	p := printer.New(o.Out)
	return p.PrintResources(context.Background(), resourceList, printer.ResourcePrintOptions{
		Columns:         o.Fields(),
		SortBy:          o.SortBy(),
		OmitHeader:      o.NoHeader(),
		DescendingOrder: o.Descending(),
	})
}
