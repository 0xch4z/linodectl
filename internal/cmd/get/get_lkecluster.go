package get

import (
	"context"

	"github.com/0xch4z/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/0xch4z/linodectl/internal/cmd/util"
	"github.com/0xch4z/linodectl/internal/printer"
	"github.com/0xch4z/linodectl/internal/resource/lkecluster"
	"github.com/0xch4z/linodectl/internal/resource/resourceref"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type GetLKEClusterOptions struct {
	refs resourceref.List

	genericoptions.PaginationFlags
	genericoptions.ProfileFlags
	genericoptions.PrinterFlags
	lkecluster.FilterFlags
	cmdutil.IOStreams
}

func NewGetLKEClusterOptions(ioStreams cmdutil.IOStreams) *GetLKEClusterOptions {
	return &GetLKEClusterOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdGetLKECluster(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewGetLKEClusterOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "lkecluster [NAME] [args...]",
		Aliases: []string{"cluster", "lkeclusters", "clusters"},
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

func (o *GetLKEClusterOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	if o.refs, err = resourceref.ListFromArgs(args); err != nil {
		return err
	}
	return nil
}

func (o *GetLKEClusterOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	// Currently, the Linode API does not support filtering by label on
	// LKE Clusters, even though the docs say otherwise. When this is fixed,
	// we should pass the string o.refs.Label(), so that if a single label
	// is specified, only that cluster is fetched.
	// If it's included in the request right now, we'll get a 400 error.
	filter := o.Filter("")

	filterBytes, err := filter.MarshalJSON()
	if err != nil {
		return err
	}

	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}

	ctx := context.Background()
	clusters, err := client.ListLKEClusters(ctx, &linodego.ListOptions{
		PageOptions: o.PageOptions(),
		PageSize:    o.PageOptions().Results,
		Filter:      string(filterBytes),
	})
	if err != nil {
		return err
	}

	if len(o.refs) > 0 {
		clusters = lkecluster.FilterByRefs(clusters, o.refs)
	}

	resourceList := lkecluster.NewList(clusters)
	p := printer.New(o.Out)
	return p.PrintResources(context.Background(), resourceList, printer.ResourcePrintOptions{
		Columns:         o.Fields(),
		SortBy:          o.SortBy(),
		OmitHeader:      o.NoHeader(),
		DescendingOrder: o.Descending(),
	})
}
