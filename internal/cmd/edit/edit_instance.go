package edit

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Charliekenney23/linodectl/internal/cli/editor"
	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/resource/instance"
	"github.com/Charliekenney23/linodectl/internal/resource/resourceref"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type EditInstanceOptions struct {
	refs resourceref.List

	genericoptions.ProfileFlags
	cmdutil.IOStreams
}

func NewEditInstanceOptions(ioStreams cmdutil.IOStreams) *EditInstanceOptions {
	return &EditInstanceOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdEditInstance(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewEditInstanceOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "instance NAME [args...]",
		Aliases: []string{"linode"},
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

func (o *EditInstanceOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	if o.refs, err = resourceref.ListFromArgs(args); err != nil {
		return err
	}
	return nil
}

func (o *EditInstanceOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	if len(o.refs) == 0 {
		return fmt.Errorf("exactly one Instance ID or label must be specified")
	}

	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}

	ctx := context.Background()
	var filterBytes []byte

	if o.refs.Label() != "" {
		filter := linodego.Filter{}
		filter.AddField(linodego.Eq, "label", o.refs.Label())
		if filterBytes, err = filter.MarshalJSON(); err != nil {
			return err
		}
	}

	instances, err := client.ListInstances(ctx, &linodego.ListOptions{
		Filter: string(filterBytes),
	})
	if err != nil {
		return err
	}

	instances = instance.FilterByRefs(instances, o.refs)

	if len(instances) != 1 {
		return fmt.Errorf("could not find instance %v", o.refs[0])
	}

	toEdit := instances[0]
	spec := instance.SpecFromObject(&toEdit)
	instanceBytes, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}

	editor := editor.NewDefaultEditor()
	specBytes, _, err := editor.EditReader("instance", toEdit.Label+".yaml", o.IOStreams, bytes.NewBuffer(instanceBytes))
	if err != nil {
		panic(err)
	}

	var updatedSpec instance.Spec
	if err := yaml.Unmarshal(specBytes, &updatedSpec); err != nil {
		return err
	}

	updateOpts, err := spec.Diff(&updatedSpec)
	if err != nil {
		return err
	}

	if _, err = client.UpdateInstance(ctx, toEdit.ID, *updateOpts); err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "instance %q updated", toEdit.Label)
	return nil
}
