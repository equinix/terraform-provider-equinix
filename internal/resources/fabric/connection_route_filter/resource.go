package connection_route_filter

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
		Description: `Fabric V4 API compatible resource allows attachment of Route Filter Polices to Fabric Connections

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/FCR/FCR-route-filters.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#route-filters`,
	}
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	connectionId := d.Get("connection_id").(string)
	connectionRouteFilter, _, err := client.RouteFiltersApi.GetConnectionRouteFilterByUuid(ctx, d.Id(), connectionId).Execute()
	if err != nil {
		log.Printf("[WARN] Route Filter Policy %s not found on Connection %s, error %s", d.Id(), connectionId, err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(connectionRouteFilter.GetUuid())
	return setConnectionRouteFilterMap(d, connectionRouteFilter)
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	connectionId := d.Get("connection_id").(string)
	routeFilterId := d.Get("route_filter_id").(string)
	direction := d.Get("direction").(string)

	start := time.Now()
	routeFilter, _, err := client.RouteFiltersApi.
		AttachConnectionRouteFilter(ctx, routeFilterId, connectionId).
		ConnectionRouteFiltersBase(
			fabricv4.ConnectionRouteFiltersBase{
				Direction: fabricv4.ConnectionRouteFiltersBaseDirection(direction),
			},
		).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	if err = d.Set("connection_id", connectionId); err != nil {
		return diag.Errorf("error setting connection_id to state %s", err)
	}
	d.SetId(routeFilter.GetUuid())

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if err = waitForStability(connectionId, d.Id(), meta, d, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for route filter (%s) to be attached to connection (%s): %s", d.Id(), connectionId, err)
	}

	return resourceRead(ctx, d, meta)
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	connectionId := d.Get("connection_id").(string)
	routeFilterId := d.Get("route_filter_id").(string)
	oldDirection, newDirection := d.GetChange("direction")
	if oldDirection.(string) == newDirection.(string) {
		return diag.Diagnostics{}
	}

	start := time.Now()
	connectionRouteFilter, _, err := client.RouteFiltersApi.
		AttachConnectionRouteFilter(ctx, routeFilterId, connectionId).
		ConnectionRouteFiltersBase(
			fabricv4.ConnectionRouteFiltersBase{
				Direction: fabricv4.ConnectionRouteFiltersBaseDirection(newDirection.(string)),
			},
		).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	if err = waitForStability(routeFilterId, d.Id(), meta, d, ctx, updateTimeout); err != nil {
		return diag.Errorf("error waiting for route filter policy (%s) on connection (%s) to be updated: %s", routeFilterId, connectionId, err)
	}

	return setConnectionRouteFilterMap(d, connectionRouteFilter)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	connectionId := d.Get("connection_id").(string)

	start := time.Now()
	_, _, err := client.RouteFiltersApi.DetachConnectionRouteFilter(ctx, d.Id(), connectionId).Execute()
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
	if err = WaitForDeletion(connectionId, d.Id(), meta, d, ctx, deleteTimeout); err != nil {
		return diag.Errorf("error waiting for route filter (%s) to be detached from connection (%s): %s", d.Id(), connectionId, err)
	}
	return diags
}

func waitForStability(connectionId, routeFilterId string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for route filter policy (%x) attachment to connection (%s) to be stable", connectionId, routeFilterId)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_ATTACHING),
		},
		Target: []string{
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_ATTACHED),
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_PENDING_BGP_CONFIGURATION),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			connectionRouteFilter, _, err := client.RouteFiltersApi.GetConnectionRouteFilterByUuid(ctx, routeFilterId, connectionId).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return connectionRouteFilter, string(connectionRouteFilter.GetAttachmentStatus()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func WaitForDeletion(connectionId, routeFilterId string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for route filter policy (%s) to be detached from connection (%s)", routeFilterId, connectionId)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_ATTACHING),
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_DETACHING),
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_PENDING_BGP_CONFIGURATION),
		},
		Target: []string{
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_DETACHED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			connectionRouteFilter, body, err := client.RouteFiltersApi.GetConnectionRouteFilterByUuid(ctx, routeFilterId, connectionId).Execute()
			if err != nil {
				if body.StatusCode >= 400 && body.StatusCode <= 499 {
					// Already deleted resource
					return connectionRouteFilter, string(fabricv4.ROUTEFILTERSTATE_DEPROVISIONED), nil
				}
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return connectionRouteFilter, string(connectionRouteFilter.GetAttachmentStatus()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
