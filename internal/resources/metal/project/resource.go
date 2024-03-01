package project

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/slices"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_metal_project",
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

	// Prepare the data for API request
	createRequest := metalv1.ProjectCreateFromRootInput{
		Name: plan.Name.ValueString(),
	}

	// Include optional fields if they are set
	if !plan.OrganizationID.IsNull() {
		createRequest.OrganizationId = plan.OrganizationID.ValueStringPointer()
	}

	// API call to create the project
	project, createResp, err := client.ProjectsApi.CreateProject(ctx).ProjectCreateFromRootInput(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project: " + equinix_errors.FriendlyErrorForMetalGo(err, createResp).Error(),
		)
		return
	}

	// Handle BGP Config if present
	if !plan.BGPConfig.IsNull() {
		bgpCreateRequest, err := expandBGPConfig(ctx, plan.BGPConfig)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating project",
				"Could not validate BGP Config: " + err.Error(),
			)
			return
		}
	
		createResp, err = client.BGPApi.RequestBgpConfig(ctx, project.GetId()).BgpConfigRequestInput(*bgpCreateRequest).Execute()
		if err != nil {
			err = equinix_errors.FriendlyErrorForMetalGo(err, createResp)
			resp.Diagnostics.AddError(
				"Error creating BGP configuration",
				"Could not create BGP configuration for project: " + err.Error(),
			)
			return
		}
	}

	// Enable Backend Transfer if True
	if plan.BackendTransfer.ValueBool() {
		pur :=  metalv1.ProjectUpdateInput{
			BackendTransferEnabled: plan.BackendTransfer.ValueBoolPointer(),
		}
		_, updateResp, err := client.ProjectsApi.UpdateProject(ctx, project.GetId()).ProjectUpdateInput(pur).Execute()
		if err != nil {
			err = equinix_errors.FriendlyErrorForMetalGo(err, updateResp)
			resp.Diagnostics.AddError(
					"Error enabling Backend Transfer",
					"Could not enable Backend Transfer for project with ID " + project.GetId() + ": " + err.Error(),
			)
			return
		}
	}

	// Use API client to get the current state of the resource
	project, diags = fetchProject(ctx, client, project.GetId())
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Fetch BGP Config if needed
	var bgpConfig *metalv1.BgpConfig
	if !plan.BGPConfig.IsNull() {
		bgpConfig, diags = fetchBGPConfig(ctx, client, project.GetId())
        diags.Append(diags...)
		if diags.HasError(){
			return
		}
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, project, bgpConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func fetchProject(ctx context.Context, client *metalv1.APIClient, projectID string) (*metalv1.Project, diag.Diagnostics) {
    var diags diag.Diagnostics

	project, apiResp, err := client.ProjectsApi.FindProjectById(ctx, projectID).Execute()
    if err != nil {
		err = equinix_errors.FriendlyErrorForMetalGo(err, apiResp)

        // Check if the Project no longer exists
        if equinix_errors.IsNotFound(err) {
            diags.AddWarning(
                "Project not found",
                fmt.Sprintf("Project (%s) not found, removing from state", projectID),
            )
        } else {
            diags.AddError(
                "Error reading project",
                "Could not read project with ID " + projectID + ": " + err.Error(),
            )
        }
        return nil, diags
    }

    return project, diags
}

func fetchBGPConfig(ctx context.Context, client *metalv1.APIClient, projectID string) (*metalv1.BgpConfig, diag.Diagnostics) {
    var diags diag.Diagnostics

	bgpConfig, _, err := client.BGPApi.FindBgpConfigByProject(ctx, projectID).Execute()
    if err != nil {
		friendlyErr := equinix_errors.FriendlyError(err)
        diags.AddError(
            "Error reading BGP configuration",
            "Could not read BGP configuration for project with ID " + projectID + ": " + friendlyErr.Error(),
        )
        return nil, diags
    }

    return bgpConfig, diags
}

func expandBGPConfig(ctx context.Context, bgpConfig fwtypes.ListNestedObjectValueOf[BGPConfigModel]) (*metalv1.BgpConfigRequestInput, error) {
	bgpConfigModel, _ := bgpConfig.ToSlice(ctx)
	bgpDeploymentType, err := metalv1.NewBgpConfigRequestInputDeploymentTypeFromValue(bgpConfigModel[0].DeploymentType.ValueString())
	if err != nil {
		return nil, err
	}
	bgpCreateRequest := metalv1.BgpConfigRequestInput{
		DeploymentType: *bgpDeploymentType,
		Asn:            int32(bgpConfigModel[0].ASN.ValueInt64()),
	}
	if !bgpConfigModel[0].MD5.IsNull() {
		bgpCreateRequest.Md5 = bgpConfigModel[0].MD5.ValueStringPointer()
	}

	return &bgpCreateRequest, nil
}


// func resourceMetalProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	client := meta.(*config.Config).NewMetalClientForSDK(d)

// 	proj, resp, err := client.ProjectsApi.FindProjectById(ctx, d.Id()).Execute()
// 	if err != nil {
// 		err = equinix_errors.FriendlyErrorForMetalGo(err, resp)

// 		// If the project somehow already destroyed, mark as successfully gone.
// 		if equinix_errors.IsNotFound(err) {
// 			d.SetId("")

// 			return nil
// 		}

// 		return diag.FromErr(err)
// 	}

// 	d.SetId(proj.GetId())
// 	if len(proj.PaymentMethod.GetHref()) != 0 {
// 		d.Set("payment_method_id", path.Base(proj.PaymentMethod.GetHref()))
// 	}
// 	d.Set("name", proj.Name)
// 	d.Set("organization_id", path.Base(proj.Organization.AdditionalProperties["href"].(string))) // spec: organization has no href
// 	d.Set("created", proj.GetCreatedAt().Format(time.RFC3339))
// 	d.Set("updated", proj.GetUpdatedAt().Format(time.RFC3339))
// 	d.Set("backend_transfer", proj.AdditionalProperties["backend_transfer_enabled"].(bool)) // No backend_transfer_enabled property in API spec

// 	bgpConf, _, err := client.BGPApi.FindBgpConfigByProject(ctx, proj.GetId()).Execute()

// 	if (err == nil) && (bgpConf != nil) {
// 		// guard against an empty struct
// 		if bgpConf.GetId() != "" {
// 			err := d.Set("bgp_config", flattenBGPConfig(bgpConf))
// 			if err != nil {
// 				return diag.FromErr(equinix_errors.FriendlyErrorForMetalGo(err, resp))
// 			}
// 		}
// 	}
// 	return nil
// }


// func resourceMetalProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	client := meta.(*config.Config).NewMetalClientForSDK(d)
// 	updateRequest := metalv1.ProjectUpdateInput{}
// 	if d.HasChange("name") {
// 		pName := d.Get("name").(string)
// 		updateRequest.Name = &pName
// 	}
// 	if d.HasChange("payment_method_id") {
// 		pPayment := d.Get("payment_method_id").(string)
// 		updateRequest.PaymentMethodId = &pPayment
// 	}
// 	if d.HasChange("backend_transfer") {
// 		pBT := d.Get("backend_transfer").(bool)
// 		updateRequest.BackendTransferEnabled = &pBT
// 	}
// 	if d.HasChange("bgp_config") {
// 		o, n := d.GetChange("bgp_config")
// 		oldarr := o.([]interface{})
// 		newarr := n.([]interface{})
// 		if len(newarr) == 1 {
// 			bgpCreateRequest, err := expandBGPConfig(d)
// 			if err != nil {
// 				return diag.FromErr(err)
// 			}

// 			resp, err := client.BGPApi.RequestBgpConfig(ctx, d.Id()).BgpConfigRequestInput(*bgpCreateRequest).Execute()
// 			if err != nil {
// 				return diag.FromErr(equinix_errors.FriendlyErrorForMetalGo(err, resp))
// 			}
// 		} else {
// 			if len(oldarr) == 1 {
// 				m := oldarr[0].(map[string]interface{})

// 				bgpConfStr := fmt.Sprintf(
// 					"bgp_config {\n"+
// 						"  deployment_type = \"%s\"\n"+
// 						"  md5 = \"%s\"\n"+
// 						"  asn = %d\n"+
// 						"}", m["deployment_type"].(string), m["md5"].(string),
// 					m["asn"].(int))

// 				return diag.Errorf("BGP Config can not be removed from a project, please add back\n%s", bgpConfStr)
// 			}
// 		}
// 	} else {
// 		_, resp, err := client.ProjectsApi.UpdateProject(ctx, d.Id()).ProjectUpdateInput(updateRequest).Execute()
// 		if err != nil {
// 			return diag.FromErr(equinix_errors.FriendlyErrorForMetalGo(err, resp))
// 		}
// 	}

// 	return resourceMetalProjectRead(ctx, d, meta)
// }

// func resourceMetalProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	client := meta.(*config.Config).NewMetalClientForSDK(d)

// 	resp, err := client.ProjectsApi.DeleteProject(ctx, d.Id()).Execute()
// 	if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
// 		return diag.FromErr(equinix_errors.FriendlyErrorForMetalGo(err, resp))
// 	}

// 	d.SetId("")
// 	return nil
// }

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve the current state
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

    // API call to delete the project
	deleteResp, err := client.ProjectsApi.DeleteProject(ctx, id).Execute()
	if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(deleteResp, err) != nil {
		err = equinix_errors.FriendlyErrorForMetalGo(err, deleteResp)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Project %s", id),
			err.Error(),
		)
	}
}