package config

import (
	"bytes"
	"testing"

	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/config"
	configmock "github.com/Charliekenney23/linodectl/internal/config/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("updates the profile correctly", func(t *testing.T) {
		p := configmock.NewMockProvider(ctrl)

		profiles := map[string]config.Profile{
			"default":        {},
			"myotherprofile": {},
			"myprofile":      {},
		}
		f := cmdutil.NewFactory(p, &config.Config{
			Profile:  "myprofile",
			Profiles: profiles,
		}, cmdutil.NoopClientFactory)
		out := bytes.NewBuffer(nil)
		streams := cmdutil.IOStreams{Out: out}

		p.EXPECT().Save(&config.Config{
			Profile:  "myotherprofile",
			Profiles: profiles,
		}).Times(1)

		cmd := NewCmdConfigSetProfile(f, streams)
		cmd.SetArgs([]string{"myotherprofile"})
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, out.String(), "")
	})

	t.Run("skips update when profile already set", func(t *testing.T) {
		p := configmock.NewMockProvider(ctrl)

		f := cmdutil.NewFactory(p, &config.Config{
			Profile:  "mine",
			Profiles: map[string]config.Profile{"default": {}, "mine": {}},
		}, cmdutil.NoopClientFactory)
		out := bytes.NewBuffer(nil)
		streams := cmdutil.IOStreams{Out: out}

		p.EXPECT().Save(gomock.Any()).Times(0)

		cmd := NewCmdConfigSetProfile(f, streams)
		cmd.SetArgs([]string{"mine"})
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, out.String(), "")
	})

	t.Run("fails to update the profile when it doesn't exist", func(t *testing.T) {
		p := configmock.NewMockProvider(ctrl)

		f := cmdutil.NewFactory(p, &config.Config{
			Profile:  "default",
			Profiles: map[string]config.Profile{"notit": {}},
		}, cmdutil.NoopClientFactory)
		out := bytes.NewBuffer(nil)
		err := bytes.NewBuffer(nil)
		streams := cmdutil.IOStreams{Out: out, Err: err}

		cmd := NewCmdConfigSetProfile(f, streams)
		cmd.SetArgs([]string{"myotherprofile"})
		assert.ErrorIs(t, cmd.Execute(), profileNotExistErr)
		assert.Equal(t, out.String(), "")
		assert.NotEmpty(t, err.String())
	})

}
