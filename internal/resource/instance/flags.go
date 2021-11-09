package instance

import (
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type FilterFlags struct {
	group  string
	image  string
	region string
	tag    string

	// TODO: add --lkecluster filter
}

// AddFlags recieves a *cobra.Command reference and binds a flag for specifying
// a preset.
func (f *FilterFlags) AddFlags(c *cobra.Command) {
	if f == nil {
		return
	}

	c.Flags().StringVarP(&f.group, "group", "g", "", "Filter for instances with this group")
	c.Flags().StringVar(&f.image, "image", "", "Filter for instances with this image")
	c.Flags().StringVar(&f.region, "region", "", "Filter for instances in this region")
	c.Flags().StringVarP(&f.tag, "tag", "t", "", "Filter for instances with this tag")
}

func (f *FilterFlags) Filter(label string) *linodego.Filter {
	filter := new(linodego.Filter)

	if label != "" {
		filter.AddField(linodego.Eq, "label", label)
	}
	if f.group != "" {
		filter.AddField(linodego.Eq, "group", f.group)
	}
	if f.image != "" {
		filter.AddField(linodego.Eq, "image", f.image)
	}
	if f.region != "" {
		filter.AddField(linodego.Eq, "region", f.region)
	}
	if f.tag != "" {
		filter.AddField(linodego.Eq, "tags", f.tag)
	}
	return filter
}
