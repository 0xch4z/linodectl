package linode

//go:generate mockgen -destination mock/mock_client.go -package mock github.com/Charliekenney23/linodectl/internal/linode Client

import (
	"context"
	"os"

	"github.com/Charliekenney23/linodectl/internal/config"
	"github.com/linode/linodego"
)

const (
	userAgent = "linodectl"
)

type Client interface {
	GetInstance(context.Context, int) (*linodego.Instance, error)
	UpdateInstance(context.Context, int, linodego.InstanceUpdateOptions) (*linodego.Instance, error)
	CreateInstance(context.Context, linodego.InstanceCreateOptions) (*linodego.Instance, error)
	DeleteInstance(context.Context, int) error
	ListInstances(context.Context, *linodego.ListOptions) ([]linodego.Instance, error)

	CreateLKECluster(context.Context, linodego.LKEClusterCreateOptions) (*linodego.LKECluster, error)
	ListLKEClusters(context.Context, *linodego.ListOptions) ([]linodego.LKECluster, error)
	ListLKEClusterPools(context.Context, int, *linodego.ListOptions) ([]linodego.LKEClusterPool, error)
	GetLKEClusterKubeconfig(context.Context, int) (*linodego.LKEClusterKubeconfig, error)
	UpdateLKECluster(context.Context, int, linodego.LKEClusterUpdateOptions) (*linodego.LKECluster, error)
	UpdateLKEClusterPool(context.Context, int, int, linodego.LKEClusterPoolUpdateOptions) (*linodego.LKEClusterPool, error)
	DeleteLKECluster(context.Context, int) error

	ListStackscripts(context.Context, *linodego.ListOptions) ([]linodego.Stackscript, error)
	DeleteStackscript(context.Context, int) error

	ListObjectStorageBuckets(context.Context, *linodego.ListOptions) ([]linodego.ObjectStorageBucket, error)
	DeleteObjectStorageBucket(context.Context, string, string) error
	CreateObjectStorageKey(context.Context, linodego.ObjectStorageKeyCreateOptions) (*linodego.ObjectStorageKey, error)
	DeleteObjectStorageKey(context.Context, int) error

	GetProfile(context.Context) (*linodego.Profile, error)
}

// *linodego.Client implements Client
var _ Client = (*linodego.Client)(nil)

func NewClient(profile config.Profile) Client {
	client := linodego.NewClient(nil)

	if profile.APIVersion != "" {
		client.SetAPIVersion(profile.APIVersion)
	}
	if profile.APIBaseURL != "" {
		client.SetBaseURL(profile.APIBaseURL)
	}

	token := profile.Token
	if token == "" {
		token = os.Getenv("LINODE_API_TOKEN")
	}

	client.SetToken(token)
	client.SetUserAgent(userAgent)
	return &client
}
