package equinix

import (
	"context"
	"fmt"
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Read:   schema.DefaultTimeout(6 * time.Minute),
		},
		ReadContext:   resourceNetworkRead,
		CreateContext: resourceNetworkCreate,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: createNetworkResourceSchema(),

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric Network\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := notificationToFabric(schemaNotifications)
	schemaLocation := d.Get("location").(*schema.Set).List()
	location := locationToFabric(schemaLocation)
	project := v4.Project{}
	schemaProject := d.Get("project").(*schema.Set).List()
	if len(schemaProject) != 0 {
		project = projectToFabric(schemaProject)
	}

	createRequest := v4.NetworkPostRequest{
		Name:          d.Get("name").(string),
		Type_:         d.Get("type").(*v4.NetworkType),
		Scope:         d.Get("type").(*v4.NetworkScope),
		Location:      &location,
		Notifications: notifications,
		Project:       &project,
	}

	fabricNetwork, _, err := client.NetworksApi.CreateNetwork(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fabricNetwork.uuid)

	if _, err = waitUntilNetworkIsProvisioned(d.Id(), meta, ctx); err != nil {
		return diag.Errorf("error waiting for Cloud Router (%s) to be created: %s", d.Id(), err)
	}

	return resourceNetworkRead(ctx, d, meta)
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	fabricNetwork, _, err := client.NetworksApi.GetNetworkByUuid(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Fabric Network %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	d.SetId(fabricNetwork.Uuid)
	return setNetworkMap(d, fabricNetwork)
}

func setNetworkMap(d *schema.ResourceData, nt v4.Network) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := setMap(d, map[string]interface{}{
		"name":          nt.Name,
		"type":          nt.Type_,
		"scope":         nt.Scope,
		"location":      locationToTerra(nt.Location),
		"notifications": notificationToTerra(nt.Notifications),
		"project":       projectToTerra(nt.Project),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	dbConn, err := waitUntilNetworkIsProvisioned(d.Id(), meta, ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.Errorf("either timed out or errored out while fetching Fabric Network for uuid %s and error %v", d.Id(), err)
	}
	// TO-DO
	update, err := getNetworkUpdateRequest(dbConn, d)
	if err != nil {
		return diag.FromErr(err)
	}
	updates := []v4.NetworkChangeOperation{update}
	_, res, err := client.NetworksApi.UpdateNetworkByUuid(ctx, updates, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error response for the Fabric Network update, response %v, error %v", res, err))
	}
	updateFg := v4.Network{}
	updateFg, err = waitForNetworkUpdateCompletion(d.Id(), meta, ctx)

	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("errored while waiting for successful Fabric Network update, response %v, error %v", res, err))
	}

	d.SetId(updateFg.Uuid)
	return setNetworkMap(d, updateFg)
}

func waitForNetworkUpdateCompletion(uuid string, meta interface{}, ctx context.Context) (v4.Network, error) {
	log.Printf("Waiting for Cloud Router update to complete, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Target: []string{string(v4.PROVISIONED_NetworkEquinixStatus)},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.NetworksApi.GetNetworkByUuid(ctx, uuid)
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
	dbConn := v4.Network{}

	if err == nil {
		dbConn = inter.(v4.Network)
	}
	return dbConn, err
}

func waitUntilNetworkIsProvisioned(uuid string, meta interface{}, ctx context.Context) (v4.Network, error) {
	log.Printf("Waiting for Fabric Network to be provisioned, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(v4.PROVISIONING_NetworkEquinixStatus),
		},
		Target: []string{
			string(v4.PROVISIONED_NetworkEquinixStatus),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.NetworksApi.GetNetworkByUuid(ctx, uuid)
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
	dbConn := v4.Network{}

	if err == nil {
		dbConn = inter.(v4.Network)
	}
	return dbConn, err
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	_, resp, err := client.NetworksApi.DeleteNetworkByUuid(ctx, d.Id())
	if err != nil {
		errors, ok := err.(v4.GenericSwaggerError).Model().([]v4.ModelError)
		if ok {
			// EQ-3040055 = There is an existing update in REQUESTED state
			if hasModelErrorCode(errors, "EQ-3040055") {
				return diags
			}
		}
		return diag.FromErr(fmt.Errorf("error response for the Fabric Network delete. Error %v and response %v", err, resp))
	}

	err = waitUntilNetworkDeprovisioned(d.Id(), meta, ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func waitUntilNetworkDeprovisioned(uuid string, meta interface{}, ctx context.Context) error {
	log.Printf("Waiting for Fabric Network to be deprovisioned, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(v4.DEPROVISIONING_NetworkEquinixStatus),
		},
		Target: []string{
			string(v4.DEPROVISIONED_NetworkEquinixStatus),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.NetworksApi.GetNetworkByUuid(ctx, uuid)
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
