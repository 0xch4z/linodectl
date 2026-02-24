package lkecluster

import (
	"context"
	"reflect"
	"time"

	"github.com/0xch4z/linodectl/internal/resource"
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
	NodePools    []NodePoolSpec                  `yaml:"node_pools"`
}

type NodePoolSpec struct {
	ID         int                               `yaml:"id"`
	Count      int                               `yaml:"count"`
	Type       string                            `yaml:"type"`
	Disks      []linodego.LKEClusterPoolDisk     `yaml:"disks"`
	Linodes    []linodego.LKEClusterPoolLinode   `yaml:"nodes"`
	Tags       []string                          `yaml:"tag"`
	Autoscaler linodego.LKEClusterPoolAutoscaler `yaml:"autoscaler"`
}

func (s *NodePoolSpec) Diff(in *NodePoolSpec) (*linodego.LKEClusterPoolUpdateOptions, error) {
	o := new(linodego.LKEClusterPoolUpdateOptions)
	if s.ID != in.ID {
		return nil, resource.NotUpdateableError{Name: "id"}
	}
	if s.Count != in.Count {
		o.Count = in.Count
	}
	if s.Type != in.Type {
		return nil, resource.NotUpdateableError{Name: "type"}
	}
	if !reflect.DeepEqual(s.Disks, in.Disks) {
		return nil, resource.NotUpdateableError{Name: "disks"}
	}
	if !reflect.DeepEqual(s.Linodes, in.Linodes) {
		return nil, resource.NotUpdateableError{Name: "linodes"}
	}
	if !reflect.DeepEqual(s.Tags, in.Tags) {
		o.Tags = &in.Tags
	}
	if !reflect.DeepEqual(s.Autoscaler, in.Autoscaler) {
		o.Autoscaler = &in.Autoscaler
	}

	if reflect.DeepEqual(*o, linodego.LKEClusterPoolUpdateOptions{}) {
		// nothing to update
		return nil, nil
	}
	return o, nil
}

func SpecFromObject(cluster *linodego.LKECluster, pools []linodego.LKEClusterPool) *Spec {
	nodePoolSpecs := make([]NodePoolSpec, len(pools))
	for i, pool := range pools {
		nodePoolSpecs[i] = NodePoolSpec{
			ID:         pool.ID,
			Count:      pool.Count,
			Type:       pool.Type,
			Disks:      pool.Disks,
			Linodes:    pool.Linodes,
			Tags:       pool.Tags,
			Autoscaler: pool.Autoscaler,
		}
	}
	return &Spec{
		ID:           cluster.ID,
		Created:      cluster.Created,
		Updated:      cluster.Updated,
		Region:       cluster.Region,
		K8sVersion:   cluster.K8sVersion,
		Status:       cluster.Status,
		Tags:         cluster.Tags,
		ControlPlane: cluster.ControlPlane,
		NodePools:    nodePoolSpecs,
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
		o.ControlPlane = &linodego.LKEClusterControlPlaneOptions{
			HighAvailability: &in.ControlPlane.HighAvailability,
		}
	}

	if reflect.DeepEqual(*o, linodego.LKEClusterUpdateOptions{}) {
		// nothing to update
		return nil, nil
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
		"id": &resource.Property{
			Getter: func(context.Context) (any, error) {
				return r.ID, nil
			},
		},
		"label": &resource.Property{
			Getter: func(context.Context) (any, error) {
				return r.Label, nil
			},
		},
		"region": &resource.Property{
			Getter: func(context.Context) (any, error) {
				return r.Region, nil
			},
		},
		"k8s_version": &resource.Property{
			Getter: func(context.Context) (any, error) {
				return r.K8sVersion, nil
			},
		},
		"status": &resource.Property{
			Getter: func(context.Context) (any, error) {
				return r.Status, nil
			},
		},
		"tags": &resource.Property{
			Getter: func(context.Context) (any, error) {
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
