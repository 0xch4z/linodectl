package edit

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Charliekenney23/linodectl/internal/cli/editor"
	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/resource/instance"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type EditInstanceOptions struct {
	Label string

	AuthorizedUsers []string
	BackupsEnabled  bool
	Group           string
	Image           string
	PoweredOff      bool
	PrivateIP       bool
	Preset          string
	Region          string
	RootPass        string
	StackscriptData string
	StackscriptID   int
	SwapSize        int
	Tags            []string
	Type            string

	AuthorizeMe bool

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

	cmd.Flags().StringSliceVar(&o.AuthorizedUsers, "authorized-users", o.AuthorizedUsers, "Users to authorize for this instance.")
	cmd.Flags().BoolVar(&o.BackupsEnabled, "enable-backups", false, "If true, backups will be enabled.")
	cmd.Flags().StringVarP(&o.Group, "group", "g", "", "The group to attribute this instance to.")
	cmd.Flags().StringVarP(&o.Image, "image", "i", "", "The image to provision this instance with.")
	cmd.Flags().BoolVar(&o.PoweredOff, "powered-off", false, "If true, the instance will not be booted.")
	cmd.Flags().BoolVar(&o.PrivateIP, "private-ip", false, "If true, a private IP will be allocated for this instance.")
	cmd.Flags().StringVar(&o.Preset, "preset", "", "The preset to use for this instance.")
	cmd.Flags().StringVar(&o.Region, "region", "", "The region to deploy this instance in.")
	cmd.Flags().StringVar(&o.RootPass, "root-pass", "", "The root pass to set on this instance.")
	cmd.Flags().IntVarP(&o.SwapSize, "swap-size", "s", 0, "The swap size for the instance.")
	cmd.Flags().StringSliceVar(&o.Tags, "tags", o.Tags, "The tags to add to this instance.")
	cmd.Flags().StringVar(&o.Type, "type", "", "The type of this instance.")

	o.ProfileFlags.AddFlags(cmd)
	return cmd
}

func (o *EditInstanceOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	if len(args) == 1 {
		o.Label = args[0]
	}

	return nil
}

func (o *EditInstanceOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}

	filter := linodego.Filter{}
	filter.AddField(linodego.Eq, "label", o.Label)
	filterBytes, err := filter.MarshalJSON()
	if err != nil {
		return err
	}

	ctx := context.Background()
	instances, err := client.ListInstances(ctx, &linodego.ListOptions{
		Filter: string(filterBytes),
	})
	if err != nil {
		return err
	}

	if len(instances) == 0 {
		return fmt.Errorf("instance %q does not exist", o.Label)
	}

	i := instances[0]
	spec := instance.SpecFromObject(&i)
	instanceBytes, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}

	editor := editor.NewDefaultEditor()
	specBytes, _, err := editor.EditReader("instance", o.Label+".yaml", o.IOStreams, bytes.NewBuffer(instanceBytes))
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

	if _, err = client.UpdateInstance(ctx, i.ID, *updateOpts); err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "instance %q updated", i.Label)
	return nil
}
