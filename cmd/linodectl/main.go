package main

import (
	"log"
	"os"

	"github.com/Charliekenney23/linodectl/internal/cmd"
	cmdutil "github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/Charliekenney23/linodectl/internal/config"
	"github.com/Charliekenney23/linodectl/internal/linode"
)

func main() {
	configProvider := config.NewProvider()
	conf, err := configProvider.Load()
	if err != nil {
		log.Fatal(err)
	}

	f := cmdutil.NewFactory(configProvider, conf, linode.NewClient)
	if err := cmd.NewRootCommand(f, os.Stdin, os.Stdout, os.Stderr).Execute(); err != nil {
		log.Fatal(err)
	}
}
