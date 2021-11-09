package config

import (
	"bytes"
	"testing"

	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestListProfiles(t *testing.T) {
	t.Run("lists all profiles", func(t *testing.T) {
		f := cmdutil.NewFactory(nil, &config.Config{
			Profile: "myprofile",
			Profiles: map[string]config.Profile{
				"default":        {},
				"myotherprofile": {},
				"myprofile":      {},
			},
		}, cmdutil.NoopClientFactory)
		out := bytes.NewBuffer(nil)
		streams := cmdutil.IOStreams{Out: out}

		cmd := NewCmdConfigListProfiles(f, streams)
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, out.String(), `  default
  myotherprofile
* myprofile
`)
	})

	t.Run("shows default config when none specified", func(t *testing.T) {
		f := cmdutil.NewFactory(nil, &config.Config{
			Profiles: map[string]config.Profile{
				"not-used":      {},
				"also-not-used": {},
			},
		}, cmdutil.NoopClientFactory)
		out := bytes.NewBuffer(nil)
		streams := cmdutil.IOStreams{Out: out}

		cmd := NewCmdConfigListProfiles(f, streams)
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, out.String(), `  also-not-used
* default
  not-used
`)
	})

	t.Run("shows default config when config is empty", func(t *testing.T) {
		f := cmdutil.NewFactory(nil, &config.Config{}, cmdutil.NoopClientFactory)
		out := bytes.NewBuffer(nil)
		streams := cmdutil.IOStreams{Out: out}

		cmd := NewCmdConfigListProfiles(f, streams)
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, out.String(), `* default
`)
	})
}
