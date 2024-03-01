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

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
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
				ValidateFunc: validation.StringInSlice([]string{"COLO", "VD", "VG", "SP", "IGW", "SUBNET", "CLOUD_ROUTER", "NETWORK"}, true),
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
	ctx = context.WithValue(ctx, fabricv4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	conType := fabricv4.ConnectionType(d.Get("type").(string))
	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := equinix_fabric_schema.NotificationsTerraformToGo(schemaNotifications)
	schemaRedundancy := d.Get("redundancy").(*schema.Set).List()
	redundancy := connectionRedundancyTerraformToGo(schemaRedundancy)
	schemaOrder := d.Get("order").(*schema.Set).List()
	order := equinix_fabric_schema.OrderTerraformToGo(schemaOrder)
	terraConfigProject := d.Get("project").(*schema.Set).List()
	project := equinix_fabric_schema.ProjectTerraformToGo(terraConfigProject)
	additionalInfoTerraConfig := d.Get("additional_info").([]interface{})
	additionalInfo := additionalInfoTerraformToGo(additionalInfoTerraConfig)

	aSide := d.Get("a_side").(*schema.Set).List()
	connectionASide := connectionSideTerraformToGo(aSide)

	zSide := d.Get("z_side").(*schema.Set).List()
	connectionZSide := connectionSideTerraformToGo(zSide)

	createConnectionRequest := fabricv4.ConnectionPostRequest{
		Name:           d.Get("name").(*string),
		Type:           &conType,
		Order:          order,
		Notifications:  notifications,
		Bandwidth:      d.Get("bandwidth").(*int32),
		AdditionalInfo: additionalInfo,
		Redundancy:     redundancy,
		ASide:          connectionASide,
		ZSide:          connectionZSide,
		Project:        project,
	}

	start := time.Now()
	conn, _, err := client.ConnectionsApi.CreateConnection(ctx).ConnectionPostRequest(createConnectionRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(*conn.Uuid)

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if err = waitUntilConnectionIsCreated(d.Id(), meta, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for connection (%s) to be created: %s", d.Id(), err)
	}

	awsSecrets, hasAWSSecrets := additionalInfoContainsAWSSecrets(additionalInfoTerraConfig)
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
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	conn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, d.Id()).Execute()
	if err != nil {
		log.Printf("[WARN] Connection %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(*conn.Uuid)
	return setFabricMap(d, conn)
}

func setFabricMap(d *schema.ResourceData, conn *fabricv4.Connection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	connection := make(map[string]interface{})
	connection["name"] = conn.Name
	connection["bandwidth"] = conn.Bandwidth
	connection["href"] = conn.Href
	connection["is_remote"] = conn.IsRemote
	connection["type"] = conn.Type
	connection["state"] = conn.State
	connection["direction"] = conn.Direction
	if conn.Operation != nil {
		connection["operation"] = connectionOperationGoToTerraform(conn.Operation)
	}
	if conn.Order != nil {
		connection["order"] = equinix_fabric_schema.OrderGoToTerraform(conn.Order)
	}
	if conn.ChangeLog != nil {
		connection["change_log"] = equinix_fabric_schema.ChangeLogGoToTerraform(conn.ChangeLog)
	}
	if conn.Redundancy != nil {
		connection["redundancy"] = connectionRedundancyGoToTerraform(conn.Redundancy)
	}
	if conn.Notifications != nil {
		connection["notifications"] = equinix_fabric_schema.NotificationsGoToTerraform(conn.Notifications)
	}
	if conn.Account != nil {
		connection["account"] = equinix_fabric_schema.AccountGoToTerraform(conn.Account)
	}
	if &conn.ASide != nil {
		connection["a_side"] = connectionSideGoToTerraform(&conn.ASide)
	}
	if &conn.ZSide != nil {
		connection["z_side"] = connectionSideGoToTerraform(&conn.ZSide)
	}
	if conn.AdditionalInfo != nil {
		connection["additional_info"] = additionalInfoGoToTerraform(conn.AdditionalInfo)
	}
	if conn.Project != nil {
		connection["project"] = equinix_fabric_schema.ProjectGoToTerraform(conn.Project)
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

func connectionRedundancyTerraformToGo(redundancyTerraform []interface{}) *fabricv4.ConnectionRedundancy {
	if redundancyTerraform == nil || len(redundancyTerraform) == 0 {
		return nil
	}
	var redundancy *fabricv4.ConnectionRedundancy

	redundancyMap := redundancyTerraform[0].(map[string]interface{})
	connectionPriority := fabricv4.ConnectionPriority(redundancyMap["priority"].(string))
	redundancyGroup := redundancyMap["group"].(*string)
	redundancy = &fabricv4.ConnectionRedundancy{
		Priority: &connectionPriority,
		Group:    redundancyGroup,
	}

	return redundancy
}

func connectionRedundancyGoToTerraform(redundancy *fabricv4.ConnectionRedundancy) *schema.Set {
	if redundancy == nil {
		return nil
	}
	mappedRedundancy := make(map[string]interface{})
	mappedRedundancy["group"] = redundancy.Group
	mappedRedundancy["priority"] = string(*redundancy.Priority)
	redundancySet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: connectionRedundancySch()}),
		[]interface{}{mappedRedundancy},
	)
	return redundancySet
}

func serviceTokenTerraformToGo(serviceTokenList []interface{}) *fabricv4.ServiceToken {
	if serviceTokenList == nil || len(serviceTokenList) == 0 {
		return nil
	}

	var serviceToken *fabricv4.ServiceToken
	serviceTokenMap := serviceTokenList[0].(map[string]interface{})
	serviceTokenTypeString := serviceTokenMap["type"].(string)
	uuid := serviceTokenMap["uuid"].(string)
	serviceTokenType := fabricv4.ServiceTokenType(serviceTokenTypeString)
	serviceToken = &fabricv4.ServiceToken{Uuid: uuid, Type: &serviceTokenType}

	return serviceToken
}

func additionalInfoTerraformToGo(additionalInfoList []interface{}) []fabricv4.ConnectionSideAdditionalInfo {
	if additionalInfoList == nil || len(additionalInfoList) == 0 {
		return nil
	}

	mappedAdditionalInfoList := make([]fabricv4.ConnectionSideAdditionalInfo, len(additionalInfoList))
	for index, additionalInfo := range additionalInfoList {
		additionalInfoMap := additionalInfo.(map[string]interface{})
		key := additionalInfoMap["key"].(*string)
		value := additionalInfoMap["value"].(*string)
		additionalInfo := fabricv4.ConnectionSideAdditionalInfo{Key: key, Value: value}
		mappedAdditionalInfoList[index] = additionalInfo
	}
	return mappedAdditionalInfoList
}

func connectionSideTerraformToGo(connectionSideTerraform []interface{}) *fabricv4.ConnectionSide {
	if connectionSideTerraform == nil || len(connectionSideTerraform) == 0 {
		return nil
	}

	var connectionSide *fabricv4.ConnectionSide

	connectionSideMap := connectionSideTerraform[0].(map[string]interface{})
	accessPoint := connectionSideMap["access_point"].(*schema.Set).List()
	serviceTokenRequest := connectionSideMap["service_token"].(*schema.Set).List()
	additionalInfoRequest := connectionSideMap["additional_info"].([]interface{})
	if len(accessPoint) != 0 {
		ap := accessPointTerraformToGo(accessPoint)
		connectionSide = &fabricv4.ConnectionSide{AccessPoint: ap}
	}
	if len(serviceTokenRequest) != 0 {
		serviceToken := serviceTokenTerraformToGo(serviceTokenRequest)
		connectionSide = &fabricv4.ConnectionSide{ServiceToken: serviceToken}
	}
	if len(additionalInfoRequest) != 0 {
		accessPointAdditionalInfo := additionalInfoTerraformToGo(additionalInfoRequest)
		connectionSide = &fabricv4.ConnectionSide{AdditionalInfo: accessPointAdditionalInfo}
	}

	return connectionSide
}

func accessPointTerraformToGo(accessPointTerraform []interface{}) *fabricv4.AccessPoint {
	if accessPointTerraform == nil || len(accessPointTerraform) == 0 {
		return nil
	}

	var accessPoint *fabricv4.AccessPoint
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
		accessPoint.AuthenticationKey = &authenticationKey
	}
	providerConnectionId := accessPointMap["provider_connection_id"].(string)
	if providerConnectionId != "" {
		accessPoint.ProviderConnectionId = &providerConnectionId
	}
	sellerRegion := accessPointMap["seller_region"].(string)
	if sellerRegion != "" {
		accessPoint.SellerRegion = &sellerRegion
	}
	peeringTypeRaw := accessPointMap["peering_type"].(string)
	if peeringTypeRaw != "" {
		peeringType := fabricv4.PeeringType(peeringTypeRaw)
		accessPoint.PeeringType = &peeringType
	}
	cloudRouterRequest := accessPointMap["router"].(*schema.Set).List()
	if len(cloudRouterRequest) == 0 {
		log.Print("[DEBUG] The router attribute was not used, attempting to revert to deprecated gateway attribute")
		cloudRouterRequest = accessPointMap["gateway"].(*schema.Set).List()
	}

	if len(cloudRouterRequest) != 0 {
		cloudRouter := cloudRouterTerraformToGo(cloudRouterRequest)
		if *cloudRouter.Uuid != "" {
			accessPoint.Router = cloudRouter
		}
	}
	apt := fabricv4.AccessPointType(typeVal)
	accessPoint.Type = &apt
	if len(portList) != 0 {
		port := portTerraformToGo(portList)
		if *port.Uuid != "" {
			accessPoint.Port = port
		}
	}

	if len(networkList) != 0 {
		network := networkTerraformToGo(networkList)
		if network.Uuid != "" {
			accessPoint.Network = network
		}
	}
	linkProtocolList := accessPointMap["link_protocol"].(*schema.Set).List()

	if len(linkProtocolList) != 0 {
		linkProtocol := linkProtocolTerraformToGo(linkProtocolList)
		if linkProtocol.Type != nil {
			accessPoint.LinkProtocol = linkProtocol
		}
	}

	if len(profileList) != 0 {
		serviceProfile := simplifiedServiceProfileTerraformToGo(profileList)
		if *serviceProfile.Uuid != "" {
			accessPoint.Profile = serviceProfile
		}
	}

	if len(locationList) != 0 {
		location := equinix_fabric_schema.LocationTerraformToGo(locationList)
		accessPoint.Location = location
	}

	if len(virtualDeviceList) != 0 {
		virtualDevice := virtualDeviceTerraformToGo(virtualDeviceList)
		accessPoint.VirtualDevice = virtualDevice
	}

	if len(interfaceList) != 0 {
		interface_ := interfaceTerraformToGo(interfaceList)
		accessPoint.Interface = interface_
	}

	return accessPoint
}

func cloudRouterTerraformToGo(cloudRouterRequest []interface{}) *fabricv4.CloudRouter {
	if cloudRouterRequest == nil || len(cloudRouterRequest) == 0 {
		return nil
	}
	var cloudRouterMapped *fabricv4.CloudRouter
	cloudRouterMap := cloudRouterRequest[0].(map[string]interface{})
	uuid := cloudRouterMap["uuid"].(*string)
	cloudRouterMapped = &fabricv4.CloudRouter{Uuid: uuid}

	return cloudRouterMapped
}

func linkProtocolTerraformToGo(linkProtocolList []interface{}) *fabricv4.SimplifiedLinkProtocol {
	if linkProtocolList == nil || len(linkProtocolList) == 0 {
		return nil
	}

	var linkProtocol *fabricv4.SimplifiedLinkProtocol
	lpMap := linkProtocolList[0].(map[string]interface{})
	lpType := lpMap["type"].(string)
	lpVlanSTag := lpMap["vlan_s_tag"].(*int32)
	lpVlanTag := lpMap["vlan_tag"].(*int32)
	lpVlanCTag := lpMap["vlan_c_tag"].(*int32)
	lpt := fabricv4.LinkProtocolType(lpType)
	linkProtocol = &fabricv4.SimplifiedLinkProtocol{Type: &lpt, VlanSTag: lpVlanSTag, VlanTag: lpVlanTag, VlanCTag: lpVlanCTag}

	return linkProtocol
}

func networkTerraformToGo(networkList []interface{}) *fabricv4.SimplifiedNetwork {
	if networkList == nil || len(networkList) == 0 {
		return nil
	}
	var network *fabricv4.SimplifiedNetwork
	networkListMap := networkList[0].(map[string]interface{})
	uuid := networkListMap["uuid"].(string)
	network = &fabricv4.SimplifiedNetwork{Uuid: uuid}
	return network
}

func simplifiedServiceProfileTerraformToGo(profileList []interface{}) *fabricv4.SimplifiedServiceProfile {
	if profileList == nil || len(profileList) == 0 {
		return nil
	}

	var serviceProfile *fabricv4.SimplifiedServiceProfile
	profileListMap := profileList[0].(map[string]interface{})
	profileListType := profileListMap["type"].(string)
	profileType := fabricv4.ServiceProfileTypeEnum(profileListType)
	uuid := profileListMap["uuid"].(*string)
	serviceProfile = &fabricv4.SimplifiedServiceProfile{Uuid: uuid, Type: &profileType}
	return serviceProfile
}

func virtualDeviceTerraformToGo(virtualDeviceList []interface{}) *fabricv4.VirtualDevice {
	if virtualDeviceList == nil || len(virtualDeviceList) == 0 {
		return nil
	}

	var virtualDevice *fabricv4.VirtualDevice
	virtualDeviceMap := virtualDeviceList[0].(map[string]interface{})
	href := virtualDeviceMap["href"].(*string)
	typeString := virtualDeviceMap["type"].(string)
	type_ := fabricv4.VirtualDeviceType(typeString)
	uuid := virtualDeviceMap["uuid"].(*string)
	name := virtualDeviceMap["name"].(*string)
	virtualDevice = &fabricv4.VirtualDevice{Href: href, Type: &type_, Uuid: uuid, Name: name}

	return virtualDevice
}

func interfaceTerraformToGo(interfaceList []interface{}) *fabricv4.Interface {
	if interfaceList == nil || len(interfaceList) == 0 {
		return nil
	}

	var interface_ *fabricv4.Interface
	interfaceMap := interfaceList[0].(map[string]interface{})
	uuid := interfaceMap["uuid"].(*string)
	typeString := interfaceMap["type"].(string)
	type_ := fabricv4.InterfaceType(typeString)
	id := interfaceMap["id"].(*int32)
	interface_ = &fabricv4.Interface{Type: &type_, Uuid: uuid, Id: id}

	return interface_
}

func connectionOperationGoToTerraform(operation *fabricv4.ConnectionOperation) *schema.Set {
	if operation == nil {
		return nil
	}
	operations := []*fabricv4.ConnectionOperation{operation}
	mappedOperations := make([]interface{}, len(operations))
	for _, operation := range operations {
		mappedOperation := make(map[string]interface{})
		mappedOperation["provider_status"] = string(*operation.ProviderStatus)
		mappedOperation["equinix_status"] = string(*operation.EquinixStatus)
		if operation.Errors != nil {
			mappedOperation["errors"] = equinix_fabric_schema.ErrorGoToTerraform(operation.Errors)
		}
		mappedOperations = append(mappedOperations, mappedOperation)
	}
	operationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: operationSch()}),
		mappedOperations,
	)
	return operationSet
}

func serviceTokenGoToTerraform(serviceToken *fabricv4.ServiceToken) *schema.Set {
	if serviceToken == nil {
		return nil
	}
	mappedServiceToken := make(map[string]interface{})
	if serviceToken.Type != nil {
		mappedServiceToken["type"] = string(*serviceToken.Type)
	}
	mappedServiceToken["href"] = serviceToken.Href
	mappedServiceToken["uuid"] = serviceToken.Uuid

	serviceTokenSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: serviceTokenSch()}),
		[]interface{}{mappedServiceToken},
	)
	return serviceTokenSet
}

func connectionSideGoToTerraform(connectionSide *fabricv4.ConnectionSide) *schema.Set {
	mappedConnectionSide := make(map[string]interface{})
	serviceTokenSet := serviceTokenGoToTerraform(connectionSide.ServiceToken)
	if serviceTokenSet != nil {
		mappedConnectionSide["service_token"] = serviceTokenSet
	}
	mappedConnectionSide["access_point"] = accessPointGoToTerraform(connectionSide.AccessPoint)
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
			"key":   additionalInfo.Key,
			"value": additionalInfo.Value,
		}
	}
	return mappedAdditionalInfo
}

func cloudRouterGoToTerraform(cloudRouter *fabricv4.CloudRouter) *schema.Set {
	if cloudRouter == nil {
		return nil
	}
	mappedCloudRouter := make(map[string]interface{})
	mappedCloudRouter["uuid"] = cloudRouter.Uuid
	mappedCloudRouter["href"] = cloudRouter.Href

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
	mappedVirtualDevice["name"] = virtualDevice.Name
	mappedVirtualDevice["href"] = virtualDevice.Href
	mappedVirtualDevice["type"] = virtualDevice.Type
	mappedVirtualDevice["uuid"] = virtualDevice.Uuid

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
	mappedMInterface["id"] = int(*mInterface.Id)
	mappedMInterface["type"] = mInterface.Type
	mappedMInterface["uuid"] = mInterface.Uuid

	mInterfaceSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointInterface()}),
		[]interface{}{mappedMInterface})
	return mInterfaceSet
}

func accessPointGoToTerraform(accessPoint *fabricv4.AccessPoint) *schema.Set {
	mappedAccessPoint := make(map[string]interface{})
	if accessPoint.Type != nil {
		mappedAccessPoint["type"] = string(*accessPoint.Type)
	}
	if accessPoint.Account != nil {
		mappedAccessPoint["account"] = equinix_fabric_schema.AccountGoToTerraform(accessPoint.Account)
	}
	if accessPoint.Location != nil {
		mappedAccessPoint["location"] = equinix_fabric_schema.LocationGoToTerraform(accessPoint.Location)
	}
	if accessPoint.Port != nil {
		mappedAccessPoint["port"] = portGoToTerraform(accessPoint.Port)
	}
	if accessPoint.Profile != nil {
		mappedAccessPoint["profile"] = simplifiedServiceProfileGoToTerraform(accessPoint.Profile)
	}
	if accessPoint.Router != nil {
		mappedAccessPoint["router"] = cloudRouterGoToTerraform(accessPoint.Router)
		mappedAccessPoint["gateway"] = cloudRouterGoToTerraform(accessPoint.Router)
	}
	if accessPoint.LinkProtocol != nil {
		mappedAccessPoint["link_protocol"] = linkedProtocolGoToTerraform(accessPoint.LinkProtocol)
	}
	if accessPoint.VirtualDevice != nil {
		mappedAccessPoint["virtual_device"] = virtualDeviceGoToTerraform(accessPoint.VirtualDevice)
	}
	if accessPoint.Interface != nil {
		mappedAccessPoint["interface"] = interfaceGoToTerraform(accessPoint.Interface)
	}
	mappedAccessPoint["seller_region"] = accessPoint.SellerRegion
	if accessPoint.PeeringType != nil {
		mappedAccessPoint["peering_type"] = string(*accessPoint.PeeringType)
	}
	mappedAccessPoint["authentication_key"] = accessPoint.AuthenticationKey
	mappedAccessPoint["provider_connection_id"] = accessPoint.ProviderConnectionId

	accessPointSet := schema.NewSet(
		schema.HashResource(accessPointSch()),
		[]interface{}{mappedAccessPoint},
	)
	return accessPointSet
}

func linkedProtocolGoToTerraform(linkedProtocol *fabricv4.SimplifiedLinkProtocol) *schema.Set {

	mappedLinkedProtocol := make(map[string]interface{})
	mappedLinkedProtocol["type"] = string(*linkedProtocol.Type)
	mappedLinkedProtocol["vlan_tag"] = int(*linkedProtocol.VlanTag)
	mappedLinkedProtocol["vlan_s_tag"] = int(*linkedProtocol.VlanSTag)
	mappedLinkedProtocol["vlan_c_tag"] = int(*linkedProtocol.VlanCTag)

	linkedProtocolSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: accessPointLinkProtocolSch()}),
		[]interface{}{mappedLinkedProtocol})
	return linkedProtocolSet
}

func simplifiedServiceProfileGoToTerraform(profile *fabricv4.SimplifiedServiceProfile) *schema.Set {

	mappedProfile := make(map[string]interface{})
	mappedProfile["href"] = profile.Href
	mappedProfile["type"] = string(*profile.Type)
	mappedProfile["name"] = profile.Name
	mappedProfile["uuid"] = profile.Uuid
	mappedProfile["access_point_type_configs"] = accessPointTypeConfigToTerra(profile.AccessPointTypeConfigs)

	profileSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: serviceProfileSch()}),
		[]interface{}{mappedProfile},
	)
	return profileSet
}

func apiConfigGoToTerraform(apiConfig *fabricv4.ApiConfig) *schema.Set {

	mappedApiConfig := make(map[string]interface{})
	mappedApiConfig["api_available"] = apiConfig.ApiAvailable
	mappedApiConfig["equinix_managed_vlan"] = apiConfig.EquinixManagedVlan
	mappedApiConfig["bandwidth_from_api"] = apiConfig.BandwidthFromApi
	mappedApiConfig["integration_id"] = apiConfig.IntegrationId
	mappedApiConfig["equinix_managed_port"] = apiConfig.EquinixManagedPort

	apiConfigSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createApiConfigSch()}),
		[]interface{}{mappedApiConfig})
	return apiConfigSet
}

func authenticationKeyGoToTerraform(authenticationKey *fabricv4.AuthenticationKey) *schema.Set {
	mappedAuthenticationKey := make(map[string]interface{})
	mappedAuthenticationKey["required"] = authenticationKey.Required
	mappedAuthenticationKey["label"] = authenticationKey.Label
	mappedAuthenticationKey["description"] = authenticationKey.Description

	apiConfigSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createAuthenticationKeySch()}),
		[]interface{}{mappedAuthenticationKey})
	return apiConfigSet
}

func supportedBandwidthsGoToTerraform(supportedBandwidths []int32) []interface{} {
	if supportedBandwidths == nil {
		return nil
	}
	mappedSupportedBandwidths := make([]interface{}, len(supportedBandwidths))
	for index, bandwidth := range supportedBandwidths {
		mappedSupportedBandwidths[index] = int(bandwidth)
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
