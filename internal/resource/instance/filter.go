package instance

import (
	"context"
	"fmt"

	"github.com/0xch4z/linodectl/internal/linode"
	"github.com/0xch4z/linodectl/internal/resource/resourceref"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type FilterFlags struct {
	group      string
	image      string
	region     string
	tag        string
	lkeCluster int
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
	c.Flags().IntVar(&f.lkeCluster, "lke-cluster", 0, "Filter for instances associated with this cluster ID")
}

func (f *FilterFlags) LKECluster() int {
	return f.lkeCluster
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

// FilterLKECluster filters for LKE Cluster Pool Nodes a slice of linodego.Instance, a client
// and an LKE Cluster ID.
func FilterLKECluster(ctx context.Context, client linode.Client, id int, instances []linodego.Instance) ([]linodego.Instance, error) {
	clusterPools, err := client.ListLKEClusterPools(ctx, id, &linodego.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pools for LKE Cluster %d: %w", id, err)
	}

	clusterInstanceIDs := make(map[int]struct{})
	for _, clusterPool := range clusterPools {
		for _, node := range clusterPool.Linodes {
			clusterInstanceIDs[node.InstanceID] = struct{}{}
		}
	}

	clusterInstances := make([]linodego.Instance, 0, len(clusterInstanceIDs))
	for _, instance := range instances {
		if _, ok := clusterInstanceIDs[instance.ID]; ok {
			clusterInstances = append(clusterInstances, instance)
		}
	}
	return clusterInstances, nil
}

// FilterByRefs filters for the referenced Linode Instances.
func FilterByRefs(instances []linodego.Instance, refs resourceref.List) (r []linodego.Instance) {
	labels, ids := refs.Identifiers()
	for _, instance := range instances {
		if _, ok := ids[instance.ID]; ok {
			r = append(r, instance)
			continue
		}
		if _, ok := labels[instance.Label]; ok {
			r = append(r, instance)
			continue
		}
	}
	return r
}
