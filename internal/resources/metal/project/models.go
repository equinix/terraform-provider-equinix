package project

import (
	"context"
	"path"
	"time"

	"github.com/equinix/equinix-sdk-go/services/metalv1"

	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID              types.String                                     `tfsdk:"id"`
	Name            types.String                                     `tfsdk:"name"`
	Created         types.String                                     `tfsdk:"created"`
	Updated         types.String                                     `tfsdk:"updated"`
	BackendTransfer types.Bool                                       `tfsdk:"backend_transfer"`
	PaymentMethodID types.String                                     `tfsdk:"payment_method_id"`
	OrganizationID  types.String                                     `tfsdk:"organization_id"`
	BGPConfig       fwtypes.ListNestedObjectValueOf[BGPConfigModel] `tfsdk:"bgp_config"`
}

type BGPConfigModel struct {
	DeploymentType types.String `tfsdk:"deployment_type"`
	ASN            types.Int64  `tfsdk:"asn"`
	MD5            types.String `tfsdk:"md5"`
	Status         types.String `tfsdk:"status"`
	MaxPrefix      types.Int64  `tfsdk:"max_prefix"`
}

func (m *ResourceModel) parse(ctx context.Context, project *metalv1.Project, bgpConfig *metalv1.BgpConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	m.ID              = types.StringValue(project.GetId()) 
	m.Name            = types.StringValue(project.GetId())
	m.Created         = types.StringValue(project.GetCreatedAt().Format(time.RFC3339))
	m.Updated         = types.StringValue(project.GetUpdatedAt().Format(time.RFC3339))
	m.BackendTransfer = types.BoolValue(project.AdditionalProperties["backend_transfer_enabled"].(bool)) // No backend_transfer_enabled property in API spec

	m.OrganizationID  = types.StringValue(path.Base(project.Organization.AdditionalProperties["href"].(string)))

	if len(project.PaymentMethod.GetHref()) != 0 {
		m.PaymentMethodID = types.StringValue(path.Base(project.PaymentMethod.GetHref()))
	}

	// Handle BGP Config if present	
	m.BGPConfig = parseBGPConfig(ctx, bgpConfig)

	return diags
}

func parseBGPConfig(ctx context.Context, bgpConfig *metalv1.BgpConfig) fwtypes.ListNestedObjectValueOf[BGPConfigModel] {
	if bgpConfig != nil {
		bgpConfigResourceModel := make([]BGPConfigModel, 1)
		bgpConfigResourceModel[0] = BGPConfigModel{
			DeploymentType: types.StringValue(string(bgpConfig.GetDeploymentType())),
			ASN:            types.Int64Value(int64(bgpConfig.GetAsn())),
			MD5:            types.StringValue(bgpConfig.GetMd5()),
			Status:         types.StringValue(string(bgpConfig.GetStatus())),
			MaxPrefix:      types.Int64Value(int64(bgpConfig.GetMaxPrefix())),
		}
		return fwtypes.NewListNestedObjectValueOfValueSlice[BGPConfigModel](ctx, bgpConfigResourceModel)
	}
	return fwtypes.NewListNestedObjectValueOfNull[BGPConfigModel](ctx)
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