package equinix

import (
	"context"
	"log"
	"slices"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
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
	l2Types = []string{"layer2-individual", "layer2-bonded"}
	l3Types = []string{"layer3", "hybrid", "hybrid-bonded"}
)

func resourceMetalPort() *schema.Resource {
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
	cpr, _, err := getClientPortResource(d, meta)
	if err != nil {
		return diag.FromErr(equinix_errors.FriendlyError(err))
	}

	for _, f := range [](func(*ClientPortResource) error){
		portSanityChecks,
		batchVlans(ctx, start, true),
		makeDisbond,
		convertToL2,
		makeBond,
		convertToL3,
		batchVlans(ctx, start, false),
		updateNativeVlan,
	} {
		if err := f(cpr); err != nil {
			return diag.FromErr(equinix_errors.FriendlyError(err))
		}
	}

	return resourceMetalPortRead(ctx, d, meta)
}

func resourceMetalPortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	port, err := getPortByResourceData(d, client)
	if err != nil {
		if equinix_errors.IsNotFound(err) || equinix_errors.IsForbidden(err) {
			log.Printf("[WARN] Port (%s) not accessible, removing from state", d.Id())
			d.SetId("")

			return nil
		}
		return diag.FromErr(err)
	}
	m := map[string]interface{}{
		"port_id":           port.ID,
		"type":              port.Type,
		"name":              port.Name,
		"network_type":      port.NetworkType,
		"mac":               port.Data.MAC,
		"bonded":            port.Data.Bonded,
		"disbond_supported": port.DisbondOperationSupported,
	}
	l2 := slices.Contains(l2Types, port.NetworkType)
	l3 := slices.Contains(l3Types, port.NetworkType)

	if l2 {
		m["layer2"] = true
	}
	if l3 {
		m["layer2"] = false
	}

	if port.NativeVirtualNetwork != nil {
		m["native_vlan_id"] = port.NativeVirtualNetwork.ID
	}

	vlans := []string{}
	vxlans := []int{}
	for _, n := range port.AttachedVirtualNetworks {
		vlans = append(vlans, n.ID)
		vxlans = append(vxlans, n.VXLAN)
	}
	m["vlan_ids"] = vlans
	m["vxlan_ids"] = vxlans

	if port.Bond != nil {
		m["bond_id"] = port.Bond.ID
		m["bond_name"] = port.Bond.Name
	}

	d.SetId(port.ID)
	return diag.FromErr(equinix_schema.SetMap(d, m))
}

func resourceMetalPortDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resetRaw, resetOk := d.GetOk("reset_on_delete")
	if resetOk && resetRaw.(bool) {
		start := time.Now()
		cpr, resp, err := getClientPortResource(d, meta)
		if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
			return diag.FromErr(err)
		}

		// to reset the port to defaults we iterate through helpers (used in
		// create/update), some of which rely on resource state. reuse those helpers by
		// setting ephemeral state.
		port := resourceMetalPort()
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
		for _, f := range [](func(*ClientPortResource) error){
			batchVlans(ctx, start, true),
			makeBond,
			convertToL3,
		} {
			if err := f(cpr); err != nil {
				return diag.FromErr(err)
			}
		}
		// TODO(displague) error or warn?
		if warn := portProperlyDestroyed(cpr.Port); warn != nil {
			log.Printf("[WARN] %s\n", warn)
		}
	}
	return nil
}
