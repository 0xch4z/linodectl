package config

import (
	"github.com/0xch4z/linodectl/internal/strutil"
)

type SSHConfig struct {
	User    string `yaml:"user,omitempty"`
	KeyPath string `yaml:"key_path,omitempty"`
}

type Profile struct {
	APIVersion string `yaml:"api_version,omitempty"`
	APIBaseURL string `yaml:"api_base_url,omitempty"`
	Region     string `yaml:"region,omitempty"`
	Token      string `yaml:"token,omitempty"`
}

type Config struct {
	Profile string `yaml:"profile,omitempty"`

	Profiles  map[string]Profile `yaml:"profiles,omitempty"`
	Instances InstanceConfig     `yaml:"instances,omitempty"`
	SSH       SSHConfig          `yaml:"ssh,omitempty"`
}

func (c Config) CurrentProfile() (*Profile, bool) {
	profileName := strutil.Fallback(c.Profile, "default")
	profile, ok := c.Profiles[profileName]
	return &profile, ok
}
