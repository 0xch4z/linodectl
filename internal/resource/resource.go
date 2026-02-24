package resource

import (
	"context"
	"sync"
)

// PropertyMap represents a mapping of property names to getters and setters.
type PropertyMap map[string]*Property

// Resource represents a Linode resource
type Resource interface {
	Properties() PropertyMap
}

// Resource represents a list of Linode resources.
type List interface {
	Meta() Meta
	Items() []Resource
}

type Meta interface {
	Name() string
	DefaultColumns() []string
	RequiredColumns() []string
}

type Property struct {
	// Getter gets the property
	Getter func(context.Context) (any, error)

	// Setter sets the property. Setter is not required
	Setter func(context.Context, any) error

	// value is the cached result of Getter
	value any

	sync.Once
}

// GetWithCache
func (p *Property) GetWithCache(ctx context.Context) (any, error) {
	var err error
	p.Do(func() {
		p.value, err = p.Getter(ctx)
	})
	return p.value, err
}

// IsSettable returns whether or not a property has a setter.
func (p *Property) IsSettable() bool {
	return p.Setter != nil
}
