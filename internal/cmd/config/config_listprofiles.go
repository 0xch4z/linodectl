package config

import (
	"fmt"
	"sort"

	cmdutil "github.com/0xch4z/linodectl/internal/cmd/util"
	"github.com/0xch4z/linodectl/internal/strutil"
	"github.com/spf13/cobra"
)

func NewCmdConfigListProfiles(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	return &cobra.Command{
		Use:   "list-profiles",
		Short: "List all profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := f.Config()
			currentProfile := strutil.Fallback(config.Profile, "default")
			profileSeen := false
			profileNames := make([]string, 0, len(config.Profiles)+1)

			for profile := range config.Profiles {
				if currentProfile == profile {
					profileSeen = true
				}
				profileNames = append(profileNames, profile)
			}

			if !profileSeen {
				// profile does not actually exist but is infered
				profileNames = append(profileNames, currentProfile)
			}

			sort.Strings(profileNames)

			for _, profile := range profileNames {
				if profile == currentProfile {
					_, _ = fmt.Fprintf(ioStreams.Out, "* %s\n", profile)
					continue
				}
				_, _ = fmt.Fprintf(ioStreams.Out, "  %s\n", profile)
			}

			return nil
		},
	}
}
