package metal_connection

import (
	"context"
	"fmt"
    "strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-framework/diag"
)

type MetalConnectionResourceModel struct {
    ID               types.String `tfsdk:"id"`
    Name             types.String   `tfsdk:"name"`
    Facility         types.String   `tfsdk:"facility"`
    Metro            types.String   `tfsdk:"metro"`
    Redundancy       types.String   `tfsdk:"redundancy"`
    ContactEmail     types.String   `tfsdk:"contact_email"`
    Type             types.String   `tfsdk:"type"`
    ProjectID        types.String   `tfsdk:"project_id"`
    Speed            types.String   `tfsdk:"speed"`
    Description      types.String   `tfsdk:"description"`
    Mode             types.String   `tfsdk:"mode"`
    Tags             types.List     `tfsdk:"tags"` // List of strings
    Vlans            types.List     `tfsdk:"vlans"` // List of ints
    ServiceTokenType types.String   `tfsdk:"service_token_type"`
    OrganizationID   types.String   `tfsdk:"organization_id"`
    Status           types.String   `tfsdk:"status"`
    Token            types.String   `tfsdk:"token"`
    Ports            types.List     `tfsdk:"ports"` // List of Port
    ServiceTokens    types.List     `tfsdk:"service_tokens"` // List of ServiceToken
}

type Port struct {
    ID                types.String `tfsdk:"id"`
    Name              types.String `tfsdk:"name"`
    Role              types.String `tfsdk:"role"`
    Speed             types.Int64  `tfsdk:"speed"`
    Status            types.String `tfsdk:"status"`
    LinkStatus        types.String `tfsdk:"link_status"`
    VirtualCircuitIDs types.List   `tfsdk:"virtual_circuit_ids"` // List of String
}

type ServiceToken struct {
    ID              types.String `tfsdk:"id"`
    MaxAllowedSpeed types.String `tfsdk:"max_allowed_speed"`
    Role            types.String `tfsdk:"role"`
    State           types.String `tfsdk:"state"`
    Type            types.String `tfsdk:"type"`
}

func (rm *MetalConnectionResourceModel) parse(ctx context.Context, conn *packngo.Connection) diag.Diagnostics {
    var diags diag.Diagnostics

    rm.OrganizationID = types.StringValue(conn.Organization.ID)
    rm.Name = types.StringValue(conn.Name)
    rm.Facility = types.StringValue(conn.Facility.Code)
    rm.Metro = types.StringValue(conn.Metro.Code)
    rm.Description = types.StringValue(conn.Description)
    rm.ContactEmail = types.StringValue(conn.ContactEmail)
    rm.Status = types.StringValue(conn.Status)
    rm.Redundancy = types.StringValue(string(conn.Redundancy))
    rm.Token = types.StringValue(conn.Token)
    rm.Type = types.StringValue(string(conn.Type))
    rm.Mode = types.StringValue(string(*conn.Mode))
    tags, diags := types.ListValueFrom(ctx, types.StringType, conn.Tags)
    if diags.HasError() {
        return diags
    }
    rm.Tags = tags

    // Parse Speed
    speed := "0"
    var err error
    if conn.Speed > 0 {
		speed, err = speedUintToStr(conn.Speed)
		if err != nil {
			diags.AddError(
                fmt.Sprintf("Failed to convert Speed (%d) to string", conn.Speed),
                err.Error(),
            )
            return diags
		}
	}
    rm.Speed = types.StringValue(speed)

    // Parse Project ID
	// fix the project id get when it's added straight to the Connection API resource
	// https://github.com/packethost/packngo/issues/317
	if conn.Type == packngo.ConnectionShared {
		rm.ProjectID = types.StringValue(conn.Ports[0].VirtualCircuits[0].Project.ID)
	}

    // Parse Service Token Type
    tokenType := ""
    if len(conn.Tokens) > 0 {
        tokenType = string(conn.Tokens[0].ServiceTokenType)
    }
    rm.ServiceTokenType = types.StringValue(tokenType)

    // Parse Vlans
    diags = rm.parseConnectionVlans(ctx, conn)
    if diags.HasError() {
		return diags
	}

    // Parse Ports
    diags = rm.parseConnectionPorts(ctx, conn.Ports)
    if diags.HasError() {
		return diags
	}

    // Parse ServiceTokens
    connServiceTokens := make([]ServiceToken, len(conn.Tokens))
    for i, token := range conn.Tokens {
        speed, err := speedUintToStr(token.MaxAllowedSpeed)
        if err != nil {
            var diags diag.Diagnostics
			diags.AddError(
                fmt.Sprintf("Failed to convert token MaxAllowedSpeed (%d) to string", token.MaxAllowedSpeed),
                err.Error(),
            )
            return diags
		}
        connServiceTokens[i] = ServiceToken{
            ID:              types.StringValue(token.ID),
            MaxAllowedSpeed: types.StringValue(speed),
            Role:            types.StringValue(string(token.Role)),
            State:           types.StringValue(token.State),
            Type:            types.StringValue(string(token.ServiceTokenType)),
        }
    }
    serviceTokens, diags := types.ListValueFrom(ctx, ServiceTokensObjectType, connServiceTokens) 
    if diags.HasError() {
		return diags
	}
    rm.ServiceTokens = serviceTokens

    return diags
}

func (rm *MetalConnectionResourceModel) parseConnectionPorts(ctx context.Context, cps []packngo.ConnectionPort) diag.Diagnostics {
    ret := make([]Port, len(cps))
    order := map[packngo.ConnectionPortRole]int{
        packngo.ConnectionPortPrimary:   0,
        packngo.ConnectionPortSecondary: 1,
    }

    for _, p := range cps {
        // Parse VirtualCircuits
        portVcIDs := make([]string, len(p.VirtualCircuits))
        for i, vc := range p.VirtualCircuits {
            portVcIDs[i] = vc.ID
        }
        vcIDs, diags := types.ListValueFrom(ctx, types.StringType, portVcIDs)
        if diags.HasError() {
            return diags
        }
        connPort := Port{
            ID:                types.StringValue(p.ID),
            Name:              types.StringValue(p.Name),
            Role:              types.StringValue(string(p.Role)),
            Speed:             types.Int64Value(int64(p.Speed)),
            Status:            types.StringValue(p.Status),
            LinkStatus:        types.StringValue(p.LinkStatus),
            VirtualCircuitIDs: vcIDs,
        }

        // Sort the ports by role, asserting the API always returns primary for len of 1 responses
        ret[order[p.Role]] = connPort
    }

    ports, diags := types.ListValueFrom(ctx, types.StringType, ret)
    if diags.HasError() {
        return diags
    }
    rm.Ports = ports
    return nil
}

func (rm *MetalConnectionResourceModel) parseConnectionVlans(ctx context.Context, conn *packngo.Connection) diag.Diagnostics {
	var ret []int

	if conn.Type == packngo.ConnectionShared {
		order := map[packngo.ConnectionPortRole]int{
			packngo.ConnectionPortPrimary:   0,
			packngo.ConnectionPortSecondary: 1,
		}

		rawVlans := make([]int, len(conn.Ports))
		for _, p := range conn.Ports {
			rawVlans[order[p.Role]] = p.VirtualCircuits[0].VNID
		}

		for _, v := range rawVlans {
			if v > 0 {
				ret = append(ret, v)
			}
		}
	}
	vlans, diags := types.ListValueFrom(ctx, types.StringType, ret)
    if diags.HasError() {
        return diags
    }
    rm.Vlans = vlans
    return nil
}

func speedStrToUint(speed string) (uint64, error) {
	allowedStrings := []string{}
	for _, allowedSpeed := range allowedSpeeds {
		if allowedSpeed.Str == speed {
			return allowedSpeed.Int, nil
		}
		allowedStrings = append(allowedStrings, allowedSpeed.Str)
	}
	return 0, fmt.Errorf("invalid speed string: %s. Allowed strings: %s", speed, strings.Join(allowedStrings, ", "))
}

func speedUintToStr(speed uint64) (string, error) {
	allowedUints := []uint64{}
	for _, allowedSpeed := range allowedSpeeds {
		if speed == allowedSpeed.Int {
			return allowedSpeed.Str, nil
		}
		allowedUints = append(allowedUints, allowedSpeed.Int)
	}
	return "", fmt.Errorf("%d is not allowed speed value. Allowed values: %v", speed, allowedUints)
}