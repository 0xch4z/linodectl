package genericoptions

import "github.com/spf13/cobra"

// ProfileFlags provides flags for specifying a profile.
type ProfileFlags struct {
	profileName string
}

// AddFlags recieves a *cobra.Command reference and binds a flag for specifying
// a profile.
func (f *ProfileFlags) AddFlags(c *cobra.Command) {
	if f == nil {
		return
	}

	c.PersistentFlags().StringVarP(&f.profileName, "profile", "p", "", "The profile to use for communicating with the Linode API")
}

// ProfileName gets the name of the specified profile from the flag.
func (f *ProfileFlags) ProfileName() string {
	return f.profileName
}
