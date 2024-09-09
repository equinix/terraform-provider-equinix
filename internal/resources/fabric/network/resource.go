package network

import (
	"context"
	"log"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"

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
		ReadContext:   resourceFabricNetworkRead,
		CreateContext: resourceFabricNetworkCreate,
		UpdateContext: resourceFabricNetworkUpdate,
		DeleteContext: resourceFabricNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: fabricNetworkResourceSchema(),

		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric Network

Additional documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/Fabric/IMPLEMENTATION/fabric-networks-implement.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#fabric-networks`,
	}
}

func resourceFabricNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	createRequest := fabricv4.NetworkPostRequest{}
	createRequest.SetName(d.Get("name").(string))

	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := equinix_fabric_schema.NotificationsTerraformToGo(schemaNotifications)
	createRequest.SetNotifications(notifications)

	if schemaLocation, ok := d.GetOk("location"); ok {
		location := equinix_fabric_schema.LocationTerraformToGo(schemaLocation.(*schema.Set).List())
		createRequest.SetLocation(location)
	}

	schemaProject := d.Get("project").(*schema.Set).List()
	project := equinix_fabric_schema.ProjectTerraformToGo(schemaProject)
	createRequest.SetProject(project)

	netType := fabricv4.NetworkType(d.Get("type").(string))
	createRequest.SetType(netType)

	netScope := fabricv4.NetworkScope(d.Get("scope").(string))
	createRequest.SetScope(netScope)

	start := time.Now()
	fabricNetwork, _, err := client.NetworksApi.CreateNetwork(ctx).NetworkPostRequest(createRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(fabricNetwork.Uuid)

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if _, err = waitUntilFabricNetworkIsProvisioned(d.Id(), meta, d, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for Network (%s) to be created: %s", d.Id(), err)
	}

	return resourceFabricNetworkRead(ctx, d, meta)
}

func resourceFabricNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	fabricNetwork, _, err := client.NetworksApi.GetNetworkByUuid(ctx, d.Id()).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(fabricNetwork.Uuid)
	return setFabricNetworkMap(d, fabricNetwork)
}

func resourceFabricNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	start := time.Now()
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	dbConn, waitUntilNetworkProvisionedErr := waitUntilFabricNetworkIsProvisioned(d.Id(), meta, d, ctx, updateTimeout)
	if waitUntilNetworkProvisionedErr != nil {
		return diag.Errorf("either timed out or errored out while fetching Fabric Network for uuid %s and error %v", d.Id(), waitUntilNetworkProvisionedErr)
	}
	update, getUpdateRequestErr := getFabricNetworkUpdateRequest(dbConn, d)
	if getUpdateRequestErr != nil {
		return diag.Errorf("error retrieving intended updates from network config: %v", getUpdateRequestErr)
	}
	updates := []fabricv4.NetworkChangeOperation{update}
	_, _, err := client.NetworksApi.UpdateNetworkByUuid(ctx, d.Id()).NetworkChangeOperation(updates).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	updateTimeout = d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	updateNetwork, waitForNetworkUpdateCompletionErr := waitForFabricNetworkUpdateCompletion(d.Id(), meta, d, ctx, updateTimeout)

	if waitForNetworkUpdateCompletionErr != nil {
		return diag.Errorf("errored while waiting for successful Fabric Network update error %v", waitForNetworkUpdateCompletionErr)
	}

	d.SetId(updateNetwork.Uuid)
	return setFabricNetworkMap(d, updateNetwork)
}

func waitForFabricNetworkUpdateCompletion(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.Network, error) {
	log.Printf("Waiting for Network update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{string(fabricv4.NETWORKEQUINIXSTATUS_PROVISIONED)},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.NetworksApi.GetNetworkByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.Operation.EquinixStatus), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbConn *fabricv4.Network

	if err == nil {
		dbConn = inter.(*fabricv4.Network)
	}
	return dbConn, err
}

func waitUntilFabricNetworkIsProvisioned(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.Network, error) {
	log.Printf("Waiting for Fabric Network to be provisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.NETWORKEQUINIXSTATUS_PROVISIONING),
		},
		Target: []string{
			string(fabricv4.NETWORKEQUINIXSTATUS_PROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.NetworksApi.GetNetworkByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.Operation.EquinixStatus), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbConn *fabricv4.Network

	if err == nil {
		dbConn = inter.(*fabricv4.Network)
	}
	return dbConn, err
}

func resourceFabricNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	start := time.Now()
	_, _, err := client.NetworksApi.DeleteNetworkByUuid(ctx, d.Id()).Execute()
	if err != nil {
		if genericError, ok := err.(*fabricv4.GenericOpenAPIError); ok {
			if fabricErrs, ok := genericError.Model().([]fabricv4.Error); ok {
				// EQ-3040055 = There is an existing update in REQUESTED state
				if equinix_errors.HasErrorCode(fabricErrs, "EQ-3040055") {
					return diags
				}
			}
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	deleteTimeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)
	err = WaitUntilFabricNetworkDeprovisioned(d.Id(), meta, d, ctx, deleteTimeout)
	if err != nil {
		return diag.Errorf("API call failed while waiting for resource deletion. Error %v", err)
	}
	return diags
}

func WaitUntilFabricNetworkDeprovisioned(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for Fabric Network to be deprovisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.NETWORKEQUINIXSTATUS_DEPROVISIONING),
		},
		Target: []string{
			string(fabricv4.NETWORKEQUINIXSTATUS_DEPROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.NetworksApi.GetNetworkByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.Operation.EquinixStatus), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
