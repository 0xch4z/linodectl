package bucket

import (
	"github.com/Charliekenney23/linodectl/internal/resource/resourceref"
	"github.com/linode/linodego"
)

// FilterByRefs filters for the referenced OBJ Buckets.
func FilterByRefs(buckets []linodego.ObjectStorageBucket, refs resourceref.List) (r []linodego.ObjectStorageBucket) {
	labels, _ := refs.Identifiers()
	for _, bucket := range buckets {
		if _, ok := labels[bucket.Label]; ok {
			r = append(r, bucket)
			continue
		}
	}
	return r
}
