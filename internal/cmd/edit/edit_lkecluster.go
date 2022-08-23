package edit

import (
	"bytes"
	"context"
	"fmt"

	"github.com/0xch4z/linodectl/internal/cli/editor"
	"github.com/0xch4z/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/0xch4z/linodectl/internal/cmd/util"
	"github.com/0xch4z/linodectl/internal/resource/lkecluster"
	"github.com/0xch4z/linodectl/internal/resource/resourceref"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type EditLKEClusterOptions struct {
	refs resourceref.List

	genericoptions.ProfileFlags
	cmdutil.IOStreams
}

func NewEditLKEClusterOptions(ioStreams cmdutil.IOStreams) *EditLKEClusterOptions {
	return &EditLKEClusterOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdEditLKECluster(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewEditLKEClusterOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "lkecluster NAME [args...]",
		Aliases: []string{"cluster"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(f, ioStreams, args); err != nil {
				return err
			}
			return o.Run(f, cmd)
		},
	}

	o.ProfileFlags.AddFlags(cmd)
	return cmd
}

func (o *EditLKEClusterOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	if o.refs, err = resourceref.ListFromArgs(args); err != nil {
		return err
	}
	return nil
}

func (o *EditLKEClusterOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	if len(o.refs) == 0 {
		return fmt.Errorf("exactly one LKE Cluster ID or label must be specified")
	}

	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}

	ctx := context.Background()
	var filterBytes []byte

	clusters, err := client.ListLKEClusters(ctx, &linodego.ListOptions{
		Filter: string(filterBytes),
	})
	if err != nil {
		return err
	}

	clusters = lkecluster.FilterByRefs(clusters, o.refs)

	if len(clusters) != 1 {
		return fmt.Errorf("could not find instance %v", o.refs[0])
	}

	toEdit := clusters[0]
	pools, err := client.ListLKEClusterPools(ctx, toEdit.ID, &linodego.ListOptions{})
	if err != nil {
		return err
	}

	spec := lkecluster.SpecFromObject(&toEdit, pools)
	instanceBytes, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}

	editor := editor.NewDefaultEditor()
	specBytes, _, err := editor.EditReader("lkecluster", toEdit.Label+".yaml", o.IOStreams, bytes.NewBuffer(instanceBytes))
	if err != nil {
		panic(err)
	}

	var updatedSpec lkecluster.Spec
	if err := yaml.Unmarshal(specBytes, &updatedSpec); err != nil {
		return err
	}

	clusterPoolUpdates := make(map[int]*linodego.LKEClusterPoolUpdateOptions)
	for i, poolSpec := range spec.NodePools {
		if updatedSpec.NodePools == nil || len(updatedSpec.NodePools) < i+1 {
			return fmt.Errorf("missing nodepool %d", poolSpec.ID)
		}

		updatedPoolSpec := updatedSpec.NodePools[i]
		updateOpts, err := poolSpec.Diff(&updatedPoolSpec)
		if err != nil {
			return err
		}
		clusterPoolUpdates[poolSpec.ID] = updateOpts
	}

	updateOpts, err := spec.Diff(&updatedSpec)
	if err != nil {
		return err
	}

	for id, update := range clusterPoolUpdates {
		if update == nil {
			fmt.Fprintf(o.Out, "LKE Node Pool %d not updated\n", id)
			continue
		}
		if _, err := client.UpdateLKEClusterPool(ctx, toEdit.ID, id, *update); err != nil {
			return fmt.Errorf("failed to update LKE Node Pool %d: %w", id, err)
		}
		fmt.Fprintf(o.Out, "LKE Node Pool %d updated\n", id)
	}

	if updateOpts != nil {
		if _, err = client.UpdateLKECluster(ctx, toEdit.ID, *updateOpts); err != nil {
			return fmt.Errorf("failed to update LKE Cluster %d: %w", toEdit.ID, err)
		}
		fmt.Fprintf(o.Out, "LKE Cluster %q updated\n", toEdit.Label)
	} else {
		fmt.Fprintf(o.Out, "LKE Cluster %q not updated\n", toEdit.Label)
	}

	return nil
}
