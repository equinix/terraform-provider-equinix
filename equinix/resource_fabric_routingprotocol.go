package equinix

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"strings"
	"time"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var connId = "temp connection uuid" // fixme: get connectionId

func resourceFabricRoutingProtocol() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Read:   schema.DefaultTimeout(6 * time.Minute),
		},
		ReadContext:   resourceFabricRoutingProtocolRead,
		CreateContext: resourceFabricRoutingProtocolCreate,
		//UpdateContext: resourceFabricRoutingProtocolUpdate,
		DeleteContext: resourceFabricRoutingProtocolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: createFabricRoutingProtocolResourceSchema(),

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric connection\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
	}
}

func resourceFabricRoutingProtocolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	fabricRoutingProtocol, _, err := client.RoutingProtocolsApi.GetConnectionRoutingProtocolByUuid(ctx, d.Id(), connId)
	if err != nil {
		log.Printf("[WARN] Routing Protocol %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(err)
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
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	schemaBgpIpv4 := d.Get("bgp_ipv4").(*schema.Set).List()
	bgpIpv4 := routingProtocolBgpIpv4ToFabric(schemaBgpIpv4)
	schemaBgpIpv6 := d.Get("bgp_ipv6").(*schema.Set).List()
	bgpIpv6 := routingProtocolBgpIpv6ToFabric(schemaBgpIpv6)
	schemaDirectIpv4 := d.Get("direct_ipv4").(*schema.Set).List()
	directIpv4 := routingProtocolDirectIpv4ToFabric(schemaDirectIpv4)
	schemaDirectIpv6 := d.Get("direct_ipv6").(*schema.Set).List()
	DirectIpv6 := routingProtocolDirectIpv6ToFabric(schemaDirectIpv6)
	schemaBfd := d.Get("bfd").(*schema.Set).List()
	bfd := routingProtocolBfdToFabric(schemaBfd)

	var createRequest = v4.RoutingProtocolBase{
		Type_: d.Get("type").(string),
		OneOfRoutingProtocolBase: v4.OneOfRoutingProtocolBase{
			RoutingProtocolBgpType: v4.RoutingProtocolBgpType{
				Type_:       d.Get("type").(string),
				Name:        d.Get("name").(string),
				BgpIpv4:     &bgpIpv4,
				BgpIpv6:     &bgpIpv6,
				CustomerAsn: d.Get("customer_asn").(int64),
				EquinixAsn:  d.Get("equinix_asn").(int64),
				BgpAuthKey:  d.Get("bgp_auth_key").(string),
				Bfd:         &bfd,
			},
			RoutingProtocolDirectType: v4.RoutingProtocolDirectType{
				Type_:      d.Get("type").(string),
				Name:       d.Get("name").(string),
				DirectIpv4: &directIpv4,
				DirectIpv6: &DirectIpv6,
			},
		},
	}
	fabricRoutingProtocol, _, err := client.RoutingProtocolsApi.CreateConnectionRoutingProtocol(ctx, createRequest, connId)
	if err != nil {
		return diag.FromErr(err)
	}

	switch fabricRoutingProtocol.Type_ {
	case "BGP":
		d.SetId(fabricRoutingProtocol.RoutingProtocolBgpData.Uuid)
	case "DIRECT":
		d.SetId(fabricRoutingProtocol.RoutingProtocolDirectData.Uuid)
	}

	if err = waitUntilRoutingProtocolIsProvisioned(d.Id(), connId, meta, ctx); err != nil {
		return diag.Errorf("error waiting for RP (%s) to be created: %s", d.Id(), err)
	}

	return resourceFabricRoutingProtocolRead(ctx, d, meta)
}

func resourceFabricRoutingProtocolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	_, resp, err := client.RoutingProtocolsApi.DeleteConnectionRoutingProtocolByUuid(ctx, d.Id(), connId)
	if err != nil {
		errors, ok := err.(v4.GenericSwaggerError).Model().([]v4.ModelError)
		if ok {
			// EQ-3142509 = Connection already deleted
			if hasModelErrorCode(errors, "EQ-3142509") {
				return diags
			}
		}
		return diag.FromErr(fmt.Errorf("error response for the routing protocol delete. Error %v and response %v", err, resp))
	}

	err = waitUntilRoutingProtocolDeprovisioned(d.Id(), connId, meta, ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func setFabricRoutingProtocolMap(d *schema.ResourceData, rp v4.RoutingProtocolData) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := error(nil) // fixme: doesnt look right
	if rp.Type_ == "BGP" {
		err = setMap(d, map[string]interface{}{
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
			"change_log":   changeLogToTerra(rp.RoutingProtocolBgpData.Changelog),
		})
	} else if rp.Type_ == "DIRECT" {
		err = setMap(d, map[string]interface{}{
			"name":        rp.RoutingProtocolDirectData.Name,
			"href":        rp.RoutingProtocolDirectData.Href,
			"type":        rp.RoutingProtocolDirectData.Type_,
			"state":       rp.RoutingProtocolDirectData.State,
			"operation":   routingProtocolOperationToTerra(rp.RoutingProtocolDirectData.Operation),
			"direct_ipv4": routingProtocolDirectConnectionIpv4ToTerra(rp.DirectIpv4),
			"direct_ipv6": routingProtocolDirectConnectionIpv6ToTerra(rp.DirectIpv6),
			"change":      routingProtocolChangeToTerra(rp.RoutingProtocolDirectData.Change),
			"change_log":  changeLogToTerra(rp.RoutingProtocolDirectData.Changelog),
		})
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func waitUntilRoutingProtocolIsProvisioned(uuid string, connUuid string, meta interface{}, ctx context.Context) error {
	log.Printf("Waiting for routing protocol to be provisioned, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(v4.PROVISIONING_ConnectionState),
			string(v4.REPROVISIONING_ConnectionState),
		},
		Target: []string{
			string(v4.PROVISIONED_ConnectionState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.RoutingProtocolsApi.GetConnectionRoutingProtocolByUuid(ctx, uuid, connUuid)
			if err != nil {
				return "", "", err
			}
			var state string
			if dbConn.Type_ == "BGP" {
				state = dbConn.RoutingProtocolBgpData.State
			} else if dbConn.Type_ == "DIRECT" {
				state = dbConn.RoutingProtocolDirectData.State
			}
			return dbConn, state, nil
		},
		Timeout:    5 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func waitUntilRoutingProtocolDeprovisioned(uuid string, connUuid string, meta interface{}, ctx context.Context) error {
	log.Printf("Waiting for routing protocol to be deprovisioned, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(v4.DEPROVISIONING_ConnectionState),
		},
		Target: []string{
			string(v4.DEPROVISIONED_ConnectionState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.RoutingProtocolsApi.GetConnectionRoutingProtocolByUuid(ctx, uuid, connUuid)
			if err != nil {
				return "", "", err
			}
			var state string
			if dbConn.Type_ == "BGP" {
				state = dbConn.RoutingProtocolBgpData.State
			} else if dbConn.Type_ == "DIRECT" {
				state = dbConn.RoutingProtocolDirectData.State
			}
			return dbConn, state, nil

		},
		Timeout:    5 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
