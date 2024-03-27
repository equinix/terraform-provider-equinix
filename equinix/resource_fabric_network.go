package equinix

import (
	"context"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"log"
	"time"
)

func fabricNetworkChangeSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Absolute URL that returns the details of the given change.\nExample: https://api.equinix.com/fabric/v4/networks/92dc376a-a932-43aa-a6a2-c806dedbd784",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Asset change request identifier.",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Asset instance change request type.: NETWORK_CREATION, NETWORK_UPDATE, NETWORK_DELETION",
		},
	}
}
func fabricNetworkOperationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"equinix_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Progress towards provisioning a given asset.",
		},
	}
}
func fabricNetworkProjectSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Customer project identifier",
		},
	}
}
func fabricNetworkResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Network URI information",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 24),
			Description:  "Fabric Network name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned network identifier",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Network overall state",
		},
		"scope": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Fabric Network scope",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"IPWAN", "EPLAN", "EVPLAN"}, true),
			Description:  "Supported Network types - EVPLAN, EPLAN, IPWAN",
		},
		"location": {
			Type:        schema.TypeSet,
			Computed:    true,
			Optional:    true,
			Description: "Fabric Network location",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.LocationSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Fabric Network project",
			Elem: &schema.Resource{
				Schema: fabricNetworkProjectSch(),
			},
		},
		"operation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Network operation information that is associated with this Fabric Network",
			Elem: &schema.Resource{
				Schema: fabricNetworkOperationSch(),
			},
		},
		"change": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Information on asset change operation",
			Elem: &schema.Resource{
				Schema: fabricNetworkChangeSch(),
			},
		},
		"notifications": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Preferences for notifications on Fabric Network configuration or status changes",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.NotificationSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "A permanent record of asset creation, modification, or deletion",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ChangeLogSch(),
			},
		},
		"connections_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of connections associated with this network",
		},
	}
}
func resourceFabricNetwork() *schema.Resource {
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

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric Network",
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

func fabricNetworkOperationGoToTerraform(operation *fabricv4.NetworkOperation) *schema.Set {
	if operation == nil {
		return nil
	}
	mappedOperation := make(map[string]interface{})
	mappedOperation["equinix_status"] = string(*operation.EquinixStatus)

	operationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: fabricNetworkOperationSch()}),
		[]interface{}{mappedOperation},
	)
	return operationSet
}
func simplifiedFabricNetworkChangeGoToTerraform(networkChange *fabricv4.SimplifiedNetworkChange) *schema.Set {

	mappedChange := make(map[string]interface{})
	mappedChange["href"] = networkChange.GetHref()
	mappedChange["type"] = string(networkChange.GetType())
	mappedChange["uuid"] = networkChange.GetUuid()

	changeSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: fabricNetworkChangeSch()}),
		[]interface{}{mappedChange},
	)
	return changeSet
}

func setFabricNetworkMap(d *schema.ResourceData, nt *fabricv4.Network) diag.Diagnostics {
	diags := diag.Diagnostics{}
	operation := nt.GetOperation()
	change := nt.GetChange()
	location := nt.GetLocation()
	notifications := nt.GetNotifications()
	project := nt.GetProject()
	changeLog := nt.GetChangeLog()
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"name":              nt.GetName(),
		"href":              nt.GetHref(),
		"uuid":              nt.GetUuid(),
		"type":              string(nt.GetType()),
		"scope":             string(nt.GetScope()),
		"state":             string(nt.GetState()),
		"operation":         fabricNetworkOperationGoToTerraform(&operation),
		"change":            simplifiedFabricNetworkChangeGoToTerraform(&change),
		"location":          equinix_fabric_schema.LocationGoToTerraform(&location),
		"notifications":     equinix_fabric_schema.NotificationsGoToTerraform(notifications),
		"project":           equinix_fabric_schema.ProjectGoToTerraform(&project),
		"change_log":        equinix_fabric_schema.ChangeLogGoToTerraform(&changeLog),
		"connections_count": nt.GetConnectionsCount(),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func getFabricNetworkUpdateRequest(network *fabricv4.Network, d *schema.ResourceData) (fabricv4.NetworkChangeOperation, error) {
	changeOps := fabricv4.NetworkChangeOperation{}
	existingName := network.GetName()
	updateNameVal := d.Get("name").(string)

	log.Printf("existing name %s, Update Name Request %s ", existingName, updateNameVal)

	if existingName != updateNameVal {
		changeOps = fabricv4.NetworkChangeOperation{Op: "replace", Path: "/name", Value: updateNameVal}
	} else {
		return changeOps, fmt.Errorf("nothing to update for the Fabric Network: %s", existingName)
	}
	return changeOps, nil
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
	stateConf := &resource.StateChangeConf{
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
	stateConf := &resource.StateChangeConf{
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
	stateConf := &resource.StateChangeConf{
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
