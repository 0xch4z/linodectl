package instance

import (
	"context"
	"net"
	"reflect"
	"time"

	"github.com/Charliekenney23/linodectl/internal/ptr"
	"github.com/Charliekenney23/linodectl/internal/resource"
	"github.com/linode/linodego"
)

type Spec struct {
	ID              int                      `yaml:"id"`
	Created         *time.Time               `yaml:"created"`
	Updated         *time.Time               `yaml:"updated"`
	Region          string                   `yaml:"region"`
	Alerts          *linodego.InstanceAlert  `yaml:"alerts"`
	Backups         *linodego.InstanceBackup `yaml:"backups"`
	Image           string                   `yaml:"image"`
	Group           string                   `yaml:"group"`
	IPv4            []*net.IP                `yaml:"ipv4"`
	IPv6            string                   `yaml:"ipv6"`
	Label           string                   `yaml:"label"`
	Type            string                   `yaml:"type"`
	Status          linodego.InstanceStatus  `yaml:"status"`
	Hypervisor      string                   `yaml:"hypervisor"`
	Specs           *linodego.InstanceSpec   `yaml:"specs"`
	WatchdogEnabled bool                     `yaml:"watchdogEnabled"`
	Tags            []string                 `yaml:"tags"`
}

func SpecFromObject(instance *linodego.Instance) *Spec {
	return &Spec{
		ID:              instance.ID,
		Created:         instance.Created,
		Updated:         instance.Updated,
		Region:          instance.Region,
		Alerts:          instance.Alerts,
		Backups:         instance.Backups,
		Image:           instance.Image,
		Group:           instance.Group,
		IPv4:            instance.IPv4,
		IPv6:            instance.IPv6,
		Label:           instance.Label,
		Type:            instance.Type,
		Status:          instance.Status,
		Hypervisor:      instance.Hypervisor,
		Specs:           instance.Specs,
		WatchdogEnabled: instance.WatchdogEnabled,
		Tags:            instance.Tags,
	}
}

func (s *Spec) Diff(in *Spec) (*linodego.InstanceUpdateOptions, error) {
	o := new(linodego.InstanceUpdateOptions)
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
	if !reflect.DeepEqual(s.Alerts, in.Alerts) {
		o.Alerts = in.Alerts
	}
	if !reflect.DeepEqual(s.Backups, in.Backups) {
		o.Backups = in.Backups
	}
	if s.Image != in.Image {
		return nil, resource.NotUpdateableError{Name: "image"}
	}
	if s.Group != in.Group {
		o.Group = in.Group
	}
	if !reflect.DeepEqual(s.IPv4, in.IPv4) {
		return nil, resource.NotUpdateableError{Name: "ipv4"}
	}
	if s.IPv6 != in.IPv6 {
		return nil, resource.NotUpdateableError{Name: "ipv6"}
	}
	if s.Label != in.Label {
		o.Label = in.Label
	}
	if s.Type != in.Type {
		return nil, resource.NotUpdateableError{Name: "type"}
	}
	if s.Status != in.Status {
		return nil, resource.NotUpdateableError{Name: "status"}
	}
	if s.Hypervisor != in.Hypervisor {
		return nil, resource.NotUpdateableError{Name: "hypervisor"}
	}
	if !reflect.DeepEqual(s.Specs, in.Specs) {
		return nil, resource.NotUpdateableError{Name: "specs"}
	}
	if s.WatchdogEnabled != in.WatchdogEnabled {
		o.WatchdogEnabled = ptr.Bool(in.WatchdogEnabled)
	}
	if !reflect.DeepEqual(s.Tags, in.Tags) {
		o.Tags = &in.Tags
	}
	return o, nil
}

type Meta struct{}

var _ resource.Meta = (*Meta)(nil)

func (Meta) Name() string {
	return "instance"
}

func (Meta) RequiredColumns() []string {
	return []string{"id", "label"}
}

func (Meta) DefaultColumns() []string {
	return []string{"region", "type", "image", "status", "ipv4"}
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

type Resource linodego.Instance

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
		"image": {
			Getter: func(context.Context) (interface{}, error) {
				return r.Image, nil
			},
		},
		"type": {
			Getter: func(context.Context) (interface{}, error) {
				return r.Type, nil
			},
		},
		"region": {
			Getter: func(c context.Context) (interface{}, error) {
				return r.Region, nil
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
		"ipv4": {
			Getter: func(c context.Context) (interface{}, error) {
				return r.IPv4, nil
			},
		},
		"ipv6": {
			Getter: func(c context.Context) (interface{}, error) {
				return r.IPv6, nil
			},
		},
	}
}

func NewList(instances []linodego.Instance) List {
	resources := make([]resource.Resource, len(instances))
	for i, instance := range instances {
		resources[i] = Resource(instance)
	}
	return List{items: resources}
}
