package acceptance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6to5server"
)

var ProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
	"equinix": func() (tfprotov5.ProviderServer, error) {
		ctx := context.Background()

		frameworkServer := providerserver.NewProtocol6(TestAccFrameworkProvider)

		providers := []func() tfprotov5.ProviderServer{
			func() tfprotov5.ProviderServer {
				downgradedServer, err := tf6to5server.DowngradeServer(ctx, frameworkServer)

				if err != nil {
					panic(err)
				}
				return downgradedServer
			},
			TestAccProviders["equinix"].GRPCProvider,
		}

		muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	},
}
