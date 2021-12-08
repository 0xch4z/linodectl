# linodectl

![linodectl edit demo](/assets/demo.gif)

A command line interface to manage your Linode infrastructure.

## Usage

```
❯ linodectl
linodectl manages Linode resources

Usage:
  linodectl [flags]
  linodectl [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  config      Manage configurations
  create      Create a Linode resource
  delete      Delete a resource
  edit        Edit a resource
  get         Get a resource
  help        Help about any command
  whoami      Introspect

Flags:
  -h, --help             help for linodectl
  -p, --profile string   The profile to use for communicating with the Linode API

Use "linodectl [command] --help" for more information about a command.
```

See documentation [here](./docs/linodectl.md)
