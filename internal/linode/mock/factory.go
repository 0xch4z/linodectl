package mock

import (
	"github.com/0xch4z/linodectl/internal/cmd/util"
	"github.com/0xch4z/linodectl/internal/config"
	"github.com/0xch4z/linodectl/internal/linode"
	"go.uber.org/mock/gomock"
)

func ClientFactory(ctrl *gomock.Controller) (util.ClientFactoryFunc, *MockClient) {
	client := NewMockClient(ctrl)
	return func(p config.Profile) linode.Client {
		return client
	}, client
}
