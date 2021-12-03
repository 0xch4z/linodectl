package bucket

import (
	"context"

	"github.com/Charliekenney23/linodectl/internal/resource"
	"github.com/linode/linodego"
)

type Meta struct{}

var _ resource.Meta = (*Meta)(nil)

func (Meta) Name() string {
	return "bucket"
}

func (Meta) RequiredColumns() []string {
	return []string{"label"}
}

func (Meta) DefaultColumns() []string {
	return []string{"cluster"}
}

func (l List) Items() []resource.Resource {
	return l.items
}

func (l List) Meta() resource.Meta {
	return l.meta
}

type List struct {
	meta  Meta
	items []resource.Resource
}

var _ resource.List = (*List)(nil)

type Resource linodego.ObjectStorageBucket

var _ resource.Resource = (*Resource)(nil)

func (r Resource) Properties() resource.PropertyMap {
	return resource.PropertyMap{
		"label": {
			Getter: func(context.Context) (interface{}, error) {
				return r.Label, nil
			},
		},
		"cluster": {
			Getter: func(c context.Context) (interface{}, error) {
				return r.Cluster, nil
			},
		},
		"hostname": {
			Getter: func(context.Context) (interface{}, error) {
				return r.Hostname, nil
			},
		},
	}
}

func NewList(buckets []linodego.ObjectStorageBucket) List {
	resources := make([]resource.Resource, len(buckets))
	for i, cluster := range buckets {
		resources[i] = Resource(cluster)
	}
	return List{items: resources}
}
