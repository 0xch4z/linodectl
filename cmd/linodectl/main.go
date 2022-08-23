package main

import (
	"log"
	"os"

	"github.com/0xch4z/linodectl/internal/cmd"
	cmdutil "github.com/0xch4z/linodectl/internal/cmd/util"
	"github.com/0xch4z/linodectl/internal/config"
	"github.com/0xch4z/linodectl/internal/linode"
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
