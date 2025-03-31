package precisiontime

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
)

func NewDataSourceByEptServiceID() datasource.DataSource {
	return &DataSourceByEptServiceID{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_precision_time_service",
			},
		),
	}
}

type DataSourceByEptServiceID struct {
	framework.BaseDataSource
}

func (r *DataSourceByEptServiceID) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSingleEptServiceSchema(ctx)
}

func (r *DataSourceByEptServiceID) Read(
	ctx context.Context,
	request datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data dataSourceByIDModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	eptServiceID := data.EptServiceID.ValueString()

	ept, _, err := client.PrecisionTimeApi.GetTimeServicesById(ctx, eptServiceID).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving ept service data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	response.Diagnostics.Append(data.parse(ctx, ept)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
