package route_filter_rule

import (
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func buildCreateRequest(d *schema.ResourceData) fabricv4.RouteFilterRulesBase {
	routeFilterRuleRequest := fabricv4.RouteFilterRulesBase{}

	prefix := d.Get("prefix").(string)
	routeFilterRuleRequest.SetPrefix(prefix)

	prefixMatch := d.Get("prefix_match").(string)
	if prefixMatch != "" {
		routeFilterRuleRequest.SetPrefixMatch(prefixMatch)
	}

	name := d.Get("name").(string)
	if name != "" {
		routeFilterRuleRequest.SetName(name)
	}

	description := d.Get("description").(string)
	if description != "" {
		routeFilterRuleRequest.SetDescription(description)
	}

	return routeFilterRuleRequest
}

func buildUpdateRequest(d *schema.ResourceData) []fabricv4.RouteFilterRulesPatchRequestItem {
	patches := make([]fabricv4.RouteFilterRulesPatchRequestItem, 0)
	oldName, newName := d.GetChange("name")
	if oldName.(string) != newName.(string) {
		patches = append(patches, fabricv4.RouteFilterRulesPatchRequestItem{
			Op:    "replace",
			Path:  "/name",
			Value: newName.(string),
		})
	}

	oldPrefix, newPrefix := d.GetChange("prefix")
	if oldPrefix.(string) != newPrefix.(string) {
		patches = append(patches, fabricv4.RouteFilterRulesPatchRequestItem{
			Op:    "replace",
			Path:  "/prefix",
			Value: newPrefix.(string),
		})
	}

	oldPrefixMatch, newPrefixMatch := d.GetChange("prefix_match")
	if oldPrefixMatch.(string) != newPrefixMatch.(string) {
		patches = append(patches, fabricv4.RouteFilterRulesPatchRequestItem{
			Op:    "replace",
			Path:  "/prefixMatch",
			Value: newPrefixMatch.(string),
		})
	}

	return patches
}

func routeFilterRuleResponseMap(data *fabricv4.RouteFilterRulesData) map[string]interface{} {
	routeFilterRuleMap := make(map[string]interface{})
	routeFilterRuleMap["type"] = string(data.GetType())
	routeFilterRuleMap["name"] = data.GetName()
	routeFilterRuleMap["description"] = data.GetDescription()
	routeFilterRuleMap["href"] = data.GetHref()
	routeFilterRuleMap["uuid"] = data.GetUuid()
	routeFilterRuleMap["state"] = string(data.GetState())
	routeFilterRuleMap["prefix"] = data.GetPrefix()
	routeFilterRuleMap["prefix_match"] = data.GetPrefixMatch()
	routeFilterRuleMap["action"] = data.GetAction()
	if data.Change != nil {
		change := data.GetChange()
		routeFilterRuleMap["change"] = changeGoToTerraform(&change)
	}
	if data.Changelog != nil {
		changeLog := data.GetChangelog()
		routeFilterRuleMap["change_log"] = equinix_fabric_schema.ChangeLogGoToTerraform(&changeLog)
	}

	return routeFilterRuleMap
}

func changeGoToTerraform(change *fabricv4.RouteFilterRulesChange) *schema.Set {
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

func paginationGoToTerraform(pagination *fabricv4.Pagination) *schema.Set {
	if pagination == nil {
		return nil
	}
	mappedPagination := make(map[string]interface{})
	mappedPagination["offset"] = pagination.GetOffset()
	mappedPagination["limit"] = pagination.GetLimit()
	mappedPagination["total"] = pagination.GetTotal()
	mappedPagination["next"] = pagination.GetNext()
	mappedPagination["previous"] = pagination.GetPrevious()

	return schema.NewSet(
		schema.HashResource(paginationSchema()),
		[]interface{}{mappedPagination},
	)
}
