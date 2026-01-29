// Package advertised_route implements datasource for advertised route
package advertised_route

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

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

	var data dataSourceSearchAdvertisedRoutesModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	var tffilter FilterModel

	diags := data.Filter.As(ctx, &tffilter, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return
	}
	values := []string{}
	if len(tffilter.Values) > 0 {
		for _, strVal := range tffilter.Values {
			if !strVal.IsNull() && !strVal.IsUnknown() {
				values = append(values, strVal.ValueString())
			}
		}
	}

	// propertyValue := fabricv4.RouteFiltersSearchFilterItemProperty(tffilter.Property.ValueString()) ////
	propertyValue := tffilter.Property.ValueString()

	filterItem := fabricv4.ConnectionRouteEntrySimpleExpression{
		Property: &propertyValue, ///////
	}

	if !tffilter.Operator.IsNull() && !tffilter.Operator.IsUnknown() {
		filterItem.Operator = tffilter.Operator.ValueStringPointer()
	}

	if len(values) > 0 {
		filterItem.Values = values
	}
	filterEntry := fabricv4.ConnectionRouteEntryFilter{
		ConnectionRouteEntrySimpleExpression: &filterItem,
	}

	filter := fabricv4.ConnectionRouteEntryFilters{
		And: []fabricv4.ConnectionRouteEntryFilter{
			filterEntry,
		},
	}

	var tfpagination paginationModel
	diags = data.Pagination.As(ctx, &tfpagination, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return
	}
	offset := tfpagination.Offset.ValueInt32()
	limit := tfpagination.Limit.ValueInt32()
	if limit == 0 {
		limit = 20
	}

	pagination := fabricv4.PaginationRequest{
		Offset: &offset,
		Limit:  &limit,
	}

	var tfsort sortModel
	diags = data.Sort.As(ctx, &tfsort, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return
	}
	direction := tfsort.Direction.ValueString()
	property := tfsort.Property.ValueString()

	pValue := fabricv4.ConnectionRouteEntrySortBy(property)
	dValue := fabricv4.ConnectionRouteEntrySortDirection(direction)

	sort := fabricv4.ConnectionRouteSortCriteria{
		Property:  &pValue,
		Direction: &dValue,
	}

	advertisedRoutesSearch := fabricv4.ConnectionRouteSearchRequest{
		Filter:     &filter,
		Pagination: &pagination,
		Sort:       []fabricv4.ConnectionRouteSortCriteria{sort},
	}
	connectionID := data.ConnectionID.ValueString()
	advertisedRoutes, _, err := client.CloudRoutersApi.SearchConnectionAdvertisedRoutes(ctx, connectionID).ConnectionRouteSearchRequest(advertisedRoutesSearch).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving advertised routes data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, advertisedRoutes)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
