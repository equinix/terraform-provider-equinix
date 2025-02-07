package metro

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewDataSourceMetroCode() datasource.DataSource {
	return &DataSourceMetroCode{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_metro",
			},
		),
	}
}

type DataSourceMetroCode struct {
	framework.BaseDataSource
}

func (r *DataSourceMetroCode) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSingleMetroSchema(ctx)
}

func (r *DataSourceMetroCode) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DataSourceByCodeModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	metroCode := data.MetroCode.ValueString()
	metroByCode, _, err := client.MetrosApi.GetMetroByCode(ctx, metroCode).Execute()
	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("Get By Metro Code API Error", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, metroByCode)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
