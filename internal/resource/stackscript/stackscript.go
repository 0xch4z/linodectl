package stackscript

import (
	"context"

	"github.com/0xch4z/linodectl/internal/resource"
	"github.com/linode/linodego"
)

type Meta struct{}

var _ resource.Meta = (*Meta)(nil)

func (Meta) Name() string {
	return "stackscript"
}

func (Meta) RequiredColumns() []string {
	return []string{"id", "label"}
}

func (Meta) DefaultColumns() []string {
	return []string{"mine", "is_public", "images"}
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

type Resource linodego.Stackscript

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
		"mine": &resource.Property{
			Getter: func(context.Context) (any, error) {
				return r.Mine, nil
			},
		},
		"is_public": &resource.Property{
			Getter: func(context.Context) (any, error) {
				return r.IsPublic, nil
			},
		},
		"images": &resource.Property{
			Getter: func(context.Context) (any, error) {
				return r.Images, nil
			},
		},
	}
}

func NewList(stackScripts []linodego.Stackscript) List {
	resources := make([]resource.Resource, len(stackScripts))
	for i, cluster := range stackScripts {
		resources[i] = Resource(cluster)
	}
	return List{items: resources}
}
