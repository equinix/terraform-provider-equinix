package equinix

import (
	"context"
	"fmt"
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
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
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := equinix_fabric_schema.NotificationsToFabric(schemaNotifications)
	schemaLocation := d.Get("location").(*schema.Set).List()
	location := equinix_fabric_schema.LocationToFabric(schemaLocation)
	schemaProject := d.Get("project").(*schema.Set).List()
	project := equinix_fabric_schema.ProjectToFabric(schemaProject)
	netType := v4.NetworkType(d.Get("type").(string))
	netScope := v4.NetworkScope(d.Get("scope").(string))

	createRequest := v4.NetworkPostRequest{
		Name:          d.Get("name").(string),
		Type_:         &netType,
		Scope:         &netScope,
		Location:      &location,
		Notifications: notifications,
		Project:       project,
	}

	start := time.Now()
	fabricNetwork, _, err := client.NetworksApi.CreateNetwork(ctx, createRequest)
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(fabricNetwork.Uuid)

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if _, err = waitUntilFabricNetworkIsProvisioned(d.Id(), meta, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for Network (%s) to be created: %s", d.Id(), err)
	}

	return resourceFabricNetworkRead(ctx, d, meta)
}

func resourceFabricNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	fabricNetwork, _, err := client.NetworksApi.GetNetworkByUuid(ctx, d.Id())
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(fabricNetwork.Uuid)
	return setFabricNetworkMap(d, fabricNetwork)
}
func fabricNetworkOperationToTerra(operation *v4.NetworkOperation) *schema.Set {
	if operation == nil {
		return nil
	}
	operations := []*v4.NetworkOperation{operation}
	mappedOperations := make([]interface{}, len(operations))
	for _, operation := range operations {
		mappedOperation := make(map[string]interface{})
		mappedOperation["equinix_status"] = string(*operation.EquinixStatus)
		mappedOperations = append(mappedOperations, mappedOperation)
	}

	operationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: fabricNetworkOperationSch()}),
		mappedOperations,
	)
	return operationSet
}
func simplifiedFabricNetworkChangeToTerra(networkChange *v4.SimplifiedNetworkChange) *schema.Set {
	changes := []*v4.SimplifiedNetworkChange{networkChange}
	mappedChanges := make([]interface{}, len(changes))
	for _, change := range changes {
		mappedChange := make(map[string]interface{})
		mappedChange["href"] = change.Href
		mappedChange["type"] = string(*change.Type_)
		mappedChange["uuid"] = change.Uuid
		mappedChanges = append(mappedChanges, mappedChange)
	}

	changeSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: fabricNetworkChangeSch()}),
		mappedChanges,
	)
	return changeSet
}

func setFabricNetworkMap(d *schema.ResourceData, nt v4.Network) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"name":              nt.Name,
		"href":              nt.Href,
		"uuid":              nt.Uuid,
		"type":              nt.Type_,
		"scope":             nt.Scope,
		"state":             nt.State,
		"operation":         fabricNetworkOperationToTerra(nt.Operation),
		"change":            simplifiedFabricNetworkChangeToTerra(nt.Change),
		"location":          equinix_fabric_schema.LocationToTerra(nt.Location),
		"notifications":     equinix_fabric_schema.NotificationsToTerra(nt.Notifications),
		"project":           equinix_fabric_schema.ProjectToTerra(nt.Project),
		"change_log":        equinix_fabric_schema.ChangeLogToTerra(nt.ChangeLog),
		"connections_count": nt.ConnectionsCount,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func getFabricNetworkUpdateRequest(network v4.Network, d *schema.ResourceData) (v4.NetworkChangeOperation, error) {
	changeOps := v4.NetworkChangeOperation{}
	existingName := network.Name
	updateNameVal := d.Get("name")

	log.Printf("existing name %s, Update Name Request %s ", existingName, updateNameVal)

	if existingName != updateNameVal {
		changeOps = v4.NetworkChangeOperation{Op: "replace", Path: "/name", Value: &updateNameVal}
	} else {
		return changeOps, fmt.Errorf("nothing to update for the Fabric Network: %s", existingName)
	}
	return changeOps, nil
}
func resourceFabricNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	start := time.Now()
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	dbConn, err := waitUntilFabricNetworkIsProvisioned(d.Id(), meta, ctx, updateTimeout)
	if err != nil {
		return diag.Errorf("either timed out or errored out while fetching Fabric Network for uuid %s and error %v", d.Id(), err)
	}
	update, err := getFabricNetworkUpdateRequest(dbConn, d)
	if err != nil {
		return diag.Errorf("error retrieving intended updates from network config: %v", err)
	}
	updates := []v4.NetworkChangeOperation{update}
	_, res, err := client.NetworksApi.UpdateNetworkByUuid(ctx, updates, d.Id())
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	updateFg := v4.Network{}
	updateTimeout = d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	updateFg, err = waitForFabricNetworkUpdateCompletion(d.Id(), meta, ctx, updateTimeout)

	if err != nil {
		return diag.Errorf("errored while waiting for successful Fabric Network update, response %v, error %v", res, err)
	}

	d.SetId(updateFg.Uuid)
	return setFabricNetworkMap(d, updateFg)
}

func waitForFabricNetworkUpdateCompletion(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) (v4.Network, error) {
	log.Printf("Waiting for Network update to complete, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Target: []string{string(v4.PROVISIONED_NetworkEquinixStatus)},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.NetworksApi.GetNetworkByUuid(ctx, uuid)
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
	dbConn := v4.Network{}

	if err == nil {
		dbConn = inter.(v4.Network)
	}
	return dbConn, err
}

func waitUntilFabricNetworkIsProvisioned(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) (v4.Network, error) {
	log.Printf("Waiting for Fabric Network to be provisioned, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(v4.PROVISIONING_NetworkEquinixStatus),
		},
		Target: []string{
			string(v4.PROVISIONED_NetworkEquinixStatus),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.NetworksApi.GetNetworkByUuid(ctx, uuid)
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
	dbConn := v4.Network{}

	if err == nil {
		dbConn = inter.(v4.Network)
	}
	return dbConn, err
}

func resourceFabricNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	start := time.Now()
	_, _, err := client.NetworksApi.DeleteNetworkByUuid(ctx, d.Id())
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

	deleteTimeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)
	err = WaitUntilFabricNetworkDeprovisioned(d.Id(), meta, ctx, deleteTimeout)
	if err != nil {
		return diag.Errorf("API call failed while waiting for resource deletion. Error %v", err)
	}
	return diags
}

func WaitUntilFabricNetworkDeprovisioned(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for Fabric Network to be deprovisioned, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(v4.DEPROVISIONING_NetworkEquinixStatus),
		},
		Target: []string{
			string(v4.DEPROVISIONED_NetworkEquinixStatus),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.NetworksApi.GetNetworkByUuid(ctx, uuid)
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
