package config

import (
	"fmt"

	cmdutil "github.com/0xch4z/linodectl/internal/cmd/util"
	"github.com/0xch4z/linodectl/internal/strutil"
	"github.com/spf13/cobra"
)

const notConfigured = "[not configured]"

func wrapProfileValue(s string) string {
	return strutil.Fallback(s, notConfigured)
}

func NewCmdConfigGetProfile(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	return &cobra.Command{
		Use:   "get-profile",
		Short: "Get current profile",
		RunE: func(_ *cobra.Command, _ []string) error {
			config := f.Config()

			profile, _ := config.CurrentProfile()
			_, _ = fmt.Fprintf(ioStreams.Out, "name: %s\n", strutil.Fallback(config.Profile, "default"))
			_, _ = fmt.Fprintf(ioStreams.Out, "apiVersion: %s\n", strutil.Fallback(profile.APIVersion, "$LINODE_API_VERSION"))
			_, _ = fmt.Fprintf(ioStreams.Out, "apiBaseURL: %s\n", strutil.Fallback(profile.APIBaseURL, "$LINODE_URL"))
			_, _ = fmt.Fprintf(ioStreams.Out, "token: %s\n", strutil.Fallback(strutil.Mask(profile.Token, '*'), "$LINODE_API_TOKEN"))
			_, _ = fmt.Fprintf(ioStreams.Out, "region: %s\n", wrapProfileValue(profile.Region))

			return nil
		},
	}
}
