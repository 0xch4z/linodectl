package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/Charliekenney23/linodectl/internal/cmd"
	"github.com/Charliekenney23/linodectl/internal/cmd/util"
	"github.com/spf13/cobra/doc"
)

func main() {
	rootcmd := cmd.NewRootCommand(util.NewFactory(nil, nil, nil), os.Stdin, ioutil.Discard, ioutil.Discard)
	if err := doc.GenMarkdownTree(rootcmd, "./docs"); err != nil {
		log.Fatal(err)
	}
}
