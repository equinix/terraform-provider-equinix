package port

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"slices"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

func NewResource() resource.Resource {
	r := &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_metal_port",
			},
		),
	}

	r.SetDefaultUpdateTimeout(30 * time.Minute)
	r.SetDefaultDeleteTimeout(30 * time.Minute)

	return r
}

type Resource struct {
	framework.BaseResource
	framework.WithTimeouts
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	s := resourceSchema(ctx)
	if s.Blocks == nil {
		s.Blocks = make(map[string]schema.Block)
	}

	s.Blocks["timeouts"] = timeouts.Block(ctx, timeouts.opts{
		Create: true,
		Update: true,
		Delete: true,
	})
	resp.Schema = s
}



func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	return r.Read(ctx, d, meta)

}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state resourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	
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
