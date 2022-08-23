package delete

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/0xch4z/linodectl/internal/cli/genericoptions"
	cmdutil "github.com/0xch4z/linodectl/internal/cmd/util"
	"github.com/0xch4z/linodectl/internal/obj"
	"github.com/0xch4z/linodectl/internal/resource/bucket"
	"github.com/0xch4z/linodectl/internal/resource/resourceref"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type DeleteBucketOptions struct {
	refs resourceref.List

	force bool

	genericoptions.PaginationFlags
	genericoptions.ProfileFlags
	genericoptions.PrinterFlags
	cmdutil.IOStreams
}

func NewDeleteBucketOptions(ioStreams cmdutil.IOStreams) *DeleteBucketOptions {
	return &DeleteBucketOptions{
		IOStreams: ioStreams,
	}
}
func NewCmdDeleteBucket(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewDeleteBucketOptions(ioStreams)

	cmd := &cobra.Command{
		Use:     "bucket [NAME] [args...]",
		Aliases: []string{"buckets"},
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
	cmd.Flags().BoolVar(&o.force, "force", false, "If specified, the bucket will be emptied before deletion.")
	return cmd
}

func (o *DeleteBucketOptions) Complete(f cmdutil.Factory, ioStreams cmdutil.IOStreams, args []string) (err error) {
	if o.refs, err = resourceref.ListFromArgs(args); err != nil {
		return err
	}
	return nil
}

func (o *DeleteBucketOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	if len(o.refs) == 0 {
		// we can't just delete every bucket
		return errors.New("at least one Bucket label is required")
	}

	client, err := f.Client(o.ProfileName())
	if err != nil {
		return err
	}
	ctx := context.Background()

	buckets, err := client.ListObjectStorageBuckets(ctx, &linodego.ListOptions{
		PageOptions: o.PageOptions(),
		PageSize:    o.PageOptions().Results,
	})
	if err != nil {
		return err
	}

	toDelete := bucket.FilterByRefs(buckets, o.refs)
	if len(toDelete) == 0 {
		return nil
	}

	var keypair *linodego.ObjectStorageKey
	if o.force {
		key, cleanup, err := obj.CreateTempKeyPair(ctx, client, toDelete)
		if err != nil {
			return fmt.Errorf("failed to create temporary Object Storage key for bucket deletion: %s", err)
		}
		defer cleanup()
		keypair = key
	}

	for _, bucket := range toDelete {
		if o.force {
			// ensure bucket is empty before attempting deletion
			conn := obj.BuildS3Conn(bucket.Cluster, keypair)

			iter := s3manager.NewDeleteListIterator(conn, &s3.ListObjectsInput{
				Bucket: aws.String(bucket.Label),
			})
			if err := s3manager.NewBatchDeleteWithClient(conn).Delete(context.TODO(), iter); err != nil {
				return err
			}
		}

		if err := client.DeleteObjectStorageBucket(context.Background(), bucket.Cluster, bucket.Label); err != nil {
			if linodeErr, ok := err.(*linodego.Error); ok && strings.HasSuffix(linodeErr.Message, "Please delete all objects and try again.") {
				return fmt.Errorf("bucket %q is not empty. Use --force to empty before deletion.", bucket.Label)
			}
			return fmt.Errorf("failed to delete Bucket %q: %w", bucket.Label, err)
		}
		fmt.Fprintf(o.Out, "Bucket %q deleted...\n", bucket.Label)
	}
	return nil
}
