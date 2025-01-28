package route_filter

import (
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func buildCreateRequest(d *schema.ResourceData) fabricv4.RouteFiltersBase {
	routeFilterRequest := fabricv4.RouteFiltersBase{}

	typeConfig := d.Get("type").(string)
	routeFilterRequest.SetType(fabricv4.RouteFiltersBaseType(typeConfig))

	nameConfig := d.Get("name").(string)
	routeFilterRequest.SetName(nameConfig)

	projectConfig := d.Get("project").(*schema.Set).List()
	project := projectTerraformToGo(projectConfig)
	routeFilterRequest.SetProject(project)

	description := d.Get("description").(string)
	if description != "" {
		routeFilterRequest.SetDescription(description)
	}

	return routeFilterRequest
}

func buildUpdateRequest(d *schema.ResourceData) []fabricv4.RouteFiltersPatchRequestItem {
	patches := make([]fabricv4.RouteFiltersPatchRequestItem, 0)
	oldName, newName := d.GetChange("name")
	if oldName.(string) != newName.(string) {
		patches = append(patches, fabricv4.RouteFiltersPatchRequestItem{
			Op:    "replace",
			Path:  "/name",
			Value: newName.(string),
		})
	}

	oldDescription, newDescription := d.GetChange("description")
	if oldDescription.(string) != newDescription.(string) {
		patches = append(patches, fabricv4.RouteFiltersPatchRequestItem{
			Op:    "replace",
			Path:  "/description",
			Value: newDescription.(string),
		})
	}

	oldProjectId, newProjectId := d.GetChange("project.0.project_id")
	if oldProjectId.(string) != newProjectId.(string) {
		patches = append(patches, fabricv4.RouteFiltersPatchRequestItem{
			Op:    "replace",
			Path:  "/project/projectId",
			Value: newProjectId.(string),
		})
	}

	return patches
}

func buildSearchRequest(d *schema.ResourceData) fabricv4.RouteFiltersSearchBase {
	searchRequest := fabricv4.RouteFiltersSearchBase{}

	schemaFilters := d.Get("filter").([]interface{})
	filter := filtersTerraformToGo(schemaFilters)
	searchRequest.SetFilter(filter)

	if schemaPagination, ok := d.GetOk("pagination"); ok {
		pagination := paginationTerraformToGo(schemaPagination.(*schema.Set).List())
		searchRequest.SetPagination(pagination)
	}

	if schemaSort, ok := d.GetOk("sort"); ok {
		sort := sortTerraformToGo(schemaSort.([]interface{}))
		searchRequest.SetSort(sort)
	}
	return searchRequest
}

func setRouteFilterMap(d *schema.ResourceData, routeFilter *fabricv4.RouteFiltersData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	routeFilterMap := routeFilterResponseMap(routeFilter)
	err := equinix_schema.SetMap(d, routeFilterMap)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func setRouteFiltersData(d *schema.ResourceData, routeFilters *fabricv4.RouteFiltersSearchResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}
	mappedRouteFilters := make([]map[string]interface{}, len(routeFilters.Data))
	pagination := routeFilters.GetPagination()
	if routeFilters.Data != nil {
		for index, routeFilter := range routeFilters.Data {
			mappedRouteFilters[index] = routeFilterResponseMap(&routeFilter)
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

func routeFilterResponseMap(data *fabricv4.RouteFiltersData) map[string]interface{} {
	routeFilterMap := make(map[string]interface{})
	routeFilterMap["type"] = string(data.GetType())
	routeFilterMap["name"] = data.GetName()
	if data.Project != nil {
		project := data.GetProject()
		routeFilterMap["project"] = projectGoToTerraform(&project)
	}
	routeFilterMap["description"] = data.GetDescription()
	routeFilterMap["href"] = data.GetHref()
	routeFilterMap["uuid"] = data.GetUuid()
	routeFilterMap["state"] = string(data.GetState())
	routeFilterMap["not_matched_rule_action"] = data.GetNotMatchedRuleAction()
	routeFilterMap["connections_count"] = data.GetConnectionsCount()
	routeFilterMap["rules_count"] = data.GetRulesCount()
	if data.Change != nil {
		change := data.GetChange()
		routeFilterMap["change"] = changeGoToTerraform(&change)
	}
	if data.Changelog != nil {
		changeLog := data.GetChangelog()
		routeFilterMap["change_log"] = equinix_fabric_schema.ChangeLogGoToTerraform(&changeLog)
	}

	return routeFilterMap
}

func projectTerraformToGo(projectTerraform []interface{}) fabricv4.Project {
	if projectTerraform == nil {
		return fabricv4.Project{}
	}

	project := fabricv4.Project{}
	projectMap := projectTerraform[0].(map[string]interface{})
	project.SetProjectId(projectMap["project_id"].(string))
	return project
}

func projectGoToTerraform(project *fabricv4.RouteFiltersDataProject) *schema.Set {
	if project == nil {
		return nil
	}

	mappedProject := make(map[string]interface{})
	mappedProject["project_id"] = project.GetProjectId()
	mappedProject["href"] = project.GetHref()
	projectSet := schema.NewSet(
		schema.HashResource(projectSch()),
		[]interface{}{mappedProject},
	)
	return projectSet
}

func changeGoToTerraform(change *fabricv4.RouteFiltersChange) *schema.Set {
	if change == nil {
		return nil
	}

	mappedChange := make(map[string]interface{})
	mappedChange["href"] = change.GetHref()
	mappedChange["type"] = string(change.GetType())
	mappedChange["uuid"] = change.GetUuid()
	changeSet := schema.NewSet(
		schema.HashResource(changeSch()),
		[]interface{}{mappedChange},
	)

	return changeSet
}

func filtersTerraformToGo(filters []interface{}) fabricv4.RouteFiltersSearchBaseFilter {
	if filters == nil {
		return fabricv4.RouteFiltersSearchBaseFilter{}
	}

	searchFiltersList := make([]fabricv4.RouteFiltersSearchFilterItem, 0)

	for _, filter := range filters {
		filterMap := filter.(map[string]interface{})
		filterItem := fabricv4.RouteFiltersSearchFilterItem{}
		if property, ok := filterMap["property"]; ok {
			filterItem.SetProperty(fabricv4.RouteFiltersSearchFilterItemProperty(property.(string)))
		}
		if operator, ok := filterMap["operator"]; ok {
			filterItem.SetOperator(operator.(string))
		}
		if values, ok := filterMap["values"]; ok {
			stringValues := converters.IfArrToStringArr(values.([]interface{}))
			filterItem.SetValues(stringValues)
		}
		searchFiltersList = append(searchFiltersList, filterItem)
	}

	searchFilters := fabricv4.RouteFiltersSearchBaseFilter{}
	searchFilters.SetAnd(searchFiltersList)

	return searchFilters
}

func paginationTerraformToGo(pagination []interface{}) fabricv4.Pagination {
	if pagination == nil {
		return fabricv4.Pagination{}
	}
	paginationRequest := fabricv4.Pagination{}
	for _, page := range pagination {
		pageMap := page.(map[string]interface{})
		if offset, ok := pageMap["offset"]; ok {
			paginationRequest.SetOffset(int32(offset.(int)))
		}
		if limit, ok := pageMap["limit"]; ok {
			paginationRequest.SetLimit(int32(limit.(int)))
		}
		if total, ok := pageMap["total"]; ok {
			paginationRequest.SetTotal(int32(total.(int)))
		}
	}

	return paginationRequest
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

func sortTerraformToGo(sort []interface{}) []fabricv4.SortItem {
	if sort == nil {
		return []fabricv4.SortItem{}
	}
	sortItems := make([]fabricv4.SortItem, len(sort))
	for index, item := range sort {
		sortItem := fabricv4.SortItem{}
		pageMap := item.(map[string]interface{})
		if direction, ok := pageMap["direction"]; ok {
			sortItem.SetDirection(fabricv4.SortItemDirection(direction.(string)))
		}
		if property, ok := pageMap["property"]; ok {
			sortItem.SetProperty(fabricv4.SortItemProperty(property.(string)))
		}
		sortItems[index] = sortItem
	}
	return sortItems
}
