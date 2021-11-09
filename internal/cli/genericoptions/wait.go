package genericoptions

import (
	"github.com/spf13/cobra"
)

type WaitFlags struct {
	wait    bool
	timeout string
}

func (f *WaitFlags) AddFlags(c *cobra.Command) {
	if f == nil {
		return
	}

	c.Flags().BoolVar(&f.wait, "wait", false, "Wait for resource to reach the desired state")
	c.Flags().StringVar(&f.timeout, "timeout", f.timeout, "The amount of time to wait for the resource to reach the desired state")
}
