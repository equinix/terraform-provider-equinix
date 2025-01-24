package metro

import (
	"context"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func NewDataSourceMetros() datasource.DataSource {
	return &DataSourceMetros{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_metros",
			},
		),
	}
}

type DataSourceMetros struct {
	framework.BaseDataSource
}

func (r *DataSourceMetros) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAllMetroSchema(ctx)
}

func (r *DataSourceMetros) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var allMetrosData DataSourceAllMetrosModel
	var pagination PaginationModel
	response.Diagnostics.Append(request.Config.Get(ctx, &allMetrosData)...)
	if response.Diagnostics.HasError() {
		return
	}

	diags := allMetrosData.Pagination.As(ctx, &pagination, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return
	}
	offset := pagination.Offset.ValueInt32()
	limit := pagination.Limit.ValueInt32()
	presence := allMetrosData.Presence.ValueString()
	if limit == 0 {
		limit = 20
	}

	metroRequest := client.MetrosApi.GetMetros(ctx).
		Limit(limit).
		Offset(offset)
	if presence != "" {
		metroRequest.Presence(fabricv4.Presence(presence))
	}
	metros, _, err := metroRequest.Execute()

	if len(metros.GetData()) < 1 {
		response.Diagnostics.AddError("no data retrieved by Get All Metros data source",
			"either the account does not have any metros data to pull or the combination of limit and offset needs to be updated")
		return
	}

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("Get All Metros API Error", equinix_errors.FormatFabricError(err).Error())
		diag.FromErr(equinix_errors.FormatFabricError(err))
		return
	}

	response.Diagnostics.Append(allMetrosData.parse(ctx, metros)...)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &allMetrosData)...)
}
