package gateway

import (
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/packethost/packngo"

	"context"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name:   "equinix_metal_gateway",
				Schema: &dataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	framework.BaseDataSource
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Retrieve values from plan
	var data DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

    // Extract the ID of the resource from the state
	id := data.GatewayID.ValueString()

    // API call to get the Metal Gateway
    includes := &packngo.GetOptions{Includes: []string{"project", "ip_reservation", "virtual_network", "vrf"}}
    gw, _, err := client.MetalGateways.Get(id, includes)
    if err != nil {
        err = equinix_errors.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error reading Metal Gateway",
            "Could not read Metal Gateway with ID " + id + ": " + err.Error(),
        )
        return
    }

	// Set state to fully populated data
	resp.Diagnostics.Append(data.parse(gw)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
