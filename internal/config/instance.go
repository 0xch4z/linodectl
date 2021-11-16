package config

import "github.com/linode/linodego"

type InstanceConfig struct {
	Presets map[string]InstancePreset `yaml:"presets,omitempty"`

	AutoAuthorizeMe bool `yaml:"auto_authorize_me,omitempty"`
}

type InstanceInterface struct {
	IPAMAddress string                          `yaml:"ipam_address,omitempty"`
	Label       string                          `yaml:"label,omitempty"`
	Purpose     linodego.ConfigInterfacePurpose `yaml:"purpose,omitempty"`
}

type InstancePreset struct {
	AuthorizedUsers []string                       `yaml:"autorized_users,omitempty"`
	Group           string                         `yaml:"group,omitempty"`
	Image           string                         `yaml:"image,omitempty"`
	Interfaces      []map[string]InstanceInterface `yaml:"interfaces,omitempty"`
	Region          string                         `yaml:"region,omitempty"`
	RootPass        string                         `yaml:"root_pass,omitempty"`
	StackscriptData map[string]interface{}         `yaml:"stackscript_data,omitempty"`
	StackscriptID   int                            `yaml:"stackscript_id,omitempty"`
	SwapSize        int                            `yaml:"swap_size,omitempty"`
	Tags            []string                       `yaml:"tags,omitempty"`
	Type            string                         `yaml:"type,omitempty"`
}
