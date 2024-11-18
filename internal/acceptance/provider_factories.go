package acceptance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"equinix": func() (tfprotov6.ProviderServer, error) {
		ctx := context.Background()

		sdkv2Provider, err := tf5to6server.UpgradeServer(ctx, TestAccProvider.GRPCProvider)
		if err != nil {
			return nil, err
		}

		providers := []func() tfprotov6.ProviderServer{
			func() tfprotov6.ProviderServer { return sdkv2Provider },
			providerserver.NewProtocol6(TestAccFrameworkProvider),
		}

		muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	},
}
