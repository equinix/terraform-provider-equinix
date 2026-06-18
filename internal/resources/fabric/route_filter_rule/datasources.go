// Package route_filter_rule provides resources and data sources for managing Equinix Fabric route filter rules.
package route_filter_rule

import (
	"context"
	"fmt"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSource returns the schema.Resource for fetching a Fabric route filter rule by UUID.
func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema:      dataSourceByUUIDSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch route filter for a given UUID

Additional Documentation:
* Getting Started: https://docs.equinix.com/fabric-cloud-router/bgp/fcr-route-filters/
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Route-Filter-Rules`,
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	uuid := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceRead(ctx, d, meta)
}

// DataSourceGetRules returns the schema.Resource for fetching all route filter rules for a given route filter.
func DataSourceGetRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGetRules,
		Schema:      dataSourceRulesForRouteFilterSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch route filter for a given search data set

Additional Documentation:
* Getting Started: https://docs.equinix.com/fabric-cloud-router/bgp/fcr-route-filters/
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Route-Filter-Rules`,
	}
}

func dataSourceGetRules(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	routeFilterID := d.Get("route_filter_id").(string)

	filters := d.Get("filter").([]any)

	toReq := func(f map[string]any) *fabricv4.RouteFilterRuleSimpleExpression {
		property := f["property"].(string)
		operator := f["operator"].(string)
		values := f["values"].([]any)
		return &fabricv4.RouteFilterRuleSimpleExpression{
			Property: &property,
			Operator: &operator,
			Values: func(vs []any) []string {
				result := []string{}
				for _, v := range vs {
					result = append(result, v.(string))
				}
				return result
			}(values),
		}
	}

	innerRuleFilters := []fabricv4.RouteFilterRuleExpression{}
	for _, filter := range filters {
		innerRuleFilters = append(innerRuleFilters, fabricv4.RouteFilterRuleExpression{
			RouteFilterRuleSimpleExpression: toReq(filter.(map[string]any))})
	}

	rulesFilter := fabricv4.RouteFilterRulesFilter{}

	if outerOperator := d.Get("outer_operator"); outerOperator != nil {
		outerOperator := outerOperator.(string)

		switch outerOperator {
		case "AND":
			rulesFilter.RouteFilterRuleAndExpression = &fabricv4.RouteFilterRuleAndExpression{And: innerRuleFilters}
		case "OR":
			rulesFilter.RouteFilterRuleOrExpression = &fabricv4.RouteFilterRuleOrExpression{Or: innerRuleFilters}
		}
	}

	searchRequest := fabricv4.NewRouteFilterRulesSearchRequestWithDefaults()
	searchRequest.SetFilter(rulesFilter)

	sortInputs := d.Get("sort").([]any)
	reqSortParam := []fabricv4.RouteFilterRuleSortCriteria{}
	for _, sortInput := range sortInputs {
		sortItem := fabricv4.NewRouteFilterRuleSortCriteriaWithDefaults()
		sortInput := sortInput.(map[string]any)

		if direction, ok := sortInput["direction"]; ok {
			sortItem.SetDirection(fabricv4.RouteFilterRuleSortDirection(direction.(string)))
		}
		if property, ok := sortInput["property"]; ok {
			sortItem.SetProperty(fabricv4.RouteFilterRuleSortBy(property.(string)))
		}

		reqSortParam = append(reqSortParam, *sortItem)
	}
	searchRequest.SetSort(reqSortParam)

	searchRequest.SetPagination(*fabricv4.NewPaginationRequestWithDefaults())

	limit := d.Get("limit").(int)
	searchRequest.Pagination.SetLimit(int32(limit))

	offset := d.Get("offset").(int)
	searchRequest.Pagination.SetOffset(int32(offset))

	req := client.RouteFilterRulesApi.SearchRouteFilterRules(ctx, routeFilterID).
		RouteFilterRulesSearchRequest(*searchRequest)

	searchRouteFilterRules, _, err := req.Execute()

	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	if len(searchRouteFilterRules.Data) < 1 {
		return diag.FromErr(fmt.Errorf("no records are found for the route filter (%s) - %d , please change the search criteria", routeFilterID, len(searchRouteFilterRules.Data)))
	}

	d.SetId(routeFilterID)
	if err := d.Set("route_filter_id", routeFilterID); err != nil {
		return diag.Errorf("error setting route_filter_id to state %s", err)
	}
	return setRouteFilterRulesData(d, searchRouteFilterRules)
}
