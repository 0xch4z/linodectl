package mock

import (
	"github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/config"
	"github.com/Charliekenney23/linodectl/internal/linode"
	"github.com/golang/mock/gomock"
)

func ClientFactory(ctrl *gomock.Controller) (util.ClientFactoryFunc, *MockClient) {
	client := NewMockClient(ctrl)
	return func(p config.Profile) linode.Client {
		return client
	}, client
}
