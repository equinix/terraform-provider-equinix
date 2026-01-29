// Package advertised_route implements datasource for advertised route
package advertised_route

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// NewDataSourceAdvertisedRoutes creates a new data source for Advertised Routes
func NewDataSourceAdvertisedRoutes() datasource.DataSource {
	return &DataSourceAllAdvertisedRoutes{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_advertised_routes",
			},
		),
	}
}

// DataSourceAllAdvertisedRoutes datasource represents advertised routes
type DataSourceAllAdvertisedRoutes struct {
	framework.BaseDataSource
}

// Schema returns the advertised routes datasource schema
func (r *DataSourceAllAdvertisedRoutes) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAdvertisedRoutesSchema(ctx)
}

func (r *DataSourceAllAdvertisedRoutes) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var searchAdvertisedRoutesData dataSourceSearchAdvertisedRoutesModel
	var pagination paginationModel
	response.Diagnostics.Append(request.Config.Get(ctx, &searchAdvertisedRoutesData)...)
	if response.Diagnostics.HasError() {
		return
	}

	diags := searchAdvertisedRoutesData.Pagination.As(ctx, &pagination, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return
	}
	
	connectionID := searchAdvertisedRoutesData.ConnectionID.ValueString()
	advertisedRoutes, _, err := client.CloudRoutersApi.SearchConnectionAdvertisedRoutes(ctx, connectionID).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving advertised routes data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(searchAdvertisedRoutesData.parse(ctx, advertisedRoutes)...)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &searchAdvertisedRoutesData)...)
}
