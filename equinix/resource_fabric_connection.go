package equinix

import (
	"context"
	"fmt"
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"time"
)

func resourceFabricConnection() *schema.Resource {

	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
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

		Description: "Resource allows creation and management of Equinix Fabric	layer 2 connections",
	}
}

func resourceFabricConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf(" inside fabric connection create resource ")
	client := meta.(*Config).fabricClient

	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	conType := v4.ConnectionType(d.Get("type").(string))
	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := notificationToFabric(schemaNotifications)
	schemaRedundancy := d.Get("redundancy").(interface{}).(*schema.Set).List()
	red := redundancyToFabric(schemaRedundancy)
	schemaOrder := d.Get("order").(interface{}).(*schema.Set).List()
	order := orderToFabric(schemaOrder)
	aside := d.Get("a_side").(interface{}).(*schema.Set).List()
	connectionASide := v4.ConnectionSide{}
	for _, as := range aside {
		asideMap := as.(map[string]interface{})
		accessPoint := asideMap["access_point"].(interface{}).(*schema.Set).List()
		invitationRequest := asideMap["invitation"].(interface{}).(*schema.Set).List()
		serviceTokenRequest := asideMap["service_token"].(interface{}).(*schema.Set).List()
		additionalInfoRequest := asideMap["additional_info"].([]interface{})

		if len(accessPoint) != 0 {
			ap := accessPointToFabric(accessPoint)
			connectionASide = v4.ConnectionSide{AccessPoint: &ap}
		}
		if len(invitationRequest) != 0 {
			mappedInvitation := invitationToFabric(invitationRequest)
			connectionASide = v4.ConnectionSide{Invitation: &mappedInvitation}
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

	zside := d.Get("z_side").(interface{}).(*schema.Set).List()
	connectionZSide := v4.ConnectionSide{}
	for _, as := range zside {
		zsideMap := as.(map[string]interface{})
		accessPoint := zsideMap["access_point"].(interface{}).(*schema.Set).List()
		invitationRequest := zsideMap["invitation"].(interface{}).(*schema.Set).List()
		serviceTokenRequest := zsideMap["service_token"].(interface{}).(*schema.Set).List()
		additionalInfoRequest := zsideMap["additional_info"].([]interface{})

		if len(accessPoint) != 0 {
			ap := accessPointToFabric(accessPoint)
			connectionZSide = v4.ConnectionSide{AccessPoint: &ap}
		}
		if len(invitationRequest) != 0 {
			mappedInvitation := invitationToFabric(invitationRequest)
			connectionZSide = v4.ConnectionSide{Invitation: &mappedInvitation}
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
	}

	conn, _, err := client.ConnectionsApi.CreateConnection(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(conn.Uuid)
	return resourceFabricConnectionRead(ctx, d, meta)
}

func resourceFabricConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf(" inside fabric connection read resource ")
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	conn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, d.Id(), nil)
	//TODO check for 404 or 403 from the API and see if we need specific handling
	if err != nil {
		log.Printf("[WARN] Connection %s not found ", d.Id())
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(conn.Uuid)
	return setFabricMap(d, conn)
}

func setFabricMap(d *schema.ResourceData, conn v4.Connection) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := setMap(d, map[string]interface{}{
		"name":          conn.Name,
		"bandwidth":     conn.Bandwidth,
		"href":          conn.Href,
		"description":   conn.Description,
		"is_remote":     conn.IsRemote,
		"type":          conn.Type_,
		"state":         conn.State,
		"direction":     conn.Direction,
		"operation":     operationToTerra(conn.Operation),
		"order":         orderMappingToTerra(conn.Order),
		"change_log":    changeLogToTerra(conn.ChangeLog),
		"redundancy":    redundancyToTerra(conn.Redundancy),
		"notifications": notificationToTerra(conn.Notifications),
		"account":       accountToTerra(conn.Account),
		"a_side":        connectionSideToTerra(conn.ASide),
		"z_side":        connectionSideToTerra(conn.ZSide),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceFabricConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	uuid := d.Id()
	if uuid == "" {
		//TODO added an error and return
		return diag.Errorf("No connection found for the value uuid %v ", uuid)
	}

	dbConn := v4.Connection{}
	var err error
	counter := 0
	//TODO Terraform retry helpers to use
	for dbConn.State == nil || "ACTIVE" != *dbConn.State {
		time.Sleep(30 * time.Second)
		dbConn, _, err = client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
		if err != nil {
			//TODO handle error and add to the diag rather returning it directly
			return diag.Errorf(" Error while fetching connection for uuid %s and error %v", uuid, err)
		}
		if counter >= 4 {
			break
		}
		counter++
	}

	update, err := getUpdateRequest(dbConn, d)
	if err != nil {
		return diag.FromErr(err)
	}
	updates := []v4.ConnectionChangeOperation{update}
	_, res, err := client.ConnectionsApi.UpdateConnectionByUuid(ctx, updates, uuid)
	if err != nil {
		//TODO handle response
		fmt.Errorf(" Error response for the connection update, response %v, error %v", res, err)
		return diag.FromErr(err)
	}

	time.Sleep(1 * time.Minute)
	updatedConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
	log.Printf(" inside fabric connection update resource ")
	if err != nil {
		//TODO handle response
		fmt.Errorf(" Error response while getting updated conection, response %v, error %v", res, err)
		return diag.FromErr(err)
	}
	d.SetId(updatedConn.Uuid)
	return setFabricMap(d, updatedConn)
}

func resourceFabricConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf(" inside fabric connection delete resource ")
	diags := diag.Diagnostics{}
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	uuid := d.Id()
	if uuid == "" {
		//TODO added an error and return
		return diag.Errorf("No uuid found %v ", uuid)
	}
	_, resp, err := client.ConnectionsApi.DeleteConnectionByUuid(ctx, uuid)
	//TODO check if delete is returing 404 in case if API did not find an active connction for the ID
	if err != nil {
		fmt.Errorf(" Error response for the connection delete error %v and response %v", err, resp)
		return diag.FromErr(err)
	}

	return diags
}
