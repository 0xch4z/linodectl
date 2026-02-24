package instance

import (
	"context"
	"net/http"
	"testing"

	"github.com/0xch4z/linodectl/internal/linode/mock"
	"go.uber.org/mock/gomock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestFilterLKECluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	defer ctrl.Finish()

	t.Run("filters for associated LKE Cluster nodes", func(t *testing.T) {
		clusterID := 39291
		client.EXPECT().ListLKEClusterPools(gomock.Any(), clusterID, gomock.Any()).
			Times(1).
			Return([]linodego.LKEClusterPool{
				{ID: 2039, Linodes: []linodego.LKEClusterPoolLinode{
					{InstanceID: 3029}, {InstanceID: 2919}, {InstanceID: 2921},
				}},
				{ID: 2033, Linodes: []linodego.LKEClusterPoolLinode{
					{InstanceID: 3219}, {InstanceID: 4838},
				}},
			}, nil)

		instances, err := FilterLKECluster(context.TODO(), client, clusterID, []linodego.Instance{
			{ID: 2029}, {ID: 2031}, {ID: 2919}, {ID: 2921}, {ID: 3291}, {ID: 512},
		})
		assert.NoError(t, err)
		assert.Equal(t, []linodego.Instance{{ID: 2919}, {ID: 2921}}, instances)
	})

	t.Run("throws error on generic API error", func(t *testing.T) {
		linodeErr := linodego.Error{Code: http.StatusBadGateway}
		client.EXPECT().ListLKEClusterPools(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, linodeErr)

		instances, err := FilterLKECluster(context.TODO(), client, 245, []linodego.Instance{{ID: 123}})
		assert.ErrorIs(t, err, linodeErr)
		assert.Nil(t, instances)
	})
}
