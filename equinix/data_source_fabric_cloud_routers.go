package equinix

import (
	"context"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func readFabricCloudRouterResourceSchemaUpdated() map[string]*schema.Schema {
	sch := readFabricCloudRouterResourceSchema()
	sch["uuid"].Computed = true
	sch["uuid"].Optional = false
	sch["uuid"].Required = false
	return sch
}

func readFabricCloudRouterSearchSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of Cloud Routers",
			Elem: &schema.Resource{
				Schema: readFabricCloudRouterResourceSchemaUpdated(),
			},
		},
		"filter": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Filters for the Data Source Search Request. Maximum of 8 total filters.",
			MaxItems:    10,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The API response property which you want to filter your request on. Can be one of the following: \"/project/projectId\", \"/name\", \"/uuid\", \"/state\", \"/location/metroCode\", \"/location/metroName\", \"/package/code\", \"/*\"",
						ValidateFunc: validation.StringInSlice([]string{"/project/projectId", "/name", "/uuid", "/state", "/location/metroCode", "/location/metroName", "/package/code", "/*"}, true),
					},
					"operator": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Possible operators to use on the filter property. Can be one of the following: [= - equal, != - not equal, > - greater than, >= - greater than or equal to, < - less than, <= - less than or equal to, [NOT] BETWEEN - (not) between, [NOT] LIKE - (not) like, [NOT] IN - (not) in",
						ValidateFunc: validation.StringInSlice([]string{"=", "!=", ">", ">=", "<", "<=", "[NOT] BETWEEN", "[NOT] LIKE", "[NOT] IN"}, true),
					},
					"values": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "The values that you want to apply the property+operator combination to in order to filter your data search",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"or": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Boolean flag indicating whether this filter is included in the OR group. There can only be one OR group and it can have a maximum of 3 filters. The OR group only counts as 1 of the 8 possible filters",
					},
				},
			},
		},
		"pagination": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Pagination details for the Data Source Search Request",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"offset": {
						Type:        schema.TypeInt,
						Optional:    true,
						Default:     0,
						Description: "The page offset for the pagination request. Index of the first element. Default is 0.",
					},
					"limit": {
						Type:        schema.TypeInt,
						Optional:    true,
						Default:     20,
						Description: "Number of elements to be requested per page. Number must be between 1 and 100. Default is 20",
					},
				},
			},
		},
		"sort": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Filters for the Data Source Search Request",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"direction": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "DESC",
						Description:  "The sorting direction. Can be one of: [DESC, ASC], Defaults to DESC",
						ValidateFunc: validation.StringInSlice([]string{"DESC", "ASC"}, true),
					},
					"property": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "/changeLog/updatedDateTime",
						Description:  "The property name to use in sorting. Can be one of the following: [/name, /uuid, /state, /location/metroCode, /location/metroName, /package/code, /changeLog/createdDateTime, /changeLog/updatedDateTime], Defaults to /changeLog/updatedDateTime",
						ValidateFunc: validation.StringInSlice([]string{"/name", "/uuid", "/state", "/location/metroCode", "/location/metroName", "/package/code", "/changeLog/createdDateTime", "/changeLog/updatedDateTime"}, true),
					},
				},
			},
		},
	}
}

func dataSourceFabricGetCloudRouters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricGetCloudRoutersRead,
		Schema:      readFabricCloudRouterSearchSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch Fabric Cloud Routers matching custom search criteria",
	}
}

func dataSourceFabricGetCloudRoutersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceFabricCloudRoutersSearch(ctx, d, meta)
}

func resourceFabricCloudRoutersSearch(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	cloudRouterSearchRequest := fabricv4.CloudRouterSearchRequest{}

	schemaFilters := d.Get("filter").([]interface{})
	filters, err := cloudRouterFiltersTerraformToGo(schemaFilters)
	if err != nil {
		return diag.FromErr(err)
	}

	cloudRouterSearchRequest.SetFilter(filters)

	if schemaPagination, ok := d.GetOk("pagination"); ok {
		pagination := cloudRouterPaginationTerraformToGo(schemaPagination.(*schema.Set).List())
		cloudRouterSearchRequest.SetPagination(pagination)
	}

	if schemaSort, ok := d.GetOk("sort"); ok {
		sort := cloudRouterSortTerraformToGo(schemaSort.([]interface{}))
		cloudRouterSearchRequest.SetSort(sort)
	}

	cloudRouters, _, err := client.CloudRoutersApi.SearchCloudRouters(ctx).CloudRouterSearchRequest(cloudRouterSearchRequest).Execute()

	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	if len(cloudRouters.Data) < 1 {
		return diag.FromErr(fmt.Errorf("no records are found for the cloud router search criteria provided - %d , please change the search criteria", len(cloudRouters.Data)))
	}

	d.SetId(cloudRouters.Data[0].GetUuid())
	return setFabricCloudRoutersData(d, cloudRouters)
}

func cloudRouterFiltersTerraformToGo(filters []interface{}) (fabricv4.CloudRouterFilters, error) {
	if filters == nil || len(filters) == 0 {
		return fabricv4.CloudRouterFilters{}, fmt.Errorf("no filters passed to filtersTerraformToGoMethod")
	}
	cloudRouterFiltersList := make([]fabricv4.CloudRouterFilter, 0)
	cloudRouterOrFilter := fabricv4.CloudRouterOrFilter{}

	for _, filter := range filters {
		filterMap := filter.(map[string]interface{})
		cloudRouterFilter := fabricv4.CloudRouterFilter{}
		filterExpression := fabricv4.CloudRouterSimpleExpression{}
		if property, ok := filterMap["property"]; ok {
			filterExpression.SetProperty(property.(string))
		}
		if operator, ok := filterMap["operator"]; ok {
			filterExpression.SetOperator(operator.(string))
		}
		if values, ok := filterMap["values"]; ok {
			stringValues := converters.IfArrToStringArr(values.([]interface{}))
			filterExpression.SetValues(stringValues)
		}

		// If the parent has any contents then all the children schema properties will be included in the map even
		// if they aren't given a value. Still need to check for empty string for the value because of this.
		if orGroup, ok := filterMap["or"]; ok && orGroup.(bool) {
			orValues := cloudRouterOrFilter.GetOr()
			orValues = append(orValues, filterExpression)
			if len(orValues) > 3 {
				return fabricv4.CloudRouterFilters{}, fmt.Errorf("too many OR group filters passed. Passed %d but can only have a maximum of 3", len(orValues))
			}
			cloudRouterOrFilter.SetOr(orValues)
		} else {
			cloudRouterFilter.CloudRouterSimpleExpression = &filterExpression
			cloudRouterFiltersList = append(cloudRouterFiltersList, cloudRouterFilter)
		}
	}

	if orGroupHasValues := cloudRouterOrFilter.GetOr(); len(orGroupHasValues) > 0 {
		cloudRouterFilter := fabricv4.CloudRouterFilter{}
		cloudRouterFilter.CloudRouterOrFilter = &cloudRouterOrFilter
		cloudRouterFiltersList = append(cloudRouterFiltersList, cloudRouterFilter)
	}

	cloudRouterFilters := fabricv4.CloudRouterFilters{}
	cloudRouterFilters.SetAnd(cloudRouterFiltersList)

	if len(cloudRouterFilters.GetAnd()) > 8 {
		return fabricv4.CloudRouterFilters{}, fmt.Errorf("too many filters are applied to the data source. The maximum is 8 and %d were provided. Please reduce your filter count to 8", len(cloudRouterFilters.GetAnd()))
	}

	return cloudRouterFilters, nil
}

func cloudRouterPaginationTerraformToGo(pagination []interface{}) fabricv4.PaginationRequest {
	if pagination == nil || len(pagination) == 0 {
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

func cloudRouterSortTerraformToGo(sort []interface{}) []fabricv4.CloudRouterSortCriteria {
	if sort == nil || len(sort) == 0 {
		return []fabricv4.CloudRouterSortCriteria{}
	}
	sortCriteria := make([]fabricv4.CloudRouterSortCriteria, len(sort))
	for index, item := range sort {
		sortItem := fabricv4.CloudRouterSortCriteria{}
		pageMap := item.(map[string]interface{})
		if direction, ok := pageMap["direction"]; ok {
			sortItem.SetDirection(fabricv4.CloudRouterSortDirection(direction.(string)))
		}
		if property, ok := pageMap["property"]; ok {
			sortItem.SetProperty(fabricv4.CloudRouterSortBy(property.(string)))
		}
		sortCriteria[index] = sortItem
	}
	return sortCriteria
}

func setFabricCloudRoutersData(d *schema.ResourceData, cloudRouters *fabricv4.SearchResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}
	mappedCloudRouters := make([]map[string]interface{}, len(cloudRouters.Data))
	if cloudRouters.Data != nil {
		for index, cloudRouter := range cloudRouters.Data {
			mappedCloudRouters[index] = fabricCloudRouterMap(&cloudRouter)
		}
	} else {
		mappedCloudRouters = nil
	}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"data": mappedCloudRouters,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
