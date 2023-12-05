package equinix

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFabricConnection() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Read:   schema.DefaultTimeout(6 * time.Minute),
		},
		ReadContext:   resourceFabricConnectionRead,
		CreateContext: resourceFabricConnectionCreate,
		UpdateContext: resourceFabricConnectionUpdate,
		DeleteContext: resourceFabricConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: createFabricConnectionResourceSchema(),

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric connection\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
	}
}

func resourceFabricConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	conType := v4.ConnectionType(d.Get("type").(string))
	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := notificationToFabric(schemaNotifications)
	schemaRedundancy := d.Get("redundancy").(*schema.Set).List()
	red := redundancyToFabric(schemaRedundancy)
	schemaOrder := d.Get("order").(*schema.Set).List()
	order := orderToFabric(schemaOrder)
	aside := d.Get("a_side").(*schema.Set).List()
	projectReq := d.Get("project").(*schema.Set).List()
	project := projectToFabric(projectReq)
	additionalInfoTerraConfig := d.Get("additional_info").([]interface{})
	additionalInfo := additionalInfoTerraToGo(additionalInfoTerraConfig)
	connectionASide := v4.ConnectionSide{}
	for _, as := range aside {
		asideMap := as.(map[string]interface{})
		accessPoint := asideMap["access_point"].(*schema.Set).List()
		serviceTokenRequest := asideMap["service_token"].(*schema.Set).List()
		additionalInfoRequest := asideMap["additional_info"].([]interface{})

		if len(accessPoint) != 0 {
			ap := accessPointToFabric(accessPoint)
			connectionASide = v4.ConnectionSide{AccessPoint: &ap}
		}
		if len(serviceTokenRequest) != 0 {
			mappedServiceToken, err := serviceTokenToFabric(serviceTokenRequest)
			if err != nil {
				return diag.FromErr(err)
			}
			connectionASide = v4.ConnectionSide{ServiceToken: &mappedServiceToken}
		}
		if len(additionalInfoRequest) != 0 {
			mappedAdditionalInfo := additionalInfoTerraToGo(additionalInfoRequest)
			connectionASide = v4.ConnectionSide{AdditionalInfo: mappedAdditionalInfo}
		}
	}

	zside := d.Get("z_side").(*schema.Set).List()
	connectionZSide := v4.ConnectionSide{}
	for _, as := range zside {
		zsideMap := as.(map[string]interface{})
		accessPoint := zsideMap["access_point"].(*schema.Set).List()
		serviceTokenRequest := zsideMap["service_token"].(*schema.Set).List()
		additionalInfoRequest := zsideMap["additional_info"].([]interface{})
		if len(accessPoint) != 0 {
			ap := accessPointToFabric(accessPoint)
			connectionZSide = v4.ConnectionSide{AccessPoint: &ap}
		}
		if len(serviceTokenRequest) != 0 {
			mappedServiceToken, err := serviceTokenToFabric(serviceTokenRequest)
			if err != nil {
				return diag.FromErr(err)
			}
			connectionZSide = v4.ConnectionSide{ServiceToken: &mappedServiceToken}
		}
		if len(additionalInfoRequest) != 0 {
			mappedAdditionalInfo := additionalInfoTerraToGo(additionalInfoRequest)
			connectionZSide = v4.ConnectionSide{AdditionalInfo: mappedAdditionalInfo}
		}
	}

	createRequest := v4.ConnectionPostRequest{
		Name:           d.Get("name").(string),
		Type_:          &conType,
		Order:          &order,
		Notifications:  notifications,
		Bandwidth:      int32(d.Get("bandwidth").(int)),
		AdditionalInfo: additionalInfo,
		Redundancy:     &red,
		ASide:          &connectionASide,
		ZSide:          &connectionZSide,
		Project:        &project,
	}

	conn, _, err := client.ConnectionsApi.CreateConnection(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(conn.Uuid)

	if err = waitUntilConnectionIsCreated(d.Id(), meta, ctx); err != nil {
		return diag.Errorf("error waiting for connection (%s) to be created: %s", d.Id(), err)
	}

	awsSecrets, hasAWSSecrets := additionalInfoContainsAWSSecrets(additionalInfoTerraConfig)
	if hasAWSSecrets {
		patchChangeOperation := []v4.ConnectionChangeOperation{
			{
				Op:    "add",
				Path:  "",
				Value: map[string]interface{}{"additionalInfo": awsSecrets},
			},
		}

		_, _, patchErr := client.ConnectionsApi.UpdateConnectionByUuid(ctx, patchChangeOperation, conn.Uuid)
		if patchErr != nil {
			return diag.FromErr(err)
		}

		if _, statusChangeErr := waitForConnectionProviderStatusChange(d.Id(), meta, ctx); err != nil {
			return diag.Errorf("error waiting for AWS Approval for connection %s: %v", d.Id(), statusChangeErr)
		}
	}

	return resourceFabricConnectionRead(ctx, d, meta)
}

func additionalInfoContainsAWSSecrets(info []interface{}) ([]interface{}, bool) {
	var awsSecrets []interface{}

	for _, item := range info {
		if value, _ := item.(map[string]interface{})["key"]; value == "accessKey" {
			awsSecrets = append(awsSecrets, item)
		}

		if value, _ := item.(map[string]interface{})["key"]; value == "secretKey" {
			awsSecrets = append(awsSecrets, item)
		}
	}

	return awsSecrets, len(awsSecrets) == 2
}

func resourceFabricConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	conn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, d.Id(), nil)
	if err != nil {
		log.Printf("[WARN] Connection %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	d.SetId(conn.Uuid)
	return setFabricMap(d, conn)
}

func setFabricMap(d *schema.ResourceData, conn v4.Connection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"name":      conn.Name,
		"bandwidth": conn.Bandwidth,
		"href":      conn.Href,
		// TODO v4.ConnectionPostRequest doesn't have a "description" field,
		// so it always returns empty because it was never in the API, that produces an inconsistency
		// "description":     conn.Description,
		"is_remote":       conn.IsRemote,
		"type":            conn.Type_,
		"state":           conn.State,
		"direction":       conn.Direction,
		"operation":       operationToTerra(conn.Operation),
		"order":           orderMappingToTerra(conn.Order),
		"change_log":      changeLogToTerra(conn.ChangeLog),
		"redundancy":      redundancyToTerra(conn.Redundancy),
		"notifications":   notificationToTerra(conn.Notifications),
		"account":         accountToTerra(conn.Account),
		"a_side":          connectionSideToTerra(conn.ASide),
		"z_side":          connectionSideToTerra(conn.ZSide),
		"additional_info": additionalInfoToTerra(conn.AdditionalInfo),
		"project":         projectToTerra(conn.Project),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceFabricConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	dbConn, err := verifyConnectionCreated(d.Id(), meta, ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.Errorf("either timed out or errored out while fetching connection for uuid %s: error -> %v", d.Id(), err)
	}

	diags := diag.Diagnostics{}
	updateRequests, err := getUpdateRequests(dbConn, d)
	if err != nil {
		diags = append(diags, diag.Diagnostic{Severity: 1, Summary: err.Error()})
		return diags
	}
	updatedConn := dbConn

	for _, update := range updateRequests {
		_, _, err := client.ConnectionsApi.UpdateConnectionByUuid(ctx, update, d.Id())
		if err != nil {
			diags = append(diags, diag.Diagnostic{Severity: 2, Summary: fmt.Sprintf("connectionn property update request error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", err, update)})
			continue
		}

		var waitFunction func(uuid string, meta interface{}, ctx context.Context) (v4.Connection, error)
		if update[0].Op == "replace" {
			// Update type is either name or bandwidth
			waitFunction = waitForConnectionUpdateCompletion
		} else if update[0].Op == "add" {
			// Update type is aws secret additionalInfo
			waitFunction = waitForConnectionProviderStatusChange
		}

		conn, err := waitFunction(d.Id(), meta, ctx)

		if err != nil {
			if !strings.Contains(err.Error(), "500") {
				d.SetId("")
			}
			diags = append(diags, diag.Diagnostic{Severity: 2, Summary: fmt.Sprintf("connection property update completion timeout error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", err, update)})
		} else {
			updatedConn = conn
		}
	}

	d.SetId(updatedConn.Uuid)
	return append(diags, setFabricMap(d, updatedConn)...)
}

func waitForConnectionUpdateCompletion(uuid string, meta interface{}, ctx context.Context) (v4.Connection, error) {
	log.Printf("[DEBUG] Waiting for connection update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
			if err != nil {
				return "", "", err
			}
			updatableState := ""
			if dbConn.Change.Status == "COMPLETED" {
				updatableState = dbConn.Change.Status
			}
			return dbConn, updatableState, nil
		},
		Timeout:    3 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	dbConn := v4.Connection{}

	if err == nil {
		dbConn = inter.(v4.Connection)
	}
	return dbConn, err
}

func waitUntilConnectionIsCreated(uuid string, meta interface{}, ctx context.Context) error {
	log.Printf("Waiting for connection to be created, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(v4.PROVISIONING_ConnectionState),
		},
		Target: []string{
			string(v4.PENDING_ConnectionState),
			string(v4.PROVISIONED_ConnectionState),
			string(v4.ACTIVE_ConnectionState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
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

func waitForConnectionProviderStatusChange(uuid string, meta interface{}, ctx context.Context) (v4.Connection, error) {
	log.Printf("DEBUG: wating for provider status to update. Connection uuid: %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(v4.PENDING_APPROVAL_ProviderStatus),
		},
		Target: []string{
			string(v4.PROVISIONING_ProviderStatus),
			string(v4.PROVISIONED_ProviderStatus),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
			if err != nil {
				return "", "", err
			}
			return dbConn, string(*dbConn.Operation.ProviderStatus), nil
		},
		Timeout:    5 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	dbConn := v4.Connection{}

	if err == nil {
		dbConn = inter.(v4.Connection)
	}
	return dbConn, err
}

func verifyConnectionCreated(uuid string, meta interface{}, ctx context.Context) (v4.Connection, error) {
	log.Printf("Waiting for connection to be in created state, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{
			string(v4.ACTIVE_ConnectionState),
			string(v4.PROVISIONED_ConnectionState),
			string(v4.PENDING_ConnectionState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
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
	dbConn := v4.Connection{}

	if err == nil {
		dbConn = inter.(v4.Connection)
	}
	return dbConn, err
}

func resourceFabricConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	_, _, err := client.ConnectionsApi.DeleteConnectionByUuid(ctx, d.Id())
	if err != nil {
		errors, ok := err.(v4.GenericSwaggerError).Model().([]v4.ModelError)
		if ok {
			// EQ-3142509 = Connection already deleted
			if hasModelErrorCode(errors, "EQ-3142509") {
				return diags
			}
		}
		return diag.FromErr(fmt.Errorf("error response for the connection delete: %v", err))
	}

	err = waitUntilConnectionDeprovisioned(d.Id(), meta, ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func waitUntilConnectionDeprovisioned(uuid string, meta interface{}, ctx context.Context) error {
	log.Printf("Waiting for connection to be deprovisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(v4.DEPROVISIONING_ConnectionState),
		},
		Target: []string{
			string(v4.DEPROVISIONED_ConnectionState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
			if err != nil {
				return "", "", err
			}
			return dbConn, string(*dbConn.State), nil
		},
		Timeout:    6 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
