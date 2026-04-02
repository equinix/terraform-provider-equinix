// Package connection_route_filter provides resources and data sources for managing route filter attachments to Fabric connections.
package connection_route_filter

import (
	"context"
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSource returns the schema.Resource for fetching a route filter attachment to a Fabric connection by UUID.
func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema:      dataSourceByUUIDSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch route filter policy attachment to a fabric connection

Additional Documentation:
* Getting Started: https://docs.equinix.com/fabric-cloud-router/bgp/fcr-route-filters/
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Route-Filter-Rules`,
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid := d.Get("route_filter_id").(string)
	d.SetId(uuid)
	return resourceRead(ctx, d, meta)
}

// DataSourceGetAllRules returns the schema.Resource for fetching all route filter policies attached to a Fabric connection.
func DataSourceGetAllRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGetAllFilters,
		Schema:      dataSourceAllFiltersSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch all route filter policies attached to a fabric connection

Additional Documentation:
* Getting Started: https://docs.equinix.com/fabric-cloud-router/bgp/fcr-route-filters/
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Route-Filter-Rules`,
	}
}

func dataSourceGetAllFilters(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	connectionID := d.Get("connection_id").(string)
	connectionRouteFilters, _, err := client.RouteFiltersApi.GetConnectionRouteFilters(ctx, connectionID).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	if len(connectionRouteFilters.Data) < 1 {
		return diag.FromErr(fmt.Errorf("no records are found for the connection (%s) - %d , please change the search criteria", connectionID, len(connectionRouteFilters.Data)))
	}
	d.SetId(connectionID)
	return setConnectionRouteFilterData(d, connectionRouteFilters)
}
