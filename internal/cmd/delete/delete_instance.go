package delete

import (
	"context"
	"errors"
	"fmt"

	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/resource/instance"
	"github.com/Charliekenney23/linodectl/internal/resource/resourceref"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type DeleteInstanceOptions struct {
	refs resourceref.List

	genericoptions.PaginationFlags
	genericoptions.ProfileFlags
	instance.FilterFlags
	cmdutil.IOStreams
}

func NewDeleteInstanceOptions(ioStreams cmdutil.IOStreams) *DeleteInstanceOptions {
	return &DeleteInstanceOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdDeleteInstance(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewDeleteInstanceOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "instance [NAME] [args...]",
		Aliases: []string{"instances", "linode", "linodes"},
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.Complete(f, ioStreams, args); err != nil {
				panic(err)
			}
			if err := o.Run(f, cmd); err != nil {
				panic(err)
			}
		},
	}

	o.PaginationFlags.AddFlags(cmd)
	o.ProfileFlags.AddFlags(cmd)
	o.FilterFlags.AddFlags(cmd)
	return cmd
}

func (o *DeleteInstanceOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	if o.refs, err = resourceref.ListFromArgs(args); err != nil {
		return err
	}
	return nil
}

func (o *DeleteInstanceOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	filter := o.Filter(o.refs.Label())

	if len(o.refs) == 0 && len(filter.Children) == 0 && o.LKECluster() == 0 {
		// we can't just delete every instance
		return errors.New("no filters provided")
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

	instances, err := client.ListInstances(context.Background(), &linodego.ListOptions{
		PageOptions: o.PageOptions(),
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

	for _, instance := range instances {
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			return err
		}
		fmt.Fprintf(o.Out, "Instance %q (%d) deleted...\n", instance.Label, instance.ID)
	}
	return nil
}
