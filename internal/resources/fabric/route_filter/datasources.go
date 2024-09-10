package route_filter

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
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceRead(ctx, d, meta)
}

func DataSourceSearch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSearch,
		Schema:      dataSourceSearchSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch route filter for a given search data set",
	}
}

func dataSourceSearch(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	searchRequest := buildSearchRequest(d)

	routeFilters, _, err := client.RouteFiltersApi.SearchRouteFilters(ctx).RouteFiltersSearchBase(searchRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	if len(routeFilters.Data) < 1 {
		return diag.FromErr(fmt.Errorf("no records are found for the route filter search criteria provided - %d , please change the search criteria", len(routeFilters.Data)))
	}

	d.SetId(routeFilters.Data[0].GetUuid())
	return setRouteFiltersData(d, routeFilters)
}
