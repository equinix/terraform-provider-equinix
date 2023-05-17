package equinix

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
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
			mappedServiceToken := serviceTokenToFabric(serviceTokenRequest)
			connectionASide = v4.ConnectionSide{ServiceToken: &mappedServiceToken}
		}
		if len(additionalInfoRequest) != 0 {
			mappedAdditionalInfo := additionalInfoToFabric(additionalInfoRequest)
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
			mappedServiceToken := serviceTokenToFabric(serviceTokenRequest)
			connectionZSide = v4.ConnectionSide{ServiceToken: &mappedServiceToken}
		}
		if len(additionalInfoRequest) != 0 {
			mappedAdditionalInfo := additionalInfoToFabric(additionalInfoRequest)
			connectionZSide = v4.ConnectionSide{AdditionalInfo: mappedAdditionalInfo}
		}
	}

	createRequest := v4.ConnectionPostRequest{
		Name:          d.Get("name").(string),
		Type_:         &conType,
		Order:         &order,
		Notifications: notifications,
		Bandwidth:     int32(d.Get("bandwidth").(int)),
		Redundancy:    &red,
		ASide:         &connectionASide,
		ZSide:         &connectionZSide,
		Project:       &project,
	}

	conn, _, err := client.ConnectionsApi.CreateConnection(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(conn.Uuid)

	if err = waitUntilConnectionIsProvisioned(d.Id(), meta, ctx); err != nil {
		return diag.Errorf("error waiting for connection (%s) to be created: %s", d.Id(), err)
	}

	return resourceFabricConnectionRead(ctx, d, meta)
}

func resourceFabricConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
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
	err := setMap(d, map[string]interface{}{
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
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	dbConn, err := waitUntilConnectionIsActive(d.Id(), meta, ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.Errorf("either timed out or errored out while fetching connection for uuid %s and error %v", d.Id(), err)
	}
	update, err := getUpdateRequest(dbConn, d)
	if err != nil {
		return diag.FromErr(err)
	}
	updates := []v4.ConnectionChangeOperation{update}
	_, res, err := client.ConnectionsApi.UpdateConnectionByUuid(ctx, updates, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error response for the connection update, response %v, error %v", res, err))
	}
	updatedConn := v4.Connection{}
	updatedConn, err = waitForConnectionUpdateCompletion(d.Id(), meta, ctx)

	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("errored while waiting for successful connection update, response %v, error %v", res, err))
	}

	d.SetId(updatedConn.Uuid)
	return setFabricMap(d, updatedConn)
}

func waitForConnectionUpdateCompletion(uuid string, meta interface{}, ctx context.Context) (v4.Connection, error) {
	log.Printf("Waiting for connection update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
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
		Timeout:    2 * time.Minute,
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

func waitUntilConnectionIsProvisioned(uuid string, meta interface{}, ctx context.Context) error {
	log.Printf("Waiting for connection to be provisioned, uuid %s", uuid)
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
			client := meta.(*Config).fabricClient
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

func waitUntilConnectionIsActive(uuid string, meta interface{}, ctx context.Context) (v4.Connection, error) {
	log.Printf("Waiting for connection to be in active state, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{
			string(v4.ACTIVE_ConnectionState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
			if err != nil {
				return "", "", err
			}
			updatableState := ""
			if *dbConn.State == v4.ACTIVE_ConnectionState {
				updatableState = string(*dbConn.State)
			}
			return dbConn, updatableState, nil
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
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	_, resp, err := client.ConnectionsApi.DeleteConnectionByUuid(ctx, d.Id())
	if err != nil {
		errors, ok := err.(v4.GenericSwaggerError).Model().([]v4.ModelError)
		if ok {
			// EQ-3142509 = Connection already deleted
			if hasModelErrorCode(errors, "EQ-3142509") {
				return diags
			}
		}
		return diag.FromErr(fmt.Errorf("error response for the connection delete. Error %v and response %v", err, resp))
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
			client := meta.(*Config).fabricClient
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
