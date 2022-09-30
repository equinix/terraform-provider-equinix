package main

import (
	"context"
	"flag"
	"log"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	opts := &plugin.ServeOpts{ProviderFunc: equinix.Provider}

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/equinix/equinix", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
