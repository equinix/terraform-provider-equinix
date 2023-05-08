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

func resourceFabricGateway() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Read:   schema.DefaultTimeout(6 * time.Minute),
		},
		ReadContext:   resourceFabricGatewayRead,
		CreateContext: resourceFabricGatewayCreate,
		UpdateContext: resourceFabricGatewayUpdate,
		DeleteContext: resourceFabricGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: createFabricGatewayResourceSchema(),

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric Gateway\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
	}
}

func resourceFabricGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := notificationToFabric(schemaNotifications)
	schemaOrder := d.Get("order").(*schema.Set).List()
	order := orderToFabric(schemaOrder)
	schemaAccount := d.Get("account").(*schema.Set).List()
	account := accountToFabricGateway(schemaAccount)
	schemaLocation := d.Get("location").(*schema.Set).List()
	location := locationToFabricGateway(schemaLocation)
	project := v4.Project{}
	schemaProject := d.Get("project").(*schema.Set).List()
	if len(schemaProject) != 0 {
		project = projectToFabricGateway(schemaProject)
	}
	schemaPackage := d.Get("package").(*schema.Set).List()
	packages := packageToFabricGateway(schemaPackage)

	createRequest := v4.FabricGatewayPostRequest{
		Name:          d.Get("name").(string),
		Type_:         d.Get("type").(string),
		Order:         &order,
		Location:      &location,
		Notifications: notifications,
		Package_:      &packages,
		Account:       &account,
		Project:       &project,
	}

	fg, _, err := client.GatewaysApi.CreateGateway(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fg.Uuid)

	if _, err = waitUntilFGIsProvisioned(d.Id(), meta, ctx); err != nil {
		return diag.Errorf("error waiting for FG (%s) to be created: %s", d.Id(), err)
	}

	return resourceFabricGatewayRead(ctx, d, meta)
}

func resourceFabricGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	fabricGateway, _, err := client.GatewaysApi.GetGatewayByUuid(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Fabric Gateway %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	d.SetId(fabricGateway.Uuid)
	return setFabricGatewayMap(d, fabricGateway)
}

func setFabricGatewayMap(d *schema.ResourceData, fg v4.FabricGateway) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := setMap(d, map[string]interface{}{
		"name":          fg.Name,
		"href":          fg.Href,
		"type":          fg.Type_,
		"state":         fg.State,
		"package":       fabricGatewayPackageToTerra(fg.Package_),
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

func resourceFabricGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	dbConn, err := waitUntilFGIsProvisioned(d.Id(), meta, ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.Errorf("either timed out or errored out while fetching Fabric Gateway for uuid %s and error %v", d.Id(), err)
	}
	// TO-DO
	update, err := getFabricGatewayUpdateRequest(dbConn, d)
	if err != nil {
		return diag.FromErr(err)
	}
	updates := []v4.FabricGatewayChangeOperation{update}
	_, res, err := client.GatewaysApi.UpdateGatewayByUuid(ctx, updates, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error response for the Fabric Gateway update, response %v, error %v", res, err))
	}
	updateFg := v4.FabricGateway{}
	updateFg, err = waitForFGUpdateCompletion(d.Id(), meta, ctx)

	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("errored while waiting for successful Fabric Gateway update, response %v, error %v", res, err))
	}

	d.SetId(updateFg.Uuid)
	return setFabricGatewayMap(d, updateFg)
}

func waitForFGUpdateCompletion(uuid string, meta interface{}, ctx context.Context) (v4.FabricGateway, error) {
	log.Printf("Waiting for FG update to complete, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Target: []string{string(v4.PROVISIONED_FabricGatewayAccessPointState)},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.GatewaysApi.GetGatewayByUuid(ctx, uuid)
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
	dbConn := v4.FabricGateway{}

	if err == nil {
		dbConn = inter.(v4.FabricGateway)
	}
	return dbConn, err
}

func waitUntilFGIsProvisioned(uuid string, meta interface{}, ctx context.Context) (v4.FabricGateway, error) {
	log.Printf("Waiting for FG to be provisioned, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(v4.PROVISIONING_FabricGatewayAccessPointState),
		},
		Target: []string{
			string(v4.PENDING_INTERFACE_CONFIGURATION_EquinixStatus),
			string(v4.PROVISIONED_FabricGatewayAccessPointState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.GatewaysApi.GetGatewayByUuid(ctx, uuid)
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
	dbConn := v4.FabricGateway{}

	if err == nil {
		dbConn = inter.(v4.FabricGateway)
	}
	return dbConn, err
}

func resourceFabricGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	resp, err := client.GatewaysApi.DeleteGatewayByUuid(ctx, d.Id())
	if err != nil {
		errors, ok := err.(v4.GenericSwaggerError).Model().([]v4.ModelError)
		if ok {
			// EQ-3040055 = There is an existing update in REQUESTED state
			if hasModelErrorCode(errors, "EQ-3040055") {
				return diags
			}
		}
		return diag.FromErr(fmt.Errorf("error response for the Fabric Gateway delete. Error %v and response %v", err, resp))
	}

	err = waitUntilFGDeprovisioned(d.Id(), meta, ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func waitUntilFGDeprovisioned(uuid string, meta interface{}, ctx context.Context) error {
	log.Printf("Waiting for Fabric Gateway to be deprovisioned, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(v4.DEPROVISIONING_FabricGatewayAccessPointState),
		},
		Target: []string{
			string(v4.DEPROVISIONED_FabricGatewayAccessPointState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.GatewaysApi.GetGatewayByUuid(ctx, uuid)
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
