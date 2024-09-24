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
		Schema:      dataSourceBaseSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch route filter for a given UUID",
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
		Description: "Fabric V4 API compatible data resource that allow user to fetch route filter for a given search data set",
	}
}

func dataSourceGetAllRules(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	routeFilterId := d.Get("route_filter_id").(string)

	routeFilterRules, _, err := client.RouteFilterRulesApi.GetRouteFilterRules(ctx, routeFilterId).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	if len(routeFilterRules.Data) < 1 {
		return diag.FromErr(fmt.Errorf("no records are found for the route filter (%s) - %d , please change the search criteria", routeFilterId, len(routeFilterRules.Data)))
	}

	d.SetId(routeFilterId)
	return setRouteFilterRulesData(d, routeFilterRules)
}
