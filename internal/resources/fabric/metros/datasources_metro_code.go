package metros

import (
	"context"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSingleMetroSchema(ctx)
}

// READ function for GET Metro Code data source
func (r *DataSourceMetroCode) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data MetroModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	metroCode := data.MetroCode.ValueString()

	metroByCode, _, err := client.MetrosApi.GetMetroByCode(ctx, metroCode).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		diag.FromErr(equinix_errors.FormatFabricError(err))
		return
	}

	response.Diagnostics.Append(data.parseDataSourceByMetroCode(ctx, metroByCode)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
