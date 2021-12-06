package obj

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Charliekenney23/linodectl/internal/linode"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/linode/linodego"
)

func CreateTempKeyPair(ctx context.Context, client linode.Client, buckets []linodego.ObjectStorageBucket) (*linodego.ObjectStorageKey, func(), error) {
	label := fmt.Sprintf("linodectl-tmp%d", time.Now().Unix())

	accessPermissions := make([]linodego.ObjectStorageKeyBucketAccess, len(buckets))
	for i, bucket := range buckets {
		accessPermissions[i] = linodego.ObjectStorageKeyBucketAccess{
			Cluster:     bucket.Cluster,
			BucketName:  bucket.Label,
			Permissions: "read_write",
		}
	}
	key, err := client.CreateObjectStorageKey(ctx, linodego.ObjectStorageKeyCreateOptions{
		Label:        label,
		BucketAccess: &accessPermissions,
	})
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := client.DeleteObjectStorageKey(ctx, key.ID); err != nil {
			log.Printf("failed to delete temporary object storage key %q (%d): %v", key.Label, key.ID, err)
		}
	}
	return key, cleanup, nil
}

func BuildS3Conn(cluster string, key *linodego.ObjectStorageKey) *s3.S3 {
	endpoint := fmt.Sprintf("https://%s.linodeobjects.com", cluster)

	sess, _ := session.NewSession(&aws.Config{
		// This region is hardcoded strictly for preflight validation purposes.
		// The real region is in the endpoint URL.
		Region:      aws.String("us-east-1"),
		Endpoint:    aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(key.AccessKey, key.SecretKey, ""),
	})
	return s3.New(sess)
}
