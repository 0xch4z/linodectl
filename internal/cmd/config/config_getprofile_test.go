package config

import (
	"testing"

	"github.com/0xch4z/linodectl/internal/cmd/test"
	"github.com/0xch4z/linodectl/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	t.Run("outputs profile", func(t *testing.T) {
		cmd, i := test.Command(t, NewCmdConfig)
		i.Config.Profiles["default"] = config.Profile{
			APIVersion: "v4beta",
			Token:      "bogus",
			Region:     "us-east",
		}

		cmd.SetArgs([]string{"get-profile"})
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, i.Streams.Out.String(), `name: default
apiVersion: v4beta
apiBaseURL: $LINODE_URL
token: *****
region: us-east
`)
	})
}
