package config

import (
	"bytes"
	"testing"

	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	t.Run("outputs profile", func(t *testing.T) {
		f := cmdutil.NewFactory(nil, &config.Config{
			Profiles: map[string]config.Profile{
				"default": {
					APIVersion: "v4beta",
					Token:      "bogus",
					Region:     "us-east",
				},
			},
		}, cmdutil.NoopClientFactory)
		out := bytes.NewBuffer(nil)
		streams := cmdutil.IOStreams{Out: out}

		cmd := NewCmdConfig(f, streams)
		cmd.SetArgs([]string{"get-profile"})
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, out.String(), `name: default
apiVersion: v4beta
apiBaseURL: $LINODE_URL
token: *****
region: us-east
`)
	})

	t.Run("outputs nothing when no profile", func(t *testing.T) {
		f := cmdutil.NewFactory(nil, &config.Config{}, cmdutil.NoopClientFactory)
		out := bytes.NewBuffer(nil)
		streams := cmdutil.IOStreams{Out: out}

		cmd := NewCmdConfig(f, streams)
		cmd.SetArgs([]string{"get-profile"})
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, out.String(), "")
	})
}
