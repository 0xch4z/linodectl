package config

import (
	"bytes"

	"github.com/Charliekenney23/linodectl/internal/cli/editor"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func NewCmdConfigEdit(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	return &cobra.Command{
		Use: "edit",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := f.Config()
			configBytes, err := yaml.Marshal(conf)
			if err != nil {
				return err
			}

			e := editor.NewDefaultEditor()
			if configBytes, _, err = e.EditReader("", "config.yaml", ioStreams, bytes.NewBuffer(configBytes)); err != nil {
				return err
			}

			if err := yaml.Unmarshal(configBytes, &conf); err != nil {
				return err
			}
			return f.ConfigProvider().Save(conf)
		},
	}
}
