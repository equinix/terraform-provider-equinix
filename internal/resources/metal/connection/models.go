package connection

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
)

type ResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Facility         types.String `tfsdk:"facility"`
	Metro            types.String `tfsdk:"metro"`
	Redundancy       types.String `tfsdk:"redundancy"`
	ContactEmail     types.String `tfsdk:"contact_email"`
	Type             types.String `tfsdk:"type"`
	ProjectID        types.String `tfsdk:"project_id"`
	Speed            types.String `tfsdk:"speed"`
	Description      types.String `tfsdk:"description"`
	Mode             types.String `tfsdk:"mode"`
	Tags             types.List   `tfsdk:"tags"`  // List of strings
	Vlans            types.List   `tfsdk:"vlans"` // List of ints
	ServiceTokenType types.String `tfsdk:"service_token_type"`
	OrganizationID   types.String `tfsdk:"organization_id"`
	Status           types.String `tfsdk:"status"`
	Token            types.String `tfsdk:"token"`
	Ports            types.List   `tfsdk:"ports"`          // List of Port
	ServiceTokens    types.List   `tfsdk:"service_tokens"` // List of ServiceToken
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

func (m *ResourceModel) parse(ctx context.Context, conn *packngo.Connection) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(conn.ID)
	m.OrganizationID = types.StringValue(conn.Organization.ID)
	m.Name = types.StringValue(conn.Name)
	m.Facility = types.StringValue(conn.Facility.Code)

	// TODO(ocobles) we were using "StateFunc: converters.ToLowerIf" for "metro" field in the sdkv2
	// version of this resource. StateFunc doesn't exist in terraform and it requires implementation
	// of bespoke logic before storing state. To ensure backward compatibility we ignore lower/upper
	// case diff for now, but we may want to require input upper case
	if !strings.EqualFold(m.Metro.ValueString(), conn.Metro.Code) {
		m.Metro = types.StringValue(conn.Metro.Code)
	}

	// TODO(ocobles) API returns "" when description was not provided
	// To ensure backward compatibility we ignore null/empty diff
	if !m.Description.IsNull() || (m.Description.IsNull() && conn.Description != "") {
		m.Description = types.StringValue(conn.Description)
	}

	m.ContactEmail = types.StringValue(conn.ContactEmail)
	m.Status = types.StringValue(conn.Status)
	m.Redundancy = types.StringValue(string(conn.Redundancy))
	m.Token = types.StringValue(conn.Token)
	m.Type = types.StringValue(string(conn.Type))

	if conn.Mode != nil {
		m.Mode = types.StringValue(string(*conn.Mode))
	}

	if !m.Tags.IsNull() || (m.Tags.IsNull() && conn.Tags != nil && len(conn.Tags) > 0) {
		tags, diags := types.ListValueFrom(ctx, types.StringType, conn.Tags)
		if diags.HasError() {
			return diags
		}
		m.Tags = tags
	}

	// Parse Service Token Type
	tokenType := ""
	if len(conn.Tokens) > 0 {
		tokenType = string(conn.Tokens[0].ServiceTokenType)
	}
	m.ServiceTokenType = types.StringValue(tokenType)

	// Parse Speed
	if !(tokenType == "z_side" && conn.Type == packngo.ConnectionShared) {
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
		m.Speed = types.StringValue(speed)
	}

	// Parse Project ID
	// fix the project id get when it's added straight to the Connection API resource
	// https://github.com/packethost/packngo/issues/317
	if conn.Type == packngo.ConnectionShared {
		// Note: we were using conn.Ports[0].VirtualCircuits[0].Project.ID in the sdkv2 version but
		// it is empty and in framework that produces an unexpected new value.
		m.ProjectID = types.StringValue(path.Base(conn.Ports[0].VirtualCircuits[0].Project.URL))
	}

	// Parse Vlans
	diags = m.parseConnectionVlans(ctx, conn)
	if diags.HasError() {
		return diags
	}

	// Parse Ports
	diags = m.parseConnectionPorts(ctx, conn.Ports)
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
	m.ServiceTokens = serviceTokens

	return diags
}

func (m *ResourceModel) parseConnectionPorts(ctx context.Context, cps []packngo.ConnectionPort) diag.Diagnostics {
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

	ports, diags := types.ListValueFrom(ctx, PortsObjectType, ret)
	if diags.HasError() {
		return diags
	}
	m.Ports = ports
	return nil
}

func (m *ResourceModel) parseConnectionVlans(ctx context.Context, conn *packngo.Connection) diag.Diagnostics {
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
	vlans, diags := types.ListValueFrom(ctx, types.Int64Type, ret)
	if diags.HasError() {
		return diags
	}
	m.Vlans = vlans
	return nil
}
