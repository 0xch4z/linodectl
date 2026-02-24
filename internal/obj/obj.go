package obj

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/0xch4z/linodectl/internal/linode"
	"github.com/linode/linodego"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

func BuildS3Conn(cluster string, key *linodego.ObjectStorageKey) (*minio.Client, error) {
	endpoint := fmt.Sprintf("%s.linodeobjects.com", cluster)

	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(key.AccessKey, key.SecretKey, ""),
		Secure: true,
	})
}
