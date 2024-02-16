package main

import (
	"context"
	"flag"
	"log"

	sdkprovider "github.com/equinix/terraform-provider-equinix/equinix"
	frameworkprovider "github.com/equinix/terraform-provider-equinix/internal/provider"
	"github.com/equinix/terraform-provider-equinix/version"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {

	ctx := context.Background()

	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	upgradedSdkProvider, err := tf5to6server.UpgradeServer(
		context.Background(),
		sdkprovider.Provider().GRPCProvider,
	)

	providers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedSdkProvider
		},

		// Example terraform-plugin-framework provider
		providerserver.NewProtocol6(frameworkprovider.CreateFrameworkProvider(version.ProviderVersion)),
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
