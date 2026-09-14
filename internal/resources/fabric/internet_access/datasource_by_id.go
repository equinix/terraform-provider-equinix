package internetaccess

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewDataSourceByInternetAccessServiceID() datasource.DataSource {
	return &DataSourceByID{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_internet_access_service",
			},
		),
	}
}

type DataSourceByID struct {
	framework.BaseDataSource
}

func (r *DataSourceByID) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSingleInternetAccessServiceSchema(ctx)
}

func (r *DataSourceByID) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DataSourceByIDModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	svc, err := getInternetAccessServiceRaw(ctx, client, data.InternetAccessServiceID.ValueString())
	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving internet access service", err.Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, svc)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
