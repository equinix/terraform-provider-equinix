package project

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func fetchBGPConfig(ctx context.Context, client *metalv1.APIClient, projectID string) (*metalv1.BgpConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	bgpConfig, _, err := client.BGPApi.FindBgpConfigByProject(ctx, projectID).Execute()
	if err != nil {
		friendlyErr := equinix_errors.FriendlyError(err)
		diags.AddError(
			"Error reading BGP configuration",
			"Could not read BGP configuration for project with ID "+projectID+": "+friendlyErr.Error(),
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

func handleBGPConfigChanges(ctx context.Context, client *metalv1.APIClient, plan, state *ResourceModel, projectID string) (*metalv1.BgpConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var bgpConfig *metalv1.BgpConfig

	if plan.BGPConfig.IsNull() && state.BGPConfig.IsNull() {
		return bgpConfig, nil
	}

	bgpAdded := !plan.BGPConfig.IsNull() && state.BGPConfig.IsNull()
	bgpChanged := !plan.BGPConfig.IsNull() && !state.BGPConfig.IsNull() && !plan.BGPConfig.Equal(state.BGPConfig)

	if bgpAdded || bgpChanged {
		// Create BGP Config
		bgpCreateRequest, err := expandBGPConfig(ctx, plan.BGPConfig)
		if err != nil {
			diags.AddError(
				"Error creating project",
				"Could not validate BGP Config: "+err.Error(),
			)
			return nil, diags
		}
		createResp, err := client.BGPApi.RequestBgpConfig(ctx, projectID).BgpConfigRequestInput(*bgpCreateRequest).Execute()
		if err != nil {
			err = equinix_errors.FriendlyErrorForMetalGo(err, createResp)
			diags.AddError(
				"Error creating BGP configuration",
				"Could not create BGP configuration for project: "+err.Error(),
			)
			return nil, diags
		}
		// Fetch the newly created BGP Config
		bgpConfig, diags = fetchBGPConfig(ctx, client, projectID)
		diags.Append(diags...)
	} else { // assuming already exists
		// Fetch the existing BGP Config
		bgpConfig, diags = fetchBGPConfig(ctx, client, projectID)
		diags.Append(diags...)
	}

	return bgpConfig, diags
}
