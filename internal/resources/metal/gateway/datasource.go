package gateway

import (
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

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
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	// Retrieve values from plan
	var data DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := data.GatewayID.ValueString()

	// API call to get the Metal Gateway
	includes := []string{"project", "ip_reservation", "virtual_network", "vrf"}
	gw, _, err := client.MetalGatewaysApi.FindMetalGatewayById(ctx, id).Include(includes).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Metal Gateway",
			"Could not read Metal Gateway with ID "+id+": "+err.Error(),
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
