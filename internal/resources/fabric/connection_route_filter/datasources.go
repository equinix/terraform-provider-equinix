package connection_route_filter

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
		Description: `Fabric V4 API compatible data resource that allow user to fetch route filter policy attachment to a fabric connection

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/FCR/FCR-route-filters.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#route-filter-rules`,
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid := d.Get("route_filter_id").(string)
	d.SetId(uuid)
	return resourceRead(ctx, d, meta)
}

func DataSourceGetAllRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGetAllFilters,
		Schema:      dataSourceAllFiltersSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch all route filter policies attached to a fabric connection

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/FCR/FCR-route-filters.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#route-filter-rules`,
	}
}

func dataSourceGetAllFilters(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	connectionId := d.Get("connection_id").(string)
	connectionRouteFilters, _, err := client.RouteFiltersApi.GetConnectionRouteFilters(ctx, connectionId).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	if len(connectionRouteFilters.Data) < 1 {
		return diag.FromErr(fmt.Errorf("no records are found for the connection (%s) - %d , please change the search criteria", connectionId, len(connectionRouteFilters.Data)))
	}
	d.SetId(connectionId)
	return setConnectionRouteFilterData(d, connectionRouteFilters)
}
