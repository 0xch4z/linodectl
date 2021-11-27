package config

import (
	"errors"

	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/spf13/cobra"
)

const setProfileExamples = `  # Set profile to "my-profile"
  linodectl config set-profile my-profile`

var profileNotExistErr = errors.New("profile does not exist")

func NewCmdConfigSetProfile(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-profile NAME",
		Short:   "Set default profile",
		Example: setProfileExamples,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := f.Config()

			profileName := args[0]
			if profileName == config.Profile {
				// this is already the configured profile
				return nil
			}

			if _, ok := config.Profiles[profileName]; !ok {
				return profileNotExistErr
			}

			// save the specified profile
			config.Profile = profileName
			return f.ConfigProvider().Save(config)
		},
	}
	cmd.SetErr(ioStreams.Err)

	return cmd
}
