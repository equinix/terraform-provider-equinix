package main

import (
	"context"
	"flag"
	"log"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/provider"
	"github.com/equinix/terraform-provider-equinix/version"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs --rendered-provider-name=Equinix
func main() {

	ctx := context.Background()

	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

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
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debugMode {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/equinix/equinix",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
