package acceptance

import (
	"context"
	"log"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/provider"
	"github.com/equinix/terraform-provider-equinix/version"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"equinix": func() (tfprotov6.ProviderServer, error) {
		ctx := context.Background()

		sdkv2Provider, err := tf5to6server.UpgradeServer(ctx, equinix.Provider().GRPCProvider)
		if err != nil {
			log.Fatal(err)
		}

		sdkv2ProviderFunc := func() tfprotov6.ProviderServer { return sdkv2Provider }
		frameworkProvider := providerserver.NewProtocol6(
			provider.CreateFrameworkProvider(version.ProviderVersion))

		providers := []func() tfprotov6.ProviderServer{
			sdkv2ProviderFunc,
			frameworkProvider,
		}

		muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	},
}
