package metal_virtual_circuit

import (
    "context"
    "regexp"
    "strconv"
    "log"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

type MetalVirtualCircuitResourceModel struct {
    ID            types.String   `tfsdk:"id"`
    ConnectionID  types.String   `tfsdk:"connection_id"`
    ProjectID     types.String   `tfsdk:"project_id"`
    PortID        types.String   `tfsdk:"port_id"`
    Name          types.String   `tfsdk:"name"`
    Description   types.String   `tfsdk:"description"`
    Speed         types.String   `tfsdk:"speed"`
    Tags          types.List     `tfsdk:"tags"`
    NniVLAN       types.Int64    `tfsdk:"nni_vlan"`
    VlanID        types.String   `tfsdk:"vlan_id"`
    VrfID         types.String   `tfsdk:"vrf_id"`
    PeerASN       types.Int64    `tfsdk:"peer_asn"`
    Subnet        types.String   `tfsdk:"subnet"`
    MetalIP       types.String   `tfsdk:"metal_ip"`
    CustomerIP    types.String   `tfsdk:"customer_ip"`
    MD5           types.String   `tfsdk:"md5"`
    Vnid          types.Int64    `tfsdk:"vnid"`
    NniVnid       types.Int64    `tfsdk:"nni_vnid"`
    Status        types.String   `tfsdk:"status"`
    Timeouts      timeouts.Value `tfsdk:"timeouts"`
}

func (rm *MetalVirtualCircuitResourceModel) parse(ctx context.Context, vc *packngo.VirtualCircuit) diag.Diagnostics {
    var diags diag.Diagnostics

    rm.ID          = types.StringValue(vc.ID)
    rm.ProjectID   = types.StringValue(vc.Project.ID)
    rm.Name        = types.StringValue(vc.Name)
    rm.Description = types.StringValue(vc.Description)
    rm.Speed       = types.StringValue(strconv.Itoa(vc.Speed))
    rm.Status      = types.StringValue(string(vc.Status))
    rm.NniVLAN     = types.Int64Value(int64(vc.NniVLAN)) 
    rm.Vnid        = types.Int64Value(int64(vc.VNID))
    rm.NniVnid     =  types.Int64Value(int64(vc.NniVNID))
    rm.PeerASN     = types.Int64Value(int64(vc.PeerASN))
    rm.Subnet      = types.StringValue(vc.Subnet)
    rm.MetalIP     = types.StringValue(vc.MetalIP)
    rm.CustomerIP  = types.StringValue(vc.CustomerIP)
    rm.MD5         = types.StringValue(vc.MD5)

    // TODO: use API field from VC responses when available The regexp is
	// optimistic, not guaranteed. This affects resource imports. "port" is not
	// in the Includes above to assure the Href needed below.
	connectionID := types.StringNull() // vc.Connection.ID is not available yet
	portID := ""       // vc.Port.ID would be available with ?include=port
	connectionRe := regexp.MustCompile("/connections/([0-9a-z-]+)/ports/([0-9a-z-]+)")
	matches := connectionRe.FindStringSubmatch(vc.Port.Href.Href)
	if len(matches) == 3 {
		connectionID = types.StringValue(matches[1])
		portID = matches[2]
	} else {
		log.Printf("[DEBUG] Could not parse connection and port ID from port href %s", vc.Port.Href.Href)
	}
    rm.ConnectionID = connectionID
    rm.PortID = types.StringValue(portID)

    // VRF ID
    if vc.VRF != nil {
        rm.VrfID = types.StringValue(vc.VRF.ID)
    }

    // VLAN ID
    if vc.VirtualNetwork != nil {
        rm.VlanID = types.StringValue(vc.VirtualNetwork.ID)
    }

    // Handling tags as a list
    tags, diags := types.ListValueFrom(ctx, types.StringType, vc.Tags)
    if diags.HasError() {
        return diags
    }
    rm.Tags = tags

    return diags
}
