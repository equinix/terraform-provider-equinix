package connection

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
	"strings"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricConnectionRead,
		Schema:      readFabricConnectionResourceSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch connection for a given UUID",
	}
}

func dataSourceFabricConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricConnectionRead(ctx, d, meta)
}

func DataSourceSearch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricConnectionSearch,
		Schema:      readFabricConnectionSearchSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch connection for a given UUID",
	}
}

func dataSourceFabricConnectionSearch(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceFabricConnectionSearch(ctx, d, meta)
}

func resourceFabricConnectionSearch(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	connectionSearchRequest := fabricv4.SearchRequest{}

	schemaFilters := d.Get("filter").([]interface{})
	schemaOuterOperator := d.Get("outer_operator").(string)
	filter, err := connectionFiltersTerraformToGo(schemaFilters, schemaOuterOperator)
	if err != nil {
		return diag.FromErr(err)
	}

	connectionSearchRequest.SetFilter(filter)

	if schemaPagination, ok := d.GetOk("pagination"); ok {
		pagination := connectionPaginationTerraformToGo(schemaPagination.(*schema.Set).List())
		connectionSearchRequest.SetPagination(pagination)
	}

	if schemaSort, ok := d.GetOk("sort"); ok {
		sort := connectionSortTerraformToGo(schemaSort.([]interface{}))
		connectionSearchRequest.SetSort(sort)
	}

	connections, _, err := client.ConnectionsApi.SearchConnections(ctx).SearchRequest(connectionSearchRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	if len(connections.Data) < 1 {
		return diag.FromErr(fmt.Errorf("no records are found for the connection search criteria provided - %d , please change the search criteria", len(connections.Data)))
	}

	d.SetId(connections.Data[0].GetUuid())
	return setConnectionsData(d, connections)
}

func setConnectionsData(d *schema.ResourceData, connections *fabricv4.ConnectionSearchResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}
	mappedConnections := make([]map[string]interface{}, len(connections.Data))
	if connections.Data != nil {
		for index, connection := range connections.Data {
			mappedConnections[index] = connectionMap(&connection)
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

func connectionFiltersTerraformToGo(filters []interface{}, outerOperator string) (fabricv4.Expression, error) {
	if filters == nil || len(filters) == 0 {
		return fabricv4.Expression{}, fmt.Errorf("no filters passed to filtersTerraformToGoMethod")
	}
	outerExpression := fabricv4.Expression{}
	expressions := make([]fabricv4.Expression, 0)
	groups := make(map[string]fabricv4.Expression)

	for _, filter := range filters {
		filterMap := filter.(map[string]interface{})
		expression := fabricv4.Expression{}
		if property, ok := filterMap["property"]; ok {
			expression.SetProperty(fabricv4.SearchFieldName(property.(string)))
		}
		if operator, ok := filterMap["operator"]; ok {
			expression.SetOperator(fabricv4.ExpressionOperator(operator.(string)))
		}
		if values, ok := filterMap["values"]; ok {
			stringValues := converters.IfArrToStringArr(values.([]interface{}))
			expression.SetValues(stringValues)
		}

		// If the parent has any contents then all the children schema properties will be included in the map even
		// if they aren't given a value. Still need to check for empty string for the value because of this.
		if groupInterface, ok := filterMap["group"]; ok && groupInterface.(string) != "" {
			group := groupInterface.(string)
			groupExpression := fabricv4.Expression{}
			if _, ok := groups[group]; ok {
				groupExpression = groups[group]
				var expressionList []fabricv4.Expression
				if strings.HasPrefix(group, "AND_") {
					expressionList = groupExpression.GetAnd()
					expressionList = append(expressionList, expression)
					groupExpression.SetAnd(expressionList)
				} else if strings.HasPrefix(group, "OR_") {
					expressionList = groupExpression.GetOr()
					expressionList = append(expressionList, expression)
					groupExpression.SetOr(expressionList)
				}
			} else {
				expressionList := make([]fabricv4.Expression, 1)
				expressionList[0] = expression
				if strings.HasPrefix(group, "AND_") {
					groupExpression.SetAnd(expressionList)
				} else if strings.HasPrefix(group, "OR_") {
					groupExpression.SetOr(expressionList)
				}
			}
			groups[group] = groupExpression
		} else {
			expressions = append(expressions, expression)
		}
	}

	for _, value := range groups {
		expressions = append(expressions, value)
	}

	if outerOperator == "AND" {
		outerExpression.SetAnd(expressions)
	} else if outerOperator == "OR" {
		outerExpression.SetOr(expressions)
	}

	return outerExpression, nil
}

func connectionPaginationTerraformToGo(pagination []interface{}) fabricv4.PaginationRequest {
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

func connectionSortTerraformToGo(sort []interface{}) []fabricv4.SortCriteria {
	if sort == nil || len(sort) == 0 {
		return []fabricv4.SortCriteria{}
	}
	sortCriteria := make([]fabricv4.SortCriteria, len(sort))
	for index, item := range sort {
		sortItem := fabricv4.SortCriteria{}
		pageMap := item.(map[string]interface{})
		if direction, ok := pageMap["direction"]; ok {
			sortItem.SetDirection(fabricv4.SortDirection(direction.(string)))
		}
		if property, ok := pageMap["property"]; ok {
			sortItem.SetProperty(fabricv4.SortBy(property.(string)))
		}
		sortCriteria[index] = sortItem
	}
	return sortCriteria
}
