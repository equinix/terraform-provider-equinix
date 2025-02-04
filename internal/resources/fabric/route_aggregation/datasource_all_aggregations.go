package route_aggregation

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewDataSourceAllRouteAggregation() datasource.DataSource {
	return &DataSourceAllRouteAggregations{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_route_aggregations",
			},
		),
	}
}

type DataSourceAllRouteAggregations struct {
	framework.BaseDataSource
}

func (r *DataSourceAllRouteAggregations) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAllRouteAggregationsSchema(ctx)
}

func (r *DataSourceAllRouteAggregations) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DatsSourceAllRouteAggregationsModel

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

	propertyValue := fabricv4.RouteFiltersSearchFilterItemProperty(tffilter.Property.ValueString())

	filterItem := fabricv4.RouteAggregationsSearchFilterItem{
		Property: &propertyValue,
	}

	if !tffilter.Operator.IsNull() && !tffilter.Operator.IsUnknown() {
		filterItem.Operator = tffilter.Operator.ValueStringPointer()
	}

	if len(values) > 0 {
		filterItem.Values = values
	}

	filter := fabricv4.RouteAggregationsSearchBaseFilter{
		And: []fabricv4.RouteAggregationsSearchFilterItem{filterItem},
	}

	var tfpagination PaginationModel
	diags = data.Pagination.As(ctx, &tfpagination, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return
	}
	offset := tfpagination.Offset.ValueInt32()
	limit := tfpagination.Limit.ValueInt32()
	if limit == 0 {
		limit = 20
	}

	pagination := fabricv4.Pagination{
		Offset: &offset,
		Limit:  limit,
	}

	var tfsort SortModel
	diags = data.Sort.As(ctx, &tfsort, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return
	}
	direction := tfsort.Direction.ValueString()
	property := tfsort.Property.ValueString()

	pValue := fabricv4.RouteAggregationSortItemProperty(property)
	dValue := fabricv4.SortItemDirection(direction)

	sort := fabricv4.RouteAggregationSortItem{
		Property:  &pValue,
		Direction: &dValue,
	}

	routeAggregationsSearch := fabricv4.RouteAggregationsSearchBase{
		Filter:     &filter,
		Pagination: &pagination,
		Sort:       []fabricv4.RouteAggregationSortItem{sort},
	}

	routeAggregations, _, err := client.RouteAggregationsApi.SearchRouteAggregations(ctx).RouteAggregationsSearchBase(routeAggregationsSearch).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving route aggregations data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, routeAggregations)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
