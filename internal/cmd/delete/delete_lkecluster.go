package delete

import (
	"context"
	"errors"
	"fmt"

	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/resource/lkecluster"
	"github.com/Charliekenney23/linodectl/internal/resource/resourceref"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type DeleteLKEClusterOptions struct {
	refs resourceref.List

	genericoptions.PaginationFlags
	genericoptions.ProfileFlags
	genericoptions.PrinterFlags
	lkecluster.FilterFlags
	cmdutil.IOStreams
}

func NewDeleteLKEClusterOptions(ioStreams cmdutil.IOStreams) *DeleteLKEClusterOptions {
	return &DeleteLKEClusterOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdDeleteLKECluster(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewDeleteLKEClusterOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "lkecluster [NAME] [args...]",
		Aliases: []string{"lkeclusters", "cluster", "clusters"},
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

func (o *DeleteLKEClusterOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	if o.refs, err = resourceref.ListFromArgs(args); err != nil {
		return err
	}
	return nil
}

func (o *DeleteLKEClusterOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	// Currently, the Linode API does not support filtering by label on
	// LKE Clusters, even though the docs say otherwise. When this is fixed,
	// we should pass the string o.refs.Label(), so that if a single label
	// is specified, only that cluster is fetched.
	// If it's included in the request right now, we'll get a 400 error.
	filter := o.Filter("")

	if len(o.refs) == 0 && len(filter.Children) == 0 {
		// we can't just delete every cluster
		return errors.New("at least one cluster or filter is required")
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
	clusters, err := client.ListLKEClusters(ctx, &linodego.ListOptions{
		PageOptions: o.PageOptions(),
		PageSize:    o.PageOptions().Results,
		Filter:      string(filterBytes),
	})
	if err != nil {
		return fmt.Errorf("failed to get LKE Clusters: %w", err)
	}

	if len(o.refs) > 0 {
		clusters = lkecluster.FilterByRefs(clusters, o.refs)
	}

	for _, cluster := range clusters {
		if err := client.DeleteLKECluster(context.Background(), cluster.ID); err != nil {
			return err
		}
		fmt.Fprintf(o.Out, "LKE Cluster %q (%d) deleted...\n", cluster.Label, cluster.ID)
	}
	return nil
}
