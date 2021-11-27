package stackscript

import (
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type FilterFlags struct {
	mine bool
}

// AddFlags recieves a *cobra.Command reference and binds a flag for specifying
// a preset.
func (f *FilterFlags) AddFlags(c *cobra.Command) {
	if f == nil {
		return
	}

	c.Flags().BoolVar(&f.mine, "mine", false, "Filter for stackscripts owned by you")
}

func (f *FilterFlags) Filter(label string) *linodego.Filter {
	filter := new(linodego.Filter)

	if label != "" {
		filter.AddField(linodego.Eq, "label", label)
	}

	if f.mine {
		filter.AddField(linodego.Eq, "mine", true)
	}
	return filter
}
