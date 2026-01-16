package port

import (
	"context"
	"fmt"
	"time"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/batch"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

type tfResource struct {
	framework.BaseResource
	framework.WithTimeouts
}

// NewResource returns the TF resource representing device network ports.
func NewResource() resource.Resource {
	r := &tfResource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_metal_port",
			},
		),
	}

	r.SetDefaultCreateTimeout(20 * time.Minute)
	r.SetDefaultUpdateTimeout(20 * time.Minute)
	r.SetDefaultUpdateTimeout(20 * time.Minute)

	return r
}

// Schema implements resource.Resource.
// Subtle: this method shadows the method (BaseResource).Schema of Resource.BaseResource.
func (r *tfResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	s := resourceSchema(ctx)

	if s.Blocks == nil {
		s.Blocks = make(map[string]schema.Block)
	}

	resp.Schema = s
}

// Create implements resource.Resource.
func (r *tfResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state, plan resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)
	id := plan.PortID.ValueString()

	port, _, err := client.PortsApi.FindPortById(ctx, id).
		Include([]string{"native_virtual_network", "virtual_networks"}).
		Execute()

	if err != nil {
		if equinix_errors.IsNotFound(err) {
			resp.Diagnostics.AddWarning("Metal Port",
				fmt.Sprintf("[WARN] Metal Port (%s) not found, removing from state", id),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading Metal Port",
			fmt.Sprintf("Could not read Metal Port with ID %s: %v", id, err),
		)
		return
	}

	ops, diags := plan.ToExecutionPlan(ctx, port)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	_, diags = r.performOperations(ctx, client, ops, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	port, diags = r.refreshPort(ctx, client, plan.PortID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(state.parse(ctx, port)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read implements resource.Resource.
func (r *tfResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	id := state.PortID.ValueString()

	port, _, err := client.PortsApi.FindPortById(ctx, id).
		Include([]string{"virtual_networks", "native_virtual_network"}).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed reading Metal Port",
			fmt.Sprintf("Could not find port with id %s: %s", id, err),
		)
		return
	}

	resp.Diagnostics.Append(state.parse(ctx, port)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (r *tfResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)
	id := state.PortID.ValueString()

	port, _, err := client.PortsApi.FindPortById(ctx, id).
		Include([]string{"native_virtual_network", "virtual_networks"}).
		Execute()

	if err != nil {
		if equinix_errors.IsNotFound(err) {
			resp.Diagnostics.AddWarning("Metal Port",
				fmt.Sprintf("[WARN] Metal Port (%s) not found, removing from state", id),
			)
			req.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading Metal Port",
			fmt.Sprintf("Could not read Metal Port with ID %s: %v", id, err),
		)
		return
	}

	ops, diags := plan.ToExecutionPlan(ctx, port)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	batch, diags := r.updateVlanAssigments(ctx, state, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.waitForBatch(ctx, client, batch)...)

	_, diags = r.performOperations(ctx, client, ops, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	port, diags = r.refreshPort(ctx, client, state.PortID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(state.parse(ctx, port)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete implements resource.Resource.
func (r *tfResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state, actual resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the reset_on_delete option is not set, just let terraform cleanup the resource from the state.
	if state.ResetOnDelete.IsNull() {
		return
	}

	// If the reset_on_delete option _is_ set, but is false, let Terraform clean up the resource from the state.
	if !state.ResetOnDelete.ValueBool() {
		return
	}

	// If we're here, we have reset_on_delete = true, so we'll perform operations to reset the port.
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)
	id := state.PortID.ValueString()

	port, _, err := client.PortsApi.FindPortById(ctx, id).
		Include([]string{"native_virtual_network", "virtual_networks"}).
		Execute()

	if err != nil {
		if equinix_errors.IsNotFound(err) {
			resp.Diagnostics.AddWarning("Metal Port",
				fmt.Sprintf("[WARN] Metal Port (%s) not found, removing from state", id),
			)
			req.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading Metal Port",
			fmt.Sprintf("Could not read Metal Port with ID %s: %v", id, err),
		)
		return
	}

	ops := []string{}

	resp.Diagnostics.Append(actual.parse(ctx, port)...)
	if resp.Diagnostics.HasError() {
		return
	}

	wantBonded := state.Bonded.ValueBool()
	isBonded := actual.Bonded.ValueBool()
	if wantBonded && !isBonded {
		ops = append(ops, "bond")
	}

	wantLayer3 := state.Layer2.IsNull() || !state.Layer2.ValueBool()
	portIsLayer2 := actual.Layer2.ValueBool()
	if wantLayer3 && portIsLayer2 {
		ops = append(ops, "convertToLayer3")
	}

	port, diags := r.performOperations(ctx, client, ops, state)
	resp.Diagnostics.Append(diags...)

	if err := ProperlyDestroyed(port); err != nil {
		resp.Diagnostics.AddError(
			"Failed to reset port state",
			err.Error(),
		)
		return
	}
}

func (r *tfResource) performOperations(ctx context.Context, client *metalv1.APIClient, ops []string, state resourceModel) (*metalv1.Port, diag.Diagnostics) {
	var diags diag.Diagnostics
	var newPort *metalv1.Port
	var err error
	portID := state.PortID.ValueString()

	for _, op := range ops {
		switch op {
		case "disbond":
			newPort, _, err = client.PortsApi.DisbondPort(ctx, portID).Execute()
		case "bond":
			newPort, _, err = client.PortsApi.BondPort(ctx, portID).Execute()
		case "toLayer2":
			input := metalv1.PortAssignInput{}
			newPort, _, err = client.PortsApi.ConvertLayer2(ctx, portID).PortAssignInput(input).Execute()
		case "toLayer3":
			input := metalv1.PortConvertLayer3Input{
				RequestIps: []metalv1.PortConvertLayer3InputRequestIpsInner{
					{AddressFamily: metalv1.PtrInt32(4), Public: metalv1.PtrBool(true)},
					{AddressFamily: metalv1.PtrInt32(4), Public: metalv1.PtrBool(false)},
					{AddressFamily: metalv1.PtrInt32(6), Public: metalv1.PtrBool(true)},
				},
			}
			newPort, _, err = client.PortsApi.ConvertLayer3(ctx, portID).PortConvertLayer3Input(input).Execute()
		case "removeNativeVlan":
			newPort, _, err = client.PortsApi.DeleteNativeVlan(ctx, portID).Execute()
		case "assignNativeVlan":
			vlan := state.NativeVlanID.ValueString()
			newPort, _, err = client.PortsApi.AssignNativeVlan(ctx, portID).Vnid(vlan).Execute()
		}

		if err != nil {
			diags.AddError(
				"Failed to modify Port",
				fmt.Sprintf("Port %s failed to execute '%s' operation: %s", portID, op, err),
			)
		}
	}

	return newPort, diags
}

func (r *tfResource) refreshPort(ctx context.Context, client *metalv1.APIClient, portID string) (*metalv1.Port, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, _, err := client.PortsApi.FindPortById(ctx, portID).
		Include([]string{"native_virtual_network", "virtual_networks"}).
		Execute()
	if err != nil {
		diags.AddError("Failed to refresh port data", fmt.Sprintf("refreshing port encountered: %s", err))
		return nil, diags
	}

	return p, diags
}

// updateVlanAssignments takes in the state and the plan and
// determines the payload to the VLAN batch assignment endpoint.
// It determines which VLANs should be added and which should be removed,
// without touching any common between the currents state and the desired state.
// The function returns a batchID, and does not wait for completion.
func (r *tfResource) updateVlanAssigments(ctx context.Context, state resourceModel, plan *resourceModel) (*batch.VlanBatch, diag.Diagnostics) {
	var diags diag.Diagnostics
	var b *batch.VlanBatch

	if plan == nil {
		return b, diags
	}

	vlanIDsSet := !plan.VLANIDs.IsUnknown()
	vxlanIDsSet := !plan.VXLANIDs.IsUnknown()

	if vlanIDsSet {
		assignedVlans := []string{}
		diags = state.VLANIDs.ElementsAs(ctx, &assignedVlans, false)

		if diags.HasError() {
			return b, diags
		}

		desiredVlans := []string{}
		plan.VLANIDs.ElementsAs(ctx, &desiredVlans, false)
		if diags.HasError() {
			return b, diags
		}

		b = batch.NewVlanBatch(state.PortID.ValueString())
		toAdd := setDifference(desiredVlans, assignedVlans)
		toRemove := setDifference(assignedVlans, desiredVlans)

		for _, vlan := range toAdd {
			b.AddAssignment(vlan)
		}

		for _, vlan := range toRemove {
			b.RemoveAssignment(vlan)
		}

		return b, diags
	}

	if vxlanIDsSet {
		assignedVlans := []int32{}
		diags = state.VXLANIDs.ElementsAs(ctx, &assignedVlans, false)

		if diags.HasError() {
			return b, diags
		}

		desiredVlans := []int32{}
		plan.VLANIDs.ElementsAs(ctx, &desiredVlans, false)
		if diags.HasError() {
			return b, diags
		}

		b = batch.NewVlanBatch(state.PortID.ValueString())
		toAdd := setDifference(desiredVlans, assignedVlans)
		toRemove := setDifference(assignedVlans, desiredVlans)

		for _, vlan := range toAdd {
			b.AddAssignment(fmt.Sprintf("%d", vlan))
		}

		for _, vlan := range toRemove {
			b.RemoveAssignment(fmt.Sprintf("%d", vlan))
		}

		return b, diags
	}

	return b, diags
}

func (r *tfResource) waitForBatch(ctx context.Context, client *metalv1.APIClient, b *batch.VlanBatch) diag.Diagnostics {
	var diags diag.Diagnostics
	_, _, err := b.Execute(ctx, client)
	if err != nil {
		diags.AddError(
			"Failed to wait for VLAN batch to complete",
			fmt.Sprintf("Batch encountered error: %s", err),
		)
	}

	return diags

}

// setDifference returns the elements LEFT that are not in RIGHT.
func setDifference[T comparable](left []T, right []T) []T {
	result := []T{}

	rightSet := make(map[T]bool, len(right))
	for _, v := range right {
		rightSet[v] = true
	}

	for _, v := range left {
		_, exists := rightSet[v]
		if !exists {
			result = append(result, v)
		}
	}

	return result
}

// ProperlyDestroyed does some state checking, we don't actually destroy the port
// since that's a resource tied to the lifetime of an instance.
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
