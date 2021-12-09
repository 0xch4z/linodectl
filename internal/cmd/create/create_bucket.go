package create

import (
	"context"
	"fmt"

	"github.com/Charliekenney23/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/printer"
	"github.com/Charliekenney23/linodectl/internal/resource/bucket"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

const (
	bucketExamples = `	# Create an Object Storage Bucket in us-east-1 cluster
  linodectl create bucket mybucket --cluster us-east-1

  # Create an Object Storage Bucket with CORS enabled
  linodectl create bucket mystaticsite --cluster us-east-1 --cors

  # Create an Object Storage Bucket with a default public-read-write ACL
  linodectl create bucket securebucket --cluster us-east-1 --acl public-read-write`

	bucketACLUsage = "The Access Control Level of the bucket. One of private, public-read, " +
		"authenticated-read, public-read-write (defaults to private)."
)

type CreateBucketOptions struct {
	Label string

	ACL         string
	Cluster     string
	CorsEnabled bool

	genericoptions.ProfileFlags
	genericoptions.PrinterFlags
	cmdutil.IOStreams
}

func NewCreateBucketOptions(ioStreams cmdutil.IOStreams) *CreateBucketOptions {
	return &CreateBucketOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdCreateBucket(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewCreateBucketOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "bucket NAME [args...]",
		Short:   "Create a Bucket",
		Aliases: []string{"bucket"},
		Example: bucketExamples,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(f, ioStreams, args); err != nil {
				return err
			}
			return o.Run(f, cmd)
		},
	}

	cmd.Flags().StringVar(&o.ACL, "acl", o.ACL, bucketACLUsage)
	cmd.Flags().StringVar(&o.Cluster, "cluster", o.Cluster, "The cluster to create this bucket in.")
	cmd.Flags().BoolVar(&o.CorsEnabled, "cors", false, "If true, CORS will be enabled for all origins.")

	o.ProfileFlags.AddFlags(cmd)
	o.PrinterFlags.AddFlags(cmd)
	return cmd
}

func (o *CreateBucketOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	o.Label, err = NameFromCommandArgs(args)
	if err != nil {
		return err
	}
	return nil
}

func (o *CreateBucketOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	options := linodego.ObjectStorageBucketCreateOptions{
		Label:       o.Label,
		Cluster:     o.Cluster,
		ACL:         linodego.ObjectStorageACL(o.ACL),
		CorsEnabled: &o.CorsEnabled,
	}

	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}

	objBucket, err := client.CreateObjectStorageBucket(context.Background(), options)
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	resList := bucket.NewList([]linodego.ObjectStorageBucket{*objBucket})
	p := printer.New(o.Out)
	return p.PrintResources(context.Background(), resList, printer.ResourcePrintOptions{
		Columns:         o.Fields(),
		SortBy:          o.SortBy(),
		OmitHeader:      o.NoHeader(),
		DescendingOrder: o.Descending(),
	})
}
