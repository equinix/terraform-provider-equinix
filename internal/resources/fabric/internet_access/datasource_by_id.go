package internetaccess

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// NewDataSourceByInternetAccessID creates a new data source for fetching a single Internet Access service by UUID
func NewDataSourceByInternetAccessID() datasource.DataSource {
	return &DataSourceByInternetAccessID{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_internet_access_service",
			},
		),
	}
}

// DataSourceByInternetAccessID datasource for fetching a single Internet Access service
type DataSourceByInternetAccessID struct {
	framework.BaseDataSource
}

// Schema returns the data source schema
func (r *DataSourceByInternetAccessID) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSingleInternetAccessSchema(ctx)
}

// Read fetches an Internet Access service by its UUID
func (r *DataSourceByInternetAccessID) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DataSourceByIDModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	service, _, err := client.InternetAccessServicesApi.GetEiaService(ctx, data.InternetAccessServiceID.ValueString()).Execute()
	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving internet access service", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, service)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
