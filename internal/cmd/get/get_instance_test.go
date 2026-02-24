package get

import (
	"strconv"
	"testing"

	"github.com/0xch4z/linodectl/internal/cmd/test"
	"github.com/0xch4z/linodectl/internal/printer"
	"github.com/0xch4z/linodectl/internal/resource/instance"
	"go.uber.org/mock/gomock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestGetInstance(t *testing.T) {
	t.Run("gets instance by label", func(t *testing.T) {
		cmd, i := test.Command(t, NewCmdGetInstance)
		linodeInstance := linodego.Instance{ID: 123, Label: "my-instance"}
		i.Client.EXPECT().ListInstances(gomock.Any(), gomock.Any()).Times(1).
			Return([]linodego.Instance{
				linodeInstance,
				{ID: 124, Label: "my-other-instance"},
			}, nil)

		expectedOut := test.MockPrintResources(instance.NewList([]linodego.Instance{
			linodeInstance,
		}), printer.ResourcePrintOptions{})

		cmd.SetArgs([]string{linodeInstance.Label})
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, i.Streams.Out.String(), expectedOut)
	})

	t.Run("get instance by id", func(t *testing.T) {
		cmd, i := test.Command(t, NewCmdGetInstance)
		linodeInstance := linodego.Instance{ID: 38818, Label: "staging jenkins"}
		i.Client.EXPECT().ListInstances(gomock.Any(), gomock.Any()).Times(1).
			Return([]linodego.Instance{
				{ID: 38928, Label: "dev-jenkins"},
				{ID: 28912, Label: "prod-jenkins"},
				linodeInstance,
			}, nil)

		expectedOut := test.MockPrintResources(instance.NewList([]linodego.Instance{
			linodeInstance,
		}), printer.ResourcePrintOptions{})

		cmd.SetArgs([]string{strconv.Itoa(linodeInstance.ID)})
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, i.Streams.Out.String(), expectedOut)
	})

}
