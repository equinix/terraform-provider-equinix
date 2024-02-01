package equinix

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FabricConnectionResourceSchema() map[string]*schema.Schema {
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
				Schema: equinix_schema.OrderSch(),
			},
		},
		"notifications": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Preferences for notifications on connection configuration or status changes",
			Elem: &schema.Resource{
				Schema: equinix_schema.NotificationSch(),
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
			Description: "Project information",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_schema.ProjectSch(),
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
				Schema: equinix_schema.AccountSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures connection lifecycle change information",
			Elem: &schema.Resource{
				Schema: equinix_schema.ChangeLogSch(),
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
				ValidateFunc: validation.StringInSlice([]string{"COLO", "VD", "VG", "SP", "IGW", "SUBNET", "CLOUD_ROUTER", "NETWORK"}, true),
				Description:  "Access point type - COLO, VD, VG, SP, IGW, SUBNET, CLOUD_ROUTER, NETWORK",
			},
			"account": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Account",
				Elem: &schema.Resource{
					Schema: equinix_schema.AccountSch(),
				},
			},
			"location": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Access point location",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: equinix_schema.LocationSch(),
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
				Optional:    true,
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
				Schema: accessPointTypeConfigSch(),
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

func accessPointTypeConfigSch() map[string]*schema.Schema {
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
				Schema: equinix_schema.ErrorSch(),
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
			Description: "Redundancy group identifier (UUID of primary connection)",
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
		Schema: FabricConnectionResourceSchema(),

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric connection",
	}
}

func resourceFabricConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	conType := v4.ConnectionType(d.Get("type").(string))
	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := equinix_schema.NotificationsToFabric(schemaNotifications)
	schemaRedundancy := d.Get("redundancy").(*schema.Set).List()
	red := connectionRedundancyToFabric(schemaRedundancy)
	schemaOrder := d.Get("order").(*schema.Set).List()
	order := equinix_schema.OrderToFabric(schemaOrder)
	aside := d.Get("a_side").(*schema.Set).List()
	projectReq := d.Get("project").(*schema.Set).List()
	project := equinix_schema.ProjectToFabric(projectReq)
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
		return diag.FromErr(equinix_errors.FormatFabricError(err))
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
			return diag.FromErr(equinix_errors.FormatFabricError(patchErr))
		}

		if _, statusChangeErr := waitForConnectionProviderStatusChange(d.Id(), meta, ctx); statusChangeErr != nil {
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
		"order":           equinix_schema.OrderToTerra(conn.Order),
		"change_log":      equinix_schema.ChangeLogToTerra(conn.ChangeLog),
		"redundancy":      connectionRedundancyToTerra(conn.Redundancy),
		"notifications":   equinix_schema.NotificationsToTerra(conn.Notifications),
		"account":         accountToTerra(conn.Account),
		"a_side":          connectionSideToTerra(conn.ASide),
		"z_side":          connectionSideToTerra(conn.ZSide),
		"additional_info": additionalInfoToTerra(conn.AdditionalInfo),
		"project":         equinix_schema.ProjectToTerra(conn.Project),
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
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("connection property update request error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", equinix_errors.FormatFabricError(err), update)})
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
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("connection property update completion timeout error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", err, update)})
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
				return "", "", equinix_errors.FormatFabricError(err)
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
				return "", "", equinix_errors.FormatFabricError(err)
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
				return "", "", equinix_errors.FormatFabricError(err)
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
			if equinix_errors.HasModelErrorCode(errors, "EQ-3142509") {
				return diags
			}
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	err = WaitUntilConnectionDeprovisioned(d.Id(), meta, ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func WaitUntilConnectionDeprovisioned(uuid string, meta interface{}, ctx context.Context) error {
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
		Timeout:    6 * time.Minute,
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
