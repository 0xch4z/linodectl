package lkecluster

import (
	"context"

	"github.com/Charliekenney23/linodectl/internal/resource"
	"github.com/linode/linodego"
)

type Meta struct{}

var _ resource.Meta = (*Meta)(nil)

func (Meta) Name() string {
	return "lkecluster"
}

func (Meta) RequiredColumns() []string {
	return []string{"id", "label"}
}

func (Meta) DefaultColumns() []string {
	return []string{"region", "k8s_version"}
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

type Resource linodego.LKECluster

var _ resource.Resource = (*Resource)(nil)

func (r Resource) Properties() resource.PropertyMap {
	return resource.PropertyMap{
		"id": {
			Getter: func(context.Context) (interface{}, error) {
				return r.ID, nil
			},
		},
		"label": {
			Getter: func(context.Context) (interface{}, error) {
				return r.Label, nil
			},
		},
		"region": {
			Getter: func(c context.Context) (interface{}, error) {
				return r.Region, nil
			},
		},
		"k8s_version": {
			Getter: func(context.Context) (interface{}, error) {
				return r.K8sVersion, nil
			},
		},
		"status": {
			Getter: func(c context.Context) (interface{}, error) {
				return r.Status, nil
			},
		},
		"tags": {
			Getter: func(c context.Context) (interface{}, error) {
				return r.Tags, nil
			},
		},
	}
}

func NewList(clusters []linodego.LKECluster) List {
	resources := make([]resource.Resource, len(clusters))
	for i, cluster := range clusters {
		resources[i] = Resource(cluster)
	}
	return List{items: resources}
}
