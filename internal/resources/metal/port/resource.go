package port

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"slices"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
Race conditions:
 - assigning and removing the same VLAN in the same terraform run
 - Bonding a bond port where underlying eth port has vlans assigned, and those vlans are being removed in the same terraform run
*/

var (
	l2Types = []metalv1.PortNetworkType{"layer2-individual", "layer2-bonded"}
	l3Types = []metalv1.PortNetworkType{"layer3", "hybrid", "hybrid-bonded"}
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		ReadWithoutTimeout: resourceMetalPortRead,
		// Create and Update are the same func
		CreateContext: resourceMetalPortUpdate,
		UpdateContext: resourceMetalPortUpdate,
		DeleteContext: resourceMetalPortDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"port_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the port to lookup",
				ForceNew:    true,
			},
			"bonded": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Flag indicating whether the port should be bonded",
			},
			"layer2": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag indicating whether the port is in layer2 (or layer3) mode. The `layer2` flag can be set only for bond ports.",
			},
			"native_vlan_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "UUID of native VLAN of the port",
			},
			"vxlan_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Description:   "VLAN VXLAN ids to attach (example: [1000])",
				Elem:          &schema.Schema{Type: schema.TypeInt},
				ConflictsWith: []string{"vlan_ids"},
			},
			"vlan_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Description:   "UUIDs VLANs to attach. To avoid jitter, use the UUID and not the VXLAN",
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"vxlan_ids"},
			},
			"reset_on_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Behavioral setting to reset the port to default settings (layer3 bonded mode without any vlan attached) before delete/destroy",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the port to look up, e.g. bond0, eth1",
			},
			"network_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "One of layer2-bonded, layer2-individual, layer3, hybrid and hybrid-bonded. This attribute is only set on bond ports.",
			},
			"disbond_supported": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating whether the port can be removed from a bond",
			},
			"bond_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the bond port",
			},
			"bond_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the bond port",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port type",
			},
			"mac": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MAC address of the port",
			},
		},
	}
}

func resourceMetalPortUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	start := time.Now()
	cpr, _, err := getClientPortResource(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, f := range [](func(context.Context, *ClientPortResource) error){
		portSanityChecks,
		batchVlans(start, true),
		makeDisbond,
		convertToL2,
		makeBond,
		convertToL3,
		batchVlans(start, false),
		updateNativeVlan,
	} {
		if err := f(ctx, cpr); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMetalPortRead(ctx, d, meta)
}

func resourceMetalPortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	port, resp, err := getPortByResourceData(ctx, d, client)
	if err != nil {
		if resp != nil && slices.Contains([]int{http.StatusNotFound, http.StatusForbidden}, resp.StatusCode) {
			log.Printf("[WARN] Port (%s) not accessible, removing from state", d.Id())
			d.SetId("")

			return nil
		}
		return diag.FromErr(err)
	}
	m := map[string]interface{}{
		"port_id":           port.GetId(),
		"type":              port.GetType(),
		"name":              port.GetName(),
		"network_type":      port.GetNetworkType(),
		"mac":               port.Data.GetMac(),
		"bonded":            port.Data.GetBonded(),
		"disbond_supported": port.GetDisbondOperationSupported(),
	}
	l2 := slices.Contains(l2Types, port.GetNetworkType())
	l3 := slices.Contains(l3Types, port.GetNetworkType())

	if l2 {
		m["layer2"] = true
	}
	if l3 {
		m["layer2"] = false
	}

	if port.NativeVirtualNetwork != nil {
		m["native_vlan_id"] = port.NativeVirtualNetwork.GetId()
	}

	vlans := []string{}
	vxlans := []int{}
	for _, n := range port.VirtualNetworks {
		vlans = append(vlans, n.GetId())
		vxlans = append(vxlans, int(n.GetVxlan()))
	}
	m["vlan_ids"] = vlans
	m["vxlan_ids"] = vxlans

	if port.Bond != nil {
		m["bond_id"] = port.Bond.GetId()
		m["bond_name"] = port.Bond.GetName()
	}

	d.SetId(port.GetId())
	return diag.FromErr(equinix_schema.SetMap(d, m))
}

func resourceMetalPortDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resetRaw, resetOk := d.GetOk("reset_on_delete")
	if resetOk && resetRaw.(bool) {
		start := time.Now()
		cpr, resp, err := getClientPortResource(ctx, d, meta)
		if err != nil {
			if resp != nil && !slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, resp.StatusCode) {
				return diag.FromErr(err)
			}
		}

		// to reset the port to defaults we iterate through helpers (used in
		// create/update), some of which rely on resource state. reuse those helpers by
		// setting ephemeral state.
		port := Resource()
		copy := port.Data(d.State())
		cpr.Resource = copy
		if err = equinix_schema.SetMap(cpr.Resource, map[string]interface{}{
			"layer2":         false,
			"bonded":         true,
			"native_vlan_id": nil,
			"vlan_ids":       []string{},
			"vxlan_ids":      nil,
		}); err != nil {
			return diag.FromErr(err)
		}
		for _, f := range [](func(context.Context, *ClientPortResource) error){
			batchVlans(start, true),
			makeBond,
			convertToL3,
		} {
			if err := f(ctx, cpr); err != nil {
				return diag.FromErr(err)
			}
		}
		// TODO(displague) error or warn?
		if warn := ProperlyDestroyed(cpr.Port); warn != nil {
			log.Printf("[WARN] %s\n", warn)
		}
	}
	return nil
}

func ProperlyDestroyed(port *metalv1.Port) error {
	var errs []string
	if !port.Data.GetBonded() {
		errs = append(errs, fmt.Sprintf("port %s wasn't bonded after equinix_metal_port destroy;", port.GetId()))
	}
	if port.GetType() == "NetworkBondPort" && port.GetNetworkType() != "layer3" {
		errs = append(errs, "bond port should be in layer3 type after destroy;")
	}
	if port.NativeVirtualNetwork != nil {
		errs = append(errs, "port should not have native VLAN assigned after destroy;")
	}
	if len(port.VirtualNetworks) != 0 {
		errs = append(errs, "port should not have VLANs attached after destroy")
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", errs)
	}

	return nil
}
