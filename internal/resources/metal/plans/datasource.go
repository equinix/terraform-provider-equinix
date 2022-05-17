package plans

import (
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/equinix/terraform-provider-equinix/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func DataSource() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:               planSchema(),
		ResultAttributeName:        "plans",
		ResultAttributeDescription: "Sorted list of available server plans that match the specified filters",
		FlattenRecord:              flattenPlan,
		GetRecords:                 getPlans,
	}

	return datalist.NewResource(dataListConfig)
}

func getPlans(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.Config).MetalClient
	opts := &packngo.ListOptions{
		Includes: []string{"available_in", "available_in_metros"},
	}
	plans, _, err := client.Plans.List(opts)
	plansIf := []interface{}{}
	for _, p := range plans {
		plansIf = append(plansIf, p)
	}
	return plansIf, err
}

func planSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "id of the plan",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "name of the plan",
		},
		"slug": {
			Type:        schema.TypeString,
			Description: "plan slug",
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description of the plan",
		},
		"line": {
			Type:        schema.TypeString,
			Description: "plan line, e.g. baremetal",
		},
		"legacy": {
			Type:        schema.TypeBool,
			Description: "flag showing if it's a legacy plan",
		},
		"class": {
			Type:        schema.TypeString,
			Description: "plan class",
		},
		"pricing_hour": {
			Type:        schema.TypeFloat,
			Description: "plan hourly price",
		},
		"pricing_month": {
			Type:        schema.TypeFloat,
			Description: "plan monthly price",
		},
		"deployment_types": {
			Type:        schema.TypeSet,
			Description: "list of deployment types, e.g. on_demand, spot_market",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"available_in": {
			Type:        schema.TypeSet,
			Description: "list of facilities where the plan is available",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"available_in_metros": {
			Type:        schema.TypeSet,
			Description: "list of metros where the plan is available",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func flattenPlan(rawPlan interface{}, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	plan, ok := rawPlan.(packngo.Plan)
	if !ok {
		return nil, fmt.Errorf("unable to convert to packngo.Plan")
	}

	facs := []string{}
	for _, f := range plan.AvailableIn {
		facs = append(facs, f.Code)
	}

	metros := []string{}
	for _, m := range plan.AvailableInMetros {
		metros = append(metros, m.Code)
	}

	flattenedFacs := schema.NewSet(schema.HashString, converters.ConvertStringArrToIfArr(facs))
	flattenedMetros := schema.NewSet(schema.HashString, converters.ConvertStringArrToIfArr(metros))
	flattenedDepTypes := schema.NewSet(schema.HashString,
		converters.ConvertStringArrToIfArr(plan.DeploymentTypes))

	flattenedPlan := map[string]interface{}{
		"id":                  plan.ID,
		"name":                plan.Name,
		"slug":                plan.Slug,
		"description":         plan.Description,
		"line":                plan.Line,
		"legacy":              plan.Legacy,
		"class":               plan.Class,
		"deployment_types":    flattenedDepTypes,
		"available_in":        flattenedFacs,
		"available_in_metros": flattenedMetros,
	}

	if plan.Pricing != nil {
		flattenedPlan["pricing_hour"] = float64(plan.Pricing.Hour)
		flattenedPlan["pricing_month"] = float64(plan.Pricing.Month)
	}

	return flattenedPlan, nil
}
