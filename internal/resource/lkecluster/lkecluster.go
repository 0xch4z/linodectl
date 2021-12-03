package lkecluster

import (
	"context"
	"reflect"
	"time"

	"github.com/Charliekenney23/linodectl/internal/resource"
	"github.com/linode/linodego"
)

type Spec struct {
	ID           int                             `yaml:"id"`
	Created      *time.Time                      `yaml:"created"`
	Updated      *time.Time                      `yaml:"updated"`
	Region       string                          `yaml:"region"`
	K8sVersion   string                          `yaml:"k8s_version"`
	Status       linodego.LKEClusterStatus       `yaml:"status"`
	Tags         []string                        `yaml:"tags"`
	ControlPlane linodego.LKEClusterControlPlane `yaml:"control_plane"`
}

func SpecFromObject(cluster *linodego.LKECluster) *Spec {
	return &Spec{
		ID:           cluster.ID,
		Created:      cluster.Created,
		Updated:      cluster.Updated,
		Region:       cluster.Region,
		K8sVersion:   cluster.K8sVersion,
		Status:       cluster.Status,
		Tags:         cluster.Tags,
		ControlPlane: cluster.ControlPlane,
	}
}

func (s *Spec) Diff(in *Spec) (*linodego.LKEClusterUpdateOptions, error) {
	o := new(linodego.LKEClusterUpdateOptions)
	if s.ID != in.ID {
		return nil, resource.NotUpdateableError{Name: "id"}
	}
	if !reflect.DeepEqual(s.Created, in.Created) {
		return nil, resource.NotUpdateableError{Name: "created"}
	}
	if !reflect.DeepEqual(s.Updated, in.Updated) {
		return nil, resource.NotUpdateableError{Name: "updated"}
	}
	if s.Region != in.Region {
		return nil, resource.NotUpdateableError{Name: "region"}
	}
	if s.K8sVersion != in.K8sVersion {
		o.K8sVersion = in.K8sVersion
	}
	if s.Status != in.Status {
		return nil, resource.NotUpdateableError{Name: "status"}
	}
	if !reflect.DeepEqual(s.Tags, in.Tags) {
		o.Tags = &in.Tags
	}
	if !reflect.DeepEqual(s.ControlPlane, in.ControlPlane) {
		o.ControlPlane = &in.ControlPlane
	}
	return o, nil
}

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
