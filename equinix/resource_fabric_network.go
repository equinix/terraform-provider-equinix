package equinix

import (
	"context"
	"fmt"
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"
)

func FabricNetworkChangeSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "href",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "UUID of Network Change",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "network change type: NETWORK_CREATION, NETWORK_UPDATE, NETWORK_DELETION",
		},
	}
}
func FabricNetworkOperationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"equinix_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Network operation status",
		},
	}
}
func FabricNetworkProjectSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Project Id",
		},
	}
}
func FabricNetworkResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Network URI information",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Fabric Network name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
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
				Schema: createLocationSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Computed:    true,
			Optional:    true,
			Description: "Fabric Network project",
			Elem: &schema.Resource{
				Schema: FabricNetworkProjectSch(),
			},
		},
		"operation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Network operation information that is associated with this Fabric Network",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: FabricNetworkOperationSch(),
			},
		},
		"change": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Change information related to this Fabric Network",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: FabricNetworkChangeSch(),
			},
		},
		"notifications": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Preferences for notifications on Fabric Network configuration or status changes",
			Elem: &schema.Resource{
				Schema: createNotificationSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures Fabric Network lifecycle change information",
			Elem: &schema.Resource{
				Schema: createChangeLogSch(),
			},
		},
	}
}
func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Read:   schema.DefaultTimeout(6 * time.Minute),
		},
		ReadContext:   resourceFabricNetworkRead,
		CreateContext: resourceFabricNetworkCreate,
		UpdateContext: resourceFabricNetworkUpdate,
		DeleteContext: resourceFabricNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: FabricNetworkResourceSchema(),

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric Network\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
	}
}

func resourceFabricNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := notificationToFabric(schemaNotifications)
	schemaLocation := d.Get("location").(*schema.Set).List()
	location := locationToFabric(schemaLocation)
	project := v4.Project{}
	schemaProject := d.Get("project").(*schema.Set).List()
	if len(schemaProject) != 0 {
		project = projectToFabric(schemaProject)
	}
	netType := v4.NetworkType(d.Get("type").(string))
	netScope := v4.NetworkScope(d.Get("scope").(string))

	createRequest := v4.NetworkPostRequest{
		Name:          d.Get("name").(string),
		Type_:         &netType,
		Scope:         &netScope,
		Location:      &location,
		Notifications: notifications,
		Project:       &project,
	}

	fabricNetwork, _, err := client.NetworksApi.CreateNetwork(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fabricNetwork.Uuid)

	if _, err = waitUntilFabricNetworkIsProvisioned(d.Id(), meta, ctx); err != nil {
		return diag.Errorf("error waiting for Network (%s) to be created: %s", d.Id(), err)
	}

	return resourceFabricNetworkRead(ctx, d, meta)
}

func resourceFabricNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	fabricNetwork, _, err := client.NetworksApi.GetNetworkByUuid(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Fabric Network %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	d.SetId(fabricNetwork.Uuid)
	return setFabricNetworkMap(d, fabricNetwork)
}
func FabricNetworkOperationToTerra(operation *v4.NetworkOperation) *schema.Set {
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
		schema.HashResource(&schema.Resource{Schema: FabricNetworkOperationSch()}),
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
		schema.HashResource(&schema.Resource{Schema: FabricNetworkChangeSch()}),
		mappedChanges,
	)
	return changeSet
}

func setFabricNetworkMap(d *schema.ResourceData, nt v4.Network) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"name":          nt.Name,
		"type":          nt.Type_,
		"scope":         nt.Scope,
		"state":         nt.State,
		"operation":     FabricNetworkOperationToTerra(nt.Operation),
		"change":        simplifiedFabricNetworkChangeToTerra(nt.Change),
		"location":      locationToTerra(nt.Location),
		"notifications": notificationToTerra(nt.Notifications),
		"project":       projectToTerra(nt.Project),
		"change_log":    changeLogToTerra(nt.ChangeLog),
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
	dbConn, err := waitUntilFabricNetworkIsProvisioned(d.Id(), meta, ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.Errorf("either timed out or errored out while fetching Fabric Network for uuid %s and error %v", d.Id(), err)
	}
	// TO-DO
	update, err := getFabricNetworkUpdateRequest(dbConn, d)
	if err != nil {
		return diag.FromErr(err)
	}
	updates := []v4.NetworkChangeOperation{update}
	_, res, err := client.NetworksApi.UpdateNetworkByUuid(ctx, updates, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error response for the Fabric Network update, response %v, error %v", res, err))
	}
	updateFg := v4.Network{}
	updateFg, err = waitForFabricNetworkUpdateCompletion(d.Id(), meta, ctx)

	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("errored while waiting for successful Fabric Network update, response %v, error %v", res, err))
	}

	d.SetId(updateFg.Uuid)
	return setFabricNetworkMap(d, updateFg)
}

func waitForFabricNetworkUpdateCompletion(uuid string, meta interface{}, ctx context.Context) (v4.Network, error) {
	log.Printf("Waiting for Network update to complete, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Target: []string{string(v4.PROVISIONED_NetworkEquinixStatus)},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.NetworksApi.GetNetworkByUuid(ctx, uuid)
			if err != nil {
				return "", "", err
			}
			return dbConn, string(*dbConn.Operation.EquinixStatus), nil
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

func waitUntilFabricNetworkIsProvisioned(uuid string, meta interface{}, ctx context.Context) (v4.Network, error) {
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
				return "", "", err
			}
			return dbConn, string(*dbConn.Operation.EquinixStatus), nil
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

func resourceFabricNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	_, resp, err := client.NetworksApi.DeleteNetworkByUuid(ctx, d.Id())
	if err != nil {
		errors, ok := err.(v4.GenericSwaggerError).Model().([]v4.ModelError)
		if ok {
			// EQ-3040055 = There is an existing update in REQUESTED state
			if equinix_errors.HasModelErrorCode(errors, "EQ-3040055") {
				return diags
			}
		}
		return diag.FromErr(fmt.Errorf("error response for the Fabric Network delete. Error %v and response %v", err, resp))
	}

	err = waitUntilFabricNetworkDeprovisioned(d.Id(), meta, ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func waitUntilFabricNetworkDeprovisioned(uuid string, meta interface{}, ctx context.Context) error {
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
				return "", "", err
			}
			return dbConn, string(*dbConn.Operation.EquinixStatus), nil
		},
		Timeout:    5 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
