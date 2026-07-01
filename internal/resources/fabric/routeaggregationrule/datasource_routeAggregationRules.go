package routeaggregationrule

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/equinix/terraform-provider-equinix/internal/slice"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewDataSourceSearchRouteAggregationRules() datasource.DataSource {
	return &DataSourceSearchRouteAggregationRules{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_route_aggregation_rules",
			},
		),
	}
}

type DataSourceSearchRouteAggregationRules struct {
	framework.BaseDataSource
}

func (r *DataSourceSearchRouteAggregationRules) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSearchRouteAggregationRulesSchema(ctx)
}

func (r *DataSourceSearchRouteAggregationRules) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)
	params := fabricv4.NewRouteAggregationRulesSearchRequest()

	var data dataSourceRouteAggregationRulesModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	var sort []sortModel
	if !data.Sort.IsNull() && !data.Sort.IsUnknown() {
		if diags := data.Sort.ElementsAs(ctx, &sort, false); diags.HasError() {
			response.Diagnostics.Append(diags...)
		}
	}

	params.SetSort(slice.Map(sort,
		func(s sortModel) fabricv4.RouteAggregationRuleSortCriteria {
			sortCriteria := fabricv4.NewRouteAggregationRuleSortCriteria()
			if !s.Direction.IsNull() {
				sortCriteria.SetDirection(
					fabricv4.RouteAggregationRuleSortDirection(s.Direction.ValueString()))
			}

			if !s.Property.IsNull() {
				sortCriteria.SetProperty(fabricv4.RouteAggregationRuleSortBy(s.Property.ValueString()))
			}

			return *sortCriteria
		}))

	var filter []filterModel
	if !data.Filter.IsNull() && !data.Filter.IsUnknown() {
		if diags := data.Filter.ElementsAs(ctx, &filter, false); diags.HasError() {
			response.Diagnostics.Append(diags...)
		}
	}

	filterParam := fabricv4.RouteAggregationRulesFilter{}

	toExp := func(f filterModel) fabricv4.RouteAggregationRuleExpression {
		return fabricv4.RouteAggregationRuleExpression{
			RouteAggregationRuleSimpleExpression: func() *fabricv4.RouteAggregationRuleSimpleExpression {
				exp := fabricv4.NewRouteAggregationRuleSimpleExpression()
				if !f.Property.IsNull() {
					exp.SetProperty(f.Property.ValueString())
				}
				if !f.Operator.IsNull() {
					exp.SetOperator(f.Operator.ValueString())
				}
				var values []string
				if diags := f.Values.ElementsAs(ctx, &values, false); diags.HasError() {
					response.Diagnostics.Append(diags...)
				}

				exp.SetValues(values)
				return exp
			}(),
		}
	}

	switch data.OuterOperator.ValueString() {
	case "AND":
		filterParam.RouteAggregationRuleAndExpression = &fabricv4.RouteAggregationRuleAndExpression{And: slice.Map(filter, toExp)}
	case "OR":
		filterParam.RouteAggregationRuleOrExpression = &fabricv4.RouteAggregationRuleOrExpression{Or: slice.Map(filter, toExp)}
	}

	params.SetFilter(filterParam)

	var tfpagination paginationModel
	if !data.Pagination.IsNull() && !data.Pagination.IsUnknown() {
		diags := data.Pagination.As(ctx, &tfpagination, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			response.Diagnostics.Append(diags...)
		}
	}

	offset := tfpagination.Offset.ValueInt32()
	limit := tfpagination.Limit.ValueInt32()
	if limit == 0 {
		limit = 20
	}

	routeAggregationID := data.RouteAggregationID.ValueString()
	routeAggregationRequest := client.RouteAggregationRulesApi.SearchRouteAggregationRules(ctx, routeAggregationID)
	pagination := fabricv4.NewPaginationRequestWithDefaults()

	if !tfpagination.Limit.IsNull() {
		pagination.SetLimit(limit)
	}
	if !tfpagination.Offset.IsNull() {
		pagination.SetOffset(offset)
	}

	params.SetPagination(*pagination)

	if response.Diagnostics.HasError() {
		return
	}

	routeAggregations, _, err := routeAggregationRequest.RouteAggregationRulesSearchRequest(*params).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving route aggregation rules data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, routeAggregations)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
