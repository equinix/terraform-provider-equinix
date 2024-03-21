package connection

import (
	"context"
	"fmt"
	"reflect"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_metal_connection",
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r *Resource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = resourceSchema(ctx)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {

	var plan ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)
	projectID := plan.ProjectID.ValueString()

	createRequest, diags := buildCreateRequest(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Create connection
	var err error
	var conn *metalv1.Interconnection

	if plan.Type.ValueString() == string(metalv1.INTERCONNECTIONTYPE_SHARED) || plan.Type.ValueString() == string(metalv1.INTERCONNECTIONTYPE_SHARED_PORT_VLAN) {
		request := client.InterconnectionsApi.CreateProjectInterconnection(ctx, projectID).
			CreateOrganizationInterconnectionRequest(createRequest)

		conn, _, err = request.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Metal Connection",
				"Could not create a shared Metal Connection: "+err.Error(),
			)
			return
		}
	}

	// Dedicated connection specific logic
	if plan.Type.ValueString() == string(metalv1.INTERCONNECTIONTYPE_DEDICATED) {

		organizationID := plan.OrganizationID.ValueString()

		// get organization ID from project
		if organizationID == "" {
			project, _, err := client.ProjectsApi.FindProjectById(ctx, projectID).
				// NB: organization.address and organization.billing_address needs to
				// be included otherwise Interconnection otherwise the response is
				// invalid against the API spec.
				Include([]string{"organization", "organization.address", "organization.billing_address"}).
				Execute()
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Failed to get Project %s", projectID),
					err.Error(),
				)
			}

			org := project.GetOrganization()
			organizationID = org.GetId()
		}

		request := client.InterconnectionsApi.CreateOrganizationInterconnection(ctx, organizationID).
			CreateOrganizationInterconnectionRequest(createRequest)

		// Create the dedicated connection
		conn, _, err = request.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Metal Connection",
				"Could not create a dedicated Metal Connection: "+err.Error(),
			)
			return
		}
	}

	// Use API client to get the current state of the resource
	conn, diags = getConnection(ctx, client, &resp.State, *conn.Id)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, conn)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// Use API client to get the current state of the resource
	conn, diags := getConnection(ctx, client, &resp.State, id)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(ctx, conn)...)
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
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	// Retrieve values from plan
	var state, plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := plan.ID.ValueString()

	// Prepare update request based on the changes
	updateRequest := metalv1.InterconnectionUpdateInput{}

	if !state.ContactEmail.Equal(plan.ContactEmail) {
		updateRequest.ContactEmail = plan.ContactEmail.ValueStringPointer()
	}
	if !state.Description.Equal(plan.Description) {
		updateRequest.Description = plan.Description.ValueStringPointer()
	}
	if !state.Mode.Equal(plan.Mode) {
		mode := metalv1.InterconnectionMode(plan.Mode.ValueString())
		updateRequest.Mode = &mode
	}

	if !state.Tags.Equal(plan.Tags) {
		tags := []string{}
		if diags := plan.Tags.ElementsAs(ctx, &tags, false); diags != nil {
			resp.Diagnostics.Append(diags...)
			return
		}
		updateRequest.Tags = tags
	}

	if !reflect.DeepEqual(updateRequest, metalv1.InterconnectionUpdateInput{}) {
		_, _, err := client.InterconnectionsApi.UpdateInterconnection(ctx, id).
			InterconnectionUpdateInput(updateRequest).
			Execute()

		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating Metal Connection",
				"Could not update Connection with ID "+id+": "+err.Error(),
			)
		}
	}

	// Use API client to get the current state of the resource
	conn, diags := getConnection(ctx, client, &resp.State, id)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(ctx, conn)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state back into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve the API client
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Retrieve the current state
	var state ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// API call to delete the Metal Connection
	deleteResp, err := client.Connections.Delete(id, true)
	if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(deleteResp, err) != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Metal Connection %s", id),
			equinix_errors.FriendlyError(err).Error(),
		)
	}
}

func buildDedicatedPortCreateRequest(ctx context.Context, plan ResourceModel, req *metalv1.CreateOrganizationInterconnectionRequest) (diags diag.Diagnostics) {
	mode, err := metalv1.NewDedicatedPortCreateInputModeFromValue(plan.Mode.ValueString())
	if err != nil {
		diags.AddError(
			"Error creating metal Connection",
			err.Error(),
		)
	}

	req.DedicatedPortCreateInput = &metalv1.DedicatedPortCreateInput{
		Type: metalv1.DEDICATEDPORTCREATEINPUTTYPE_DEDICATED,

		Name:       plan.Name.ValueString(),
		Metro:      plan.Metro.ValueString(),
		Speed:      plan.Speed.ValueStringPointer(),
		Mode:       mode,
		Redundancy: plan.Redundancy.ValueString(),
	}

	if email := plan.ContactEmail.ValueString(); email != "" {
		req.DedicatedPortCreateInput.ContactEmail = &email
	}
	if description := plan.Description.ValueString(); description != "" {
		req.DedicatedPortCreateInput.Description = &description
	}
	if project := plan.ProjectID.ValueString(); project != "" {
		req.DedicatedPortCreateInput.Project = &project
	}
	if facility := plan.Facility.ValueString(); facility != "" {
		req.DedicatedPortCreateInput.FacilityId = &facility
	}

	tagDiags := getPlanTags(ctx, plan, &req.DedicatedPortCreateInput.Tags)
	diags.Append(tagDiags...)

	return
}

func buildVLANFabricVCCreateRequest(ctx context.Context, plan ResourceModel, req *metalv1.CreateOrganizationInterconnectionRequest) diag.Diagnostics {
	diags := validateSharedConnection(plan)

	project := plan.ProjectID.ValueString()

	req.VlanFabricVcCreateInput = &metalv1.VlanFabricVcCreateInput{
		Type: metalv1.VLANFABRICVCCREATEINPUTTYPE_SHARED,

		Name:             plan.Name.ValueString(),
		Project:          &project,
		Metro:            plan.Metro.ValueString(),
		Redundancy:       plan.Redundancy.ValueString(),
		ServiceTokenType: metalv1.VlanFabricVcCreateInputServiceTokenType(plan.ServiceTokenType.ValueString()),
		Speed:            plan.Speed.ValueStringPointer(),
	}

	if email := plan.ContactEmail.ValueString(); email != "" {
		req.VlanFabricVcCreateInput.ContactEmail = &email
	}
	if description := plan.Description.ValueString(); description != "" {
		req.VlanFabricVcCreateInput.Description = &description
	}
	if facility := plan.Facility.ValueString(); facility != "" {
		req.VlanFabricVcCreateInput.FacilityId = &facility
	}

	vlansDiags := plan.Vlans.ElementsAs(ctx, &req.VlanFabricVcCreateInput.Vlans, true)
	diags.Append(vlansDiags...)

	tagDiags := getPlanTags(ctx, plan, &req.VlanFabricVcCreateInput.Tags)
	diags.Append(tagDiags...)

	return diags
}

func buildVRFFabricVCCreateRequest(ctx context.Context, plan ResourceModel, req *metalv1.CreateOrganizationInterconnectionRequest) diag.Diagnostics {
	diags := validateSharedConnection(plan)

	project := plan.ProjectID.ValueString()

	req.VrfFabricVcCreateInput = &metalv1.VrfFabricVcCreateInput{
		Type: metalv1.VLANFABRICVCCREATEINPUTTYPE_SHARED,

		Name:             plan.Name.ValueString(),
		Project:          &project,
		Metro:            plan.Metro.ValueString(),
		Redundancy:       plan.Redundancy.ValueString(),
		ServiceTokenType: metalv1.VlanFabricVcCreateInputServiceTokenType(plan.ServiceTokenType.ValueString()),
		Speed:            plan.Speed.ValueStringPointer(),
	}

	if email := plan.ContactEmail.ValueString(); email != "" {
		req.VrfFabricVcCreateInput.ContactEmail = &email
	}
	if description := plan.Description.ValueString(); description != "" {
		req.VrfFabricVcCreateInput.Description = &description
	}
	if facility := plan.Facility.ValueString(); facility != "" {
		req.VrfFabricVcCreateInput.FacilityId = &facility
	}

	vrfDiags := plan.Vrfs.ElementsAs(ctx, &req.VrfFabricVcCreateInput.Vrfs, false)
	diags.Append(vrfDiags...)

	tagDiags := getPlanTags(ctx, plan, &req.VrfFabricVcCreateInput.Tags)
	diags.Append(tagDiags...)

	return diags
}

func buildSharedPortVCVLANCreateRequest(ctx context.Context, plan ResourceModel, req *metalv1.CreateOrganizationInterconnectionRequest) diag.Diagnostics {
	diags := validateSharedConnection(plan)

	project := plan.ProjectID.ValueString()

	req.SharedPortVCVlanCreateInput = &metalv1.SharedPortVCVlanCreateInput{
		Type: metalv1.SHAREDPORTVCVLANCREATEINPUTTYPE_SHARED_PORT_VLAN,

		Name:    plan.Name.ValueString(),
		Project: project,
		Metro:   plan.Metro.ValueString(),
		Speed:   plan.Speed.ValueStringPointer(),
	}

	if email := plan.ContactEmail.ValueString(); email != "" {
		req.SharedPortVCVlanCreateInput.ContactEmail = &email
	}
	if description := plan.Description.ValueString(); description != "" {
		req.SharedPortVCVlanCreateInput.Description = &description
	}

	vlansDiags := plan.Vlans.ElementsAs(ctx, &req.SharedPortVCVlanCreateInput.Vlans, true)
	diags.Append(vlansDiags...)

	tagDiags := getPlanTags(ctx, plan, &req.SharedPortVCVlanCreateInput.Tags)
	diags.Append(tagDiags...)

	return diags
}

func getPlanTags(ctx context.Context, plan ResourceModel, tags *[]string) diag.Diagnostics {
	if len(plan.Tags.Elements()) != 0 {
		return plan.Tags.ElementsAs(context.Background(), tags, false)
	}
	return diag.Diagnostics{}
}

func validateSharedConnection(plan ResourceModel) (diags diag.Diagnostics) {
	// ensure project ID is set
	if plan.ProjectID.ValueString() == "" {
		diags.AddAttributeError(
			path.Root("project_id"),
			"Missing project_id",
			"project_id is required for 'shared' connection type",
		)
	}

	// ensure standard connection mode
	if plan.Mode.ValueString() == string(metalv1.INTERCONNECTIONMODE_TUNNEL) {
		diags.AddAttributeError(
			path.Root("mode"),
			"Wrong mode",
			"tunnel mode is not supported for 'shared' connections",
		)
	}

	entity := "vlan"
	entityLen := len(plan.Vlans.Elements())

	if entityLen == 0 {
		// using vrfs
		entity = "vrf"
		entityLen = len(plan.Vrfs.Elements())
	}

	targetLen := 1
	if plan.Redundancy.ValueString() == string(metalv1.INTERCONNECTIONREDUNDANCY_REDUNDANT) {
		targetLen = 2
	}

	// check redundancy
	if entityLen != targetLen {
		diags.AddAttributeError(
			path.Root(entity),
			fmt.Sprintf("Wrong number of %ss", entity),
			fmt.Sprintf("shared primary connections must have 1 %[1]s, shared redundant connections must 2 %[1]ss", entity),
		)
	}

	return
}

func buildCreateRequest(ctx context.Context, plan ResourceModel) (request metalv1.CreateOrganizationInterconnectionRequest, diags diag.Diagnostics) {
	hasVlans := len(plan.Vlans.Elements()) != 0
	hasVrfs := len(plan.Vrfs.Elements()) != 0
	hasSharedPortVlans := len(plan.Vlans.Elements()) != 0

	connType := metalv1.InterconnectionType(plan.Type.ValueString())

	if hasVlans && hasVrfs {
		// vlans and vrfs are mutually exclusive
		diags.AddError(
			"Cannot specify vlans and vrfs together",
			"Shared connections must specify either vlans or vrfs",
		)
		return
	}

	if connType == metalv1.INTERCONNECTIONTYPE_DEDICATED && (hasVlans || hasVrfs) {
		diags.AddAttributeError(
			path.Root("type"),
			"Cannot specify vlans or vrfs",
			"Dedicated connections must not specify vlans or vrfs",
		)

		return
	} else if connType == metalv1.INTERCONNECTIONTYPE_SHARED && !(hasVlans || hasVrfs) {
		diags.AddAttributeError(
			path.Root("type"),
			"Must specify either vlans or vrfs",
			"Shared connections must specify either vlans or vrfs",
		)

		return
	} else if connType == metalv1.INTERCONNECTIONTYPE_SHARED_PORT_VLAN && (hasSharedPortVlans) {
		diags.AddAttributeError(
			path.Root("type"),
			"Must specify vlans",
			"Port Shared connections must specify vlans",
		)

		return
	}

	// ensure speed is valid if specified
	if speed := plan.Speed.ValueString(); speed != "" {
		if err := validateSpeedStr(speed); err != nil {
			diags.AddAttributeError(
				path.Root("mode"),
				"Invalid speed",
				err.Error(),
			)
		}
	}

	var requestFunc func(context.Context, ResourceModel, *metalv1.CreateOrganizationInterconnectionRequest) diag.Diagnostics

	switch {
	case hasVlans:
		requestFunc = buildVLANFabricVCCreateRequest

	case hasVrfs:
		requestFunc = buildVRFFabricVCCreateRequest

	case hasSharedPortVlans:
		requestFunc = buildSharedPortVCVLANCreateRequest

	default:
		// has to be a dedicated connection
		requestFunc = buildDedicatedPortCreateRequest
	}

	return request, requestFunc(ctx, plan, &request)
}

func getConnection(ctx context.Context, client *metalv1.APIClient, state *tfsdk.State, id string) (*metalv1.Interconnection, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Use API client to get the current state of the resource
	conn, _, err := client.InterconnectionsApi.GetInterconnection(ctx, id).
		// NB: organization.address and organization.billing_address needs to
		// be included otherwise Interconnection otherwise the response is
		// invalid against the API spec.
		Include([]string{"service_tokens", "organization", "organization.address", "organization.billing_address", "facility", "metro", "project"}).
		Execute()

	if err != nil {
		// If the Metal Connection is not found, remove it from the state
		if equinix_errors.IsNotFound(err) {
			diags.AddWarning(
				"Metal Connection",
				fmt.Sprintf("[WARN] Connection (%s) not found, removing from state", id),
			)
			state.RemoveResource(ctx)
			return nil, diags
		}

		diags.AddError(
			"Error reading Metal Connection",
			"Could not read Metal Connection with ID "+id+": "+err.Error(),
		)
		return nil, diags
	}
	return conn, diags
}
