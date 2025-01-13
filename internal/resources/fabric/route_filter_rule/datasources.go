package route_filter_rule

import (
	"context"
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema:      dataSourceByUUIDSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch route filter for a given UUID

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/FCR/FCR-route-filters.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#route-filter-rules`,
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceRead(ctx, d, meta)
}

func DataSourceGetAllRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGetAllRules,
		Schema:      dataSourceAllRulesForRouteFilterSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch route filter for a given search data set

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/FCR/FCR-route-filters.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#route-filter-rules`,
	}
}

func dataSourceGetAllRules(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	routeFilterId := d.Get("route_filter_id").(string)
	getRouteFilterRulesRequest := client.RouteFilterRulesApi.GetRouteFilterRules(ctx, routeFilterId)

	limit := d.Get("limit").(int)
	if limit != 0 {
		getRouteFilterRulesRequest.Limit(int32(limit))
		err := d.Set("limit", limit)
		if err != nil {
			return diag.Errorf("error setting limit to state %s", err)
		}
	}
	offset := d.Get("offset").(int)
	if offset != 0 {
		getRouteFilterRulesRequest.Offset(int32(offset))
		err := d.Set("offset", offset)
		if err != nil {
			return diag.Errorf("error setting offset to state %s", err)
		}
	}

	routeFilterRules, _, err := getRouteFilterRulesRequest.Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	if len(routeFilterRules.Data) < 1 {
		return diag.FromErr(fmt.Errorf("no records are found for the route filter (%s) - %d , please change the search criteria", routeFilterId, len(routeFilterRules.Data)))
	}

	d.SetId(routeFilterId)
	err = d.Set("route_filter_id", routeFilterId)
	if err != nil {
		return diag.Errorf("error setting route_filter_id to state %s", err)
	}
	return setRouteFilterRulesData(d, routeFilterRules)
}
