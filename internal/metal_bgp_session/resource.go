package metal_bgp_session

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
	"github.com/equinix/terraform-provider-equinix/internal/helper"
)

type BgpSessionResourceModel struct {
    ID            types.String `tfsdk:"id"`
    DeviceID      types.String `tfsdk:"device_id"`
    AddressFamily types.String `tfsdk:"address_family"`
    DefaultRoute  types.Bool   `tfsdk:"default_route"`
    Status        types.String `tfsdk:"status"`
}

func (rm *BgpSessionResourceModel) parse(bgpSession *packngo.BGPSession) diag.Diagnostics {
	var diags diag.Diagnostics
	rm.ID = types.StringValue(bgpSession.ID)
	rm.DeviceID = types.StringValue(bgpSession.Device.ID)
	rm.AddressFamily = types.StringValue(bgpSession.AddressFamily)

	defaultRouteValue := false
	if bgpSession.DefaultRoute != nil {
		defaultRouteValue = *bgpSession.DefaultRoute
	}
	rm.DefaultRoute = types.BoolValue(defaultRouteValue)

	rm.Status = types.StringValue(bgpSession.Status)
	return diags
}


func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_metal_bgp_session",
				Schema: &bgpSessionResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
){
	// Create an instance of the BgpSessionResourceModel to hold the planned state
	var plan BgpSessionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the data for API request
	createRequest := packngo.CreateBGPSessionRequest{
		AddressFamily: plan.AddressFamily.ValueString(),
		DefaultRoute:  plan.DefaultRoute.ValueBoolPointer(),
	}

	// Retrieve the API client from the provider meta
	// r.Meta.AddModuleToMetalUserAgent(d)
	client := r.Meta.Metal

	// API call to create the BGP session
	bgpSession, _, err := client.BGPSessions.Create(plan.DeviceID.ValueString(), createRequest)
	if err != nil {
		err = helper.FriendlyError(err)
		resp.Diagnostics.AddError(
			"Error creating BGP session",
			"Could not create BGP session: " + err.Error(),
		)
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(bgpSession)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Retrieve the current state
	var state BgpSessionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client
	// r.Meta.AddModuleToMetalUserAgent(d)
	client := r.Meta.Metal

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// API call to get the current state of the BGP session
	bgpSession, _, err := client.BGPSessions.Get(id, nil)
	if err != nil {
		err = helper.FriendlyError(err)

		// Check if the BGP session no longer exists
		if helper.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"BGP session",
				fmt.Sprintf("[WARN] BGP session (%s) not found, removing from state", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading BGP session",
			"Could not read BGP session with ID " + id + ": " + err.Error(),
		)
		return
	}

	// Parse the API response into the Terraform state
	resp.Diagnostics.Append(state.parse(bgpSession)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
    ctx context.Context,
    req resource.UpdateRequest,
    resp *resource.UpdateResponse,
) {
	// This resource does not support updates
}

func (r *Resource) Delete(
	ctx context.Context, 
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve the current state
	var state BgpSessionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client
	// r.Meta.AddModuleToMetalUserAgent(d)
	client := r.Meta.Metal

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// API call to delete the BGP session
	deleteResp, err := client.BGPSessions.Delete(id)
	if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
		err = helper.FriendlyError(err)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete BGP session %s", id),
			err.Error(),
		)
	}
}

var bgpSessionResourceSchema = schema.Schema{
    Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
            Description: "The unique identifier for this BGP session",
            Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
        },
        "device_id": schema.StringAttribute{
            Description: "ID of device",
            Required:    true,
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.RequiresReplace(),
            },
        },
        "address_family": schema.StringAttribute{
            Description: "ipv4 or ipv6",
            Required:    true,
            PlanModifiers: []planmodifier.String{
                stringplanmodifier.RequiresReplace(),
            },
            Validators: []validator.String{
                stringvalidator.OneOf("ipv4", "ipv6"),
            },
        },
        "default_route": schema.BoolAttribute{
            Description: "Boolean flag to set the default route policy. False by default",
            Optional:    true,
            Default:     booldefault.StaticBool(false),
            PlanModifiers: []planmodifier.Bool{
                boolplanmodifier.RequiresReplace(),
            },
        },
        "status": schema.StringAttribute{
            Description: "Status of the session - up or down",
            Computed:    true,
        },
    },
}
