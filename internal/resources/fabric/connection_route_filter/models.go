package connection_route_filter

import (
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setConnectionRouteFilterMap(d *schema.ResourceData, connectionRouteFilter *fabricv4.ConnectionRouteFilterData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	routeFilterMap := connectionRouteFilterResponseMap(connectionRouteFilter)
	err := equinix_schema.SetMap(d, routeFilterMap)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func setConnectionRouteFilterData(d *schema.ResourceData, connectionRouteFilters *fabricv4.GetAllConnectionRouteFiltersResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}
	mappedRouteFilters := make([]map[string]interface{}, len(connectionRouteFilters.Data))
	pagination := connectionRouteFilters.GetPagination()
	if connectionRouteFilters.Data != nil {
		for index, routeFilter := range connectionRouteFilters.Data {
			mappedRouteFilters[index] = connectionRouteFilterResponseMap(&routeFilter)
		}
	} else {
		mappedRouteFilters = nil
	}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"data":       mappedRouteFilters,
		"pagination": paginationGoToTerraform(&pagination),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func connectionRouteFilterResponseMap(data *fabricv4.ConnectionRouteFilterData) map[string]interface{} {
	connectionRouteFilterMap := make(map[string]interface{})
	connectionRouteFilterMap["href"] = data.GetHref()
	connectionRouteFilterMap["type"] = string(data.GetType())
	connectionRouteFilterMap["uuid"] = data.GetUuid()
	connectionRouteFilterMap["attachment_status"] = string(data.GetAttachmentStatus())
	connectionRouteFilterMap["direction"] = string(data.GetDirection())

	return connectionRouteFilterMap
}

func paginationGoToTerraform(pagination *fabricv4.Pagination) *schema.Set {
	if pagination == nil {
		return nil
	}
	mappedPagination := make(map[string]interface{})
	mappedPagination["offset"] = int(pagination.GetOffset())
	mappedPagination["limit"] = int(pagination.GetLimit())
	mappedPagination["total"] = int(pagination.GetTotal())
	mappedPagination["next"] = pagination.GetNext()
	mappedPagination["previous"] = pagination.GetPrevious()

	return schema.NewSet(
		schema.HashResource(paginationSchema()),
		[]interface{}{mappedPagination},
	)
}
