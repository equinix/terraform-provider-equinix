package stream

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewDataSourceByStreamID() datasource.DataSource {
	return &DataSourceByStreamID{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_stream",
			},
		),
	}
}

type DataSourceByStreamID struct {
	framework.BaseDataSource
}

func (r *DataSourceByStreamID) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSingleStreamSchema(ctx)
}

func (r *DataSourceByStreamID) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	// Retrieve values from plan
	var data DataSourceByIDModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	streamID := data.StreamID.ValueString()

	// Use API client to get the current state of the resource
	stream, _, err := client.StreamsApi.GetStreamByUuid(ctx, streamID).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving stream data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	response.Diagnostics.Append(data.parse(ctx, stream)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
