package equinix

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudRouter() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Read:   schema.DefaultTimeout(6 * time.Minute),
		},
		ReadContext:   resourceCloudRouterRead,
		CreateContext: resourceCloudRouterCreate,
		UpdateContext: resourceCloudRouterUpdate,
		DeleteContext: resourceCloudRouterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: createCloudRouterResourceSchema(),

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric Cloud Router\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
	}
}

func resourceCloudRouterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := notificationToFabric(schemaNotifications)
	schemaAccount := d.Get("account").(*schema.Set).List()
	account := accountToCloudRouter(schemaAccount)
	schemaLocation := d.Get("location").(*schema.Set).List()
	location := locationToCloudRouter(schemaLocation)
	project := v4.Project{}
	schemaProject := d.Get("project").(*schema.Set).List()
	if len(schemaProject) != 0 {
		project = projectToCloudRouter(schemaProject)
	}
	schemaPackage := d.Get("package").(*schema.Set).List()
	packages := packageToCloudRouter(schemaPackage)

	createRequest := v4.CloudRouterPostRequest{
		Name:          d.Get("name").(string),
		Type_:         d.Get("type").(string),
		Location:      &location,
		Notifications: notifications,
		Package_:      &packages,
		Account:       &account,
		Project:       &project,
	}

	if v, ok := d.GetOk("order"); ok {
		order := orderToFabric(v.(*schema.Set).List())
		createRequest.Order = &order
	}

	fcr, _, err := client.CloudRoutersApi.CreateCloudRouter(ctx, createRequest)
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(fcr.Uuid)

	if _, err = waitUntilCloudRouterIsProvisioned(d.Id(), meta, ctx); err != nil {
		return diag.Errorf("error waiting for Cloud Router (%s) to be created: %s", d.Id(), err)
	}

	return resourceCloudRouterRead(ctx, d, meta)
}

func resourceCloudRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	CloudRouter, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Fabric Cloud Router %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(CloudRouter.Uuid)
	return setCloudRouterMap(d, CloudRouter)
}

func setCloudRouterMap(d *schema.ResourceData, fcr v4.CloudRouter) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"name":                  fcr.Name,
		"href":                  fcr.Href,
		"type":                  fcr.Type_,
		"state":                 fcr.State,
		"package":               cloudRouterPackageToTerra(fcr.Package_),
		"location":              locationCloudRouterToTerra(fcr.Location),
		"change_log":            changeLogToTerra(fcr.ChangeLog),
		"account":               accountCloudRouterToTerra(fcr.Account),
		"notifications":         notificationToTerra(fcr.Notifications),
		"project":               projectToTerra(fcr.Project),
		"equinix_asn":           fcr.EquinixAsn,
		"bgp_ipv4_routes_count": fcr.BgpIpv4RoutesCount,
		"bgp_ipv6_routes_count": fcr.BgpIpv6RoutesCount,
		"connections_count":     fcr.ConnectionsCount,
		"order":                 orderToTerra(fcr.Order),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceCloudRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	dbConn, err := waitUntilCloudRouterIsProvisioned(d.Id(), meta, ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.Errorf("either timed out or errored out while fetching Fabric Cloud Router for uuid %s and error %v", d.Id(), err)
	}
	// TO-DO
	update, err := getCloudRouterUpdateRequest(dbConn, d)
	if err != nil {
		return diag.FromErr(err)
	}
	updates := []v4.CloudRouterChangeOperation{update}
	_, _, err = client.CloudRoutersApi.UpdateCloudRouterByUuid(ctx, updates, d.Id())
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	updateFg := v4.CloudRouter{}
	updateFg, err = waitForCloudRouterUpdateCompletion(d.Id(), meta, ctx)

	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("errored while waiting for successful Fabric Cloud Router update, error %v", err))
	}

	d.SetId(updateFg.Uuid)
	return setCloudRouterMap(d, updateFg)
}

func waitForCloudRouterUpdateCompletion(uuid string, meta interface{}, ctx context.Context) (v4.CloudRouter, error) {
	log.Printf("Waiting for Cloud Router update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{string(v4.PROVISIONED_CloudRouterAccessPointState)},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, uuid)
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.State), nil
		},
		Timeout:    2 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	dbConn := v4.CloudRouter{}

	if err == nil {
		dbConn = inter.(v4.CloudRouter)
	}
	return dbConn, err
}

func waitUntilCloudRouterIsProvisioned(uuid string, meta interface{}, ctx context.Context) (v4.CloudRouter, error) {
	log.Printf("Waiting for Cloud Router to be provisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(v4.PROVISIONING_CloudRouterAccessPointState),
		},
		Target: []string{
			string(v4.PROVISIONED_CloudRouterAccessPointState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, uuid)
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.State), nil
		},
		Timeout:    5 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	dbConn := v4.CloudRouter{}

	if err == nil {
		dbConn = inter.(v4.CloudRouter)
	}
	return dbConn, err
}

func resourceCloudRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	_, err := client.CloudRoutersApi.DeleteCloudRouterByUuid(ctx, d.Id())
	if err != nil {
		errors, ok := err.(v4.GenericSwaggerError).Model().([]v4.ModelError)
		if ok {
			// EQ-3040055 = There is an existing update in REQUESTED state
			if equinix_errors.HasModelErrorCode(errors, "EQ-3040055") {
				return diags
			}
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	err = waitUntilCloudRouterDeprovisioned(d.Id(), meta, ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func waitUntilCloudRouterDeprovisioned(uuid string, meta interface{}, ctx context.Context) error {
	log.Printf("Waiting for Fabric Cloud Router to be deprovisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(v4.DEPROVISIONING_CloudRouterAccessPointState),
		},
		Target: []string{
			string(v4.DEPROVISIONED_CloudRouterAccessPointState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, uuid)
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.State), nil
		},
		Timeout:    5 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
