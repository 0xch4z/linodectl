package whoami

import (
	"testing"

	"github.com/0xch4z/linodectl/internal/cmd/test"
	"github.com/golang/mock/gomock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestWhoami(t *testing.T) {
	t.Run("gets username", func(t *testing.T) {
		cmd, i := test.Command(t, NewCmdWhoami)
		username := "someonesname"
		i.Client.EXPECT().GetProfile(gomock.Any()).Times(1).Return(&linodego.Profile{Username: username}, nil)
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, i.Streams.Out.String(), username+"\n")
	})
}
