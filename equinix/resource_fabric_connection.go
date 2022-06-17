package equinix

import (
	"context"
	"fmt"
	v4 "github.com/equinix/terraform-provider-equinix/internal/apis/fabric/v4"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"sort"
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
	fmt.Println(" inside fabric connection create resource ")
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
		companyProfileRequest := asideMap["company_profile"].(interface{}).(*schema.Set).List()
		natRequest := asideMap["nat"].(interface{}).(*schema.Set).List()
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

		if len(companyProfileRequest) != 0 {
			mappedCompanyProfile := companyProfileToFabric(companyProfileRequest)
			connectionASide = v4.ConnectionSide{CompanyProfile: &mappedCompanyProfile}
		}

		if len(natRequest) != 0 {
			mappedNat := natToFabric(natRequest)
			connectionASide = v4.ConnectionSide{Nat: &mappedNat}
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
		companyProfileRequest := zsideMap["company_profile"].(interface{}).(*schema.Set).List()
		natRequest := zsideMap["nat"].(interface{}).(*schema.Set).List()
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

		if len(companyProfileRequest) != 0 {
			mappedCompanyProfile := companyProfileToFabric(companyProfileRequest)
			connectionZSide = v4.ConnectionSide{CompanyProfile: &mappedCompanyProfile}
		}

		if len(natRequest) != 0 {
			mappedNat := natToFabric(natRequest)
			connectionZSide = v4.ConnectionSide{Nat: &mappedNat}
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
	conn, _, err := client.ConnectionsApi.CreateConnection(ctx, createRequest, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(conn.Uuid)
	return resourceFabricConnectionRead(ctx, d, meta)
}

func resourceFabricConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fmt.Println(" inside fabric connection read resource ")
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	conn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, d.Id(), nil)
	if err != nil {
		log.Printf("[WARN] Connection %s not found ", d.Id())
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(conn.Uuid)
	return setFabricMap(d, conn)
}

func setFabricMap(d *schema.ResourceData, conn v4.Connection) diag.Diagnostics {
	errs := &multierror.Error{}
	diags := diag.Diagnostics{}
	var err error

	//TODO need clarity if we need to handle for each "Set" error
	err = d.Set("name", conn.Name)

	err = d.Set("bandwidth", conn.Bandwidth)

	err = d.Set("href", conn.Href)

	err = d.Set("platform_uuid", conn.PlatformUuid)

	err = d.Set("description", conn.Description)

	err = d.Set("bandwidth", conn.Bandwidth)

	err = d.Set("is_remote", conn.IsRemote)

	err = d.Set("type", conn.Type_)

	err = d.Set("state", conn.State)

	err = d.Set("direction", conn.Direction)

	//TODO need to see if we need implementation
	//err = d.Set("change", TBD)
	err = d.Set("operation", operationToTerra(conn.Operation))

	if conn.Order != nil {
		err = d.Set("order", orderMappingToTerra(conn.Order))
	}

	err = d.Set("change_log", changeLogToTerra(conn.ChangeLog))

	err = d.Set("redundancy", redundancyToTerra(conn.Redundancy))

	err = d.Set("notifications", notificationToTerra(conn.Notifications))

	if conn.Account != nil {
		err = d.Set("account", accountToTerra(conn.Account))
	}

	err = d.Set("tags", conn.Tags)

	err = d.Set("a_side", connectionSideToTerra(conn.ASide))

	err = d.Set("z_side", connectionSideToTerra(conn.ZSide))
	if err != nil {
		//err = fmt.Errorf("connection z_side  %v", err)
		errs = multierror.Append(errs, err)
		err = nil
	}
	if errs.Len() > 0 {
		sort.Sort(errs)
		return diag.FromErr(errs)
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
	cunter := 0
	for dbConn.State == nil || "ACTIVE" != *dbConn.State {
		time.Sleep(30 * time.Second)
		dbConn, _, err = client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
		if err != nil {
			//TODO handle error and add to the diag rather returning it directly
			return diag.Errorf(" Error while fetching connection for uuid %s and error %v", uuid, err)
		}
		if cunter >= 4 {
			break
		}
		cunter++
	}

	update, err := getUpdateRequest(dbConn, d)
	if err != nil {
		return diag.FromErr(err)
	}
	updates := []v4.ConnectionChangeOperation{update}
	_, res, err := client.ConnectionsApi.UpdateConnectionByUuid(ctx, updates, uuid, nil)
	if err != nil {
		//TODO handle response
		fmt.Errorf(" Error response for the connection update, response %v, error %v", res, err)
		return diag.FromErr(err)
	}

	time.Sleep(1 * time.Minute)
	updatedConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
	fmt.Println(" inside fabric connection update resource ")
	if err != nil {
		//TODO handle response
		fmt.Errorf(" Error response while getting updated conection, response %v, error %v", res, err)
		return diag.FromErr(err)
	}
	d.SetId(updatedConn.Uuid)
	return setFabricMap(d, updatedConn)
}

func resourceFabricConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fmt.Println(" inside fabric connection delete resource ")
	diags := diag.Diagnostics{}
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	uuid := d.Id()
	if uuid == "" {
		//TODO added an error and return
		return diag.Errorf("No uuid found %v ", uuid)
	}
	conn, resp, err := client.ConnectionsApi.DeleteConnectionByUuid(ctx, uuid, nil)

	if err != nil {
		fmt.Errorf(" Error response for the connection delete error %v and response %v", err, resp)
		return diag.FromErr(err)
	} else {
		fmt.Println(" Connection Details ")
		fmt.Println(conn)
	}

	return diags
}
