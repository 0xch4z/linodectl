package genericoptions

import (
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
)

type PaginationFlags struct {
	limit int
	page  int
}

func (f *PaginationFlags) AddFlags(c *cobra.Command) {
	if f == nil {
		return
	}

	c.Flags().IntVar(&f.limit, "limit", f.limit, "The number of resources to query for")
	c.Flags().IntVar(&f.page, "page", f.page, "The page number of resources to fetch")
}

func (f *PaginationFlags) PageOptions() *linodego.PageOptions {
	return &linodego.PageOptions{
		Page:    f.page,
		Results: f.limit,
	}
}
