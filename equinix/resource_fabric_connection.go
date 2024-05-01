package equinix

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func fabricConnectionResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"EVPL_VC", "EPL_VC", "IP_VC", "IPWAN_VC", "ACCESS_EPL_VC", "EVPLAN_VC", "EPLAN_VC", "EIA_VC", "IA_VC", "EC_VC"}, false),
			Description:  "Defines the connection type like EVPL_VC, EPL_VC, IPWAN_VC, IP_VC, ACCESS_EPL_VC, EVPLAN_VC, EPLAN_VC, EIA_VC, IA_VC, EC_VC",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 24),
			Description:  "Connection name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"order": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Description: "Order details",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.OrderSch(),
			},
		},
		"notifications": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Preferences for notifications on connection configuration or status changes",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.NotificationSch(),
			},
		},
		"bandwidth": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Connection bandwidth in Mbps",
		},
		//"geo_scope": {
		//	Type:         schema.TypeString,
		//	Optional:     true,
		//	ValidateFunc: validation.StringInSlice([]string{"CANADA", "CONUS"}, false),
		//	Description:  "Geographic boundary types",
		//},
		"redundancy": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Connection Redundancy Configuration",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: connectionRedundancySch(),
			},
		},
		"a_side": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Requester or Customer side connection configuration object of the multi-segment connection",
			MaxItems:    1,
			Elem:        connectionSideSch(),
			Set:         schema.HashResource(accessPointSch()),
		},
		"z_side": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Destination or Provider side connection configuration object of the multi-segment connection",
			MaxItems:    1,
			Elem:        connectionSideSch(),
			Set:         schema.HashResource(accessPointSch()),
		},
		"project": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Description: "Project information",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ProjectSch(),
			},
		},
		"additional_info": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Connection additional information",
			Elem: &schema.Schema{
				Type: schema.TypeMap,
			},
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection URI information",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned connection identifier",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Customer-provided connection description",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection overall state",
		},
		"operation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Connection type-specific operational data",
			Elem: &schema.Resource{
				Schema: operationSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Customer account information that is associated with this connection",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.AccountSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures connection lifecycle change information",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ChangeLogSch(),
			},
		},
		"is_remote": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Connection property derived from access point locations",
		},
		"direction": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection directionality from the requester point of view",
		},
	}
}

func connectionSideSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service_token": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "For service token based connections, Service tokens authorize users to access protected resources and services. Resource owners can distribute the tokens to trusted partners and vendors, allowing selected third parties to work directly with Equinix network assets",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: serviceTokenSch(),
				},
			},
			"access_point": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Point of access details",
				MaxItems:    1,
				Elem:        accessPointSch(),
			},
			"additional_info": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Connection side additional information",
				Elem: &schema.Resource{
					Schema: additionalInfoSch(),
				},
			},
		},
	}
}

func serviceTokenSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"VC_TOKEN"}, true),
			Description:  "Token type - VC_TOKEN",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "An absolute URL that is the subject of the link's context",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix-assigned service token identifier",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service token description",
		},
	}
}

func accessPointSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"COLO", "VD", "VG", "SP", "IGW", "SUBNET", "CLOUD_ROUTER", "NETWORK", "METAL_NETWORK"}, true),
				Description:  "Access point type - COLO, VD, VG, SP, IGW, SUBNET, CLOUD_ROUTER, NETWORK",
			},
			"account": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Account",
				Elem: &schema.Resource{
					Schema: equinix_fabric_schema.AccountSch(),
				},
			},
			"location": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Access point location",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: equinix_fabric_schema.LocationSch(),
				},
			},
			"port": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Port access point information",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: portSch(),
				},
			},
			"profile": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Service Profile",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: serviceProfileSch(),
				},
			},
			"gateway": {
				Type:        schema.TypeSet,
				Optional:    true,
				Deprecated:  "use router attribute instead; gateway is no longer a part of the supported backend",
				Description: "**Deprecated** `gateway` Use `router` attribute instead",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: cloudRouterSch(),
				},
			},
			"router": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Cloud Router access point information that replaces `gateway`",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: cloudRouterSch(),
				},
			},
			"link_protocol": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Connection link protocol",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: accessPointLinkProtocolSch(),
				},
			},
			"virtual_device": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Virtual device",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: accessPointVirtualDeviceSch(),
				},
			},
			"interface": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Virtual device interface",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: accessPointInterface(),
				},
			},
			"network": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "network access point information",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: networkSch(),
				},
			},
			"seller_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access point seller region",
			},
			"peering_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PRIVATE", "MICROSOFT", "PUBLIC", "MANUAL"}, true),
				Description:  "Peering Type- PRIVATE,MICROSOFT,PUBLIC, MANUAL",
			},
			"authentication_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Authentication key for provider based connections",
			},
			"provider_connection_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provider assigned Connection Id",
			},
		},
	}
}

func serviceProfileSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service Profile URI response attribute",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"L2_PROFILE", "L3_PROFILE", "ECIA_PROFILE", "ECMC_PROFILE", "IA_PROFILE"}, true),
			Description:  "Service profile type - L2_PROFILE, L3_PROFILE, ECIA_PROFILE, ECMC_PROFILE, IA_PROFILE",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Customer-assigned service profile name",
		},
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Equinix assigned service profile identifier",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "User-provided service description",
		},
		"access_point_type_configs": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Access point config information",
			Elem: &schema.Resource{
				Schema: connectionAccessPointTypeConfigSch(),
			},
		},
	}
}

func cloudRouterSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix-assigned virtual gateway identifier",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
	}
}

func accessPointLinkProtocolSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Type of the link protocol - UNTAGGED, DOT1Q, QINQ, EVPN_VXLAN",
			ValidateFunc: validation.StringInSlice([]string{"UNTAGGED", "DOT1Q", "QINQ", "EVPN_VXLAN"}, true),
		},
		"vlan_tag": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Vlan Tag information, vlanTag value specified for DOT1Q connections",
		},
		"vlan_s_tag": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Vlan Provider Tag information, vlanSTag value specified for QINQ connections",
		},
		"vlan_c_tag": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Vlan Customer Tag information, vlanCTag value specified for QINQ connections",
		},
	}
}

func accessPointVirtualDeviceSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix-assigned Virtual Device identifier",
		},
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Virtual Device type",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Customer-assigned Virtual Device Name",
		},
	}
}

func accessPointInterface() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix-assigned interface identifier",
		},
		"id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "id",
		},
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Interface type",
		},
	}
}

func networkSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix-assigned Network identifier",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
	}
}

func portSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Equinix-assigned Port identifier",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port name",
		},
		"redundancy": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Redundancy Information",
			Elem: &schema.Resource{
				Schema: PortRedundancySch(),
			},
		},
	}
}

func connectionAccessPointTypeConfigSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of access point type config - VD, COLO",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned access point type config identifier",
		},
	}
}

func additionalInfoSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Additional information key",
		},
		"value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Additional information value",
		},
	}
}

func operationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"provider_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection provider readiness status",
		},
		"equinix_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection status",
		},
		"errors": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Errors occurred",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ErrorSch(),
			},
		},
	}
}

func connectionRedundancySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Redundancy group identifier (Use the redundancy.0.group UUID of primary connection; e.g. one(equinix_fabric_connection.primary_port_connection.redundancy).group or equinix_fabric_connection.primary_port_connection.redundancy.0.group)",
		},
		"priority": {
			Type:         schema.TypeString,
			Computed:     true,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"PRIMARY", "SECONDARY"}, true),
			Description:  "Connection priority in redundancy group - PRIMARY, SECONDARY",
		},
	}
}

func resourceFabricConnection() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
		},
		ReadContext:   resourceFabricConnectionRead,
		CreateContext: resourceFabricConnectionCreate,
		UpdateContext: resourceFabricConnectionUpdate,
		DeleteContext: resourceFabricConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: fabricConnectionResourceSchema(),

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric connection",
	}
}

func resourceFabricConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)

	createConnectionRequest := fabricv4.ConnectionPostRequest{}

	name := d.Get("name").(string)
	createConnectionRequest.SetName(name)

	conType := d.Get("type").(string)
	createConnectionRequest.SetType(fabricv4.ConnectionType(conType))

	if orderSchema, ok := d.GetOk("order"); ok {
		order := equinix_fabric_schema.OrderTerraformToGo(orderSchema.(*schema.Set).List())
		createConnectionRequest.SetOrder(order)
	}

	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := equinix_fabric_schema.NotificationsTerraformToGo(schemaNotifications)
	createConnectionRequest.SetNotifications(notifications)

	bandwidth := d.Get("bandwidth").(int)
	createConnectionRequest.SetBandwidth(int32(bandwidth))

	if schemaRedundancy, ok := d.GetOk("redundancy"); ok {
		redundancy := connectionRedundancyTerraformToGo(schemaRedundancy.(*schema.Set).List())
		createConnectionRequest.SetRedundancy(redundancy)
	}

	if terraConfigProject, ok := d.GetOk("project"); ok {
		project := equinix_fabric_schema.ProjectTerraformToGo(terraConfigProject.(*schema.Set).List())
		createConnectionRequest.SetProject(project)
	}

	aSide := d.Get("a_side").(*schema.Set).List()
	connectionASide := connectionSideTerraformToGo(aSide)
	createConnectionRequest.SetASide(connectionASide)

	zSide := d.Get("z_side").(*schema.Set).List()
	connectionZSide := connectionSideTerraformToGo(zSide)
	createConnectionRequest.SetZSide(connectionZSide)

	additionalInfoTerraConfig, ok := d.GetOk("additional_info")
	if ok {
		zSideAccessPoint := connectionZSide.GetAccessPoint()
		zSideAccessPointServiceProfile := zSideAccessPoint.GetProfile()
		serviceProfile, _, _ := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, zSideAccessPointServiceProfile.GetUuid()).Execute()
		customFields := serviceProfile.GetCustomFields()

		if len(customFields) != 0 {
			additionalInfo := additionalInfoTerraformToGo(additionalInfoTerraConfig.([]interface{}))
			createConnectionRequest.SetAdditionalInfo(additionalInfo)
		}
	}

	start := time.Now()
	conn, _, err := client.ConnectionsApi.CreateConnection(ctx).ConnectionPostRequest(createConnectionRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(conn.GetUuid())

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if err = waitUntilConnectionIsCreated(d.Id(), meta, d, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for connection (%s) to be created: %s", d.Id(), err)
	}

	awsSecrets, hasAWSSecrets := additionalInfoContainsAWSSecrets(additionalInfoTerraConfig.([]interface{}))
	if hasAWSSecrets {
		patchChangeOperation := []fabricv4.ConnectionChangeOperation{
			{
				Op:    "add",
				Path:  "",
				Value: map[string]interface{}{"additionalInfo": awsSecrets},
			},
		}

		_, _, patchErr := client.ConnectionsApi.UpdateConnectionByUuid(ctx, *conn.Uuid).ConnectionChangeOperation(patchChangeOperation).Execute()
		if patchErr != nil {
			return diag.FromErr(equinix_errors.FormatFabricError(patchErr))
		}

		createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
		if _, statusChangeErr := waitForConnectionProviderStatusChange(d.Id(), meta, d, ctx, createTimeout); statusChangeErr != nil {
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
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	conn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, d.Id()).Execute()
	if err != nil {
		log.Printf("[WARN] Connection %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(conn.GetUuid())
	return setFabricMap(d, conn)
}

func setFabricMap(d *schema.ResourceData, conn *fabricv4.Connection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	connection := make(map[string]interface{})
	connection["name"] = conn.GetName()
	connection["bandwidth"] = conn.GetBandwidth()
	connection["href"] = conn.GetHref()
	connection["is_remote"] = conn.GetIsRemote()
	connection["type"] = string(conn.GetType())
	connection["state"] = string(conn.GetState())
	connection["direction"] = conn.GetDirection()
	if conn.Operation != nil {
		operation := conn.GetOperation()
		connection["operation"] = connectionOperationGoToTerraform(&operation)
	}
	if conn.Order != nil {
		order := conn.GetOrder()
		connection["order"] = equinix_fabric_schema.OrderGoToTerraform(&order)
	}
	if conn.ChangeLog != nil {
		changeLog := conn.GetChangeLog()
		connection["change_log"] = equinix_fabric_schema.ChangeLogGoToTerraform(&changeLog)
	}
	if conn.Redundancy != nil {
		redundancy := conn.GetRedundancy()
		connection["redundancy"] = connectionRedundancyGoToTerraform(&redundancy)
	}
	if conn.Notifications != nil {
		notifications := conn.GetNotifications()
		connection["notifications"] = equinix_fabric_schema.NotificationsGoToTerraform(notifications)
	}
	if conn.Account != nil {
		account := conn.GetAccount()
		connection["account"] = equinix_fabric_schema.AccountGoToTerraform(&account)
	}
	if &conn.ASide != nil {
		aSide := conn.GetASide()
		connection["a_side"] = connectionSideGoToTerraform(&aSide)
	}
	if &conn.ZSide != nil {
		zSide := conn.GetZSide()
		connection["z_side"] = connectionSideGoToTerraform(&zSide)
	}
	if conn.AdditionalInfo != nil {
		additionalInfo := conn.GetAdditionalInfo()
		connection["additional_info"] = additionalInfoGoToTerraform(additionalInfo)
	}
	if conn.Project != nil {
		project := conn.GetProject()
		connection["project"] = equinix_fabric_schema.ProjectGoToTerraform(&project)
	}
	err := equinix_schema.SetMap(d, connection)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceFabricConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	start := time.Now()
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	dbConn, err := verifyConnectionCreated(d.Id(), meta, d, ctx, updateTimeout)
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
		_, _, err := client.ConnectionsApi.UpdateConnectionByUuid(ctx, d.Id()).ConnectionChangeOperation(update).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("connection property update request error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", equinix_errors.FormatFabricError(err), update)})
			continue
		}

		var waitFunction func(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.Connection, error)
		if update[0].Op == "replace" {
			// Update type is either name or bandwidth
			waitFunction = waitForConnectionUpdateCompletion
		} else if update[0].Op == "add" {
			// Update type is aws secret additionalInfo
			waitFunction = waitForConnectionProviderStatusChange
		}

		updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
		conn, err := waitFunction(d.Id(), meta, d, ctx, updateTimeout)

		if err != nil {
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("connection property update completion timeout error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", err, update)})
		} else {
			updatedConn = conn
		}
	}

	d.SetId(updatedConn.GetUuid())
	return append(diags, setFabricMap(d, updatedConn)...)
}

func waitForConnectionUpdateCompletion(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.Connection, error) {
	log.Printf("[DEBUG] Waiting for connection update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			updatableState := ""
			change := dbConn.GetChange()
			status := change.GetStatus()
			if string(status) == "COMPLETED" {
				updatableState = string(status)
			}
			return dbConn, updatableState, nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbConn *fabricv4.Connection

	if err == nil {
		dbConn = inter.(*fabricv4.Connection)
	}
	return dbConn, err
}

func waitUntilConnectionIsCreated(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for connection to be created, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CONNECTIONSTATE_PROVISIONING),
		},
		Target: []string{
			string(fabricv4.CONNECTIONSTATE_PENDING),
			string(fabricv4.CONNECTIONSTATE_PROVISIONED),
			string(fabricv4.CONNECTIONSTATE_ACTIVE),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(dbConn.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func waitForConnectionProviderStatusChange(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.Connection, error) {
	log.Printf("DEBUG: wating for provider status to update. Connection uuid: %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.PROVIDERSTATUS_PENDING_APPROVAL),
			string(fabricv4.PROVIDERSTATUS_PROVISIONING),
		},
		Target: []string{
			string(fabricv4.PROVIDERSTATUS_PROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			operation := dbConn.GetOperation()
			providerStatus := operation.GetProviderStatus()
			return dbConn, string(providerStatus), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbConn *fabricv4.Connection

	if err == nil {
		dbConn = inter.(*fabricv4.Connection)
	}
	return dbConn, err
}

func verifyConnectionCreated(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.Connection, error) {
	log.Printf("Waiting for connection to be in created state, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{
			string(fabricv4.CONNECTIONSTATE_ACTIVE),
			string(fabricv4.CONNECTIONSTATE_PROVISIONED),
			string(fabricv4.CONNECTIONSTATE_PENDING),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(dbConn.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbConn *fabricv4.Connection

	if err == nil {
		dbConn = inter.(*fabricv4.Connection)
	}
	return dbConn, err
}

func resourceFabricConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	start := time.Now()
	_, _, err := client.ConnectionsApi.DeleteConnectionByUuid(ctx, d.Id()).Execute()
	if err != nil {
		if genericError, ok := err.(*fabricv4.GenericOpenAPIError); ok {
			if fabricErrs, ok := genericError.Model().([]fabricv4.Error); ok {
				// EQ-3142509 = Connection already deleted
				if equinix_errors.HasErrorCode(fabricErrs, "EQ-3142509") {
					return diags
				}
			}
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	deleteTimeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)
	err = WaitUntilConnectionDeprovisioned(d.Id(), meta, d, ctx, deleteTimeout)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func WaitUntilConnectionDeprovisioned(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for connection to be deprovisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CONNECTIONSTATE_DEPROVISIONING),
			string(fabricv4.CONNECTIONSTATE_ACTIVE),
			string(fabricv4.CONNECTIONSTATE_PENDING),
		},
		Target: []string{
			string(fabricv4.CONNECTIONSTATE_DEPROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(dbConn.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func connectionRedundancyTerraformToGo(redundancyTerraform []interface{}) fabricv4.ConnectionRedundancy {
	if redundancyTerraform == nil || len(redundancyTerraform) == 0 {
		return fabricv4.ConnectionRedundancy{}
	}
	var redundancy fabricv4.ConnectionRedundancy

	redundancyMap := redundancyTerraform[0].(map[string]interface{})
	connectionPriority := redundancyMap["priority"].(string)
	redundancyGroup := redundancyMap["group"].(string)
	redundancy.SetPriority(fabricv4.ConnectionPriority(connectionPriority))
	redundancy.SetGroup(redundancyGroup)

	return redundancy
}

func connectionRedundancyGoToTerraform(redundancy *fabricv4.ConnectionRedundancy) *schema.Set {
	if redundancy == nil {
		return nil
	}
	mappedRedundancy := make(map[string]interface{})
	mappedRedundancy["group"] = redundancy.GetGroup()
	mappedRedundancy["priority"] = string(redundancy.GetPriority())
	redundancySet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: connectionRedundancySch()}),
		[]interface{}{mappedRedundancy},
	)
	return redundancySet
}

func serviceTokenTerraformToGo(serviceTokenList []interface{}) fabricv4.ServiceToken {
	if serviceTokenList == nil || len(serviceTokenList) == 0 {
		return fabricv4.ServiceToken{}
	}

	var serviceToken fabricv4.ServiceToken
	serviceTokenMap := serviceTokenList[0].(map[string]interface{})
	serviceTokenType := serviceTokenMap["type"].(string)
	uuid := serviceTokenMap["uuid"].(string)
	serviceToken.SetType(fabricv4.ServiceTokenType(serviceTokenType))
	serviceToken.SetUuid(uuid)

	return serviceToken
}

func additionalInfoTerraformToGo(additionalInfoList []interface{}) []fabricv4.ConnectionSideAdditionalInfo {
	if additionalInfoList == nil || len(additionalInfoList) == 0 {
		return nil
	}

	mappedAdditionalInfoList := make([]fabricv4.ConnectionSideAdditionalInfo, len(additionalInfoList))
	for index, additionalInfo := range additionalInfoList {
		additionalInfoMap := additionalInfo.(map[string]interface{})
		key := additionalInfoMap["key"].(string)
		value := additionalInfoMap["value"].(string)

		additionalInfo := fabricv4.ConnectionSideAdditionalInfo{}
		additionalInfo.SetKey(key)
		additionalInfo.SetValue(value)
		mappedAdditionalInfoList[index] = additionalInfo
	}
	return mappedAdditionalInfoList
}

func connectionSideTerraformToGo(connectionSideTerraform []interface{}) fabricv4.ConnectionSide {
	if connectionSideTerraform == nil || len(connectionSideTerraform) == 0 {
		return fabricv4.ConnectionSide{}
	}

	var connectionSide fabricv4.ConnectionSide

	connectionSideMap := connectionSideTerraform[0].(map[string]interface{})
	accessPoint := connectionSideMap["access_point"].(*schema.Set).List()
	serviceTokenRequest := connectionSideMap["service_token"].(*schema.Set).List()
	additionalInfoRequest := connectionSideMap["additional_info"].([]interface{})
	if len(accessPoint) != 0 {
		ap := accessPointTerraformToGo(accessPoint)
		connectionSide.SetAccessPoint(ap)
	}
	if len(serviceTokenRequest) != 0 {
		serviceToken := serviceTokenTerraformToGo(serviceTokenRequest)
		connectionSide.SetServiceToken(serviceToken)
	}
	if len(additionalInfoRequest) != 0 {
		accessPointAdditionalInfo := additionalInfoTerraformToGo(additionalInfoRequest)
		connectionSide.SetAdditionalInfo(accessPointAdditionalInfo)
	}

	return connectionSide
}

func accessPointTerraformToGo(accessPointTerraform []interface{}) fabricv4.AccessPoint {
	if accessPointTerraform == nil || len(accessPointTerraform) == 0 {
		return fabricv4.AccessPoint{}
	}

	var accessPoint fabricv4.AccessPoint
	accessPointMap := accessPointTerraform[0].(map[string]interface{})
	portList := accessPointMap["port"].(*schema.Set).List()
	profileList := accessPointMap["profile"].(*schema.Set).List()
	locationList := accessPointMap["location"].(*schema.Set).List()
	virtualDeviceList := accessPointMap["virtual_device"].(*schema.Set).List()
	interfaceList := accessPointMap["interface"].(*schema.Set).List()
	networkList := accessPointMap["network"].(*schema.Set).List()
	typeVal := accessPointMap["type"].(string)
	authenticationKey := accessPointMap["authentication_key"].(string)
	if authenticationKey != "" {
		accessPoint.SetAuthenticationKey(authenticationKey)
	}
	providerConnectionId := accessPointMap["provider_connection_id"].(string)
	if providerConnectionId != "" {
		accessPoint.SetProviderConnectionId(providerConnectionId)
	}
	sellerRegion := accessPointMap["seller_region"].(string)
	if sellerRegion != "" {
		accessPoint.SetSellerRegion(sellerRegion)
	}
	peeringTypeRaw := accessPointMap["peering_type"].(string)
	if peeringTypeRaw != "" {
		peeringType := fabricv4.PeeringType(peeringTypeRaw)
		accessPoint.SetPeeringType(peeringType)
	}
	cloudRouterRequest := accessPointMap["router"].(*schema.Set).List()
	if len(cloudRouterRequest) == 0 {
		log.Print("[DEBUG] The router attribute was not used, attempting to revert to deprecated gateway attribute")
		cloudRouterRequest = accessPointMap["gateway"].(*schema.Set).List()
	}

	if len(cloudRouterRequest) != 0 {
		cloudRouter := cloudRouterTerraformToGo(cloudRouterRequest)
		if cloudRouter.GetUuid() != "" {
			accessPoint.SetRouter(cloudRouter)
		}
	}
	accessPoint.SetType(fabricv4.AccessPointType(typeVal))
	if len(portList) != 0 {
		port := portTerraformToGo(portList)
		if port.GetUuid() != "" {
			accessPoint.SetPort(port)
		}
	}

	if len(networkList) != 0 {
		network := networkTerraformToGo(networkList)
		if network.GetUuid() != "" {
			accessPoint.SetNetwork(network)
		}
	}
	linkProtocolList := accessPointMap["link_protocol"].(*schema.Set).List()

	if len(linkProtocolList) != 0 {
		linkProtocol := linkProtocolTerraformToGo(linkProtocolList)
		if linkProtocol.GetType().Ptr() != nil {
			accessPoint.SetLinkProtocol(linkProtocol)
		}
	}

	if len(profileList) != 0 {
		serviceProfile := simplifiedServiceProfileTerraformToGo(profileList)
		if serviceProfile.GetUuid() != "" {
			accessPoint.SetProfile(serviceProfile)
		}
	}

	if len(locationList) != 0 {
		location := equinix_fabric_schema.LocationTerraformToGo(locationList)
		accessPoint.SetLocation(location)
	}

	if len(virtualDeviceList) != 0 {
		virtualDevice := virtualDeviceTerraformToGo(virtualDeviceList)
		accessPoint.SetVirtualDevice(virtualDevice)
	}

	if len(interfaceList) != 0 {
		interface_ := interfaceTerraformToGo(interfaceList)
		accessPoint.SetInterface(interface_)
	}

	return accessPoint
}

func cloudRouterTerraformToGo(cloudRouterRequest []interface{}) fabricv4.CloudRouter {
	if cloudRouterRequest == nil || len(cloudRouterRequest) == 0 {
		return fabricv4.CloudRouter{}
	}
	var cloudRouter fabricv4.CloudRouter
	cloudRouterMap := cloudRouterRequest[0].(map[string]interface{})
	uuid := cloudRouterMap["uuid"].(string)
	cloudRouter.SetUuid(uuid)

	return cloudRouter
}

func linkProtocolTerraformToGo(linkProtocolList []interface{}) fabricv4.SimplifiedLinkProtocol {
	if linkProtocolList == nil || len(linkProtocolList) == 0 {
		return fabricv4.SimplifiedLinkProtocol{}
	}

	var linkProtocol fabricv4.SimplifiedLinkProtocol
	lpMap := linkProtocolList[0].(map[string]interface{})
	lpType := lpMap["type"].(string)
	lpVlanSTag := int32(lpMap["vlan_s_tag"].(int))
	lpVlanTag := int32(lpMap["vlan_tag"].(int))
	lpVlanCTag := int32(lpMap["vlan_c_tag"].(int))
	log.Printf("[DEBUG] linkProtocolMap: %v", lpMap)

	linkProtocol.SetType(fabricv4.LinkProtocolType(lpType))
	if lpVlanSTag != 0 {
		linkProtocol.SetVlanSTag(lpVlanSTag)
	}
	if lpVlanTag != 0 {
		linkProtocol.SetVlanTag(lpVlanTag)
	}
	if lpVlanCTag != 0 {
		linkProtocol.SetVlanCTag(lpVlanCTag)
	}

	return linkProtocol
}

func networkTerraformToGo(networkList []interface{}) fabricv4.SimplifiedNetwork {
	if networkList == nil || len(networkList) == 0 {
		return fabricv4.SimplifiedNetwork{}
	}
	var network fabricv4.SimplifiedNetwork
	networkListMap := networkList[0].(map[string]interface{})
	uuid := networkListMap["uuid"].(string)
	network.SetUuid(uuid)
	return network
}

func simplifiedServiceProfileTerraformToGo(profileList []interface{}) fabricv4.SimplifiedServiceProfile {
	if profileList == nil || len(profileList) == 0 {
		return fabricv4.SimplifiedServiceProfile{}
	}

	var serviceProfile fabricv4.SimplifiedServiceProfile
	profileListMap := profileList[0].(map[string]interface{})
	profileType := profileListMap["type"].(string)
	uuid := profileListMap["uuid"].(string)
	serviceProfile.SetType(fabricv4.ServiceProfileTypeEnum(profileType))
	serviceProfile.SetUuid(uuid)
	return serviceProfile
}

func virtualDeviceTerraformToGo(virtualDeviceList []interface{}) fabricv4.VirtualDevice {
	if virtualDeviceList == nil || len(virtualDeviceList) == 0 {
		return fabricv4.VirtualDevice{}
	}

	var virtualDevice fabricv4.VirtualDevice
	virtualDeviceMap := virtualDeviceList[0].(map[string]interface{})
	href := virtualDeviceMap["href"].(string)
	type_ := virtualDeviceMap["type"].(string)
	uuid := virtualDeviceMap["uuid"].(string)
	name := virtualDeviceMap["name"].(string)
	virtualDevice.SetHref(href)
	virtualDevice.SetType(fabricv4.VirtualDeviceType(type_))
	virtualDevice.SetUuid(uuid)
	virtualDevice.SetName(name)

	return virtualDevice
}

func interfaceTerraformToGo(interfaceList []interface{}) fabricv4.Interface {
	if interfaceList == nil || len(interfaceList) == 0 {
		return fabricv4.Interface{}
	}

	var interface_ fabricv4.Interface
	interfaceMap := interfaceList[0].(map[string]interface{})
	uuid := interfaceMap["uuid"].(string)
	type_ := interfaceMap["type"].(string)
	id := interfaceMap["id"].(int)
	interface_.SetUuid(uuid)
	interface_.SetType(fabricv4.InterfaceType(type_))
	interface_.SetId(int32(id))

	return interface_
}

func connectionOperationGoToTerraform(operation *fabricv4.ConnectionOperation) *schema.Set {
	if operation == nil {
		return nil
	}

	mappedOperation := make(map[string]interface{})
	mappedOperation["provider_status"] = string(operation.GetProviderStatus())
	mappedOperation["equinix_status"] = string(operation.GetEquinixStatus())
	if operation.Errors != nil {
		mappedOperation["errors"] = equinix_fabric_schema.ErrorGoToTerraform(operation.GetErrors())
	}
	operationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: operationSch()}),
		[]interface{}{mappedOperation},
	)
	return operationSet
}

func serviceTokenGoToTerraform(serviceToken *fabricv4.ServiceToken) *schema.Set {
	if serviceToken == nil {
		return nil
	}
	mappedServiceToken := make(map[string]interface{})
	if serviceToken.Type != nil {
		mappedServiceToken["type"] = string(serviceToken.GetType())
	}
	mappedServiceToken["href"] = serviceToken.GetHref()
	mappedServiceToken["uuid"] = serviceToken.GetUuid()

	serviceTokenSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: serviceTokenSch()}),
		[]interface{}{mappedServiceToken},
	)
	return serviceTokenSet
}

func connectionSideGoToTerraform(connectionSide *fabricv4.ConnectionSide) *schema.Set {
	mappedConnectionSide := make(map[string]interface{})
	serviceToken := connectionSide.GetServiceToken()
	serviceTokenSet := serviceTokenGoToTerraform(&serviceToken)
	if serviceTokenSet != nil {
		mappedConnectionSide["service_token"] = serviceTokenSet
	}
	accessPoint := connectionSide.GetAccessPoint()
	mappedConnectionSide["access_point"] = accessPointGoToTerraform(&accessPoint)
	connectionSideSet := schema.NewSet(
		schema.HashResource(connectionSideSch()),
		[]interface{}{mappedConnectionSide},
	)
	return connectionSideSet
}

func additionalInfoGoToTerraform(additionalInfo []fabricv4.ConnectionSideAdditionalInfo) []map[string]interface{} {
	if additionalInfo == nil {
		return nil
	}
	mappedAdditionalInfo := make([]map[string]interface{}, len(additionalInfo))
	for index, additionalInfo := range additionalInfo {
		mappedAdditionalInfo[index] = map[string]interface{}{
			"key":   additionalInfo.GetKey(),
			"value": additionalInfo.GetValue(),
		}
	}
	return mappedAdditionalInfo
}

func cloudRouterGoToTerraform(cloudRouter *fabricv4.CloudRouter) *schema.Set {
	if cloudRouter == nil {
		return nil
	}
	mappedCloudRouter := make(map[string]interface{})
	mappedCloudRouter["uuid"] = cloudRouter.GetUuid()
	mappedCloudRouter["href"] = cloudRouter.GetHref()

	linkedProtocolSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: equinix_fabric_schema.ProjectSch()}),
		[]interface{}{mappedCloudRouter})
	return linkedProtocolSet
}

func virtualDeviceGoToTerraform(virtualDevice *fabricv4.VirtualDevice) *schema.Set {
	if virtualDevice == nil {
		return nil
	}
	mappedVirtualDevice := make(map[string]interface{})
	mappedVirtualDevice["name"] = virtualDevice.GetName()
	mappedVirtualDevice["href"] = virtualDevice.GetHref()
	mappedVirtualDevice["type"] = string(virtualDevice.GetType())
	mappedVirtualDevice["uuid"] = virtualDevice.GetUuid()

	virtualDeviceSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointVirtualDeviceSch()}),
		[]interface{}{mappedVirtualDevice})
	return virtualDeviceSet
}

func interfaceGoToTerraform(mInterface *fabricv4.Interface) *schema.Set {
	if mInterface == nil {
		return nil
	}
	mappedMInterface := make(map[string]interface{})
	mappedMInterface["id"] = int(mInterface.GetId())
	mappedMInterface["type"] = string(mInterface.GetType())
	mappedMInterface["uuid"] = mInterface.GetUuid()

	mInterfaceSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointInterface()}),
		[]interface{}{mappedMInterface})
	return mInterfaceSet
}

func networkGoToTerraform(network *fabricv4.SimplifiedNetwork) *schema.Set {
	if network == nil {
		return nil
	}

	mappedNetwork := make(map[string]interface{})
	mappedNetwork["uuid"] = network.GetUuid()

	return schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: networkSch()}),
		[]interface{}{mappedNetwork})
}

func accessPointGoToTerraform(accessPoint *fabricv4.AccessPoint) *schema.Set {
	mappedAccessPoint := make(map[string]interface{})
	if accessPoint.Type != nil {
		mappedAccessPoint["type"] = string(accessPoint.GetType())
	}
	if accessPoint.Account != nil {
		account := accessPoint.GetAccount()
		mappedAccessPoint["account"] = equinix_fabric_schema.AccountGoToTerraform(&account)
	}
	if accessPoint.Location != nil {
		location := accessPoint.GetLocation()
		mappedAccessPoint["location"] = equinix_fabric_schema.LocationGoToTerraform(&location)
	}
	if accessPoint.Port != nil {
		port := accessPoint.GetPort()
		mappedAccessPoint["port"] = portGoToTerraform(&port)
	}
	if accessPoint.Profile != nil {
		profile := accessPoint.GetProfile()
		mappedAccessPoint["profile"] = simplifiedServiceProfileGoToTerraform(&profile)
	}
	if accessPoint.Router != nil {
		router := accessPoint.GetRouter()
		mappedAccessPoint["router"] = cloudRouterGoToTerraform(&router)
		mappedAccessPoint["gateway"] = cloudRouterGoToTerraform(&router)
	}
	if accessPoint.LinkProtocol != nil {
		linkProtocol := accessPoint.GetLinkProtocol()
		mappedAccessPoint["link_protocol"] = linkedProtocolGoToTerraform(&linkProtocol)
	}
	if accessPoint.VirtualDevice != nil {
		virtualDevice := accessPoint.GetVirtualDevice()
		mappedAccessPoint["virtual_device"] = virtualDeviceGoToTerraform(&virtualDevice)
	}
	if accessPoint.Interface != nil {
		interface_ := accessPoint.GetInterface()
		mappedAccessPoint["interface"] = interfaceGoToTerraform(&interface_)
	}
	if accessPoint.Network != nil {
		network := accessPoint.GetNetwork()
		mappedAccessPoint["network"] = networkGoToTerraform(&network)
	}
	mappedAccessPoint["seller_region"] = accessPoint.GetSellerRegion()
	if accessPoint.PeeringType != nil {
		mappedAccessPoint["peering_type"] = string(accessPoint.GetPeeringType())
	}
	mappedAccessPoint["authentication_key"] = accessPoint.GetAuthenticationKey()
	mappedAccessPoint["provider_connection_id"] = accessPoint.GetProviderConnectionId()

	accessPointSet := schema.NewSet(
		schema.HashResource(accessPointSch()),
		[]interface{}{mappedAccessPoint},
	)
	return accessPointSet
}

func linkedProtocolGoToTerraform(linkedProtocol *fabricv4.SimplifiedLinkProtocol) *schema.Set {

	mappedLinkedProtocol := make(map[string]interface{})
	mappedLinkedProtocol["type"] = string(linkedProtocol.GetType())
	mappedLinkedProtocol["vlan_tag"] = int(linkedProtocol.GetVlanTag())
	mappedLinkedProtocol["vlan_s_tag"] = int(linkedProtocol.GetVlanSTag())
	mappedLinkedProtocol["vlan_c_tag"] = int(linkedProtocol.GetVlanCTag())

	linkedProtocolSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointLinkProtocolSch()}),
		[]interface{}{mappedLinkedProtocol})
	return linkedProtocolSet
}

func simplifiedServiceProfileGoToTerraform(profile *fabricv4.SimplifiedServiceProfile) *schema.Set {

	mappedProfile := make(map[string]interface{})
	mappedProfile["href"] = profile.GetHref()
	mappedProfile["type"] = string(profile.GetType())
	mappedProfile["name"] = profile.GetName()
	mappedProfile["uuid"] = profile.GetName()
	mappedProfile["access_point_type_configs"] = accessPointTypeConfigGoToTerraform(profile.AccessPointTypeConfigs)

	profileSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: serviceProfileSch()}),
		[]interface{}{mappedProfile},
	)
	return profileSet
}

func getUpdateRequests(conn *fabricv4.Connection, d *schema.ResourceData) ([][]fabricv4.ConnectionChangeOperation, error) {
	var changeOps [][]fabricv4.ConnectionChangeOperation
	existingName := conn.GetName()
	existingBandwidth := int(conn.GetBandwidth())
	updateNameVal := d.Get("name").(string)
	updateBandwidthVal := d.Get("bandwidth").(int)
	additionalInfo := d.Get("additional_info").([]interface{})

	awsSecrets, hasAWSSecrets := additionalInfoContainsAWSSecrets(additionalInfo)

	if existingName != updateNameVal {
		changeOps = append(changeOps, []fabricv4.ConnectionChangeOperation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: updateNameVal,
			},
		})
	}

	if existingBandwidth != updateBandwidthVal {
		changeOps = append(changeOps, []fabricv4.ConnectionChangeOperation{
			{
				Op:    "replace",
				Path:  "/bandwidth",
				Value: updateBandwidthVal,
			},
		})
	}

	if *conn.Operation.ProviderStatus == fabricv4.PROVIDERSTATUS_PENDING_APPROVAL && hasAWSSecrets {
		changeOps = append(changeOps, []fabricv4.ConnectionChangeOperation{
			{
				Op:    "add",
				Path:  "",
				Value: map[string]interface{}{"additionalInfo": awsSecrets},
			},
		})
	}

	if len(changeOps) == 0 {
		return changeOps, fmt.Errorf("nothing to update for the connection %s", existingName)
	}

	return changeOps, nil
}
