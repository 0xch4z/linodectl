package instance

import (
	"context"

	"github.com/Charliekenney23/linodectl/internal/resource"
	"github.com/linode/linodego"
)

type Spec struct {
	AuthorizedUsers []string `mapstructure:"authorized_users"`
	BackupsEnabled  bool     `mapstructure:"backups_enabled"`
	Group           string   `mapstructure:"group"`
	Image           string   `mapstructure:"image"`
	PoweredOff      bool     `mapstructure:"powered_off"`
	PrivateIP       bool     `mapstructure:"private_ip"`
	Preset          string   `mapstructure:"preset"`
	Region          string   `mapstructure:"region"`
	RootPass        string   `mapstructure:"root_pass"`
	StackscriptData string   `mapstructure:"stackscript_data"`
	StackscriptID   int      `mapstructure:"stackscript_id"`
	SwapSize        int      `mapstructure:"swap_size"`
	Tags            []string `mapstructure:"tags"`
	Type            string   `mapstructure:"type"`
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
