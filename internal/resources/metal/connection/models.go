package connection

import (
	"context"
	"fmt"
	"path"
	"strings"

	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/packethost/packngo"
)

type ResourceModel struct {
	ID               types.String                                       `tfsdk:"id"`
	Name             types.String                                       `tfsdk:"name"`
	Facility         types.String                                       `tfsdk:"facility"`
	Metro            types.String                                       `tfsdk:"metro"`
	Redundancy       types.String                                       `tfsdk:"redundancy"`
	ContactEmail     types.String                                       `tfsdk:"contact_email"`
	Type             types.String                                       `tfsdk:"type"`
	ProjectID        types.String                                       `tfsdk:"project_id"`
	Speed            types.String                                       `tfsdk:"speed"`
	Description      types.String                                       `tfsdk:"description"`
	Mode             types.String                                       `tfsdk:"mode"`
	Tags             types.List                                         `tfsdk:"tags"`  // List of strings
	Vlans            types.List                                         `tfsdk:"vlans"` // List of ints
	ServiceTokenType types.String                                       `tfsdk:"service_token_type"`
	OrganizationID   types.String                                       `tfsdk:"organization_id"`
	Status           types.String                                       `tfsdk:"status"`
	Token            types.String                                       `tfsdk:"token"`
	Ports            fwtypes.ListNestedObjectValueOf[PortModel]         `tfsdk:"ports"`          // List of Port
	ServiceTokens    fwtypes.ListNestedObjectValueOf[ServiceTokenModel] `tfsdk:"service_tokens"` // List of ServiceToken
}

type DataSourceModel struct {
	ID               types.String                                       `tfsdk:"id"`
	ConnectionID     types.String                                       `tfsdk:"connection_id"`
	Name             types.String                                       `tfsdk:"name"`
	Facility         types.String                                       `tfsdk:"facility"`
	Metro            types.String                                       `tfsdk:"metro"`
	Redundancy       types.String                                       `tfsdk:"redundancy"`
	ContactEmail     types.String                                       `tfsdk:"contact_email"`
	Type             types.String                                       `tfsdk:"type"`
	ProjectID        types.String                                       `tfsdk:"project_id"`
	Speed            types.String                                       `tfsdk:"speed"`
	Description      types.String                                       `tfsdk:"description"`
	Mode             types.String                                       `tfsdk:"mode"`
	Tags             types.List                                         `tfsdk:"tags"`  // List of strings
	Vlans            types.List                                         `tfsdk:"vlans"` // List of ints
	ServiceTokenType types.String                                       `tfsdk:"service_token_type"`
	OrganizationID   types.String                                       `tfsdk:"organization_id"`
	Status           types.String                                       `tfsdk:"status"`
	Token            types.String                                       `tfsdk:"token"`
	Ports            fwtypes.ListNestedObjectValueOf[PortModel]         `tfsdk:"ports"`          // List of Port
	ServiceTokens    fwtypes.ListNestedObjectValueOf[ServiceTokenModel] `tfsdk:"service_tokens"` // List of ServiceToken
}

type PortModel struct {
	ID                types.String                      `tfsdk:"id"`
	Name              types.String                      `tfsdk:"name"`
	Role              types.String                      `tfsdk:"role"`
	Speed             types.Int64                       `tfsdk:"speed"`
	Status            types.String                      `tfsdk:"status"`
	LinkStatus        types.String                      `tfsdk:"link_status"`
	VirtualCircuitIDs fwtypes.ListValueOf[types.String] `tfsdk:"virtual_circuit_ids"` // List of String
}

type ServiceTokenModel struct {
	ID              types.String `tfsdk:"id"`
	ExpiresAt       types.String `tfsdk:"expires_at"`
	MaxAllowedSpeed types.String `tfsdk:"max_allowed_speed"`
	Role            types.String `tfsdk:"role"`
	State           types.String `tfsdk:"state"`
	Type            types.String `tfsdk:"type"`
}

func (m *DataSourceModel) parse(ctx context.Context, conn *packngo.Connection) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ConnectionID = types.StringValue(conn.ID)
	parseConnection(ctx, conn,
		&m.ID, &m.OrganizationID, &m.Name, &m.Facility, &m.Metro,
		&m.Description, &m.ContactEmail, &m.Status, &m.Redundancy,
		&m.Token, &m.Type, &m.Mode, &m.ServiceTokenType, &m.Speed,
		&m.ProjectID, &m.Tags, &m.Vlans, &m.Ports, &m.ServiceTokens,
	)
	return diags
}

func (m *ResourceModel) parse(ctx context.Context, conn *packngo.Connection) diag.Diagnostics {
	var diags diag.Diagnostics

	parseConnection(ctx, conn,
		&m.ID, &m.OrganizationID, &m.Name, &m.Facility, &m.Metro,
		&m.Description, &m.ContactEmail, &m.Status, &m.Redundancy,
		&m.Token, &m.Type, &m.Mode, &m.ServiceTokenType, &m.Speed,
		&m.ProjectID, &m.Tags, &m.Vlans, &m.Ports, &m.ServiceTokens,
	)
	return diags
}

func parseConnection(
	ctx context.Context,
	conn *packngo.Connection,
	id, orgID, name, facility, metro, description, contactEmail, status, redundancy,
	token, typ, mode, serviceTokenType, speed, projectID *basetypes.StringValue,
	tags, vlans *basetypes.ListValue,
	ports *fwtypes.ListNestedObjectValueOf[PortModel],
	serviceTokens *fwtypes.ListNestedObjectValueOf[ServiceTokenModel],
) diag.Diagnostics {
	var diags diag.Diagnostics

	*id = types.StringValue(conn.ID)
	*orgID = types.StringValue(conn.Organization.ID)
	*name = types.StringValue(conn.Name)
	*facility = types.StringValue(conn.Facility.Code)

	// TODO(ocobles) we were using "StateFunc: converters.ToLowerIf" for "metro" field in the sdkv2
	// version of this resource. StateFunc doesn't exist in terraform and it requires implementation
	// of bespoke logic before storing state. To ensure backward compatibility we ignore lower/upper
	// case diff for now, but we may want to require input upper case
	if !strings.EqualFold(metro.ValueString(), conn.Metro.Code) {
		*metro = types.StringValue(conn.Metro.Code)
	}

	// TODO(ocobles) API returns "" when description was not provided
	// To ensure backward compatibility we ignore null/empty diff
	if !description.IsNull() || (description.IsNull() && conn.Description != "") {
		*description = types.StringValue(conn.Description)
	}

	*contactEmail = types.StringValue(conn.ContactEmail)
	*status = types.StringValue(conn.Status)
	*redundancy = types.StringValue(string(conn.Redundancy))
	*token = types.StringValue(conn.Token)
	*typ = types.StringValue(string(conn.Type))

	*mode = types.StringValue(string(packngo.ConnectionModeStandard))
	if conn.Mode != nil {
		*mode = types.StringValue(string(*conn.Mode))
	}

	if !tags.IsNull() || (tags.IsNull() && conn.Tags != nil && len(conn.Tags) > 0) {
		connTags, diags := types.ListValueFrom(ctx, types.StringType, conn.Tags)
		if diags.HasError() {
			return diags
		}
		*tags = connTags
	}

	// Parse Service Token Type
	if len(conn.Tokens) > 0 {
		*serviceTokenType = types.StringValue(string(conn.Tokens[0].ServiceTokenType))
	}

	// Parse Speed
	connSpeed := "0"
	var err error
	if conn.Speed > 0 {
		connSpeed, err = speedUintToStr(conn.Speed)
		if err != nil {
			diags.AddError(
				fmt.Sprintf("Failed to convert Speed (%d) to string", conn.Speed),
				err.Error(),
			)
			return diags
		}
	}
	*speed = types.StringValue(connSpeed)

	// Parse Project ID
	// fix the project id get when it's added straight to the Connection API resource
	// https://github.com/packethost/packngo/issues/317
	if conn.Type == packngo.ConnectionShared {
		// Note: we were using conn.Ports[0].VirtualCircuits[0].Project.ID in the sdkv2 version but
		// it is empty and in framework that produces an unexpected new value.
		*projectID = types.StringValue(path.Base(conn.Ports[0].VirtualCircuits[0].Project.URL))
	}

	// Parse Vlans
	connVlans, diags := parseConnectionVlans(ctx, conn)
	if diags.HasError() {
		return diags
	}
	if !connVlans.IsNull() {
		*vlans = *connVlans
	}

	connPorts, diags := parseConnectionPorts(ctx, conn.Ports)
	if diags.HasError() {
		return diags
	}
	*ports = connPorts

	// Parse ServiceTokens
	connServiceTokens, diags := parseConnectionServiceTokens(ctx, conn.Tokens)
	if diags.HasError() {
		return diags
	}
	*serviceTokens = connServiceTokens

	return diags
}

func parseConnectionServiceTokens(ctx context.Context, fst []packngo.FabricServiceToken) (fwtypes.ListNestedObjectValueOf[ServiceTokenModel], diag.Diagnostics) {
	connServiceTokens := make([]ServiceTokenModel, len(fst))
	for i, token := range fst {
		speed, err := speedUintToStr(token.MaxAllowedSpeed)
		if err != nil {
			var diags diag.Diagnostics
			diags.AddError(
				fmt.Sprintf("Failed to convert token MaxAllowedSpeed (%d) to string", token.MaxAllowedSpeed),
				err.Error(),
			)
			return fwtypes.NewListNestedObjectValueOfNull[ServiceTokenModel](ctx), diags
		}
		connServiceTokens[i] = ServiceTokenModel{
			ID:              types.StringValue(token.ID),
			MaxAllowedSpeed: types.StringValue(speed),
			Role:            types.StringValue(string(token.Role)),
			State:           types.StringValue(token.State),
			Type:            types.StringValue(string(token.ServiceTokenType)),
		}
		if token.ExpiresAt != nil {
			connServiceTokens[i].ExpiresAt = types.StringValue(token.ExpiresAt.String())
		}
	}

	return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, connServiceTokens), nil
}

func parseConnectionPorts(ctx context.Context, cps []packngo.ConnectionPort) (fwtypes.ListNestedObjectValueOf[PortModel], diag.Diagnostics) {
	ret := make([]PortModel, len(cps))
	order := map[packngo.ConnectionPortRole]int{
		packngo.ConnectionPortPrimary:   0,
		packngo.ConnectionPortSecondary: 1,
	}

	for _, p := range cps {
		// Parse VirtualCircuits
		portVcIDs := make([]attr.Value, len(p.VirtualCircuits))
		for i, vc := range p.VirtualCircuits {
			portVcIDs[i] = types.StringValue(vc.ID)
		}
		vcIDs, diags := fwtypes.NewListValueOf[types.String](ctx, portVcIDs)
		if diags.HasError() {
			return fwtypes.NewListNestedObjectValueOfNull[PortModel](ctx), diags
		}
		connPort := PortModel{
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

	return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, ret), nil
}

func parseConnectionVlans(ctx context.Context, conn *packngo.Connection) (*basetypes.ListValue, diag.Diagnostics) {
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
		return nil, diags
	}
	return &vlans, nil
}
