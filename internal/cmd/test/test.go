package test

import (
	"bytes"
	"testing"

	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/config"
	configmock "github.com/Charliekenney23/linodectl/internal/config/mock"
	clientmock "github.com/Charliekenney23/linodectl/internal/linode/mock"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
)

type Streams struct {
	Out *bytes.Buffer
}

type Invocation struct {
	Client         *clientmock.MockClient
	Config         *config.Config
	ConfigProvider *configmock.MockProvider
	Streams        *Streams

	Mock *gomock.Controller
}

func Command(t *testing.T, cmdFactory cmdutil.CommandFactoryFunc) (*cobra.Command, *Invocation) {
	t.Helper()

	ctrl := gomock.NewController(t)
	clientFactory, client := clientmock.ClientFactory(ctrl)

	conf := config.DefaultConfig()
	configProvider := configmock.NewMockProvider(ctrl)
	factory := cmdutil.NewFactory(configProvider, conf, clientFactory)
	streams := &Streams{Out: bytes.NewBuffer(nil)}

	invocation := &Invocation{
		Config:         conf,
		Client:         client,
		ConfigProvider: configProvider,
		Streams:        streams,
	}

	return cmdFactory(factory, cmdutil.IOStreams{Out: streams.Out}), invocation
}
