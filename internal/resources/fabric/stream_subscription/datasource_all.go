package streamsubscription

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewDataSourceAllStreamSubscriptions() datasource.DataSource {
	return &DataSourceAllStreams{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_stream_subscriptions",
			},
		),
	}
}

type DataSourceAllStreams struct {
	framework.BaseDataSource
}

func (r *DataSourceAllStreams) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAllStreamSubscriptionsSchema(ctx)
}

func (r *DataSourceAllStreams) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	// Retrieve values from plan
	var data dataSourceAll
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	var pagination paginationModel
	diags := data.Pagination.As(ctx, &pagination, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return
	}
	offset := pagination.Offset.ValueInt32()
	limit := pagination.Limit.ValueInt32()
	if limit == 0 {
		limit = 20
	}

	// Use API client to get the current state of the resource
	streamSubscriptions, _, err := client.StreamSubscriptionsApi.GetStreamSubscriptions(ctx, data.StreamID.ValueString()).Limit(limit).Offset(offset).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving stream subscriptions data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	response.Diagnostics.Append(data.parse(ctx, streamSubscriptions)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
