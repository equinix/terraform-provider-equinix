package routeaggregation

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
	_ datasource.SchemaRequest,
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
	if diags := data.Filter.As(ctx, &tffilter, basetypes.ObjectAsOptions{}); diags.HasError() {
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

	filterItem := fabricv4.SearchSimpleExpression{
		Property: tffilter.Property.ValueString(),
	}

	filterItem.Operator = fabricv4.SearchSimpleExpressionOperator(tffilter.Operator.ValueString())

	if len(values) > 0 {
		filterItem.Values = values
	}

	filter := fabricv4.SearchFilter{
		SearchAndExpression: &fabricv4.SearchAndExpression{
			And: []fabricv4.SearchFilterExpression{
				{SearchSimpleExpression: &filterItem},
			},
		},
	}

	var tfpagination PaginationModel
	if diags := data.Pagination.As(ctx, &tfpagination, basetypes.ObjectAsOptions{}); diags.HasError() {
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

	var tfsort SortModel
	if diags := data.Sort.As(ctx, &tfsort, basetypes.ObjectAsOptions{}); diags.HasError() {
		return
	}
	direction := tfsort.Direction.ValueString()
	property := tfsort.Property.ValueString()

	pValue := fabricv4.RouteAggregationSortBy(property)
	dValue := fabricv4.RouteAggregationSortDirection(direction)

	routeAggregationsSearch := fabricv4.RouteAggregationsSearchRequest{
		Filter:     &filter,
		Pagination: &pagination,
		Sort: []fabricv4.RouteAggregationSortCriteria{
			{
				Property:  &pValue,
				Direction: &dValue,
			},
		},
	}

	routeAggregations, _, err := client.RouteAggregationsApi.SearchRouteAggregations(ctx).RouteAggregationsSearchRequest(routeAggregationsSearch).Execute()

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
