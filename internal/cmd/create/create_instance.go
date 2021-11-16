package create

import (
	"context"
	"fmt"

	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/printer"
	"github.com/Charliekenney23/linodectl/internal/ptr"
	"github.com/Charliekenney23/linodectl/internal/resource/instance"
	"github.com/Charliekenney23/linodectl/internal/strutil"
	"github.com/spf13/cobra"

	"github.com/linode/linodego"
)

type CreateInstanceOptions struct {
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

	genericoptions.PrinterFlags
	genericoptions.ProfileFlags
	cmdutil.IOStreams
}

func NewCreateInstanceOptions(ioStreams cmdutil.IOStreams) *CreateInstanceOptions {
	return &CreateInstanceOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdCreateInstance(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewCreateInstanceOptions(ioStreams)

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
	cmd.Flags().BoolVar(&o.AuthorizeMe, "authorize-me", false, "If true, this user account's keys will be authorized.")
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
	o.PrinterFlags.AddFlags(cmd)
	return cmd
}

func (o *CreateInstanceOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	o.Label, err = NameFromCommandArgs(args)
	if err != nil {
		return err
	}

	if o.Preset != "" {
		// TODO: ensure default preset here
		preset, ok := f.Config().Instances.Presets[o.Preset]
		if !ok {
			return fmt.Errorf("preset %q does not exist", o.Preset)
		}

		o.AuthorizedUsers = strutil.SliceFallback(o.AuthorizedUsers, preset.AuthorizedUsers)
		o.Tags = strutil.SliceFallback(o.Tags, preset.Tags)

		// TODO: load other bool/int defaults from preset here
		o.Image = strutil.Fallback(o.Image, preset.Image)
		o.Group = strutil.Fallback(o.Group, preset.Group)
		o.Type = strutil.Fallback(o.Type, preset.Type)
		o.RootPass = strutil.Fallback(o.RootPass, preset.RootPass)
	}

	if profile, ok := f.Config().CurrentProfile(); ok {
		o.Region = strutil.Fallback(o.Region, profile.Region)
	}
	return nil
}

func (o *CreateInstanceOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	options := linodego.InstanceCreateOptions{
		Label:           o.Label,
		AuthorizedUsers: o.AuthorizedUsers,
		BackupsEnabled:  o.BackupsEnabled,
		Booted:          ptr.Bool(!o.PoweredOff),
		Group:           o.Group,
		Image:           o.Image,
		Region:          o.Region,
		PrivateIP:       o.PrivateIP,
		RootPass:        o.RootPass,
		Tags:            o.Tags,
		Type:            o.Type,
	}

	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}
	ctx := context.Background()

	if o.AuthorizeMe {
		profile, err := client.GetProfile(ctx)
		if err != nil {
			return fmt.Errorf("failed to get profile: %w", err)
		}
		options.AuthorizedUsers = append(options.AuthorizedUsers, profile.Username)
	}

	linodeInstance, err := client.CreateInstance(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}

	resList := instance.NewList([]linodego.Instance{*linodeInstance})
	p := printer.New(o.Out)
	return p.PrintResources(context.Background(), resList, printer.ResourcePrintOptions{
		Columns:         o.Fields(),
		SortBy:          o.SortBy(),
		OmitHeader:      o.NoHeader(),
		DescendingOrder: o.Descending(),
	})
}
