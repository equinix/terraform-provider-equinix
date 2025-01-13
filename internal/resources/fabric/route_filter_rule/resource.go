package route_filter_rule

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
		},
		ReadContext:   resourceRead,
		CreateContext: resourceCreate,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceSchema(),
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric Route Filter Rule

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/FCR/FCR-route-filters.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#route-filter-rules`,
	}
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	routeFilterId := d.Get("route_filter_id").(string)
	routeFilterRule, _, err := client.RouteFilterRulesApi.GetRouteFilterRuleByUuid(ctx, routeFilterId, d.Id()).Execute()
	if err != nil {
		log.Printf("[WARN] Route Filter Rule %s not found on Route Filter %s, error %s", d.Id(), routeFilterId, err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(routeFilterRule.GetUuid())
	return setRouteFilterRuleMap(d, routeFilterRule)
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	routeFilterId := d.Get("route_filter_id").(string)
	createRequest := buildCreateRequest(d)

	start := time.Now()
	routeFilter, _, err := client.RouteFilterRulesApi.CreateRouteFilterRule(ctx, routeFilterId).RouteFilterRulesBase(createRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	err = d.Set("route_filter_id", routeFilterId)
	if err != nil {
		return diag.Errorf("error setting route_filter_id to state %s", err)
	}
	d.SetId(routeFilter.GetUuid())

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if err = waitForStability(routeFilterId, d.Id(), meta, d, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for route filter rule (%s) on route filter (%s) to be created: %s", d.Id(), routeFilterId, err)
	}

	return resourceRead(ctx, d, meta)
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	routeFilterId := d.Get("route_filter_id").(string)
	updateRequest := buildUpdateRequest(d)

	start := time.Now()
	routeFilter, _, err := client.RouteFilterRulesApi.PatchRouteFilterRuleByUuid(ctx, routeFilterId, d.Id()).RouteFilterRulesPatchRequestItem(updateRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	if err = waitForStability(routeFilterId, d.Id(), meta, d, ctx, updateTimeout); err != nil {
		return diag.Errorf("error waiting for route filter rule (%s) on route filter (%s) to be updated: %s", d.Id(), routeFilterId, err)
	}

	return setRouteFilterRuleMap(d, routeFilter)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	routeFilterId := d.Get("route_filter_id").(string)

	start := time.Now()
	_, _, err := client.RouteFilterRulesApi.DeleteRouteFilterRuleByUuid(ctx, routeFilterId, d.Id()).Execute()
	if err != nil {
		if genericError, ok := err.(*fabricv4.GenericOpenAPIError); ok {
			if fabricErrs, ok := genericError.Model().([]fabricv4.Error); ok {
				// EQ-3142509 = Connection already deleted
				if equinix_errors.HasErrorCode(fabricErrs, "EQ-3142509") {
					return diags
				}
			}
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	deleteTimeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)
	if err = WaitForDeletion(routeFilterId, d.Id(), meta, d, ctx, deleteTimeout); err != nil {
		return diag.Errorf("error waiting for route filter rule (%s) on route filter (%s) to be deleted: %s", d.Id(), routeFilterId, err)
	}
	return diags
}

func waitForStability(routeFilterId, ruleId string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for route filter rule %s on route filter %s to be stable", d.Id(), routeFilterId)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.ROUTEFILTERRULESTATE_PROVISIONING),
			string(fabricv4.ROUTEFILTERRULESTATE_REPROVISIONING),
		},
		Target: []string{
			string(fabricv4.ROUTEFILTERRULESTATE_PROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
			routeFilterRule, _, err := client.RouteFilterRulesApi.GetRouteFilterRuleByUuid(ctx, routeFilterId, ruleId).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return routeFilterRule, string(routeFilterRule.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func WaitForDeletion(routeFilterId, ruleId string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for route filter rule %s on route filter %s to be deleted", d.Id(), routeFilterId)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.ROUTEFILTERRULESTATE_PROVISIONED),
			string(fabricv4.ROUTEFILTERRULESTATE_DEPROVISIONING),
		},
		Target: []string{
			string(fabricv4.ROUTEFILTERRULESTATE_DEPROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
			routeFilterRule, body, err := client.RouteFilterRulesApi.GetRouteFilterRuleByUuid(ctx, routeFilterId, ruleId).Execute()
			if err != nil {
				if body != nil && body.StatusCode >= 400 && body.StatusCode <= 499 {
					// Already deleted resource
					return routeFilterRule, string(fabricv4.ROUTEFILTERRULESTATE_DEPROVISIONED), nil
				}
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return routeFilterRule, string(routeFilterRule.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
