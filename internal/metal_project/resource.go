package metal_project

import (
	"context"
	"path"
	"regexp"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
	"github.com/equinix/terraform-provider-equinix/internal/helper"
)

var uuidRE = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

type ProjectResourceModel struct {
	ID              types.String    `tfsdk:"id"`
	Name            types.String    `tfsdk:"name"`
	Created         types.String    `tfsdk:"created"`
	Updated         types.String    `tfsdk:"updated"`
	BackendTransfer types.Bool      `tfsdk:"backend_transfer"`
	PaymentMethodID types.String    `tfsdk:"payment_method_id"`
	OrganizationID  types.String    `tfsdk:"organization_id"`
	BGPConfig       *BGPConfigModel `tfsdk:"bgp_config"`
}

type BGPConfigModel struct {
	DeploymentType types.String `tfsdk:"deployment_type"`
	ASN            types.Int64  `tfsdk:"asn"`
	MD5            types.String `tfsdk:"md5"`
	Status         types.String `tfsdk:"status"`
	MaxPrefix      types.Int64  `tfsdk:"max_prefix"`
}

func (bgp *BGPConfigModel) equal(other *BGPConfigModel) bool {
    if bgp == nil && other == nil {
        return true
    }
    if bgp == nil || other == nil {
        return false
    }
    return bgp.DeploymentType == other.DeploymentType &&
           bgp.ASN == other.ASN &&
           bgp.MD5 == other.MD5 &&
           bgp.Status == other.Status &&
           bgp.MaxPrefix == other.MaxPrefix
}

func (rm *ProjectResourceModel) parse(project *packngo.Project, bgpConfig *packngo.BGPConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	rm.ID = types.StringValue(project.ID) 
	rm.Name = types.StringValue(project.Name)
	rm.Created = types.StringValue(project.Created)
	rm.Updated = types.StringValue(project.Updated)
	rm.BackendTransfer = types.BoolValue(project.BackendTransfer)

	if len(project.PaymentMethod.URL) != 0 {
		rm.PaymentMethodID = types.StringValue(path.Base(project.PaymentMethod.URL))
	}

	rm.OrganizationID = types.StringValue(path.Base(project.Organization.URL))

	// Handle BGP Config if present
	if bgpConfig != nil {
		rm.BGPConfig = &BGPConfigModel{
			DeploymentType: types.StringValue(bgpConfig.DeploymentType),
			ASN:            types.Int64Value(int64(bgpConfig.Asn)),
			MD5:            types.StringValue(bgpConfig.Md5),
			Status:         types.StringValue(bgpConfig.Status),
			MaxPrefix:      types.Int64Value(int64(bgpConfig.MaxPrefix)),
		}
	}

	return diags
}

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_metal_project",
				Schema: &projectResourceSchema,
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
) {
	// Create an instance of your resource model to hold the planned state
	var plan ProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the data for API request
	createRequest := packngo.ProjectCreateRequest{
		Name: plan.Name.ValueString(),
	}

	// Include optional fields if they are set
	if !plan.OrganizationID.IsNull() {
		createRequest.OrganizationID = plan.OrganizationID.ValueString()
	}

	// Retrieve the API client
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// API call to create the project
	project, _, err := client.Projects.Create(&createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project: " + err.Error(),
		)
		return
	}

	// Handle BGP Config if present
	if plan.BGPConfig != nil {
		bgpCreateRequest := packngo.CreateBGPConfigRequest{
			DeploymentType: plan.BGPConfig.DeploymentType.ValueString(),
			Asn:            int(plan.BGPConfig.ASN.ValueInt64()),
		}
		if !plan.BGPConfig.MD5.IsNull() {
			bgpCreateRequest.Md5 = plan.BGPConfig.MD5.ValueString()
		}
		_, err := client.BGPConfig.Create(project.ID, bgpCreateRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating BGP configuration",
				"Could not create BGP configuration for project: " + err.Error(),
			)
			return
		}
	}

	// Enable Backend Transfer if True
	if plan.BackendTransfer.ValueBool() {
		pur := packngo.ProjectUpdateRequest{
			BackendTransfer: plan.BackendTransfer.ValueBoolPointer(),
		}
		project, _, err = client.Projects.Update(project.ID, &pur)
		if err != nil {
			err = helper.FriendlyError(err)
			resp.Diagnostics.AddError(
					"Error enabling Backend Transfer",
					"Could not enable Backend Transfer for project with ID " + project.ID + ": " + err.Error(),
			)
			return
		}
	}

	// Fetch BGP Config if needed
	var bgpConfig *packngo.BGPConfig
	if plan.BGPConfig != nil {
		bgpConfig, diags = fetchBGPConfig(client, project.ID)
        diags.Append(diags...)
		if diags.HasError(){
			return
		}
	}

	// Parse API response into the Terraform state
	stateDiags := (&plan).parse(project, bgpConfig)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Retrieve the current state
	var state ProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider meta
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// API call to get the current state of the project
	project, diags := fetchProject(client, id)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Fetch BGP Config if needed
	var bgpConfig *packngo.BGPConfig
	if state.BGPConfig != nil {
		bgpConfig, diags = fetchBGPConfig(client, id)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
	}

	// Parse the API response into the Terraform state
	resp.Diagnostics.Append(state.parse(project, bgpConfig)...)
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
	// Retrieve the current state and plan
	var state, plan ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// Prepare update request based on the changes
	updateRequest := &packngo.ProjectUpdateRequest{}
	if state.Name != plan.Name {
		updateRequest.Name = plan.Name.ValueStringPointer()
	}
	if state.PaymentMethodID != plan.Name {
		updateRequest.PaymentMethodID = plan.PaymentMethodID.ValueStringPointer()
	}
	if state.BackendTransfer != plan.BackendTransfer {
		updateRequest.BackendTransfer = plan.BackendTransfer.ValueBoolPointer()
	}

	// Handle BGP Config changes
	bgpConfig, diags := handleBGPConfigChanges(client, &plan, &state, id)
	resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	// NOTE (ocobles): adding this in the condition to replicate old behavior
	// but it is not clear to me if it was a mistake. I think the project
	// should be updated if has changes regardless of whether there are
	// changes to the BGP configuration or not.
	// Open discussion: https://github.com/equinix/terraform-provider-equinix/discussions/466
	var project *packngo.Project
	var err error
	if plan.BGPConfig.equal(state.BGPConfig) {
		// API call to update the project
		project, _, err = client.Projects.Update(id, updateRequest)
		if err != nil {
			friendlyErr := helper.FriendlyError(err)
			resp.Diagnostics.AddError(
				"Error updating project",
				"Could not update project with ID " + id + ": " + friendlyErr.Error(),
			)
			return
		}	
	} else {
		// Fetch the project
		project, diags = fetchProject(client, id)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(project, bgpConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the updated state back into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve the current state
    var state ProjectResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the ID of the resource from the state
    id := state.ID.ValueString()

    // API call to delete the project
	deleteResp, err := client.Projects.Delete(id)
	if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
		err = helper.FriendlyError(err)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Project %s", id),
			err.Error(),
		)
	}
}


var projectResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier for the project.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the project. The maximum length is 80 characters.",
			Required:    true,
		},
		"created": schema.StringAttribute{
			Description: "The timestamp for when the project was created",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "The timestamp for the last time the project was updated",
			Computed:    true,
		},
		"backend_transfer": schema.BoolAttribute{
			Description: "Enable or disable Backend Transfer, default is false",
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"payment_method_id": schema.StringAttribute{
			Description: "The UUID of payment method for this project.",
			Optional:    true,
			Computed:    true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(uuidRE, "must be a valid UUID"),
			},
		},
		"organization_id": schema.StringAttribute{
			Description: "The UUID of organization under which the project is created.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.RegexMatches(uuidRE, "must be a valid UUID"),
			},
		},
		"bgp_config": schema.SingleNestedAttribute{
			Description: "Optional BGP settings.",
			Optional:    true,
			Attributes: bgpConfigSchema,
		},
	},
}


var bgpConfigSchema = map[string]schema.Attribute{
	"deployment_type": schema.StringAttribute{
		Description:  "The BGP deployment type, either 'local' or 'global'.",
		Required:     true,
		Validators: []validator.String{
			stringvalidator.OneOf("local", "global"),
		},
	},
	"asn": schema.Int64Attribute{
		Description: "Autonomous System Number for local BGP deployment",
		Required:    true,
	},
	"md5": schema.StringAttribute{
		Description: "Password for BGP session in plaintext (not a checksum)",
		Sensitive:   true,
		Optional:    true,
	},
	"status": schema.StringAttribute{
		Description: "Status of BGP configuration in the project",
		Computed:    true,
	},
	"max_prefix": schema.Int64Attribute{
		Description: "The maximum number of route filters allowed per server",
		Computed:    true,
	},
}

func fetchProject(client *packngo.Client, projectID string) (*packngo.Project, diag.Diagnostics) {
    var diags diag.Diagnostics

    project, _, err := client.Projects.Get(projectID, nil)
    if err != nil {
        friendlyErr := helper.FriendlyError(err)

        // Check if the Project no longer exists
        if helper.IsNotFound(friendlyErr) {
            diags.AddWarning(
                "Project not found",
                fmt.Sprintf("Project (%s) not found, removing from state", projectID),
            )
        } else {
            diags.AddError(
                "Error reading project",
                "Could not read project with ID " + projectID + ": " + friendlyErr.Error(),
            )
        }
        return nil, diags
    }

    return project, diags
}

func fetchBGPConfig(client *packngo.Client, projectID string) (*packngo.BGPConfig, diag.Diagnostics) {
    var diags diag.Diagnostics

    bgpConfig, _, err := client.BGPConfig.Get(projectID, nil)
    if err != nil {
		friendlyErr := helper.FriendlyError(err)
        diags.AddError(
            "Error reading BGP configuration",
            "Could not read BGP configuration for project with ID " + projectID + ": " + friendlyErr.Error(),
        )
        return nil, diags
    }

    return bgpConfig, diags
}

func handleBGPConfigChanges(client *packngo.Client, plan *ProjectResourceModel, state *ProjectResourceModel, projectID string) (*packngo.BGPConfig, diag.Diagnostics) {
    var diags diag.Diagnostics
    var bgpConfig *packngo.BGPConfig

    bgpAdded := plan.BGPConfig != nil && state.BGPConfig == nil
    bgpRemoved := plan.BGPConfig == nil && state.BGPConfig != nil
    bgpChanged := plan.BGPConfig != nil && state.BGPConfig != nil && !plan.BGPConfig.equal(state.BGPConfig)

    if bgpAdded {
        // Create BGP Config
        bgpCreateRequest := packngo.CreateBGPConfigRequest{
            DeploymentType: plan.BGPConfig.DeploymentType.ValueString(),
            Asn:            int(plan.BGPConfig.ASN.ValueInt64()),
        }
        if !plan.BGPConfig.MD5.IsNull() {
            bgpCreateRequest.Md5 = plan.BGPConfig.MD5.ValueString()
        }
        _, err := client.BGPConfig.Create(projectID, bgpCreateRequest)
        if err != nil {
			friendlyErr := helper.FriendlyError(err)
            diags.AddError(
                "Error creating BGP configuration",
                "Could not create BGP configuration for project: " + friendlyErr.Error(),
            )
            return nil, diags
        }

        // Fetch the newly created BGP Config
        bgpConfig, diags = fetchBGPConfig(client, projectID)
        diags.Append(diags...)
    } else if bgpRemoved {
        bgpConfStr := fmt.Sprintf(
            "bgp_config {\n"+
                "  deployment_type = \"%s\"\n"+
                "  md5 = \"%s\"\n"+
                "  asn = %d\n"+
                "}",
            state.BGPConfig.DeploymentType.ValueString(),
            state.BGPConfig.MD5.ValueString(),
            state.BGPConfig.ASN.ValueInt64(),
        )
        diags.AddError(
            "Error removing BGP configuration",
            fmt.Sprintf("BGP Config cannot be removed from a project, please add back\n%s", bgpConfStr),
        )
    } else if bgpChanged {
        diags.AddError(
            "Error updating BGP configuration",
            "BGP configuration fields cannot be updated",
        )
    } else { // assuming already exists 
		// Fetch the existing BGP Config
        bgpConfig, diags = fetchBGPConfig(client, projectID)
        diags.Append(diags...)
	}

    return bgpConfig, diags
}
