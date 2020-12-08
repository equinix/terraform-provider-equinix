package main

import (
	"github.com/equinix/terraform-provider-metal/metal"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: metal.Provider})
}
