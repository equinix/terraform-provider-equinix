package metal_vrf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-framework/diag"
)

type MetalVRFResourceModel struct {
    ID          types.String `tfsdk:"id"`
    Name        types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Metro       types.String `tfsdk:"metro"`
    LocalASN    types.Int64  `tfsdk:"local_asn"`
    IPRanges    types.List   `tfsdk:"ip_ranges"`
    ProjectID   types.String `tfsdk:"project_id"`
}

func (rm *MetalVRFResourceModel) parse(ctx context.Context, vrf *packngo.VRF) diag.Diagnostics {
    var diags diag.Diagnostics

    rm.ID = types.StringValue(vrf.ID)
    rm.Name = types.StringValue(vrf.Name)
    rm.Description = types.StringValue(vrf.Description)
    rm.Metro = types.StringValue(vrf.Metro.Code)
    rm.LocalASN = types.Int64Value(int64(vrf.LocalASN))
    
    // Converting the IPRanges slice to a Terraform types.List
    ipRanges, diags := types.ListValueFrom(ctx, types.StringType, vrf.IPRanges)
    if diags.HasError() {
        return diags
    }
    rm.IPRanges = ipRanges

    rm.ProjectID = types.StringValue(vrf.Project.ID)

    return diags
}
