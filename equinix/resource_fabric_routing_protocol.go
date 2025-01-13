package equinix

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

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
		"outbound_as_prepend_count": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "AS path prepend count. One of: 0, 1, 3, 5",
		},
		"inbound_med": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Inbound Multi Exit Discriminator attribute",
		},
		"outbound_med": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Outbound Multi Exit Discriminator attribute",
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
		"outbound_as_prepend_count": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "AS path prepend count. One of: 0, 1, 3, 5",
		},
		"inbound_med": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Inbound Multi Exit Discriminator attribute",
		},
		"outbound_med": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Outbound Multi Exit Discriminator attribute",
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
			Computed:     true,
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
			Computed:    true,
			Description: "Routing Protocol name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
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
			Computed:    true,
			Description: "Routing Protocol Direct IPv4",
			Elem: &schema.Resource{
				Schema: createDirectConnectionIpv4Sch(),
			},
		},
		"direct_ipv6": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Description: "Routing Protocol Direct IPv6",
			Elem: &schema.Resource{
				Schema: createDirectConnectionIpv6Sch(),
			},
		},
		"bgp_ipv4": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Description: "Routing Protocol BGP IPv4",
			Elem: &schema.Resource{
				Schema: createBgpConnectionIpv4Sch(),
			},
		},
		"bgp_ipv6": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Description: "Routing Protocol BGP IPv6",
			Elem: &schema.Resource{
				Schema: createBgpConnectionIpv6Sch(),
			},
		},
		"customer_asn": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
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
			Computed:    true,
			Description: "BGP authorization key",
		},
		"as_override_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable AS number override",
		},
		"bfd": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
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
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	log.Printf("[WARN] Routing Protocol Connection uuid: %s", d.Get("connection_uuid").(string))
	fabricRoutingProtocolData, _, err := client.RoutingProtocolsApi.GetConnectionRoutingProtocolByUuid(ctx, d.Id(), d.Get("connection_uuid").(string)).Execute()
	if err != nil {
		log.Printf("[WARN] Routing Protocol %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	_ = setIdFromAPIResponse(fabricRoutingProtocolData, false, d)

	return setFabricRoutingProtocolMap(d, fabricRoutingProtocolData)
}

func resourceFabricRoutingProtocolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)

	start := time.Now()
	type_ := d.Get("type").(string)

	createRequest, err := routingProtocolPayloadFromType(type_, d)
	if err != nil {
		return diag.Errorf("error creating create request from Terraform configuration values: %s", err)
	}

	fabricRoutingProtocolData, _, err := client.RoutingProtocolsApi.CreateConnectionRoutingProtocol(ctx, d.Get("connection_uuid").(string)).RoutingProtocolBase(createRequest).Execute()

	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	_ = setIdFromAPIResponse(fabricRoutingProtocolData, false, d)

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if _, err = waitUntilRoutingProtocolIsProvisioned(d.Id(), d.Get("connection_uuid").(string), meta, d, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for RP (%s) to be created: %s", d.Id(), err)
	}

	return resourceFabricRoutingProtocolRead(ctx, d, meta)
}

func resourceFabricRoutingProtocolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)

	type_ := d.Get("type").(string)

	updateRequest, err := routingProtocolPayloadFromType(type_, d)
	if err != nil {
		return diag.Errorf("error creating update request from Terraform configuration values: %s", err)
	}

	start := time.Now()
	updatedRpResp, _, err := client.RoutingProtocolsApi.ReplaceConnectionRoutingProtocolByUuid(ctx, d.Id(), d.Get("connection_uuid").(string)).RoutingProtocolBase(updateRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	changeUuid := setIdFromAPIResponse(updatedRpResp, true, d)

	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	_, err = waitForRoutingProtocolUpdateCompletion(changeUuid, d.Id(), d.Get("connection_uuid").(string), meta, d, ctx, updateTimeout)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("timeout updating routing protocol: %v", err))
	}

	updateTimeout = d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	updatedProvisionedRpResp, err := waitUntilRoutingProtocolIsProvisioned(d.Id(), d.Get("connection_uuid").(string), meta, d, ctx, updateTimeout)
	if err != nil {
		return diag.Errorf("error waiting for RP (%s) to be replace updated: %s", d.Id(), err)
	}

	return setFabricRoutingProtocolMap(d, updatedProvisionedRpResp)
}

func resourceFabricRoutingProtocolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	start := time.Now()
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	_, _, err := client.RoutingProtocolsApi.DeleteConnectionRoutingProtocolByUuid(ctx, d.Id(), d.Get("connection_uuid").(string)).Execute()
	if err != nil {
		if genericError, ok := err.(*fabricv4.GenericOpenAPIError); ok {
			if fabricErrs, ok := genericError.Model().([]fabricv4.Error); ok {
				// EQ-3041121 = Routing Protocol already deleted
				if equinix_errors.HasErrorCode(fabricErrs, "EQ-3041121") {
					return diags
				}
			}
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	deleteTimeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)
	err = WaitUntilRoutingProtocolIsDeprovisioned(d.Id(), d.Get("connection_uuid").(string), meta, d, ctx, deleteTimeout)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}

	return diags
}

func setIdFromAPIResponse(resp *fabricv4.RoutingProtocolData, isChange bool, d *schema.ResourceData) string {
	var changeUuid string

	switch rpData := resp.GetActualInstance().(type) {
	case *fabricv4.RoutingProtocolBGPData:
		if isChange {
			change := rpData.GetChange()
			changeUuid = change.GetUuid()
		}
		d.SetId(rpData.GetUuid())
	case *fabricv4.RoutingProtocolDirectData:
		if isChange {
			change := rpData.GetChange()
			changeUuid = change.GetUuid()
		}
		d.SetId(rpData.GetUuid())
	}

	return changeUuid
}

func routingProtocolPayloadFromType(type_ string, d *schema.ResourceData) (fabricv4.RoutingProtocolBase, error) {
	payload := fabricv4.RoutingProtocolBase{}
	if type_ == "BGP" {
		bgpRP := fabricv4.RoutingProtocolBGPType{}
		bgpType := fabricv4.RoutingProtocolBGPTypeType(type_)
		bgpRP.SetType(bgpType)

		name := d.Get("name").(string)
		if name != "" {
			bgpRP.SetName(name)
		}

		customerASNSchema := d.Get("customer_asn")
		if customerASNSchema != nil {
			customerASN := int64(customerASNSchema.(int))
			bgpRP.SetCustomerAsn(customerASN)
		}

		equinixASNSchema := d.Get("equinix_asn")
		if equinixASNSchema != nil {
			equinixASN := int64(equinixASNSchema.(int))
			bgpRP.SetEquinixAsn(equinixASN)
		}

		bgpAuthKey := d.Get("bgp_auth_key").(string)
		if bgpAuthKey != "" {
			bgpRP.SetBgpAuthKey(bgpAuthKey)
		}

		schemaBgpIpv4 := d.Get("bgp_ipv4")
		if schemaBgpIpv4 != nil {
			bgpIpv4, err := routingProtocolBgpIpv4TerraformToGo(schemaBgpIpv4.(*schema.Set).List())
			if err != nil {
				return fabricv4.RoutingProtocolBase{}, err
			}
			if !reflect.DeepEqual(bgpIpv4, fabricv4.BGPConnectionIpv4{}) {
				bgpRP.SetBgpIpv4(bgpIpv4)
			}
		}

		schemaBgpIpv6 := d.Get("bgp_ipv6")
		if schemaBgpIpv6 != nil {
			bgpIpv6, err := routingProtocolBgpIpv6TerraformToGo(schemaBgpIpv6.(*schema.Set).List())
			if err != nil {
				return fabricv4.RoutingProtocolBase{}, err
			}
			if !reflect.DeepEqual(bgpIpv6, fabricv4.BGPConnectionIpv6{}) {
				bgpRP.SetBgpIpv6(bgpIpv6)
			}
		}

		asOverrideEnabled := d.Get("as_override_enabled").(bool)
		bgpRP.SetAsOverrideEnabled(asOverrideEnabled)

		bfdSchema := d.Get("bfd")
		if bfdSchema != nil {
			bfd := routingProtocolBfdTerraformToGo(bfdSchema.(*schema.Set).List())
			bgpRP.SetBfd(bfd)
		}
		payload = fabricv4.RoutingProtocolBGPTypeAsRoutingProtocolBase(&bgpRP)
	}
	if type_ == "DIRECT" {
		directRP := fabricv4.RoutingProtocolDirectType{}
		directType := fabricv4.RoutingProtocolDirectTypeType(type_)
		directRP.SetType(directType)

		name := d.Get("name").(string)
		if name != "" {
			directRP.SetName(name)
		}
		schemaDirectIpv4 := d.Get("direct_ipv4")
		if schemaDirectIpv4 != nil {
			directIpv4 := routingProtocolDirectIpv4TerraformToGo(schemaDirectIpv4.(*schema.Set).List())
			if !reflect.DeepEqual(directIpv4, fabricv4.DirectConnectionIpv4{}) {
				directRP.SetDirectIpv4(directIpv4)
			}
		}

		schemaDirectIpv6 := d.Get("direct_ipv6")
		if schemaDirectIpv6 != nil {
			directIpv6 := routingProtocolDirectIpv6TerraformToGo(schemaDirectIpv6.(*schema.Set).List())
			if !reflect.DeepEqual(directIpv6, fabricv4.DirectConnectionIpv6{}) {
				directRP.SetDirectIpv6(directIpv6)
			}
		}
		payload = fabricv4.RoutingProtocolDirectTypeAsRoutingProtocolBase(&directRP)
	}
	return payload, nil
}

func setFabricRoutingProtocolMap(d *schema.ResourceData, routingProtocolData *fabricv4.RoutingProtocolData) diag.Diagnostics {
	diags := diag.Diagnostics{}
	routingProtocol := FabricRoutingProtocolMap(routingProtocolData)
	err := equinix_schema.SetMap(d, routingProtocol)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func FabricRoutingProtocolMap(routingProtocolData *fabricv4.RoutingProtocolData) map[string]interface{} {
	routingProtocol := make(map[string]interface{})
	switch rp := routingProtocolData.GetActualInstance().(type) {
	case *fabricv4.RoutingProtocolBGPData:
		routingProtocol["name"] = rp.GetName()
		routingProtocol["href"] = rp.GetHref()
		routingProtocol["type"] = string(rp.GetType())
		routingProtocol["state"] = string(rp.GetState())
		routingProtocol["customer_asn"] = rp.GetCustomerAsn()
		routingProtocol["equinix_asn"] = rp.GetCustomerAsn()
		routingProtocol["bgp_auth_key"] = rp.GetBgpAuthKey()
		routingProtocol["as_override_enabled"] = rp.GetAsOverrideEnabled()
		if rp.Operation != nil {
			operation := rp.GetOperation()
			routingProtocol["operation"] = routingProtocolOperationGoToTerraform(&operation)
		}
		if rp.BgpIpv4 != nil {
			bgpIpv4 := rp.GetBgpIpv4()
			routingProtocol["bgp_ipv4"] = routingProtocolBgpConnectionIpv4GoToTerraform(&bgpIpv4)
		}
		if rp.BgpIpv6 != nil {
			bgpIpv6 := rp.GetBgpIpv6()
			routingProtocol["bgp_ipv6"] = routingProtocolBgpConnectionIpv6GoToTerraform(&bgpIpv6)
		}
		if rp.Bfd != nil {
			bfd := rp.GetBfd()
			routingProtocol["bfd"] = routingProtocolBfdGoToTerraform(&bfd)
		}
		if rp.Change != nil {
			change := rp.GetChange()
			routingProtocol["change"] = routingProtocolChangeGoToTerraform(&change)
		}
		if rp.Changelog != nil {
			changeLog := rp.GetChangelog()
			routingProtocol["change_log"] = equinix_fabric_schema.ChangeLogGoToTerraform(&changeLog)
		}
	case *fabricv4.RoutingProtocolDirectData:
		routingProtocol["name"] = rp.GetName()
		routingProtocol["href"] = rp.GetHref()
		routingProtocol["type"] = string(rp.GetType())
		routingProtocol["state"] = string(rp.GetState())
		if rp.Operation != nil {
			operation := rp.GetOperation()
			routingProtocol["operation"] = routingProtocolOperationGoToTerraform(&operation)
		}
		if rp.DirectIpv4 != nil {
			directIpv4 := rp.GetDirectIpv4()
			routingProtocol["direct_ipv4"] = routingProtocolDirectConnectionIpv4GoToTerraform(&directIpv4)
		}
		if rp.DirectIpv6 != nil {
			directIpv6 := rp.GetDirectIpv6()
			routingProtocol["direct_ipv6"] = routingProtocolDirectConnectionIpv6GoToTerraform(&directIpv6)
		}
		if rp.Change != nil {
			change := rp.GetChange()
			routingProtocol["change"] = routingProtocolChangeGoToTerraform(&change)
		}
		if rp.Changelog != nil {
			changeLog := rp.GetChangelog()
			routingProtocol["change_log"] = equinix_fabric_schema.ChangeLogGoToTerraform(&changeLog)
		}
	}

	return routingProtocol
}
func waitUntilRoutingProtocolIsProvisioned(uuid string, connUuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.RoutingProtocolData, error) {
	log.Printf("Waiting for routing protocol to be provisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CONNECTIONSTATE_PROVISIONING),
			string(fabricv4.CONNECTIONSTATE_REPROVISIONING),
		},
		Target: []string{
			string(fabricv4.CONNECTIONSTATE_PROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
			dbConn, _, err := client.RoutingProtocolsApi.GetConnectionRoutingProtocolByUuid(ctx, uuid, connUuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			var state string
			switch rpData := dbConn.GetActualInstance().(type) {
			case *fabricv4.RoutingProtocolBGPData:
				state = string(rpData.GetState())
			case *fabricv4.RoutingProtocolDirectData:
				state = string(rpData.GetState())
			}
			return dbConn, state, nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbConn *fabricv4.RoutingProtocolData

	if err == nil {
		dbConn = inter.(*fabricv4.RoutingProtocolData)
	}

	return dbConn, err
}

func WaitUntilRoutingProtocolIsDeprovisioned(uuid string, connUuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for routing protocol to be deprovisioned, uuid %s", uuid)

	/* check if resource is not found */
	stateConf := &retry.StateChangeConf{
		Target: []string{
			strconv.Itoa(404),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
			dbConn, resp, _ := client.RoutingProtocolsApi.GetConnectionRoutingProtocolByUuid(ctx, uuid, connUuid).Execute()
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

func waitForRoutingProtocolUpdateCompletion(rpChangeUuid string, uuid string, connUuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.RoutingProtocolChangeData, error) {
	log.Printf("Waiting for routing protocol update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
			dbConn, _, err := client.RoutingProtocolsApi.GetConnectionRoutingProtocolsChangeByUuid(ctx, connUuid, uuid, rpChangeUuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			updatableState := ""
			if dbConn.GetStatus() == "COMPLETED" {
				updatableState = dbConn.GetStatus()
			}
			return dbConn, updatableState, nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbConn *fabricv4.RoutingProtocolChangeData

	if err == nil {
		dbConn = inter.(*fabricv4.RoutingProtocolChangeData)
	}
	return dbConn, err
}

func routingProtocolDirectIpv4TerraformToGo(routingProtocolDirectIpv4Request []interface{}) fabricv4.DirectConnectionIpv4 {
	if len(routingProtocolDirectIpv4Request) == 0 {
		return fabricv4.DirectConnectionIpv4{}
	}

	rpDirectIpv4 := fabricv4.DirectConnectionIpv4{}

	directIpv4Map := routingProtocolDirectIpv4Request[0].(map[string]interface{})
	equinixIfaceIp := directIpv4Map["equinix_iface_ip"].(string)
	if equinixIfaceIp != "" {
		rpDirectIpv4.SetEquinixIfaceIp(equinixIfaceIp)
	}

	return rpDirectIpv4
}

func routingProtocolDirectIpv6TerraformToGo(routingProtocolDirectIpv6Request []interface{}) fabricv4.DirectConnectionIpv6 {
	if len(routingProtocolDirectIpv6Request) == 0 {
		return fabricv4.DirectConnectionIpv6{}
	}
	rpDirectIpv6 := fabricv4.DirectConnectionIpv6{}
	directIpv6Map := routingProtocolDirectIpv6Request[0].(map[string]interface{})
	equinixIfaceIp := directIpv6Map["equinix_iface_ip"].(string)
	if equinixIfaceIp != "" {
		rpDirectIpv6.SetEquinixIfaceIp(equinixIfaceIp)
	}

	return rpDirectIpv6
}

func routingProtocolBgpIpv4TerraformToGo(routingProtocolBgpIpv4Request []interface{}) (fabricv4.BGPConnectionIpv4, error) {
	if len(routingProtocolBgpIpv4Request) == 0 {
		return fabricv4.BGPConnectionIpv4{}, nil
	}

	rpBgpIpv4 := fabricv4.BGPConnectionIpv4{}
	bgpIpv4Map := routingProtocolBgpIpv4Request[0].(map[string]interface{})
	customerPeerIp := bgpIpv4Map["customer_peer_ip"].(string)
	if customerPeerIp != "" {
		rpBgpIpv4.SetCustomerPeerIp(customerPeerIp)
	}
	enabled := bgpIpv4Map["enabled"].(bool)
	rpBgpIpv4.SetEnabled(enabled)

	if outboundAsPrependCountStr := bgpIpv4Map["outbound_as_prepend_count"].(string); outboundAsPrependCountStr != "" {
		outboundAsPrependCount, err := strconv.ParseInt(outboundAsPrependCountStr, 10, 64)
		if err != nil {
			return fabricv4.BGPConnectionIpv4{}, fmt.Errorf("error converting outbound_as_prepend_count from string to int64: %s", err)
		}
		rpBgpIpv4.SetOutboundASPrependCount(outboundAsPrependCount)
	}
	if inboundMed := bgpIpv4Map["inbound_med"].(int); inboundMed > 0 {
		rpBgpIpv4.SetInboundMED(int64(inboundMed))
	}
	if outboundMed := bgpIpv4Map["outbound_med"].(int); outboundMed > 0 {
		rpBgpIpv4.SetOutboundMED(int64(outboundMed))
	}

	return rpBgpIpv4, nil
}

func routingProtocolBgpIpv6TerraformToGo(routingProtocolBgpIpv6Request []interface{}) (fabricv4.BGPConnectionIpv6, error) {
	if len(routingProtocolBgpIpv6Request) == 0 {
		return fabricv4.BGPConnectionIpv6{}, nil
	}

	rpBgpIpv6 := fabricv4.BGPConnectionIpv6{}
	bgpIpv6Map := routingProtocolBgpIpv6Request[0].(map[string]interface{})
	customerPeerIp := bgpIpv6Map["customer_peer_ip"].(string)
	if customerPeerIp != "" {
		rpBgpIpv6.SetCustomerPeerIp(customerPeerIp)
	}
	enabled := bgpIpv6Map["enabled"].(bool)
	rpBgpIpv6.SetEnabled(enabled)

	if outboundAsPrependCountStr := bgpIpv6Map["outbound_as_prepend_count"].(string); outboundAsPrependCountStr != "" {
		outboundAsPrependCount, err := strconv.ParseInt(outboundAsPrependCountStr, 10, 64)
		if err != nil {
			return fabricv4.BGPConnectionIpv6{}, err
		}
		rpBgpIpv6.SetOutboundASPrependCount(outboundAsPrependCount)
	}
	if inboundMed := bgpIpv6Map["inbound_med"].(int); inboundMed > 0 {
		rpBgpIpv6.SetInboundMED(int64(inboundMed))
	}
	if outboundMed := bgpIpv6Map["outbound_med"].(int); outboundMed > 0 {
		rpBgpIpv6.SetOutboundMED(int64(outboundMed))
	}

	return rpBgpIpv6, nil
}

func routingProtocolBfdTerraformToGo(routingProtocolBfdRequest []interface{}) fabricv4.RoutingProtocolBFD {
	if len(routingProtocolBfdRequest) == 0 {
		return fabricv4.RoutingProtocolBFD{}
	}

	rpBfd := fabricv4.RoutingProtocolBFD{}
	rpBfdMap := routingProtocolBfdRequest[0].(map[string]interface{})
	bfdEnabled := rpBfdMap["enabled"].(bool)
	rpBfd.SetEnabled(bfdEnabled)
	bfdInterval := rpBfdMap["interval"].(string)
	if bfdInterval != "" {
		rpBfd.SetInterval(bfdInterval)
	}

	return rpBfd
}

func routingProtocolDirectConnectionIpv4GoToTerraform(routingProtocolDirectIpv4 *fabricv4.DirectConnectionIpv4) *schema.Set {
	if routingProtocolDirectIpv4 == nil {
		return nil
	}

	mappedDirectIpv4 := map[string]interface{}{
		"equinix_iface_ip": routingProtocolDirectIpv4.GetEquinixIfaceIp(),
	}

	rpDirectIpv4Set := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createDirectConnectionIpv4Sch()}),
		[]interface{}{mappedDirectIpv4},
	)
	return rpDirectIpv4Set
}

func routingProtocolDirectConnectionIpv6GoToTerraform(routingProtocolDirectIpv6 *fabricv4.DirectConnectionIpv6) *schema.Set {
	if routingProtocolDirectIpv6 == nil {
		return nil
	}

	mappedDirectIpv6 := map[string]interface{}{
		"equinix_iface_ip": routingProtocolDirectIpv6.GetEquinixIfaceIp(),
	}

	rpDirectIpv6Set := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createDirectConnectionIpv6Sch()}),
		[]interface{}{mappedDirectIpv6},
	)
	return rpDirectIpv6Set
}

func routingProtocolBgpConnectionIpv4GoToTerraform(routingProtocolBgpIpv4 *fabricv4.BGPConnectionIpv4) *schema.Set {
	if routingProtocolBgpIpv4 == nil {
		return nil
	}

	mappedBgpIpv4 := map[string]interface{}{
		"customer_peer_ip":          routingProtocolBgpIpv4.GetCustomerPeerIp(),
		"equinix_peer_ip":           routingProtocolBgpIpv4.GetEquinixPeerIp(),
		"enabled":                   routingProtocolBgpIpv4.GetEnabled(),
		"outbound_as_prepend_count": strconv.FormatInt(routingProtocolBgpIpv4.GetOutboundASPrependCount(), 10),
		"inbound_med":               int(routingProtocolBgpIpv4.GetInboundMED()),
		"outbound_med":              int(routingProtocolBgpIpv4.GetOutboundMED()),
	}
	rpBgpIpv4Set := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createBgpConnectionIpv4Sch()}),
		[]interface{}{mappedBgpIpv4},
	)
	return rpBgpIpv4Set
}

func routingProtocolBgpConnectionIpv6GoToTerraform(routingProtocolBgpIpv6 *fabricv4.BGPConnectionIpv6) *schema.Set {
	if routingProtocolBgpIpv6 == nil {
		return nil
	}

	mappedBgpIpv6 := map[string]interface{}{
		"customer_peer_ip":          routingProtocolBgpIpv6.GetCustomerPeerIp(),
		"equinix_peer_ip":           routingProtocolBgpIpv6.GetEquinixPeerIp(),
		"enabled":                   routingProtocolBgpIpv6.GetEnabled(),
		"outbound_as_prepend_count": strconv.FormatInt(routingProtocolBgpIpv6.GetOutboundASPrependCount(), 10),
		"inbound_med":               int(routingProtocolBgpIpv6.GetInboundMED()),
		"outbound_med":              int(routingProtocolBgpIpv6.GetOutboundMED()),
	}

	rpBgpIpv6Set := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createBgpConnectionIpv6Sch()}),
		[]interface{}{mappedBgpIpv6},
	)
	return rpBgpIpv6Set
}

func routingProtocolBfdGoToTerraform(routingProtocolBfd *fabricv4.RoutingProtocolBFD) *schema.Set {
	if routingProtocolBfd == nil {
		return nil
	}

	mappedRpBfd := map[string]interface{}{
		"enabled":  routingProtocolBfd.GetEnabled(),
		"interval": routingProtocolBfd.GetInterval(),
	}

	rpBfdSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createRoutingProtocolBfdSch()}),
		[]interface{}{mappedRpBfd},
	)
	return rpBfdSet
}

func routingProtocolOperationGoToTerraform(routingProtocolOperation *fabricv4.RoutingProtocolOperation) *schema.Set {
	if routingProtocolOperation == nil {
		return nil
	}
	mappedRpOperation := make(map[string]interface{})
	errors := routingProtocolOperation.GetErrors()
	if errors != nil {
		mappedRpOperation["errors"] = equinix_fabric_schema.ErrorGoToTerraform(errors)
	}

	rpOperationSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createRoutingProtocolOperationSch()}),
		[]interface{}{mappedRpOperation},
	)
	return rpOperationSet
}

func routingProtocolChangeGoToTerraform(routingProtocolChange *fabricv4.RoutingProtocolChange) *schema.Set {
	if routingProtocolChange == nil {
		return nil
	}

	mappedRpChange := map[string]interface{}{
		"uuid": routingProtocolChange.GetUuid(),
		"type": string(routingProtocolChange.GetType()),
		"href": routingProtocolChange.GetHref(),
	}

	rpChangeSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createRoutingProtocolChangeSch()}),
		[]interface{}{mappedRpChange},
	)
	return rpChangeSet
}
