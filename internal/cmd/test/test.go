package test

import (
	"bytes"
	"context"
	"testing"

	cmdutil "github.com/0xch4z/linodectl/internal/cmd/util"
	"github.com/0xch4z/linodectl/internal/config"
	configmock "github.com/0xch4z/linodectl/internal/config/mock"
	clientmock "github.com/0xch4z/linodectl/internal/linode/mock"
	"github.com/0xch4z/linodectl/internal/printer"
	"github.com/0xch4z/linodectl/internal/resource"
	"go.uber.org/mock/gomock"
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

func MockPrintResources(resources resource.List, options printer.ResourcePrintOptions) string {
	printerBuf := bytes.NewBuffer(nil)
	p := printer.New(printerBuf)
	_ = p.PrintResources(context.TODO(), resources, options)
	return printerBuf.String()
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
