package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/packethost/terraform-provider-packet/packet"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: packet.Provider})
}
