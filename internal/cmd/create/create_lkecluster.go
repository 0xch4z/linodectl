package create

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/printer"
	"github.com/Charliekenney23/linodectl/internal/resource/lkecluster"
	"github.com/Charliekenney23/linodectl/internal/strutil"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

const (
	lkeClusterExamples = `  # Create a 1.21 Cluster with a 5 node pool of type g6-standard-2
  linodectl create lkecluster -v1.21 --region us-east --pool g6-standard-2:5

  # Create a Cluster on the latest support version of Kubernetes
  linodectl create lkecluster --region us-central --pool g6-standard-1:3

  # Create a Cluster with 2 node pools
  linodectl create lkecluster --region eu-west -v1.20 --pool g6-standard-1:3 --pool g6-standard-3:2

  # Create a Cluster with a highly available control plane
  linodectl create lkecluster --region us-east -v1.20 --ha --pool g6-standard-1:3`

	lkeClusterNodePoolUsage = "Node pool configuration in format <instance-type>:<count> (i.e. g6-standard-2:3)"
)

type CreateLKEClusterOptions struct {
	Label string

	Version string
	Region  string
	Tags    []string
	Pools   []string
	HA      bool

	NodePools []linodego.LKEClusterPoolCreateOptions

	genericoptions.ProfileFlags
	genericoptions.PrinterFlags
	cmdutil.IOStreams
}

func NewCreateLKEClusterOptions(ioStreams cmdutil.IOStreams) *CreateLKEClusterOptions {
	return &CreateLKEClusterOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdCreateLKECluister(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewCreateLKEClusterOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "lkecluster NAME [args...]",
		Short:   "Create an LKE Cluster",
		Aliases: []string{"cluster"},
		Example: lkeClusterExamples,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(f, ioStreams, args); err != nil {
				return err
			}
			return o.Run(f, cmd)
		},
	}

	cmd.Flags().StringSliceVar(&o.Pools, "pool", o.Pools, lkeClusterNodePoolUsage)
	cmd.Flags().StringVarP(&o.Version, "version", "v", "", "Version of Kubernetes to deploy")
	cmd.Flags().StringVar(&o.Region, "region", "", "Region to deploy in")
	cmd.Flags().StringSliceVar(&o.Tags, "tags", o.Tags, "Tags to add to this cluster")
	cmd.Flags().BoolVar(&o.HA, "ha", false, "If true, this cluster will be deployed with a highly available control plane")

	o.ProfileFlags.AddFlags(cmd)
	return cmd
}

func (o *CreateLKEClusterOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	o.Label, err = NameFromCommandArgs(args)
	if err != nil {
		return err
	}

	o.NodePools = make([]linodego.LKEClusterPoolCreateOptions, len(o.Pools))
	for i, pool := range o.Pools {
		args := strings.Split(pool, ":")
		if len(args) != 2 {
			return fmt.Errorf("node pool spec %q is malformed", pool)
		}

		o.NodePools[i].Type = args[0]
		if o.NodePools[i].Count, err = strconv.Atoi(args[1]); err != nil {
			return fmt.Errorf("%q is not a valid node pool count", args[1])
		}
	}

	if profile, ok := f.Config().CurrentProfile(); ok {
		o.Region = strutil.Fallback(o.Region, profile.Region)
	}
	return nil
}

func (o *CreateLKEClusterOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	options := linodego.LKEClusterCreateOptions{
		Label:      o.Label,
		K8sVersion: o.Version,
		Region:     o.Region,
		NodePools:  o.NodePools,
		Tags:       o.Tags,
		ControlPlane: &linodego.LKEClusterControlPlane{
			HighAvailability: o.HA,
		},
	}

	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}
	ctx := context.Background()

	cluster, err := client.CreateLKECluster(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to create lkecluster: %w", err)
	}

	resList := lkecluster.NewList([]linodego.LKECluster{*cluster})
	p := printer.New(o.Out)
	return p.PrintResources(context.Background(), resList, printer.ResourcePrintOptions{
		Columns:         o.Fields(),
		SortBy:          o.SortBy(),
		OmitHeader:      o.NoHeader(),
		DescendingOrder: o.Descending(),
	})
}
