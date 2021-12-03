package stackscript

import (
	"github.com/Charliekenney23/linodectl/internal/resource/resourceref"
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

// FilterByRefs filters for the referenced Linode Stackscripts.
func FilterByRefs(stackscripts []linodego.Stackscript, refs resourceref.List) (r []linodego.Stackscript) {
	labels, ids := refs.Identifiers()
	for _, ss := range stackscripts {
		if _, ok := ids[ss.ID]; ok {
			r = append(r, ss)
			continue
		}
		if _, ok := labels[ss.Label]; ok {
			r = append(r, ss)
			continue
		}
	}
	return r
}
