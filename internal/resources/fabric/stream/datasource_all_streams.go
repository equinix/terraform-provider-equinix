package stream

import (
	"context"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func NewDataSourceAllStreams() datasource.DataSource {
	return &DataSourceAllStreams{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_streams",
			},
		),
	}
}

type DataSourceAllStreams struct {
	framework.BaseDataSource
}

func (r *DataSourceAllStreams) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSingleStreamSchema(ctx)
}

func (r *DataSourceAllStreams) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	// Retrieve values from plan
	var data DataSourceAllStreamsModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	var pagination PaginationModel
	diags := data.Pagination.As(ctx, &pagination, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return
	}

	// Use API client to get the current state of the resource
	streams, _, err := client.StreamsApi.GetStreams(ctx).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		diag.FromErr(equinix_errors.FormatFabricError(err))
		return
	}

	// Set state to fully populated data
	response.Diagnostics.Append(data.parse(ctx, streams)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
