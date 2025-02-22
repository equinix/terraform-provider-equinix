package streamattachment

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewDataSourceByIDs() datasource.DataSource {
	return &DataSourceByIDs{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_stream_attachment",
			},
		),
	}
}

type DataSourceByIDs struct {
	framework.BaseDataSource
}

func (r *DataSourceByIDs) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceByIDsSchema(ctx)
}

func (r *DataSourceByIDs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	// Retrieve values from plan
	var data DataSourceByIDModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	assetID, asset, streamID := data.AssetID.ValueString(), data.Asset.ValueString(), data.StreamID.ValueString()

	attachment, _, err := client.StreamsApi.GetStreamAssetByUuid(ctx, assetID, fabricv4.Asset(asset), streamID).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving stream attachment", equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	response.Diagnostics.Append(data.parse(ctx, attachment)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
