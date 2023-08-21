package equinix

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := notificationToFabric(schemaNotifications)
	schemaOrder := d.Get("order").(*schema.Set).List()
	order := orderToFabric(schemaOrder)
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
		Order:         &order,
		Location:      &location,
		Notifications: notifications,
		Package_:      &packages,
		Account:       &account,
		Project:       &project,
	}

	fg, _, err := client.CloudRoutersApi.CreateGateway(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fg.Uuid)

	if _, err = waitUntilFGIsProvisioned(d.Id(), meta, ctx); err != nil {
		return diag.Errorf("error waiting for FG (%s) to be created: %s", d.Id(), err)
	}

	return resourceCloudRouterRead(ctx, d, meta)
}

func resourceCloudRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	CloudRouter, _, err := client.CloudRoutersApi.GetGatewayByUuid(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Fabric Cloud Router %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	d.SetId(CloudRouter.Uuid)
	return setCloudRouterMap(d, CloudRouter)
}

func setCloudRouterMap(d *schema.ResourceData, fg v4.CloudRouter) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := setMap(d, map[string]interface{}{
		"name":          fg.Name,
		"href":          fg.Href,
		"type":          fg.Type_,
		"state":         fg.State,
		"package":       CloudRouterPackageToTerra(fg.Package_),
		"location":      locationFGToTerra(fg.Location),
		"change_log":    changeLogToTerra(fg.ChangeLog),
		"account":       accountFgToTerra(fg.Account),
		"notifications": notificationToTerra(fg.Notifications),
		"project":       projectToTerra(fg.Project),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceCloudRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	dbConn, err := waitUntilFGIsProvisioned(d.Id(), meta, ctx)
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
	_, res, err := client.CloudRoutersApi.UpdateGatewayByUuid(ctx, updates, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error response for the Fabric Cloud Router update, response %v, error %v", res, err))
	}
	updateFg := v4.CloudRouter{}
	updateFg, err = waitForFGUpdateCompletion(d.Id(), meta, ctx)

	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("errored while waiting for successful Fabric Cloud Router update, response %v, error %v", res, err))
	}

	d.SetId(updateFg.Uuid)
	return setCloudRouterMap(d, updateFg)
}

func waitForFGUpdateCompletion(uuid string, meta interface{}, ctx context.Context) (v4.CloudRouter, error) {
	log.Printf("Waiting for FG update to complete, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Target: []string{string(v4.PROVISIONED_CloudRouterAccessPointState)},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.CloudRoutersApi.GetGatewayByUuid(ctx, uuid)
			if err != nil {
				return "", "", err
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

func waitUntilFGIsProvisioned(uuid string, meta interface{}, ctx context.Context) (v4.CloudRouter, error) {
	log.Printf("Waiting for FG to be provisioned, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(v4.PROVISIONED_CloudRouterAccessPointState),
		},
		Target: []string{
			string(v4.PENDING_INTERFACE_CONFIGURATION_EquinixStatus),
			string(v4.PROVISIONED_CloudRouterAccessPointState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.CloudRoutersApi.GetGatewayByUuid(ctx, uuid)
			if err != nil {
				return "", "", err
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
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	resp, err := client.CloudRoutersApi.DeleteGatewayByUuid(ctx, d.Id())
	if err != nil {
		errors, ok := err.(v4.GenericSwaggerError).Model().([]v4.ModelError)
		if ok {
			// EQ-3040055 = There is an existing update in REQUESTED state
			if hasModelErrorCode(errors, "EQ-3040055") {
				return diags
			}
		}
		return diag.FromErr(fmt.Errorf("error response for the Fabric Cloud Router delete. Error %v and response %v", err, resp))
	}

	err = waitUntilFGDeprovisioned(d.Id(), meta, ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func waitUntilFGDeprovisioned(uuid string, meta interface{}, ctx context.Context) error {
	log.Printf("Waiting for Fabric Cloud Router to be deprovisioned, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(v4.DEPROVISIONING_CloudRouterAccessPointState),
		},
		Target: []string{
			string(v4.DEPROVISIONED_CloudRouterAccessPointState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.CloudRoutersApi.GetGatewayByUuid(ctx, uuid)
			if err != nil {
				return "", "", err
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
