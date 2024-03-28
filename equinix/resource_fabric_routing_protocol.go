package equinix

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
	"strings"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createDirectConnectionIpv4Sch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"equinix_iface_ip": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Equinix side Interface IP address",
		},
	}
}

func createDirectConnectionIpv6Sch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"equinix_iface_ip": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Equinix side Interface IP address\n\n",
		},
	}
}

func createBgpConnectionIpv4Sch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"customer_peer_ip": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Customer side peering ip",
		},
		"equinix_peer_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix side peering ip",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Admin status for the BGP session",
		},
	}
}

func createBgpConnectionIpv6Sch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"customer_peer_ip": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Customer side peering ip",
		},
		"equinix_peer_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix side peering ip",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Admin status for the BGP session",
		},
	}
}

func createRoutingProtocolBfdSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Bidirectional Forwarding Detection enablement",
		},
		"interval": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     100,
			Description: "Interval range between the received BFD control packets",
		},
	}
}

func createRoutingProtocolOperationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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

func createRoutingProtocolChangeSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Uniquely identifies a change",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of change",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Routing Protocol Change URI",
		},
	}
}

func createFabricRoutingProtocolResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Connection URI associated with Routing Protocol",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Routing Protocol URI information",
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"BGP", "DIRECT"}, true),
			Description:  "Defines the routing protocol type like BGP or DIRECT",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Equinix-assigned routing protocol identifier",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Routing Protocol name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Customer-provided Fabric Routing Protocol description",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Routing Protocol overall state",
		},
		"operation": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Routing Protocol type-specific operational data",
			Elem: &schema.Resource{
				Schema: createRoutingProtocolOperationSch(),
			},
		},
		"change": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Routing Protocol configuration Changes",
			Elem: &schema.Resource{
				Schema: createRoutingProtocolChangeSch(),
			},
		},
		"direct_ipv4": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Routing Protocol Direct IPv4",
			Elem: &schema.Resource{
				Schema: createDirectConnectionIpv4Sch(),
			},
		},
		"direct_ipv6": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Routing Protocol Direct IPv6",
			Elem: &schema.Resource{
				Schema: createDirectConnectionIpv6Sch(),
			},
		},
		"bgp_ipv4": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Routing Protocol BGP IPv4",
			Elem: &schema.Resource{
				Schema: createBgpConnectionIpv4Sch(),
			},
		},
		"bgp_ipv6": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Routing Protocol BGP IPv6",
			Elem: &schema.Resource{
				Schema: createBgpConnectionIpv6Sch(),
			},
		},
		"customer_asn": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Customer-provided ASN",
		},
		"equinix_asn": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Equinix ASN",
		},
		"bgp_auth_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "BGP authorization key",
		},
		"bfd": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Bidirectional Forwarding Detection",
			Elem: &schema.Resource{
				Schema: createRoutingProtocolBfdSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures Routing Protocol lifecycle change information",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ChangeLogSch(),
			},
		},
	}
}

func resourceFabricRoutingProtocol() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
		},
		ReadContext:   resourceFabricRoutingProtocolRead,
		CreateContext: resourceFabricRoutingProtocolCreate,
		UpdateContext: resourceFabricRoutingProtocolUpdate,
		DeleteContext: resourceFabricRoutingProtocolDelete,
		Importer: &schema.ResourceImporter{
			// Custom state context function, to parse import argument as  connection_uuid/rp_uuid
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.SplitN(d.Id(), "/", 2)
				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%s), expected <conn-uuid>/<rp-uuid>", d.Id())
				}
				connectionUuid, uuid := parts[0], parts[1]
				// set set connection uuid and rp uuid as overall id of resource
				_ = d.Set("connection_uuid", connectionUuid)
				d.SetId(uuid)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: createFabricRoutingProtocolResourceSchema(),

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric connection",
	}
}

func resourceFabricRoutingProtocolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	log.Printf("[WARN] Routing Protocol Connection uuid: %s", d.Get("connection_uuid").(string))
	fabricRoutingProtocol, _, err := client.RoutingProtocolsApi.GetConnectionRoutingProtocolByUuid(ctx, d.Id(), d.Get("connection_uuid").(string))
	if err != nil {
		log.Printf("[WARN] Routing Protocol %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	switch fabricRoutingProtocol.Type_ {
	case "BGP":
		d.SetId(fabricRoutingProtocol.RoutingProtocolBgpData.Uuid)
	case "DIRECT":
		d.SetId(fabricRoutingProtocol.RoutingProtocolDirectData.Uuid)
	}

	return setFabricRoutingProtocolMap(d, fabricRoutingProtocol)
}

func resourceFabricRoutingProtocolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	schemaBgpIpv4 := d.Get("bgp_ipv4").(*schema.Set).List()
	bgpIpv4 := routingProtocolBgpIpv4ToFabric(schemaBgpIpv4)
	schemaBgpIpv6 := d.Get("bgp_ipv6").(*schema.Set).List()
	bgpIpv6 := routingProtocolBgpIpv6ToFabric(schemaBgpIpv6)
	schemaDirectIpv4 := d.Get("direct_ipv4").(*schema.Set).List()
	directIpv4 := routingProtocolDirectIpv4ToFabric(schemaDirectIpv4)
	schemaDirectIpv6 := d.Get("direct_ipv6").(*schema.Set).List()
	directIpv6 := routingProtocolDirectIpv6ToFabric(schemaDirectIpv6)
	schemaBfd := d.Get("bfd").(*schema.Set).List()
	bfd := routingProtocolBfdToFabric(schemaBfd)
	bgpAuthKey := d.Get("bgp_auth_key")
	if bgpAuthKey == nil {
		bgpAuthKey = ""
	}

	createRequest := v4.RoutingProtocolBase{}
	if d.Get("type").(string) == "BGP" {
		createRequest = v4.RoutingProtocolBase{
			Type_: d.Get("type").(string),
			OneOfRoutingProtocolBase: v4.OneOfRoutingProtocolBase{
				RoutingProtocolBgpType: v4.RoutingProtocolBgpType{
					Type_:       d.Get("type").(string),
					Name:        d.Get("name").(string),
					BgpIpv4:     &bgpIpv4,
					BgpIpv6:     &bgpIpv6,
					CustomerAsn: int64(d.Get("customer_asn").(int)),
					EquinixAsn:  int64(d.Get("equinix_asn").(int)),
					BgpAuthKey:  bgpAuthKey.(string),
					Bfd:         &bfd,
				},
			},
		}
		if bgpIpv4.CustomerPeerIp == "" {
			createRequest.BgpIpv4 = nil
		}
		if bgpIpv6.CustomerPeerIp == "" {
			createRequest.BgpIpv6 = nil
		}
		if bfd.Enabled == false {
			createRequest.Bfd = nil
		}
	}
	if d.Get("type").(string) == "DIRECT" {
		createRequest = v4.RoutingProtocolBase{
			Type_: d.Get("type").(string),
			OneOfRoutingProtocolBase: v4.OneOfRoutingProtocolBase{
				RoutingProtocolDirectType: v4.RoutingProtocolDirectType{
					Type_:      d.Get("type").(string),
					Name:       d.Get("name").(string),
					DirectIpv4: &directIpv4,
					DirectIpv6: &directIpv6,
				},
			},
		}
		if directIpv4.EquinixIfaceIp == "" {
			createRequest.DirectIpv4 = nil
		}
		if directIpv6.EquinixIfaceIp == "" {
			createRequest.DirectIpv6 = nil
		}
	}

	start := time.Now()
	fabricRoutingProtocol, _, err := client.RoutingProtocolsApi.CreateConnectionRoutingProtocol(ctx, createRequest, d.Get("connection_uuid").(string))
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	switch fabricRoutingProtocol.Type_ {
	case "BGP":
		d.SetId(fabricRoutingProtocol.RoutingProtocolBgpData.Uuid)
	case "DIRECT":
		d.SetId(fabricRoutingProtocol.RoutingProtocolDirectData.Uuid)
	}

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if _, err = waitUntilRoutingProtocolIsProvisioned(d.Id(), d.Get("connection_uuid").(string), meta, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for RP (%s) to be created: %s", d.Id(), err)
	}

	return resourceFabricRoutingProtocolRead(ctx, d, meta)
}

func resourceFabricRoutingProtocolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)

	schemaBgpIpv4 := d.Get("bgp_ipv4").(*schema.Set).List()
	bgpIpv4 := routingProtocolBgpIpv4ToFabric(schemaBgpIpv4)
	schemaBgpIpv6 := d.Get("bgp_ipv6").(*schema.Set).List()
	bgpIpv6 := routingProtocolBgpIpv6ToFabric(schemaBgpIpv6)
	schemaDirectIpv4 := d.Get("direct_ipv4").(*schema.Set).List()
	directIpv4 := routingProtocolDirectIpv4ToFabric(schemaDirectIpv4)
	schemaDirectIpv6 := d.Get("direct_ipv6").(*schema.Set).List()
	directIpv6 := routingProtocolDirectIpv6ToFabric(schemaDirectIpv6)
	schemaBfd := d.Get("bfd").(*schema.Set).List()
	bfd := routingProtocolBfdToFabric(schemaBfd)
	bgpAuthKey := d.Get("bgp_auth_key")
	if bgpAuthKey == nil {
		bgpAuthKey = ""
	}

	updateRequest := v4.RoutingProtocolBase{}
	if d.Get("type").(string) == "BGP" {
		updateRequest = v4.RoutingProtocolBase{
			Type_: d.Get("type").(string),
			OneOfRoutingProtocolBase: v4.OneOfRoutingProtocolBase{
				RoutingProtocolBgpType: v4.RoutingProtocolBgpType{
					Type_:       d.Get("type").(string),
					Name:        d.Get("name").(string),
					BgpIpv4:     &bgpIpv4,
					BgpIpv6:     &bgpIpv6,
					CustomerAsn: int64(d.Get("customer_asn").(int)),
					EquinixAsn:  int64(d.Get("equinix_asn").(int)),
					BgpAuthKey:  bgpAuthKey.(string),
					Bfd:         &bfd,
				},
			},
		}
		if bgpIpv4.CustomerPeerIp == "" {
			updateRequest.BgpIpv4 = nil
		}
		if bgpIpv6.CustomerPeerIp == "" {
			updateRequest.BgpIpv6 = nil
		}
	}
	if d.Get("type").(string) == "DIRECT" {
		updateRequest = v4.RoutingProtocolBase{
			Type_: d.Get("type").(string),
			OneOfRoutingProtocolBase: v4.OneOfRoutingProtocolBase{
				RoutingProtocolDirectType: v4.RoutingProtocolDirectType{
					Type_:      d.Get("type").(string),
					Name:       d.Get("name").(string),
					DirectIpv4: &directIpv4,
					DirectIpv6: &directIpv6,
				},
			},
		}
		if directIpv4.EquinixIfaceIp == "" {
			updateRequest.DirectIpv4 = nil
		}
		if directIpv6.EquinixIfaceIp == "" {
			updateRequest.DirectIpv6 = nil
		}
	}

	start := time.Now()
	updatedRpResp, _, err := client.RoutingProtocolsApi.ReplaceConnectionRoutingProtocolByUuid(ctx, updateRequest, d.Id(), d.Get("connection_uuid").(string))
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	var changeUuid string
	switch updatedRpResp.Type_ {
	case "BGP":
		changeUuid = updatedRpResp.RoutingProtocolBgpData.Change.Uuid
		d.SetId(updatedRpResp.RoutingProtocolBgpData.Uuid)
	case "DIRECT":
		changeUuid = updatedRpResp.RoutingProtocolDirectData.Change.Uuid
		d.SetId(updatedRpResp.RoutingProtocolDirectData.Uuid)
	}
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	_, err = waitForRoutingProtocolUpdateCompletion(changeUuid, d.Id(), d.Get("connection_uuid").(string), meta, ctx, updateTimeout)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("timeout updating routing protocol: %v", err))
	}
	updateTimeout = d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	updatedProvisionedRpResp, err := waitUntilRoutingProtocolIsProvisioned(d.Id(), d.Get("connection_uuid").(string), meta, ctx, updateTimeout)
	if err != nil {
		return diag.Errorf("error waiting for RP (%s) to be replace updated: %s", d.Id(), err)
	}

	return setFabricRoutingProtocolMap(d, updatedProvisionedRpResp)
}

func resourceFabricRoutingProtocolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	start := time.Now()
	_, _, err := client.RoutingProtocolsApi.DeleteConnectionRoutingProtocolByUuid(ctx, d.Id(), d.Get("connection_uuid").(string))
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
	err = WaitUntilRoutingProtocolIsDeprovisioned(d.Id(), d.Get("connection_uuid").(string), meta, ctx, deleteTimeout)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}

	return diags
}

func setFabricRoutingProtocolMap(d *schema.ResourceData, rp v4.RoutingProtocolData) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := error(nil)
	if rp.Type_ == "BGP" {
		err = equinix_schema.SetMap(d, map[string]interface{}{
			"name":         rp.RoutingProtocolBgpData.Name,
			"href":         rp.RoutingProtocolBgpData.Href,
			"type":         rp.RoutingProtocolBgpData.Type_,
			"state":        rp.RoutingProtocolBgpData.State,
			"operation":    routingProtocolOperationToTerra(rp.RoutingProtocolBgpData.Operation),
			"bgp_ipv4":     routingProtocolBgpConnectionIpv4ToTerra(rp.BgpIpv4),
			"bgp_ipv6":     routingProtocolBgpConnectionIpv6ToTerra(rp.BgpIpv6),
			"customer_asn": rp.CustomerAsn,
			"equinix_asn":  rp.EquinixAsn,
			"bfd":          routingProtocolBfdToTerra(rp.Bfd),
			"bgp_auth_key": rp.BgpAuthKey,
			"change":       routingProtocolChangeToTerra(rp.RoutingProtocolBgpData.Change),
			"change_log":   equinix_fabric_schema.ChangeLogToTerra(rp.RoutingProtocolBgpData.Changelog),
		})
	} else if rp.Type_ == "DIRECT" {
		err = equinix_schema.SetMap(d, map[string]interface{}{
			"name":        rp.RoutingProtocolDirectData.Name,
			"href":        rp.RoutingProtocolDirectData.Href,
			"type":        rp.RoutingProtocolDirectData.Type_,
			"state":       rp.RoutingProtocolDirectData.State,
			"operation":   routingProtocolOperationToTerra(rp.RoutingProtocolDirectData.Operation),
			"direct_ipv4": routingProtocolDirectConnectionIpv4ToTerra(rp.DirectIpv4),
			"direct_ipv6": routingProtocolDirectConnectionIpv6ToTerra(rp.DirectIpv6),
			"change":      routingProtocolChangeToTerra(rp.RoutingProtocolDirectData.Change),
			"change_log":  equinix_fabric_schema.ChangeLogToTerra(rp.RoutingProtocolDirectData.Changelog),
		})
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func waitUntilRoutingProtocolIsProvisioned(uuid string, connUuid string, meta interface{}, ctx context.Context, timeout time.Duration) (v4.RoutingProtocolData, error) {
	log.Printf("Waiting for routing protocol to be provisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(v4.PROVISIONING_ConnectionState),
			string(v4.REPROVISIONING_ConnectionState),
		},
		Target: []string{
			string(v4.PROVISIONED_ConnectionState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.RoutingProtocolsApi.GetConnectionRoutingProtocolByUuid(ctx, uuid, connUuid)
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			var state string
			if dbConn.Type_ == "BGP" {
				state = dbConn.RoutingProtocolBgpData.State
			} else if dbConn.Type_ == "DIRECT" {
				state = dbConn.RoutingProtocolDirectData.State
			}
			return dbConn, state, nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	dbConn := v4.RoutingProtocolData{}

	if err == nil {
		dbConn = inter.(v4.RoutingProtocolData)
	}

	return dbConn, err
}

func WaitUntilRoutingProtocolIsDeprovisioned(uuid string, connUuid string, meta interface{}, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for routing protocol to be deprovisioned, uuid %s", uuid)

	/* check if resource is not found */
	stateConf := &retry.StateChangeConf{
		Target: []string{
			strconv.Itoa(404),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, resp, _ := client.RoutingProtocolsApi.GetConnectionRoutingProtocolByUuid(ctx, uuid, connUuid)
			// fixme: check for error code instead?
			// ignore error for Target
			return dbConn, strconv.Itoa(resp.StatusCode), nil

		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForRoutingProtocolUpdateCompletion(rpChangeUuid string, uuid string, connUuid string, meta interface{}, ctx context.Context, timeout time.Duration) (v4.RoutingProtocolChangeData, error) {
	log.Printf("Waiting for routing protocol update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.RoutingProtocolsApi.GetConnectionRoutingProtocolsChangeByUuid(ctx, connUuid, uuid, rpChangeUuid)
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			updatableState := ""
			if dbConn.Status == "COMPLETED" {
				updatableState = dbConn.Status
			}
			return dbConn, updatableState, nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	dbConn := v4.RoutingProtocolChangeData{}

	if err == nil {
		dbConn = inter.(v4.RoutingProtocolChangeData)
	}
	return dbConn, err
}

func routingProtocolDirectIpv4ToFabric(routingProtocolDirectIpv4Request []interface{}) v4.DirectConnectionIpv4 {
	mappedRpDirectIpv4 := v4.DirectConnectionIpv4{}
	for _, str := range routingProtocolDirectIpv4Request {
		directIpv4Map := str.(map[string]interface{})
		equinixIfaceIp := directIpv4Map["equinix_iface_ip"].(string)

		mappedRpDirectIpv4 = v4.DirectConnectionIpv4{EquinixIfaceIp: equinixIfaceIp}
	}
	return mappedRpDirectIpv4
}

func routingProtocolDirectIpv6ToFabric(routingProtocolDirectIpv6Request []interface{}) v4.DirectConnectionIpv6 {
	mappedRpDirectIpv6 := v4.DirectConnectionIpv6{}
	for _, str := range routingProtocolDirectIpv6Request {
		directIpv6Map := str.(map[string]interface{})
		equinixIfaceIp := directIpv6Map["equinix_iface_ip"].(string)

		mappedRpDirectIpv6 = v4.DirectConnectionIpv6{EquinixIfaceIp: equinixIfaceIp}
	}
	return mappedRpDirectIpv6
}

func routingProtocolBgpIpv4ToFabric(routingProtocolBgpIpv4Request []interface{}) v4.BgpConnectionIpv4 {
	mappedRpBgpIpv4 := v4.BgpConnectionIpv4{}
	for _, str := range routingProtocolBgpIpv4Request {
		bgpIpv4Map := str.(map[string]interface{})
		customerPeerIp := bgpIpv4Map["customer_peer_ip"].(string)
		enabled := bgpIpv4Map["enabled"].(bool)

		mappedRpBgpIpv4 = v4.BgpConnectionIpv4{CustomerPeerIp: customerPeerIp, Enabled: enabled}
	}
	return mappedRpBgpIpv4
}

func routingProtocolBgpIpv6ToFabric(routingProtocolBgpIpv6Request []interface{}) v4.BgpConnectionIpv6 {
	mappedRpBgpIpv6 := v4.BgpConnectionIpv6{}
	for _, str := range routingProtocolBgpIpv6Request {
		bgpIpv6Map := str.(map[string]interface{})
		customerPeerIp := bgpIpv6Map["customer_peer_ip"].(string)
		enabled := bgpIpv6Map["enabled"].(bool)

		mappedRpBgpIpv6 = v4.BgpConnectionIpv6{CustomerPeerIp: customerPeerIp, Enabled: enabled}
	}
	return mappedRpBgpIpv6
}

func routingProtocolBfdToFabric(routingProtocolBfdRequest []interface{}) v4.RoutingProtocolBfd {
	mappedRpBfd := v4.RoutingProtocolBfd{}
	for _, str := range routingProtocolBfdRequest {
		rpBfdMap := str.(map[string]interface{})
		bfdEnabled := rpBfdMap["enabled"].(bool)
		bfdInterval := rpBfdMap["interval"].(string)

		mappedRpBfd = v4.RoutingProtocolBfd{Enabled: bfdEnabled, Interval: bfdInterval}
	}
	return mappedRpBfd
}

func routingProtocolDirectConnectionIpv4ToTerra(routingProtocolDirectIpv4 *v4.DirectConnectionIpv4) *schema.Set {
	if routingProtocolDirectIpv4 == nil {
		return nil
	}
	routingProtocolDirectIpv4s := []*v4.DirectConnectionIpv4{routingProtocolDirectIpv4}
	mappedDirectIpv4s := make([]interface{}, len(routingProtocolDirectIpv4s))
	for i, routingProtocolDirectIpv4 := range routingProtocolDirectIpv4s {
		mappedDirectIpv4s[i] = map[string]interface{}{
			"equinix_iface_ip": routingProtocolDirectIpv4.EquinixIfaceIp,
		}
	}
	rpDirectIpv4Set := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createDirectConnectionIpv4Sch()}),
		mappedDirectIpv4s,
	)
	return rpDirectIpv4Set
}

func routingProtocolDirectConnectionIpv6ToTerra(routingProtocolDirectIpv6 *v4.DirectConnectionIpv6) *schema.Set {
	if routingProtocolDirectIpv6 == nil {
		return nil
	}
	routingProtocolDirectIpv6s := []*v4.DirectConnectionIpv6{routingProtocolDirectIpv6}
	mappedDirectIpv6s := make([]interface{}, len(routingProtocolDirectIpv6s))
	for i, routingProtocolDirectIpv6 := range routingProtocolDirectIpv6s {
		mappedDirectIpv6s[i] = map[string]interface{}{
			"equinix_iface_ip": routingProtocolDirectIpv6.EquinixIfaceIp,
		}
	}
	rpDirectIpv6Set := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createDirectConnectionIpv6Sch()}),
		mappedDirectIpv6s,
	)
	return rpDirectIpv6Set
}

func routingProtocolBgpConnectionIpv4ToTerra(routingProtocolBgpIpv4 *v4.BgpConnectionIpv4) *schema.Set {
	if routingProtocolBgpIpv4 == nil {
		return nil
	}
	routingProtocolBgpIpv4s := []*v4.BgpConnectionIpv4{routingProtocolBgpIpv4}
	mappedBgpIpv4s := make([]interface{}, len(routingProtocolBgpIpv4s))
	for i, routingProtocolBgpIpv4 := range routingProtocolBgpIpv4s {
		mappedBgpIpv4s[i] = map[string]interface{}{
			"customer_peer_ip": routingProtocolBgpIpv4.CustomerPeerIp,
			"equinix_peer_ip":  routingProtocolBgpIpv4.EquinixPeerIp,
			"enabled":          routingProtocolBgpIpv4.Enabled,
		}
	}
	rpBgpIpv4Set := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createBgpConnectionIpv4Sch()}),
		mappedBgpIpv4s,
	)
	return rpBgpIpv4Set
}

func routingProtocolBgpConnectionIpv6ToTerra(routingProtocolBgpIpv6 *v4.BgpConnectionIpv6) *schema.Set {
	if routingProtocolBgpIpv6 == nil {
		return nil
	}
	routingProtocolBgpIpv6s := []*v4.BgpConnectionIpv6{routingProtocolBgpIpv6}
	mappedBgpIpv6s := make([]interface{}, len(routingProtocolBgpIpv6s))
	for i, routingProtocolBgpIpv6 := range routingProtocolBgpIpv6s {
		mappedBgpIpv6s[i] = map[string]interface{}{
			"customer_peer_ip": routingProtocolBgpIpv6.CustomerPeerIp,
			"equinix_peer_ip":  routingProtocolBgpIpv6.EquinixPeerIp,
			"enabled":          routingProtocolBgpIpv6.Enabled,
		}
	}
	rpBgpIpv6Set := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createBgpConnectionIpv6Sch()}),
		mappedBgpIpv6s,
	)
	return rpBgpIpv6Set
}

func routingProtocolBfdToTerra(routingProtocolBfd *v4.RoutingProtocolBfd) *schema.Set {
	if routingProtocolBfd == nil {
		return nil
	}
	routingProtocolBfds := []*v4.RoutingProtocolBfd{routingProtocolBfd}
	mappedRpBfds := make([]interface{}, len(routingProtocolBfds))
	for i, routingProtocolBfd := range routingProtocolBfds {
		mappedRpBfds[i] = map[string]interface{}{
			"enabled":  routingProtocolBfd.Enabled,
			"interval": routingProtocolBfd.Interval,
		}
	}
	rpBfdSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createRoutingProtocolBfdSch()}),
		mappedRpBfds,
	)
	return rpBfdSet
}

func routingProtocolOperationToTerra(routingProtocolOperation *v4.RoutingProtocolOperation) *schema.Set {
	if routingProtocolOperation == nil {
		return nil
	}
	routingProtocolOperations := []*v4.RoutingProtocolOperation{routingProtocolOperation}
	mappedRpOperations := make([]interface{}, len(routingProtocolOperations))
	for _, routingProtocolOperation := range routingProtocolOperations {
		mappedRpOperation := make(map[string]interface{})
		if routingProtocolOperation.Errors != nil {
			mappedRpOperation["errors"] = equinix_fabric_schema.ErrorToTerra(routingProtocolOperation.Errors)
		}
		mappedRpOperations = append(mappedRpOperations, mappedRpOperation)
	}
	rpOperationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createRoutingProtocolOperationSch()}),
		mappedRpOperations,
	)
	return rpOperationSet
}

func routingProtocolChangeToTerra(routingProtocolChange *v4.RoutingProtocolChange) *schema.Set {
	if routingProtocolChange == nil {
		return nil
	}
	routingProtocolChanges := []*v4.RoutingProtocolChange{routingProtocolChange}
	mappedRpChanges := make([]interface{}, len(routingProtocolChanges))
	for i, rpChanges := range routingProtocolChanges {
		mappedRpChanges[i] = map[string]interface{}{
			"uuid": rpChanges.Uuid,
			"type": rpChanges.Type_,
			"href": rpChanges.Href,
		}
	}
	rpChangeSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createRoutingProtocolChangeSch()}),
		mappedRpChanges,
	)
	return rpChangeSet
}
