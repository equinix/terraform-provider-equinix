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

// Resource returns the schema.Resource for managing route filter policy attachments to Fabric connections.
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
* Getting Started: https://docs.equinix.com/fabric-cloud-router/bgp/fcr-route-filters/
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Route-Filters`,
	}
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	connectionID := d.Get("connection_id").(string)
	connectionRouteFilter, _, err := client.RouteFiltersApi.GetConnectionRouteFilterByUuid(ctx, d.Id(), connectionID).Execute()
	if err != nil {
		log.Printf("[WARN] Route Filter Policy %s not found on Connection %s, error %s", d.Id(), connectionID, err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(connectionRouteFilter.GetUuid())
	return setConnectionRouteFilterMap(d, connectionRouteFilter)
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	connectionID := d.Get("connection_id").(string)
	routeFilterID := d.Get("route_filter_id").(string)
	direction := d.Get("direction").(string)

	start := time.Now()
	routeFilter, _, err := client.RouteFiltersApi.
		AttachConnectionRouteFilter(ctx, routeFilterID, connectionID).
		ConnectionRouteFiltersBase(
			fabricv4.ConnectionRouteFiltersBase{
				Direction: fabricv4.ConnectionRouteFiltersBaseDirection(direction),
			},
		).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	if err = d.Set("connection_id", connectionID); err != nil {
		return diag.Errorf("error setting connection_id to state %s", err)
	}
	d.SetId(routeFilter.GetUuid())

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if err = waitForStability(ctx, connectionID, d.Id(), meta, d, createTimeout); err != nil {
		return diag.Errorf("error waiting for route filter (%s) to be attached to connection (%s): %s", d.Id(), connectionID, err)
	}

	return resourceRead(ctx, d, meta)
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	connectionID := d.Get("connection_id").(string)
	routeFilterID := d.Get("route_filter_id").(string)
	oldDirection, newDirection := d.GetChange("direction")
	if oldDirection.(string) == newDirection.(string) {
		return diag.Diagnostics{}
	}

	start := time.Now()
	connectionRouteFilter, _, err := client.RouteFiltersApi.
		AttachConnectionRouteFilter(ctx, routeFilterID, connectionID).
		ConnectionRouteFiltersBase(
			fabricv4.ConnectionRouteFiltersBase{
				Direction: fabricv4.ConnectionRouteFiltersBaseDirection(newDirection.(string)),
			},
		).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	if err = waitForStability(ctx, routeFilterID, d.Id(), meta, d, updateTimeout); err != nil {
		return diag.Errorf("error waiting for route filter policy (%s) on connection (%s) to be updated: %s", routeFilterID, connectionID, err)
	}

	return setConnectionRouteFilterMap(d, connectionRouteFilter)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	connectionID := d.Get("connection_id").(string)

	start := time.Now()
	_, _, err := client.RouteFiltersApi.DetachConnectionRouteFilter(ctx, d.Id(), connectionID).Execute()
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
	if err = WaitForDeletion(ctx, connectionID, d.Id(), meta, d, deleteTimeout); err != nil {
		return diag.Errorf("error waiting for route filter (%s) to be detached from connection (%s): %s", d.Id(), connectionID, err)
	}
	return diags
}

func waitForStability(ctx context.Context, connectionID, routeFilterID string, meta interface{}, d *schema.ResourceData, timeout time.Duration) error {
	log.Printf("Waiting for route filter policy (%x) attachment to connection (%s) to be stable", connectionID, routeFilterID)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_ATTACHING),
		},
		Target: []string{
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_ATTACHED),
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_PENDING_BGP_CONFIGURATION),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
			connectionRouteFilter, _, err := client.RouteFiltersApi.GetConnectionRouteFilterByUuid(ctx, routeFilterID, connectionID).Execute()
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

// WaitForDeletion waits until the route filter policy is detached from the connection.
func WaitForDeletion(ctx context.Context, connectionID, routeFilterID string, meta interface{}, d *schema.ResourceData, timeout time.Duration) error {
	log.Printf("Waiting for route filter policy (%s) to be detached from connection (%s)", routeFilterID, connectionID)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_ATTACHING),
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_DETACHING),
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_PENDING_BGP_CONFIGURATION),
		},
		Target: []string{
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_DETACHED),
			string(fabricv4.ROUTEFILTERSTATE_DEPROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
			connectionRouteFilter, body, err := client.RouteFiltersApi.GetConnectionRouteFilterByUuid(ctx, routeFilterID, connectionID).Execute()
			if err != nil {
				if body != nil && body.StatusCode >= 400 && body.StatusCode <= 499 {
					// Already deleted resource - return DEPROVISIONED state
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
