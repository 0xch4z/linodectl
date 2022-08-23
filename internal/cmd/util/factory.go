package util

import (
	"errors"

	"github.com/0xch4z/linodectl/internal/config"
	"github.com/0xch4z/linodectl/internal/linode"
	"github.com/spf13/cobra"
)

var (
	profileNotExistErr = errors.New("profile does not exist")
)

type CommandFactoryFunc func(Factory, IOStreams) *cobra.Command

type Factory interface {
	ConfigProvider() config.Provider
	Client(profileName string) (linode.Client, error)
	Config() *config.Config
}

type ClientFactoryFunc func(config.Profile) linode.Client

func NoopClientFactory(config.Profile) linode.Client {
	return nil
}

type factory struct {
	config         *config.Config
	configProvider config.Provider

	clientFactory ClientFactoryFunc
}

func (f *factory) Config() *config.Config { return f.config }

func (f *factory) ConfigProvider() config.Provider { return f.configProvider }

func (f *factory) Client(profileName string) (linode.Client, error) {
	if profileName == "" {
		profileName = f.config.Profile
	}

	profile, ok := f.config.Profiles[profileName]
	if !ok {
		return nil, profileNotExistErr
	}
	return f.clientFactory(profile), nil
}

func NewFactory(configProvider config.Provider, config *config.Config, clientFactory ClientFactoryFunc) Factory {
	return &factory{
		clientFactory:  clientFactory,
		configProvider: configProvider,
		config:         config,
	}
}
