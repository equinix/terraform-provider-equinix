package equinix

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"strconv"
	"strings"
	"time"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFabricRoutingProtocols() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Read:   schema.DefaultTimeout(6 * time.Minute),
		},
		ReadContext:   resourceFabricRoutingProtocolsRead,
		CreateContext: resourceFabricRoutingProtocolsCreate,
		UpdateContext: resourceFabricRoutingProtocolsUpdate,
		DeleteContext: resourceFabricRoutingProtocolsDelete,
		Importer: &schema.ResourceImporter{
			// Custom state context function, to parse import argument as  connection_uuid/direct-rp-uuid/bgp-rp-uuid
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.SplitN(d.Id(), "/", 3)
				if len(parts) < 2 || len(parts) > 3 || parts[0] == "" || parts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%s), expected <conn-uuid>/<direct-rp-uuid> or <conn-uuid>/<direct-rp-uuid>/<bgp-rp-uuid>", d.Id())
				}
				connectionUuid, directRPUUID := parts[0], parts[1]
				idToSet := directRPUUID
				_ = d.Set("connection_uuid", connectionUuid)
				if len(parts) == 3 && parts[2] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%s), expected <conn-uuid>/<direct-rp-uuid> or <conn-uuid>/<direct-rp-uuid>/<bgp-rp-uuid>", d.Id())
				} else if len(parts) == 3 {
					bgpRPUUID := parts[2]
					idToSet += "/" + bgpRPUUID
				}

				d.SetId(idToSet)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: FabricRoutingProtocolsResourceSchema(),

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric connection\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
	}
}

type RoutingProtocols struct {
	connectionUUID        string
	directRoutingProtocol v4.RoutingProtocolDirectData
	bgpRoutingProtocol    v4.RoutingProtocolBgpData
}

func setFabricRoutingProtocolsMap(d *schema.ResourceData, rp RoutingProtocols) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := error(nil)

	err = d.Set("connection_uuid", rp.connectionUUID)

	if err != nil {
		return diag.FromErr(err)
	}

	if rp.directRoutingProtocol != (v4.RoutingProtocolDirectData{}) {
		err = setMap(d, map[string]interface{}{
			"direct_routing_protocol": directRoutingProtocolToTerra(rp.directRoutingProtocol),
		})
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if rp.bgpRoutingProtocol != (v4.RoutingProtocolBgpData{}) {
		err = setMap(d, map[string]interface{}{
			"bgp_routing_protocol": bgpRoutingProtocolToTerra(rp.bgpRoutingProtocol),
		})
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readUUIDsFromId(id string) (directRPUUID, bgpRPUUID string) {
	ids := strings.SplitN(id, "/", 2)
	if len(ids) == 2 && ids[1] != "" {
		return ids[0], ids[1]
	}
	return ids[0], ""
}

func readIdsFromBulkCreateResponse(bulkCreate v4.GetResponse) string {
	directRPUUID, bgpRPUUID := "", ""
	for _, rp := range bulkCreate.Data {
		if rp.Type_ == "DIRECT" {
			directRPUUID = rp.RoutingProtocolDirectData.Uuid
		}
		if rp.Type_ == "BGP" {
			bgpRPUUID = rp.RoutingProtocolBgpData.Uuid
		}
	}

	return directRPUUID + "/" + bgpRPUUID
}

func resourceFabricRoutingProtocolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	log.Printf("[INFO] Routing Protocols Resource Id: %s", d.Id())
	directRPUUID, bgpRPUUID := readUUIDsFromId(d.Id())
	connectionUUID := d.Get("connection_uuid").(string)
	log.Printf("[INFO] Routing Protocol Direct uuid: %s", directRPUUID)
	log.Printf("[INFO] Routing Protocol Direct uuid: %s", bgpRPUUID)
	log.Printf("[INFO] Routing Protocol Connection uuid: %s", connectionUUID)

	directRoutingProtocol, _, err := client.RoutingProtocolsApi.GetConnectionRoutingProtocolByUuid(ctx, directRPUUID, connectionUUID)
	if err != nil {
		log.Printf("[WARN] Routing Protocol %s not found , error %s", directRPUUID, err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(err)
	}

	bgpRoutingProtocol := v4.RoutingProtocolBgpData{}

	if bgpRPUUID != "" {
		bgpRoutingProtocolResponse, _, err := client.RoutingProtocolsApi.GetConnectionRoutingProtocolByUuid(ctx, bgpRPUUID, connectionUUID)
		bgpRoutingProtocol = bgpRoutingProtocolResponse.RoutingProtocolBgpData
		if err != nil {
			log.Printf("[ERROR] Routing Protocol %s not found. Error: %s", bgpRPUUID, err)
			if !strings.Contains(err.Error(), "500") {
				d.SetId("")
			}
			return diag.FromErr(err)
		}
	}

	routingProtocols := RoutingProtocols{
		connectionUUID:        connectionUUID,
		directRoutingProtocol: directRoutingProtocol.RoutingProtocolDirectData,
		bgpRoutingProtocol:    bgpRoutingProtocol,
	}

	return setFabricRoutingProtocolsMap(d, routingProtocols)
}

func resourceFabricRoutingProtocolsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)

	connectionUUID := d.Get("connection_uuid").(string)
	directRoutingProtocolConfig := d.Get("direct_routing_protocol").(*schema.Set).List()
	bgpRoutingProtocolConfig := d.Get("bgp_routing_protocol").(*schema.Set).List()
	directRoutingProtocol := directRoutingProtocolToFabric(directRoutingProtocolConfig)
	bgpRoutingProtocol := bgpRoutingProtocolToFabric(bgpRoutingProtocolConfig)

	directRoutingProtocolRequest := v4.RoutingProtocolBase{
		Type_: "DIRECT",
		OneOfRoutingProtocolBase: v4.OneOfRoutingProtocolBase{
			RoutingProtocolDirectType: directRoutingProtocol,
		},
	}

	if bgpRoutingProtocol != (v4.RoutingProtocolBgpType{}) {
		bgpRoutingProtocolRequest := v4.RoutingProtocolBase{
			Type_: "BGP",
			OneOfRoutingProtocolBase: v4.OneOfRoutingProtocolBase{
				RoutingProtocolBgpType: bgpRoutingProtocol,
			},
		}

		bulkRequest := v4.ConnectionRoutingProtocolPostRequest{
			Data: []v4.RoutingProtocolBase{
				directRoutingProtocolRequest,
				bgpRoutingProtocolRequest,
			},
		}

		bulkRPCreateResponse, _, err := client.RoutingProtocolsApi.CreateConnectionRoutingProtocolsInBulk(ctx, bulkRequest, connectionUUID)
		if err != nil {
			return diag.FromErr(err)
		}
		id := readIdsFromBulkCreateResponse(bulkRPCreateResponse)
		log.Printf("[DEBUG] ids from bulk create: %s", id)
		d.SetId(id)
		directRPUUID, bgpRPUUID := readUUIDsFromId(id)
		if _, err = waitUntilRoutingProtocolIsProvisioned(directRPUUID, connectionUUID, meta, ctx); err != nil {
			return diag.Errorf("error waiting for Routing Protocol %s to be created: %s", d.Id(), err)
		}
		if _, err = waitUntilRoutingProtocolIsProvisioned(bgpRPUUID, connectionUUID, meta, ctx); err != nil {
			return diag.Errorf("error waiting for Routing Protocol %s to be created: %s", d.Id(), err)
		}
	} else {
		fabricRoutingProtocol, _, err := client.RoutingProtocolsApi.CreateConnectionRoutingProtocol(ctx, directRoutingProtocolRequest, connectionUUID)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fabricRoutingProtocol.RoutingProtocolDirectData.Uuid)

		if _, err = waitUntilRoutingProtocolIsProvisioned(fabricRoutingProtocol.RoutingProtocolDirectData.Uuid, connectionUUID, meta, ctx); err != nil {
			return diag.Errorf("error waiting for Routing Protocol %s to be created: %s", d.Id(), err)
		}
	}

	return resourceFabricRoutingProtocolsRead(ctx, d, meta)
}

func directConfigHasChanged(old, new v4.RoutingProtocolDirectType) bool {
	if (old.DirectIpv4 != nil && new.DirectIpv4 != nil && old.DirectIpv4.EquinixIfaceIp != new.DirectIpv4.EquinixIfaceIp) ||
		(old.DirectIpv4 != nil && new.DirectIpv4 != nil && old.DirectIpv6.EquinixIfaceIp != new.DirectIpv6.EquinixIfaceIp) {
		return true
	}

	return false
}

func bgpConfigHasChanged(old, new v4.RoutingProtocolBgpType) bool {
	if (old.BgpIpv4 != nil && new.BgpIpv4 != nil && old.BgpIpv4.CustomerPeerIp != new.BgpIpv4.CustomerPeerIp) ||
		(old.BgpIpv4 != nil && new.BgpIpv4 != nil && old.BgpIpv4.EquinixPeerIp != new.BgpIpv4.EquinixPeerIp) ||
		(old.BgpIpv4 != nil && new.BgpIpv4 != nil && old.BgpIpv4.Enabled != new.BgpIpv4.Enabled) ||
		(old.BgpIpv6 != nil && new.BgpIpv6 != nil && old.BgpIpv6.CustomerPeerIp != new.BgpIpv6.CustomerPeerIp) ||
		(old.BgpIpv6 != nil && new.BgpIpv6 != nil && old.BgpIpv6.EquinixPeerIp != new.BgpIpv6.EquinixPeerIp) ||
		(old.BgpIpv6 != nil && new.BgpIpv6 != nil && old.BgpIpv6.Enabled != new.BgpIpv6.Enabled) ||
		old.CustomerAsn != new.CustomerAsn ||
		old.EquinixAsn != new.CustomerAsn ||
		old.BgpAuthKey != new.BgpAuthKey ||
		(old.Bfd != nil && new.Bfd != nil && old.Bfd.Enabled != new.Bfd.Enabled) ||
		(old.Bfd != nil && new.Bfd != nil && old.Bfd.Interval != new.Bfd.Interval) {
		return true
	}

	return false
}

func resourceFabricRoutingProtocolsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)

	directRPUUID, bgpRPUUID := readUUIDsFromId(d.Id())

	connectionUUID := d.Get("connection_uuid").(string)
	oldDirectRoutingProtocolConfig, newDirectRoutingProtocolConfig := d.GetChange("direct_routing_protocol")
	oldBgpRoutingProtocolConfig, newBgpRoutingProtocolConfig := d.GetChange("bgp_routing_protocol")

	oldDirectRoutingProtocol, directRoutingProtocol := directRoutingProtocolToFabric(oldDirectRoutingProtocolConfig.(*schema.Set).List()), directRoutingProtocolToFabric(newDirectRoutingProtocolConfig.(*schema.Set).List())
	oldBgpRoutingProtocol, bgpRoutingProtocol := bgpRoutingProtocolToFabric(oldBgpRoutingProtocolConfig.(*schema.Set).List()), bgpRoutingProtocolToFabric(newBgpRoutingProtocolConfig.(*schema.Set).List())

	directConfigChanged := directConfigHasChanged(oldDirectRoutingProtocol, directRoutingProtocol)
	bgpConfigChanged := bgpConfigHasChanged(oldBgpRoutingProtocol, bgpRoutingProtocol)

	if !directConfigChanged && !bgpConfigChanged {
		return diag.Diagnostics{}
	}

	updatedDirectRP := v4.RoutingProtocolDirectData{}

	if directConfigChanged {
		directRoutingProtocolRequest := v4.RoutingProtocolBase{
			Type_: "DIRECT",
			OneOfRoutingProtocolBase: v4.OneOfRoutingProtocolBase{
				RoutingProtocolDirectType: directRoutingProtocol,
			},
		}

		updateDirectRP, _, err := client.RoutingProtocolsApi.ReplaceConnectionRoutingProtocolByUuid(ctx, directRoutingProtocolRequest, directRPUUID, connectionUUID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error response for the routing protocol (%s) update: %v", directRPUUID, err))
		}

		_, err = waitForRoutingProtocolsUpdateCompletion(updateDirectRP.RoutingProtocolDirectData.Change.Uuid, directRPUUID, connectionUUID, meta, ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("errored while waiting for successful routing protocol (%s) update: %v", directRPUUID, err))
		}

		d.SetId(updateDirectRP.RoutingProtocolDirectData.Uuid)
		directRPUpdateProvisioned, err := waitUntilRoutingProtocolsIsProvisioned(directRPUUID, connectionUUID, meta, ctx)
		updatedDirectRP = directRPUpdateProvisioned.RoutingProtocolDirectData
		if err != nil {
			return diag.Errorf("error waiting for Routing Protocol (%s) to be update to be in provisioned state: %s", directRPUUID, err)
		}
	}

	updatedBgpRP := v4.RoutingProtocolBgpData{}

	if bgpConfigChanged {
		bgpRoutingProtocolRequest := v4.RoutingProtocolBase{
			Type_: "BGP",
			OneOfRoutingProtocolBase: v4.OneOfRoutingProtocolBase{
				RoutingProtocolBgpType: bgpRoutingProtocol,
			},
		}

		if bgpRPUUID != "" && bgpRoutingProtocol == (v4.RoutingProtocolBgpType{}) {
			// BGP removed after creation, needs to be deleted
			_, _, err := client.RoutingProtocolsApi.DeleteConnectionRoutingProtocolByUuid(ctx, bgpRPUUID, connectionUUID)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error response for the routing protocol (%s) delete in update: %v", bgpRPUUID, err))
			}

			d.SetId(directRPUUID)

			err = waitUntilRoutingProtocolsIsDeprovisioned(bgpRPUUID, connectionUUID, meta, ctx)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error waiting for routing protocol (%s) deletion in update: %v", bgpRPUUID, err))
			}
		} else {
			if bgpRPUUID != "" && bgpRoutingProtocol != (v4.RoutingProtocolBgpType{}) {
				// BGP still here, needs to be updated
				bgpRPUpdate, _, err := client.RoutingProtocolsApi.ReplaceConnectionRoutingProtocolByUuid(ctx, bgpRoutingProtocolRequest, bgpRPUUID, connectionUUID)
				if err != nil {
					return diag.FromErr(fmt.Errorf("error response for the routing protocol (%s) replace update: %v", directRPUUID, err))
				}

				_, err = waitForRoutingProtocolsUpdateCompletion(bgpRPUpdate.RoutingProtocolBgpData.Change.Uuid, bgpRPUUID, connectionUUID, meta, ctx)
				if err != nil {
					return diag.FromErr(fmt.Errorf("errored while waiting for successful routing protocol (%s) replace update: %v", bgpRPUUID, err))
				}

				d.SetId(directRPUUID + "/" + bgpRPUpdate.RoutingProtocolBgpData.Uuid)

				bgpRPUpdateProvisioned, err := waitUntilRoutingProtocolsIsProvisioned(bgpRPUUID, connectionUUID, meta, ctx)
				if err != nil {
					return diag.Errorf("error waiting for Routing Protocol (%s) to be replace updated: %s", bgpRPUUID, err)
				}
				updatedBgpRP = bgpRPUpdateProvisioned.RoutingProtocolBgpData
			} else if bgpRPUUID == "" && bgpRoutingProtocol != (v4.RoutingProtocolBgpType{}) {
				// BGP added after creation, needs to be created
				bgpRPCreate, _, err := client.RoutingProtocolsApi.CreateConnectionRoutingProtocol(ctx, bgpRoutingProtocolRequest, connectionUUID)
				if err != nil {
					return diag.FromErr(err)
				}
				d.SetId(directRPUUID + "/" + bgpRPCreate.RoutingProtocolBgpData.Uuid)

				updatedBgpRP = bgpRPCreate.RoutingProtocolBgpData
				if _, err = waitUntilRoutingProtocolIsProvisioned(bgpRPCreate.RoutingProtocolBgpData.Uuid, connectionUUID, meta, ctx); err != nil {
					return diag.Errorf("error waiting for Routing Protocol %s to be created: %s", bgpRPCreate.RoutingProtocolBgpData.Uuid, err)
				}
			}
		}
	}

	routingProtocols := RoutingProtocols{
		connectionUUID:        connectionUUID,
		directRoutingProtocol: updatedDirectRP,
		bgpRoutingProtocol:    updatedBgpRP,
	}

	return setFabricRoutingProtocolsMap(d, routingProtocols)
}

func resourceFabricRoutingProtocolsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	connectionUUID := d.Get("connection_uuid").(string)
	directRPUUID, bgpRPUUID := readUUIDsFromId(d.Id())

	if bgpRPUUID != "" {
		_, _, err := client.RoutingProtocolsApi.DeleteConnectionRoutingProtocolByUuid(ctx, bgpRPUUID, connectionUUID)
		if err != nil {
			errors, ok := err.(v4.GenericSwaggerError).Model().([]v4.ModelError)
			if ok {
				// EQ-3142509 = Connection already deleted
				if hasModelErrorCode(errors, "EQ-3142509") {
					return diag.Diagnostics{}
				}
			}
			return diag.FromErr(fmt.Errorf("error deleting routing protocol (%s): %v", bgpRPUUID, err))
		}

		err = waitUntilRoutingProtocolsIsDeprovisioned(bgpRPUUID, connectionUUID, meta, ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error waiting for routing protocol (%s) deletion: %v", bgpRPUUID, err))
		}
	}
	_, _, err := client.RoutingProtocolsApi.DeleteConnectionRoutingProtocolByUuid(ctx, directRPUUID, connectionUUID)
	if err != nil {
		errors, ok := err.(v4.GenericSwaggerError).Model().([]v4.ModelError)
		if ok {
			// EQ-3142509 = Connection already deleted
			if hasModelErrorCode(errors, "EQ-3142509") {
				return diag.Diagnostics{}
			}
		}
		return diag.FromErr(fmt.Errorf("error deleting routing protocol (%s): %v", directRPUUID, err))
	}

	err = waitUntilRoutingProtocolsIsDeprovisioned(directRPUUID, connectionUUID, meta, ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for routing protocol (%s) deletion: %v", directRPUUID, err))
	}

	return diag.Diagnostics{}
}

func waitUntilRoutingProtocolsIsProvisioned(uuid string, connUuid string, meta interface{}, ctx context.Context) (v4.RoutingProtocolData, error) {
	log.Printf("Waiting for Routing Protocol(s) %s to be provisioned", uuid)
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

	inter, err := stateConf.WaitForStateContext(ctx)
	dbConn := v4.RoutingProtocolData{}

	if err == nil {
		dbConn = inter.(v4.RoutingProtocolData)
	}

	return dbConn, err
}

func waitUntilRoutingProtocolsIsDeprovisioned(uuid string, connUuid string, meta interface{}, ctx context.Context) error {
	log.Printf("Waiting for routing protocol to be deprovisioned, uuid %s", uuid)

	/* check if resource is not found */
	stateConf := &resource.StateChangeConf{
		Target: []string{
			strconv.Itoa(404),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, resp, _ := client.RoutingProtocolsApi.GetConnectionRoutingProtocolByUuid(ctx, uuid, connUuid)
			// fixme: check for error code instead?
			// ignore error for Target
			return dbConn, strconv.Itoa(resp.StatusCode), nil
		},
		Timeout:    5 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForRoutingProtocolsUpdateCompletion(rpChangeUuid string, uuid string, connUuid string, meta interface{}, ctx context.Context) (v4.RoutingProtocolChangeData, error) {
	log.Printf("Waiting for routing protocol update to complete, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Target: []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbConn, _, err := client.RoutingProtocolsApi.GetConnectionRoutingProtocolsChangeByUuid(ctx, connUuid, uuid, rpChangeUuid)
			if err != nil {
				return "", "", err
			}
			updatableState := ""
			if dbConn.Status == "COMPLETED" {
				updatableState = dbConn.Status
			}
			return dbConn, updatableState, nil
		},
		Timeout:    2 * time.Minute,
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
