package route_filter

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
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric Route Filter Policy

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/FCR/FCR-route-filters.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#route-filters`,
	}
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	routeFilter, _, err := client.RouteFiltersApi.GetRouteFilterByUuid(ctx, d.Id()).Execute()
	if err != nil {
		log.Printf("[WARN] Route Filter %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(routeFilter.GetUuid())
	return setRouteFilterMap(d, routeFilter)
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	createRequest := buildCreateRequest(d)

	start := time.Now()
	routeFilter, _, err := client.RouteFiltersApi.CreateRouteFilter(ctx).RouteFiltersBase(createRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(routeFilter.GetUuid())

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if err = waitForStability(d.Id(), meta, d, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for route filter (%s) to be created: %s", d.Id(), err)
	}

	return resourceRead(ctx, d, meta)
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	updateRequest := buildUpdateRequest(d)

	start := time.Now()
	routeFilter, _, err := client.RouteFiltersApi.PatchRouteFilterByUuid(ctx, d.Id()).RouteFiltersPatchRequestItem(updateRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	if err = waitForStability(d.Id(), meta, d, ctx, updateTimeout); err != nil {
		return diag.Errorf("error waiting for route filter (%s) to be updated: %s", d.Id(), err)
	}

	return setRouteFilterMap(d, routeFilter)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).NewFabricClientForSDK(d)

	start := time.Now()
	_, _, err := client.RouteFiltersApi.DeleteRouteFilterByUuid(ctx, d.Id()).Execute()
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
	if err = WaitForDeletion(d.Id(), meta, d, ctx, deleteTimeout); err != nil {
		return diag.Errorf("error waiting for route filter (%s) to be deleted: %s", d.Id(), err)
	}
	return diags
}

func waitForStability(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for route filter to be stable, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.ROUTEFILTERSTATE_PROVISIONING),
			string(fabricv4.ROUTEFILTERSTATE_REPROVISIONING),
		},
		Target: []string{
			string(fabricv4.ROUTEFILTERSTATE_PROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			routeFilter, _, err := client.RouteFiltersApi.GetRouteFilterByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return routeFilter, string(routeFilter.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func WaitForDeletion(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for route filter to be deleted, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.ROUTEFILTERSTATE_PROVISIONED),
			string(fabricv4.ROUTEFILTERSTATE_DEPROVISIONING),
		},
		Target: []string{
			string(fabricv4.ROUTEFILTERSTATE_DEPROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			routeFilter, body, err := client.RouteFiltersApi.GetRouteFilterByUuid(ctx, uuid).Execute()
			if err != nil {
				if body.StatusCode >= 400 && body.StatusCode <= 499 {
					// Already deleted resource
					return routeFilter, string(fabricv4.ROUTEFILTERSTATE_DEPROVISIONED), nil
				}
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return routeFilter, string(routeFilter.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
