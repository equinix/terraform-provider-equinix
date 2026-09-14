package ipblock

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// NewDataSourceByIpBlockID creates a new data source for fetching a single IP block by UUID
func NewDataSourceByIpBlockID() datasource.DataSource {
	return &DataSourceByIpBlockID{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_ip_block",
			},
		),
	}
}

// DataSourceByIpBlockID datasource for fetching a single IP block
type DataSourceByIpBlockID struct {
	framework.BaseDataSource
}

// Schema returns the data source schema
func (r *DataSourceByIpBlockID) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSingleIpBlockSchema(ctx)
}

// Read fetches an IP block by its UUID
func (r *DataSourceByIpBlockID) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DataSourceByIDModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	ipBlock, _, err := client.IPBlocksApi.GetIpBlock(ctx, data.IpBlockID.ValueString()).Execute()
	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving ip block", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, ipBlock)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
