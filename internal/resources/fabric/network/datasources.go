package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricNetworkRead,
		Schema:      readFabricNetworkResourceSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch Fabric Network for a given UUID

Additional documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/Fabric/IMPLEMENTATION/fabric-networks-implement.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#fabric-networks`,
	}
}

func dataSourceFabricNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricNetworkRead(ctx, d, meta)
}

func DataSourceSearch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricNetworkSearch,
		Schema:      readFabricNetworkSearchSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch Fabric Network for a given UUID

Additional documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/Fabric/IMPLEMENTATION/fabric-networks-implement.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#fabric-networks`,
	}
}

func dataSourceFabricNetworkSearch(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceFabricNetworkSearch(ctx, d, meta)
}

func resourceFabricNetworkSearch(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	networkSearchRequest := fabricv4.NetworkSearchRequest{}

	schemaFilters := d.Get("filter").([]interface{})
	schemaOuterOperator := d.Get("outer_operator").(string)
	filter, err := networkFiltersTerraformToGo(schemaFilters, schemaOuterOperator)
	if err != nil {
		return diag.FromErr(err)
	}

	networkSearchRequest.SetFilter(filter)

	if schemaPagination, ok := d.GetOk("pagination"); ok {
		pagination := networkPaginationTerraformToGo(schemaPagination.(*schema.Set).List())
		networkSearchRequest.SetPagination(pagination)
	}

	if schemaSort, ok := d.GetOk("sort"); ok {
		sort := networkSortTerraformToGo(schemaSort.([]interface{}))
		networkSearchRequest.SetSort(sort)
	}

	networks, _, err := client.NetworksApi.SearchNetworks(ctx).NetworkSearchRequest(networkSearchRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	if len(networks.Data) < 1 {
		return diag.FromErr(fmt.Errorf("no records are found for the network search criteria provided - %d , please change the search criteria", len(networks.Data)))
	}

	d.SetId(networks.Data[0].GetUuid())
	return setNetworksData(d, networks)
}

func setNetworksData(d *schema.ResourceData, networks *fabricv4.NetworkSearchResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}
	mappedConnections := make([]map[string]interface{}, len(networks.Data))
	if networks.Data != nil {
		for index, network := range networks.Data {
			mappedConnections[index] = networkMap(&network)
		}
	} else {
		mappedConnections = nil
	}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"data": mappedConnections,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func networkFiltersTerraformToGo(filters []interface{}, outerOperator string) (fabricv4.NetworkFilter, error) {
	if len(filters) == 0 {
		return fabricv4.NetworkFilter{}, fmt.Errorf("no filters passed to filtersTerraformToGoMethod")
	}
	outerNetworkFilter := fabricv4.NetworkFilter{}
	networkFilters := make([]fabricv4.NetworkFilter, 0)
	groups := make(map[string]fabricv4.NetworkFilter)

	for _, filter := range filters {
		filterMap := filter.(map[string]interface{})
		networkFilter := fabricv4.NetworkFilter{}
		if property, ok := filterMap["property"]; ok {
			networkFilter.SetProperty(fabricv4.NetworkSearchFieldName(property.(string)))
		}
		if operator, ok := filterMap["operator"]; ok {
			networkFilter.SetOperator(fabricv4.NetworkFilterOperator(operator.(string)))
		}
		if values, ok := filterMap["values"]; ok {
			stringValues := converters.IfArrToStringArr(values.([]interface{}))
			networkFilter.SetValues(stringValues)
		}

		// If the parent has any contents then all the children schema properties will be included in the map even
		// if they aren't given a value. Still need to check for empty string for the value because of this.
		if groupInterface, ok := filterMap["group"]; ok && groupInterface.(string) != "" {
			group := groupInterface.(string)
			groupNetworkFilter := fabricv4.NetworkFilter{}
			if _, ok := groups[group]; ok {
				groupNetworkFilter = groups[group]
				var networkFilterList []fabricv4.NetworkFilter
				if strings.HasPrefix(group, "AND_") {
					networkFilterList = groupNetworkFilter.GetAnd()
					networkFilterList = append(networkFilterList, networkFilter)
					groupNetworkFilter.SetAnd(networkFilterList)
				} else if strings.HasPrefix(group, "OR_") {
					networkFilterList = groupNetworkFilter.GetOr()
					networkFilterList = append(networkFilterList, networkFilter)
					groupNetworkFilter.SetOr(networkFilterList)
				}
			} else {
				networkFilterList := make([]fabricv4.NetworkFilter, 1)
				networkFilterList[0] = networkFilter
				if strings.HasPrefix(group, "AND_") {
					groupNetworkFilter.SetAnd(networkFilterList)
				} else if strings.HasPrefix(group, "OR_") {
					groupNetworkFilter.SetOr(networkFilterList)
				}
			}
			groups[group] = groupNetworkFilter
		} else {
			networkFilters = append(networkFilters, networkFilter)
		}
	}

	for _, value := range groups {
		networkFilters = append(networkFilters, value)
	}

	if outerOperator == "AND" {
		outerNetworkFilter.SetAnd(networkFilters)
	} else if outerOperator == "OR" {
		outerNetworkFilter.SetOr(networkFilters)
	}

	return outerNetworkFilter, nil
}

func networkPaginationTerraformToGo(pagination []interface{}) fabricv4.PaginationRequest {
	if len(pagination) == 0 {
		return fabricv4.PaginationRequest{}
	}
	paginationRequest := fabricv4.PaginationRequest{}
	for _, page := range pagination {
		pageMap := page.(map[string]interface{})
		if offset, ok := pageMap["offset"]; ok {
			paginationRequest.SetOffset(int32(offset.(int)))
		}
		if limit, ok := pageMap["limit"]; ok {
			paginationRequest.SetLimit(int32(limit.(int)))
		}
	}

	return paginationRequest
}

func networkSortTerraformToGo(sort []interface{}) []fabricv4.NetworkSortCriteria {
	if len(sort) == 0 {
		return []fabricv4.NetworkSortCriteria{}
	}
	sortCriteria := make([]fabricv4.NetworkSortCriteria, len(sort))
	for index, item := range sort {
		sortItem := fabricv4.NetworkSortCriteria{}
		pageMap := item.(map[string]interface{})
		if direction, ok := pageMap["direction"]; ok {
			sortItem.SetDirection(fabricv4.NetworkSortDirection(direction.(string)))
		}
		if property, ok := pageMap["property"]; ok {
			sortItem.SetProperty(fabricv4.NetworkSortBy(property.(string)))
		}
		sortCriteria[index] = sortItem
	}
	return sortCriteria
}
