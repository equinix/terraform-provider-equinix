package equinix

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func fabricConnectionResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"EVPL_VC", "EPL_VC", "IP_VC", "IPWAN_VC", "ACCESS_EPL_VC", "EVPLAN_VC", "EPLAN_VC", "EIA_VC", "EC_VC"}, false),
			Description:  "Defines the connection type like EVPL_VC, EPL_VC, IPWAN_VC, IP_VC, ACCESS_EPL_VC, EVPLAN_VC, EPLAN_VC, EIA_VC, EC_VC",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 24),
			Description:  "Connection name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"order": {
			Type:        schema.TypeSet,
			Required:    true,
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
			ValidateFunc: validation.StringInSlice([]string{"L2_PROFILE", "L3_PROFILE", "ECIA_PROFILE", "ECMC_PROFILE"}, true),
			Description:  "Service profile type - L2_PROFILE, L3_PROFILE, ECIA_PROFILE, ECMC_PROFILE",
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
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	conType := v4.ConnectionType(d.Get("type").(string))
	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := equinix_fabric_schema.NotificationsToFabric(schemaNotifications)
	schemaRedundancy := d.Get("redundancy").(*schema.Set).List()
	red := connectionRedundancyToFabric(schemaRedundancy)
	schemaOrder := d.Get("order").(*schema.Set).List()
	order := equinix_fabric_schema.OrderToFabric(schemaOrder)
	aside := d.Get("a_side").(*schema.Set).List()
	terraConfigProject := d.Get("project").(*schema.Set).List()
	project := equinix_fabric_schema.ProjectToFabric(terraConfigProject)
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
		Project:        project,
	}

	start := time.Now()
	conn, _, err := client.ConnectionsApi.CreateConnection(ctx, createRequest)
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(conn.Uuid)

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if err = waitUntilConnectionIsCreated(d.Id(), meta, ctx, createTimeout); err != nil {
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
			return diag.FromErr(equinix_errors.FormatFabricError(patchErr))
		}

		createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
		if _, statusChangeErr := waitForConnectionProviderStatusChange(d.Id(), meta, ctx, createTimeout); statusChangeErr != nil {
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
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(conn.Uuid)
	return setFabricMap(d, conn)
}

func setFabricMap(d *schema.ResourceData, conn v4.Connection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	connection := make(map[string]interface{})
	connection["name"] = conn.Name
	connection["bandwidth"] = conn.Bandwidth
	connection["href"] = conn.Href
	connection["is_remote"] = conn.IsRemote
	connection["type"] = conn.Type_
	connection["state"] = conn.State
	connection["direction"] = conn.Direction
	if conn.Operation != nil {
		connection["operation"] = connectionOperationToTerra(conn.Operation)
	}
	if conn.Order != nil {
		connection["order"] = equinix_fabric_schema.OrderToTerra(conn.Order)
	}
	if conn.ChangeLog != nil {
		connection["change_log"] = equinix_fabric_schema.ChangeLogToTerra(conn.ChangeLog)
	}
	if conn.Redundancy != nil {
		connection["redundancy"] = connectionRedundancyToTerra(conn.Redundancy)
	}
	if conn.Notifications != nil {
		connection["notifications"] = equinix_fabric_schema.NotificationsToTerra(conn.Notifications)
	}
	if conn.Account != nil {
		connection["account"] = equinix_fabric_schema.AccountToTerra(conn.Account)
	}
	if conn.ASide != nil {
		connection["a_side"] = connectionSideToTerra(conn.ASide)
	}
	if conn.ZSide != nil {
		connection["z_side"] = connectionSideToTerra(conn.ZSide)
	}
	if conn.AdditionalInfo != nil {
		connection["additional_info"] = additionalInfoToTerra(conn.AdditionalInfo)
	}
	if conn.Project != nil {
		connection["project"] = equinix_fabric_schema.ProjectToTerra(conn.Project)
	}
	err := equinix_schema.SetMap(d, connection)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceFabricConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	start := time.Now()
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	dbConn, err := verifyConnectionCreated(d.Id(), meta, ctx, updateTimeout)
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
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("connection property update request error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", equinix_errors.FormatFabricError(err), update)})
			continue
		}

		var waitFunction func(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) (v4.Connection, error)
		if update[0].Op == "replace" {
			// Update type is either name or bandwidth
			waitFunction = waitForConnectionUpdateCompletion
		} else if update[0].Op == "add" {
			// Update type is aws secret additionalInfo
			waitFunction = waitForConnectionProviderStatusChange
		}

		updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
		conn, err := waitFunction(d.Id(), meta, ctx, updateTimeout)

		if err != nil {
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("connection property update completion timeout error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", err, update)})
		} else {
			updatedConn = conn
		}
	}

	d.SetId(updatedConn.Uuid)
	return append(diags, setFabricMap(d, updatedConn)...)
}

func waitForConnectionUpdateCompletion(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) (v4.Connection, error) {
	log.Printf("[DEBUG] Waiting for connection update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			updatableState := ""
			if dbConn.Change.Status == "COMPLETED" {
				updatableState = dbConn.Change.Status
			}
			return dbConn, updatableState, nil
		},
		Timeout:    timeout,
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

func waitUntilConnectionIsCreated(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) error {
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
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.State), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func waitForConnectionProviderStatusChange(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) (v4.Connection, error) {
	log.Printf("DEBUG: wating for provider status to update. Connection uuid: %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(v4.PENDING_APPROVAL_ProviderStatus),
			string(v4.PROVISIONING_ProviderStatus),
		},
		Target: []string{
			string(v4.PROVISIONED_ProviderStatus),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.Operation.ProviderStatus), nil
		},
		Timeout:    timeout,
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

func verifyConnectionCreated(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) (v4.Connection, error) {
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
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.State), nil
		},
		Timeout:    timeout,
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
	start := time.Now()
	_, _, err := client.ConnectionsApi.DeleteConnectionByUuid(ctx, d.Id())
	if err != nil {
		errors, ok := err.(v4.GenericSwaggerError).Model().([]v4.ModelError)
		if ok {
			// EQ-3142509 = Connection already deleted
			if equinix_errors.HasModelErrorCode(errors, "EQ-3142509") {
				return diags
			}
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	deleteTimeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)
	err = WaitUntilConnectionDeprovisioned(d.Id(), meta, ctx, deleteTimeout)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func WaitUntilConnectionDeprovisioned(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) error {
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
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.State), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func connectionRedundancyToFabric(schemaRedundancy []interface{}) v4.ConnectionRedundancy {
	if schemaRedundancy == nil {
		return v4.ConnectionRedundancy{}
	}
	red := v4.ConnectionRedundancy{}
	for _, r := range schemaRedundancy {
		redundancyMap := r.(map[string]interface{})
		connectionPriority := v4.ConnectionPriority(redundancyMap["priority"].(string))
		redundancyGroup := redundancyMap["group"].(string)
		red = v4.ConnectionRedundancy{
			Priority: &connectionPriority,
			Group:    redundancyGroup,
		}
	}
	return red
}

func connectionRedundancyToTerra(redundancy *v4.ConnectionRedundancy) *schema.Set {
	if redundancy == nil {
		return nil
	}
	redundancies := []*v4.ConnectionRedundancy{redundancy}
	mappedRedundancys := make([]interface{}, len(redundancies))
	for _, redundancy := range redundancies {
		mappedRedundancy := make(map[string]interface{})
		mappedRedundancy["group"] = redundancy.Group
		mappedRedundancy["priority"] = string(*redundancy.Priority)
		mappedRedundancys = append(mappedRedundancys, mappedRedundancy)
	}
	redundancySet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: connectionRedundancySch()}),
		mappedRedundancys,
	)
	return redundancySet
}

func serviceTokenToFabric(serviceTokenRequest []interface{}) (v4.ServiceToken, error) {
	mappedST := v4.ServiceToken{}
	for _, str := range serviceTokenRequest {
		stMap := str.(map[string]interface{})
		stType := stMap["type"].(string)
		uuid := stMap["uuid"].(string)
		if stType != "" {
			if stType != "VC_TOKEN" {
				return v4.ServiceToken{}, fmt.Errorf("invalid service token type in config. Must be: VC_TOKEN; Received: %s", stType)
			}
			stTypeObj := v4.ServiceTokenType(stType)
			mappedST = v4.ServiceToken{Uuid: uuid, Type_: &stTypeObj}
		} else {
			mappedST = v4.ServiceToken{Uuid: uuid}
		}

	}
	return mappedST, nil
}

func additionalInfoTerraToGo(additionalInfoRequest []interface{}) []v4.ConnectionSideAdditionalInfo {
	var mappedaiArray []v4.ConnectionSideAdditionalInfo
	for _, ai := range additionalInfoRequest {
		aiMap := ai.(map[string]interface{})
		key := aiMap["key"].(string)
		value := aiMap["value"].(string)
		mappedai := v4.ConnectionSideAdditionalInfo{Key: key, Value: value}
		mappedaiArray = append(mappedaiArray, mappedai)
	}
	return mappedaiArray
}

func accessPointToFabric(accessPointRequest []interface{}) v4.AccessPoint {
	accessPoint := v4.AccessPoint{}
	for _, ap := range accessPointRequest {
		accessPointMap := ap.(map[string]interface{})
		portList := accessPointMap["port"].(*schema.Set).List()
		profileList := accessPointMap["profile"].(*schema.Set).List()
		locationList := accessPointMap["location"].(*schema.Set).List()
		virtualdeviceList := accessPointMap["virtual_device"].(*schema.Set).List()
		interfaceList := accessPointMap["interface"].(*schema.Set).List()
		networkList := accessPointMap["network"].(*schema.Set).List()
		typeVal := accessPointMap["type"].(string)
		authenticationKey := accessPointMap["authentication_key"].(string)
		if authenticationKey != "" {
			accessPoint.AuthenticationKey = authenticationKey
		}
		providerConnectionId := accessPointMap["provider_connection_id"].(string)
		if providerConnectionId != "" {
			accessPoint.ProviderConnectionId = providerConnectionId
		}
		sellerRegion := accessPointMap["seller_region"].(string)
		if sellerRegion != "" {
			accessPoint.SellerRegion = sellerRegion
		}
		peeringTypeRaw := accessPointMap["peering_type"].(string)
		if peeringTypeRaw != "" {
			peeringType := v4.PeeringType(peeringTypeRaw)
			accessPoint.PeeringType = &peeringType
		}
		cloudRouterRequest := accessPointMap["router"].(*schema.Set).List()
		if len(cloudRouterRequest) == 0 {
			log.Print("[DEBUG] The router attribute was not used, attempting to revert to deprecated gateway attribute")
			cloudRouterRequest = accessPointMap["gateway"].(*schema.Set).List()
		}

		if len(cloudRouterRequest) != 0 {
			mappedGWr := cloudRouterToFabric(cloudRouterRequest)
			if mappedGWr.Uuid != "" {
				accessPoint.Router = &mappedGWr
			}
		}
		apt := v4.AccessPointType(typeVal)
		accessPoint.Type_ = &apt
		if len(portList) != 0 {
			port := portToFabric(portList)
			if port.Uuid != "" {
				accessPoint.Port = &port
			}
		}

		if len(networkList) != 0 {
			network := networkToFabric(networkList)
			if network.Uuid != "" {
				accessPoint.Network = &network
			}
		}
		linkProtocolList := accessPointMap["link_protocol"].(*schema.Set).List()

		if len(linkProtocolList) != 0 {
			slp := linkProtocolToFabric(linkProtocolList)
			if slp.Type_ != nil {
				accessPoint.LinkProtocol = &slp
			}
		}

		if len(profileList) != 0 {
			ssp := simplifiedServiceProfileToFabric(profileList)
			if ssp.Uuid != "" {
				accessPoint.Profile = &ssp
			}
		}

		if len(locationList) != 0 {
			sl := equinix_fabric_schema.LocationToFabric(locationList)
			accessPoint.Location = &sl
		}

		if len(virtualdeviceList) != 0 {
			vd := virtualdeviceToFabric(virtualdeviceList)
			accessPoint.VirtualDevice = &vd
		}

		if len(interfaceList) != 0 {
			il := interfaceToFabric(interfaceList)
			accessPoint.Interface_ = &il
		}

	}
	return accessPoint
}

func cloudRouterToFabric(cloudRouterRequest []interface{}) v4.CloudRouter {
	if cloudRouterRequest == nil {
		return v4.CloudRouter{}
	}
	cloudRouterMapped := v4.CloudRouter{}
	for _, crr := range cloudRouterRequest {
		crrMap := crr.(map[string]interface{})
		cruuid := crrMap["uuid"].(string)
		cloudRouterMapped = v4.CloudRouter{Uuid: cruuid}
	}
	return cloudRouterMapped
}

func linkProtocolToFabric(linkProtocolList []interface{}) v4.SimplifiedLinkProtocol {
	slp := v4.SimplifiedLinkProtocol{}
	for _, lp := range linkProtocolList {
		lpMap := lp.(map[string]interface{})
		lpType := lpMap["type"].(string)
		lpVlanSTag := lpMap["vlan_s_tag"].(int)
		lpVlanTag := lpMap["vlan_tag"].(int)
		lpVlanCTag := lpMap["vlan_c_tag"].(int)
		lpt := v4.LinkProtocolType(lpType)
		slp = v4.SimplifiedLinkProtocol{Type_: &lpt, VlanSTag: int32(lpVlanSTag), VlanTag: int32(lpVlanTag), VlanCTag: int32(lpVlanCTag)}
	}
	return slp
}

func networkToFabric(networkList []interface{}) v4.SimplifiedNetwork {
	p := v4.SimplifiedNetwork{}
	for _, pl := range networkList {
		plMap := pl.(map[string]interface{})
		uuid := plMap["uuid"].(string)
		p = v4.SimplifiedNetwork{Uuid: uuid}
	}
	return p
}

func simplifiedServiceProfileToFabric(profileList []interface{}) v4.SimplifiedServiceProfile {
	ssp := v4.SimplifiedServiceProfile{}
	for _, pl := range profileList {
		plMap := pl.(map[string]interface{})
		ptype := plMap["type"].(string)
		spte := v4.ServiceProfileTypeEnum(ptype)
		uuid := plMap["uuid"].(string)
		ssp = v4.SimplifiedServiceProfile{Uuid: uuid, Type_: &spte}

	}
	return ssp
}

func virtualdeviceToFabric(virtualdeviceList []interface{}) v4.VirtualDevice {
	vd := v4.VirtualDevice{}
	for _, ll := range virtualdeviceList {
		llMap := ll.(map[string]interface{})
		hr := llMap["href"].(string)
		tp := llMap["type"].(string)
		ud := llMap["uuid"].(string)
		na := llMap["name"].(string)
		vd = v4.VirtualDevice{Href: hr, Type_: tp, Uuid: ud, Name: na}
	}
	return vd
}

func interfaceToFabric(interfaceList []interface{}) v4.ModelInterface {
	il := v4.ModelInterface{}
	for _, ll := range interfaceList {
		llMap := ll.(map[string]interface{})
		ud := llMap["uuid"].(string)
		tp := llMap["type"].(string)
		id := llMap["id"].(int)
		il = v4.ModelInterface{Type_: tp, Uuid: ud, Id: int32(id)}
	}
	return il
}

func connectionOperationToTerra(operation *v4.ConnectionOperation) *schema.Set {
	if operation == nil {
		return nil
	}
	operations := []*v4.ConnectionOperation{operation}
	mappedOperations := make([]interface{}, len(operations))
	for _, operation := range operations {
		mappedOperation := make(map[string]interface{})
		mappedOperation["provider_status"] = string(*operation.ProviderStatus)
		mappedOperation["equinix_status"] = string(*operation.EquinixStatus)
		if operation.Errors != nil {
			mappedOperation["errors"] = equinix_fabric_schema.ErrorToTerra(operation.Errors)
		}
		mappedOperations = append(mappedOperations, mappedOperation)
	}
	operationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: operationSch()}),
		mappedOperations,
	)
	return operationSet
}

func serviceTokenToTerra(serviceToken *v4.ServiceToken) *schema.Set {
	if serviceToken == nil {
		return nil
	}
	serviceTokens := []*v4.ServiceToken{serviceToken}
	mappedServiceTokens := make([]interface{}, len(serviceTokens))
	for _, serviceToken := range serviceTokens {
		mappedServiceToken := make(map[string]interface{})
		if serviceToken.Type_ != nil {
			mappedServiceToken["type"] = string(*serviceToken.Type_)
		}
		mappedServiceToken["href"] = serviceToken.Href
		mappedServiceToken["uuid"] = serviceToken.Uuid
		mappedServiceTokens = append(mappedServiceTokens, mappedServiceToken)
	}
	serviceTokenSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: serviceTokenSch()}),
		mappedServiceTokens,
	)
	return serviceTokenSet
}

func connectionSideToTerra(connectionSide *v4.ConnectionSide) *schema.Set {
	connectionSides := []*v4.ConnectionSide{connectionSide}
	mappedConnectionSides := make([]interface{}, len(connectionSides))
	for _, connectionSide := range connectionSides {
		mappedConnectionSide := make(map[string]interface{})
		serviceTokenSet := serviceTokenToTerra(connectionSide.ServiceToken)
		if serviceTokenSet != nil {
			mappedConnectionSide["service_token"] = serviceTokenSet
		}
		mappedConnectionSide["access_point"] = accessPointToTerra(connectionSide.AccessPoint)
		mappedConnectionSides = append(mappedConnectionSides, mappedConnectionSide)
	}
	connectionSideSet := schema.NewSet(
		schema.HashResource(connectionSideSch()),
		mappedConnectionSides,
	)
	return connectionSideSet
}

func additionalInfoToTerra(additionalInfol []v4.ConnectionSideAdditionalInfo) []map[string]interface{} {
	if additionalInfol == nil {
		return nil
	}
	mappedadditionalInfol := make([]map[string]interface{}, len(additionalInfol))
	for index, additionalInfo := range additionalInfol {
		mappedadditionalInfol[index] = map[string]interface{}{
			"key":   additionalInfo.Key,
			"value": additionalInfo.Value,
		}
	}
	return mappedadditionalInfol
}

func cloudRouterToTerra(cloudRouter *v4.CloudRouter) *schema.Set {
	if cloudRouter == nil {
		return nil
	}
	cloudRouters := []*v4.CloudRouter{cloudRouter}
	mappedCloudRouters := make([]interface{}, len(cloudRouters))
	for _, cloudRouter := range cloudRouters {
		mappedCloudRouter := make(map[string]interface{})
		mappedCloudRouter["uuid"] = cloudRouter.Uuid
		mappedCloudRouter["href"] = cloudRouter.Href
		mappedCloudRouters = append(mappedCloudRouters, mappedCloudRouter)
	}
	linkedProtocolSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: equinix_fabric_schema.ProjectSch()}),
		mappedCloudRouters)
	return linkedProtocolSet
}

func virtualDeviceToTerra(virtualDevice *v4.VirtualDevice) *schema.Set {
	if virtualDevice == nil {
		return nil
	}
	virtualDevices := []*v4.VirtualDevice{virtualDevice}
	mappedVirtualDevices := make([]interface{}, len(virtualDevices))
	for _, virtualDevice := range virtualDevices {
		mappedVirtualDevice := make(map[string]interface{})
		mappedVirtualDevice["name"] = virtualDevice.Name
		mappedVirtualDevice["href"] = virtualDevice.Href
		mappedVirtualDevice["type"] = virtualDevice.Type_
		mappedVirtualDevice["uuid"] = virtualDevice.Uuid
		mappedVirtualDevices = append(mappedVirtualDevices, mappedVirtualDevice)
	}
	virtualDeviceSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointVirtualDeviceSch()}),
		mappedVirtualDevices)
	return virtualDeviceSet
}

func interfaceToTerra(mInterface *v4.ModelInterface) *schema.Set {
	if mInterface == nil {
		return nil
	}
	mInterfaces := []*v4.ModelInterface{mInterface}
	mappedMInterfaces := make([]interface{}, len(mInterfaces))
	for _, mInterface := range mInterfaces {
		mappedMInterface := make(map[string]interface{})
		mappedMInterface["id"] = int(mInterface.Id)
		mappedMInterface["type"] = mInterface.Type_
		mappedMInterface["uuid"] = mInterface.Uuid
		mappedMInterfaces = append(mappedMInterfaces, mappedMInterface)
	}
	mInterfaceSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointInterface()}),
		mappedMInterfaces)
	return mInterfaceSet
}

func accessPointToTerra(accessPoint *v4.AccessPoint) *schema.Set {
	accessPoints := []*v4.AccessPoint{accessPoint}
	mappedAccessPoints := make([]interface{}, len(accessPoints))
	for _, accessPoint := range accessPoints {
		mappedAccessPoint := make(map[string]interface{})
		if accessPoint.Type_ != nil {
			mappedAccessPoint["type"] = string(*accessPoint.Type_)
		}
		if accessPoint.Account != nil {
			mappedAccessPoint["account"] = equinix_fabric_schema.AccountToTerra(accessPoint.Account)
		}
		if accessPoint.Location != nil {
			mappedAccessPoint["location"] = equinix_fabric_schema.LocationToTerra(accessPoint.Location)
		}
		if accessPoint.Port != nil {
			mappedAccessPoint["port"] = portToTerra(accessPoint.Port)
		}
		if accessPoint.Profile != nil {
			mappedAccessPoint["profile"] = simplifiedServiceProfileToTerra(accessPoint.Profile)
		}
		if accessPoint.Router != nil {
			mappedAccessPoint["router"] = cloudRouterToTerra(accessPoint.Router)
			mappedAccessPoint["gateway"] = cloudRouterToTerra(accessPoint.Router)
		}
		if accessPoint.LinkProtocol != nil {
			mappedAccessPoint["link_protocol"] = linkedProtocolToTerra(*accessPoint.LinkProtocol)
		}
		if accessPoint.VirtualDevice != nil {
			mappedAccessPoint["virtual_device"] = virtualDeviceToTerra(accessPoint.VirtualDevice)
		}
		if accessPoint.Interface_ != nil {
			mappedAccessPoint["interface"] = interfaceToTerra(accessPoint.Interface_)
		}
		mappedAccessPoint["seller_region"] = accessPoint.SellerRegion
		if accessPoint.PeeringType != nil {
			mappedAccessPoint["peering_type"] = string(*accessPoint.PeeringType)
		}
		mappedAccessPoint["authentication_key"] = accessPoint.AuthenticationKey
		mappedAccessPoint["provider_connection_id"] = accessPoint.ProviderConnectionId
		mappedAccessPoints = append(mappedAccessPoints, mappedAccessPoint)
	}
	accessPointSet := schema.NewSet(
		schema.HashResource(accessPointSch()),
		mappedAccessPoints,
	)
	return accessPointSet
}

func linkedProtocolToTerra(linkedProtocol v4.SimplifiedLinkProtocol) *schema.Set {
	linkedProtocols := []v4.SimplifiedLinkProtocol{linkedProtocol}
	mappedLinkedProtocols := make([]interface{}, len(linkedProtocols))
	for _, linkedProtocol := range linkedProtocols {
		mappedLinkedProtocol := make(map[string]interface{})
		mappedLinkedProtocol["type"] = string(*linkedProtocol.Type_)
		mappedLinkedProtocol["vlan_tag"] = int(linkedProtocol.VlanTag)
		mappedLinkedProtocol["vlan_s_tag"] = int(linkedProtocol.VlanSTag)
		mappedLinkedProtocol["vlan_c_tag"] = int(linkedProtocol.VlanCTag)
		mappedLinkedProtocols = append(mappedLinkedProtocols, mappedLinkedProtocol)
	}
	linkedProtocolSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointLinkProtocolSch()}),
		mappedLinkedProtocols)
	return linkedProtocolSet
}

func simplifiedServiceProfileToTerra(profile *v4.SimplifiedServiceProfile) *schema.Set {
	profiles := []*v4.SimplifiedServiceProfile{profile}
	mappedProfiles := make([]interface{}, len(profiles))
	for _, profile := range profiles {
		mappedProfile := make(map[string]interface{})
		mappedProfile["href"] = profile.Href
		mappedProfile["type"] = string(*profile.Type_)
		mappedProfile["name"] = profile.Name
		mappedProfile["uuid"] = profile.Uuid
		mappedProfile["access_point_type_configs"] = accessPointTypeConfigToTerra(profile.AccessPointTypeConfigs)
		mappedProfiles = append(mappedProfiles, mappedProfile)
	}

	profileSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: serviceProfileSch()}),
		mappedProfiles,
	)
	return profileSet
}

func apiConfigToTerra(apiConfig *v4.ApiConfig) *schema.Set {
	apiConfigs := []*v4.ApiConfig{apiConfig}
	mappedApiConfigs := make([]interface{}, len(apiConfigs))
	for _, apiConfig := range apiConfigs {
		mappedApiConfig := make(map[string]interface{})
		mappedApiConfig["api_available"] = apiConfig.ApiAvailable
		mappedApiConfig["equinix_managed_vlan"] = apiConfig.EquinixManagedVlan
		mappedApiConfig["bandwidth_from_api"] = apiConfig.BandwidthFromApi
		mappedApiConfig["integration_id"] = apiConfig.IntegrationId
		mappedApiConfig["equinix_managed_port"] = apiConfig.EquinixManagedPort
		mappedApiConfigs = append(mappedApiConfigs, mappedApiConfig)
	}
	apiConfigSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createApiConfigSch()}),
		mappedApiConfigs)
	return apiConfigSet
}

func authenticationKeyToTerra(authenticationKey *v4.AuthenticationKey) *schema.Set {
	authenticationKeys := []*v4.AuthenticationKey{authenticationKey}
	mappedAuthenticationKeys := make([]interface{}, len(authenticationKeys))
	for _, authenticationKey := range authenticationKeys {
		mappedAuthenticationKey := make(map[string]interface{})
		mappedAuthenticationKey["required"] = authenticationKey.Required
		mappedAuthenticationKey["label"] = authenticationKey.Label
		mappedAuthenticationKey["description"] = authenticationKey.Description
		mappedAuthenticationKeys = append(mappedAuthenticationKeys, mappedAuthenticationKey)
	}
	apiConfigSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createAuthenticationKeySch()}),
		mappedAuthenticationKeys)
	return apiConfigSet
}

func supportedBandwidthsToTerra(supportedBandwidths *[]int32) []interface{} {
	if supportedBandwidths == nil {
		return nil
	}
	mappedSupportedBandwidths := make([]interface{}, len(*supportedBandwidths))
	for _, bandwidth := range *supportedBandwidths {
		mappedSupportedBandwidths = append(mappedSupportedBandwidths, int(bandwidth))
	}
	return mappedSupportedBandwidths
}

func getUpdateRequests(conn v4.Connection, d *schema.ResourceData) ([][]v4.ConnectionChangeOperation, error) {
	var changeOps [][]v4.ConnectionChangeOperation
	existingName := conn.Name
	existingBandwidth := int(conn.Bandwidth)
	updateNameVal := d.Get("name").(string)
	updateBandwidthVal := d.Get("bandwidth").(int)
	additionalInfo := d.Get("additional_info").([]interface{})

	awsSecrets, hasAWSSecrets := additionalInfoContainsAWSSecrets(additionalInfo)

	if existingName != updateNameVal {
		changeOps = append(changeOps, []v4.ConnectionChangeOperation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: updateNameVal,
			},
		})
	}

	if existingBandwidth != updateBandwidthVal {
		changeOps = append(changeOps, []v4.ConnectionChangeOperation{
			{
				Op:    "replace",
				Path:  "/bandwidth",
				Value: updateBandwidthVal,
			},
		})
	}

	if *conn.Operation.ProviderStatus == v4.PENDING_APPROVAL_ProviderStatus && hasAWSSecrets {
		changeOps = append(changeOps, []v4.ConnectionChangeOperation{
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
