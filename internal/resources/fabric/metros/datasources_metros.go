package metros

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

	var allMetrosData AllMetrosModel
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
	metros, _, err := client.MetrosApi.GetMetros(ctx).
		Limit(limit).
		Offset(offset).
		Presence(fabricv4.Presence(presence)).
		Execute()
	if err != nil {
		response.State.RemoveResource(ctx)
		diag.FromErr(equinix_errors.FormatFabricError(err))
		return
	}

	response.Diagnostics.Append(allMetrosData.parseDataSourceAllMetros(ctx, metros)...)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &allMetrosData)...)
}
