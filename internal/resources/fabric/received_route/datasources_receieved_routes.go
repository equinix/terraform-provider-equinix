// Package advertised_route implements datasource for advertised route
package advertised_route

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// NewDataSourceReceivedRoutes creates a new data source for Received Routes
func NewDataSourceReceivedRoutes() datasource.DataSource {
	return &DataSourceAllReceivedRoutes{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_received_routes",
			},
		),
	}
}

// DataSourceAllReceivedRoutes datasource represents Received routes
type DataSourceAllReceivedRoutes struct {
	framework.BaseDataSource
}

// Schema returns the Received routes datasource schema
func (r *DataSourceAllReceivedRoutes) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceReceivedRoutesSchema(ctx)
}

func (r *DataSourceAllReceivedRoutes) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var searchReceivedRoutesData dataSourceSearchReceivedRoutesModel
	var pagination paginationModel
	response.Diagnostics.Append(request.Config.Get(ctx, &searchReceivedRoutesData)...)
	if response.Diagnostics.HasError() {
		return
	}

	diags := searchReceivedRoutesData.Pagination.As(ctx, &pagination, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return
	}
	
	connectionID := searchReceivedRoutesData.ConnectionID.ValueString()
	receivedRoutes, _, err:= client.CloudRoutersApi.SearchConnectionReceivedRoutes(ctx, connectionID).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving Received routes data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(searchReceivedRoutesData.parse(ctx, receivedRoutes)...)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &searchReceivedRoutesData)...)
}
