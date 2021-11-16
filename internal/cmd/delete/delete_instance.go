package delete

import (
	"context"
	"errors"
	"fmt"

	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/resource/instance"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type DeleteInstanceOptions struct {
	// Label (optional) is the name of an instance to delete
	Label string

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
	if len(args) == 1 {
		o.Label = args[0]
	}

	return nil
}

func (o *DeleteInstanceOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	filter := o.Filter(o.Label)

	if len(filter.Children) == 0 {
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

	instances, err := client.ListInstances(context.Background(), &linodego.ListOptions{
		PageOptions: o.PageOptions(),
		Filter:      string(filterBytes),
	})
	if err != nil {
		return err
	}

	for _, instance := range instances {
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			return err
		}
		fmt.Fprintf(o.Out, "Instance %q (%d) deleted...\n", instance.Label, instance.ID)
	}
	return nil
}
