package get

import (
	"context"

	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/printer"
	"github.com/Charliekenney23/linodectl/internal/resource/bucket"
	"github.com/Charliekenney23/linodectl/internal/resource/resourceref"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type GetBucketOptions struct {
	refs resourceref.List

	genericoptions.PaginationFlags
	genericoptions.ProfileFlags
	genericoptions.PrinterFlags
	cmdutil.IOStreams
}

func NewGetBucketOptions(ioStreams cmdutil.IOStreams) *GetBucketOptions {
	return &GetBucketOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdGetBucket(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewGetBucketOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "bucket [NAME] [args...]",
		Aliases: []string{"buckets", "objbucket", "objbuckets"},
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
	return cmd
}

func (o *GetBucketOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	if o.refs, err = resourceref.ListFromArgs(args); err != nil {
		return err
	}
	return nil
}

func (o *GetBucketOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	// Currently, the Linode API does not support filtering by label on
	// OBJ Buckets, even though the docs say otherwise. When this is fixed,
	// we should pass the string o.refs.Label(), so that if a single label
	// is specified, only that cluster is fetched.
	// If it's included in the request right now, we'll get a 400 error.

	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}

	ctx := context.Background()
	buckets, err := client.ListObjectStorageBuckets(ctx, &linodego.ListOptions{})
	if err != nil {
		return err
	}

	if len(o.refs) > 0 {
		buckets = bucket.FilterByRefs(buckets, o.refs)
	}

	resourceList := bucket.NewList(buckets)
	p := printer.New(o.Out)
	return p.PrintResources(context.Background(), resourceList, printer.ResourcePrintOptions{
		Columns:         o.Fields(),
		SortBy:          o.SortBy(),
		OmitHeader:      o.NoHeader(),
		DescendingOrder: o.Descending(),
	})
}
