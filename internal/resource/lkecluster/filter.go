package lkecluster

import (
	"github.com/Charliekenney23/linodectl/internal/resource/resourceref"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type FilterFlags struct {
	tag string
}

// AddFlags recieves a *cobra.Command reference and binds a flag for specifying
// a preset.
func (f *FilterFlags) AddFlags(c *cobra.Command) {
	if f == nil {
		return
	}

	c.Flags().StringVarP(&f.tag, "tag", "t", "", "Filter for LKE Clusters with this tag")
}

func (f *FilterFlags) Filter(label string) *linodego.Filter {
	filter := new(linodego.Filter)

	if label != "" {
		filter.AddField(linodego.Eq, "label", label)
	}

	if f.tag != "" {
		filter.AddField(linodego.Eq, "tag", f.tag)
	}
	return filter
}

// FilterByRefs filters for the referenced LKE Clusters.
func FilterByRefs(clusters []linodego.LKECluster, refs resourceref.List) (r []linodego.LKECluster) {
	labels, ids := refs.Identifiers()
	for _, cluster := range clusters {
		if _, ok := ids[cluster.ID]; ok {
			r = append(r, cluster)
			continue
		}
		if _, ok := labels[cluster.Label]; ok {
			r = append(r, cluster)
			continue
		}
	}
	return r
}
