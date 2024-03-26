package connection

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ResourceModel struct {
	ID                types.String                                       `tfsdk:"id"`
	Name              types.String                                       `tfsdk:"name"`
	Facility          types.String                                       `tfsdk:"facility"`
	Metro             types.String                                       `tfsdk:"metro"`
	Redundancy        types.String                                       `tfsdk:"redundancy"`
	ContactEmail      types.String                                       `tfsdk:"contact_email"`
	Type              types.String                                       `tfsdk:"type"`
	ProjectID         types.String                                       `tfsdk:"project_id"`
	AuthorizationCode types.String                                       `tfsdk:"authorization_code"`
	Speed             types.String                                       `tfsdk:"speed"`
	Description       types.String                                       `tfsdk:"description"`
	Mode              types.String                                       `tfsdk:"mode"`
	Tags              types.List                                         `tfsdk:"tags"`  // List of strings
	Vlans             types.List                                         `tfsdk:"vlans"` // List of ints
	Vrfs              types.List                                         `tfsdk:"vrfs"`  // List of strings
	ServiceTokenType  types.String                                       `tfsdk:"service_token_type"`
	OrganizationID    types.String                                       `tfsdk:"organization_id"`
	Status            types.String                                       `tfsdk:"status"`
	Token             types.String                                       `tfsdk:"token"`
	Ports             fwtypes.ListNestedObjectValueOf[PortModel]         `tfsdk:"ports"`          // List of Port
	ServiceTokens     fwtypes.ListNestedObjectValueOf[ServiceTokenModel] `tfsdk:"service_tokens"` // List of ServiceToken
}

type DataSourceModel struct {
	ID                types.String                                       `tfsdk:"id"`
	ConnectionID      types.String                                       `tfsdk:"connection_id"`
	Name              types.String                                       `tfsdk:"name"`
	Facility          types.String                                       `tfsdk:"facility"`
	Metro             types.String                                       `tfsdk:"metro"`
	Redundancy        types.String                                       `tfsdk:"redundancy"`
	ContactEmail      types.String                                       `tfsdk:"contact_email"`
	Type              types.String                                       `tfsdk:"type"`
	ProjectID         types.String                                       `tfsdk:"project_id"`
	AuthorizationCode types.String                                       `tfsdk:"authorization_code"`
	Speed             types.String                                       `tfsdk:"speed"`
	Description       types.String                                       `tfsdk:"description"`
	Mode              types.String                                       `tfsdk:"mode"`
	Tags              types.List                                         `tfsdk:"tags"`  // List of strings
	Vlans             types.List                                         `tfsdk:"vlans"` // List of ints
	Vrfs              types.List                                         `tfsdk:"vrfs"`  // List of strings
	ServiceTokenType  types.String                                       `tfsdk:"service_token_type"`
	OrganizationID    types.String                                       `tfsdk:"organization_id"`
	Status            types.String                                       `tfsdk:"status"`
	Token             types.String                                       `tfsdk:"token"`
	Ports             fwtypes.ListNestedObjectValueOf[PortModel]         `tfsdk:"ports"`          // List of Port
	ServiceTokens     fwtypes.ListNestedObjectValueOf[ServiceTokenModel] `tfsdk:"service_tokens"` // List of ServiceToken
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

func (m *DataSourceModel) parse(ctx context.Context, conn *metalv1.Interconnection) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ConnectionID = types.StringPointerValue(conn.Id)

	parseConnection(ctx, conn,
		&m.ID, &m.OrganizationID, &m.Name, &m.Facility, &m.Metro,
		&m.Description, &m.ContactEmail, &m.Status, &m.Redundancy,
		&m.Token, &m.Type, &m.Mode, &m.ServiceTokenType, &m.Speed,
		&m.ProjectID, &m.AuthorizationCode, &m.Vlans, &m.Vrfs, &m.Ports, &m.ServiceTokens,
	)

	connTags, diags := types.ListValueFrom(ctx, types.StringType, conn.Tags)
	if diags.HasError() {
		return diags
	}
	m.Tags = connTags

	return diags
}

func (m *ResourceModel) parse(ctx context.Context, conn *metalv1.Interconnection) diag.Diagnostics {
	var diags diag.Diagnostics

	parseConnection(ctx, conn,
		&m.ID, &m.OrganizationID, &m.Name, &m.Facility, &m.Metro,
		&m.Description, &m.ContactEmail, &m.Status, &m.Redundancy,
		&m.Token, &m.Type, &m.Mode, &m.ServiceTokenType, &m.Speed,
		&m.ProjectID, &m.AuthorizationCode, &m.Vlans, &m.Vrfs, &m.Ports, &m.ServiceTokens,
	)

	connTags, diags := types.ListValueFrom(ctx, types.StringType, conn.Tags)
	if diags.HasError() {
		return diags
	}
	// TODO(ocobles) workaround to keep compatibility with older releases using SDKv2
	if m.Tags.IsNull() && len(conn.Tags) == 0 {
		m.Tags = types.ListNull(types.StringType)
	} else {
		m.Tags = connTags
	}

	return diags
}

// abstractVirtualCircuit represents either a metalv1.VrfVirtualCircuit or a
// metalv1.VlanVirtualCircuit
type abstractVirtualCircuit interface {
	GetId() string
	GetProject() metalv1.Project
}

func parseConnection(
	ctx context.Context,
	conn *metalv1.Interconnection,
	id, orgID, name, facility, metro, description, contactEmail, status, redundancy,
	token, typ, mode, serviceTokenType, speed, projectID, authorizationCode *basetypes.StringValue,
	vlans *basetypes.ListValue,
	vrfs *basetypes.ListValue,
	ports *fwtypes.ListNestedObjectValueOf[PortModel],
	serviceTokens *fwtypes.ListNestedObjectValueOf[ServiceTokenModel],
) diag.Diagnostics {
	var diags diag.Diagnostics

	*id = types.StringValue(conn.GetId())
	*orgID = types.StringPointerValue(conn.GetOrganization().Id)
	*name = types.StringValue(conn.GetName())
	*facility = types.StringValue(conn.Facility.GetCode())
	*description = types.StringValue(conn.GetDescription())
	*contactEmail = types.StringValue(conn.GetContactEmail())
	*status = types.StringValue(conn.GetStatus())
	*redundancy = types.StringValue(string(conn.GetRedundancy()))
	*token = types.StringValue(conn.GetToken())
	*typ = types.StringValue(string(conn.GetType()))
	*authorizationCode = types.StringValue(conn.GetAuthorizationCode())

	// TODO(ocobles) we were using "StateFunc: converters.ToLowerIf" for "metro" field in the sdkv2
	// version of this resource. StateFunc doesn't exist in terraform and it requires implementation
	// of bespoke logic before storing state. To ensure backward compatibility we ignore lower/upper
	// case diff for now, but we may want to require input upper case
	if !strings.EqualFold(metro.ValueString(), *conn.Metro.Code) {
		*metro = types.StringPointerValue(conn.GetMetro().Code)
	}

	*mode = types.StringValue(string(metalv1.INTERCONNECTIONMODE_STANDARD))
	if conn.HasMode() {
		*mode = types.StringValue(string(conn.GetMode()))
	}

	// Parse Service Token Type
	if len(conn.ServiceTokens) != 0 {
		*serviceTokenType = types.StringValue(string(conn.ServiceTokens[0].GetServiceTokenType()))
	}

	// Parse Speed
	connSpeed := "0"
	var err error
	if conn.GetSpeed() > 0 {
		connSpeed, err = speedIntToStr(conn.GetSpeed())
		if err != nil {
			diags.AddError(
				fmt.Sprintf("Failed to convert Speed (%d) to string", conn.Speed),
				err.Error(),
			)
			return diags
		}
	}
	*speed = types.StringValue(connSpeed)

	if conn.GetType() == metalv1.INTERCONNECTIONTYPE_SHARED {
		// Note: we were using conn.Ports[0].VirtualCircuits[0].Project.ID in the sdkv2 version but
		// it is empty and in framework that produces an unexpected new value.

		if len(conn.Ports) == 0 {
			diags.AddError(
				"Failed to get ports",
				"expected connection to have at least one associated port",
			)
			return diags
		}

		if len(conn.Ports[0].VirtualCircuits) == 0 {
			diags.AddError(
				"Failed to get port 0 virtual circuit",
				"expected connection port to have at least one associated virtual circuit",
			)
			return diags
		}

		vc := conn.Ports[0].VirtualCircuits[0].GetActualInstance().(abstractVirtualCircuit)
		project := vc.GetProject()

		*projectID = types.StringValue(project.GetId())

		connVlans, diags := parseConnectionVlans(ctx, conn)
		if diags.HasError() {
			return diags
		}
		if !connVlans.IsNull() && len(connVlans.Elements()) != 0 {
			*vlans = *connVlans
		}

		connVrfs, diags := parseConnectionVrfs(ctx, conn)
		if diags.HasError() {
			return diags
		}
		if !connVrfs.IsNull() && len(connVrfs.Elements()) != 0 {
			*vrfs = *connVrfs
		}
	}

	connPorts, diags := parseConnectionPorts(ctx, conn.Ports)
	if diags.HasError() {
		return diags
	}
	*ports = connPorts

	// Parse ServiceTokens
	connServiceTokens, diags := parseConnectionServiceTokens(ctx, conn.ServiceTokens)
	if diags.HasError() {
		return diags
	}
	*serviceTokens = connServiceTokens

	return diags
}

func parseConnectionServiceTokens(ctx context.Context, fst []metalv1.FabricServiceToken) (fwtypes.ListNestedObjectValueOf[ServiceTokenModel], diag.Diagnostics) {
	connServiceTokens := make([]ServiceTokenModel, len(fst))
	for i, token := range fst {
		speed, err := speedIntToStr(token.GetMaxAllowedSpeed())
		if err != nil {
			var diags diag.Diagnostics
			diags.AddError(
				fmt.Sprintf("Failed to convert token MaxAllowedSpeed (%d) to string", token.MaxAllowedSpeed),
				err.Error(),
			)
			return fwtypes.NewListNestedObjectValueOfNull[ServiceTokenModel](ctx), diags
		}

		connServiceTokens[i] = ServiceTokenModel{
			ID:              types.StringValue(token.GetId()),
			MaxAllowedSpeed: types.StringValue(speed),
			Role:            types.StringValue(string(token.GetRole())),
			State:           types.StringValue(string(token.GetState())),
			Type:            types.StringValue(string(token.GetServiceTokenType())),
		}
		if token.ExpiresAt != nil {
			connServiceTokens[i].ExpiresAt = types.StringValue(token.ExpiresAt.String())
		}
	}

	return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, connServiceTokens), nil
}

func parseConnectionPorts(ctx context.Context, cps []metalv1.InterconnectionPort) (fwtypes.ListNestedObjectValueOf[PortModel], diag.Diagnostics) {
	ret := make([]PortModel, len(cps))
	order := map[metalv1.InterconnectionPortRole]int{
		metalv1.INTERCONNECTIONPORTROLE_PRIMARY:   0,
		metalv1.INTERCONNECTIONPORTROLE_SECONDARY: 1,
	}

	for _, p := range cps {
		// Parse VirtualCircuits
		portVcIDs := make([]attr.Value, len(p.VirtualCircuits))
		for i, vc := range p.VirtualCircuits {
			vc := vc.GetActualInstance().(abstractVirtualCircuit)
			portVcIDs[i] = types.StringValue(vc.GetId())
		}
		vcIDs, diags := fwtypes.NewListValueOf[types.String](ctx, portVcIDs)
		if diags.HasError() {
			return fwtypes.NewListNestedObjectValueOfNull[PortModel](ctx), diags
		}
		connPort := PortModel{
			ID:                types.StringValue(p.GetId()),
			Name:              types.StringValue(p.GetName()),
			Role:              types.StringValue(string(p.GetRole())),
			Speed:             types.Int64Value(p.GetSpeed()),
			Status:            types.StringValue(string(p.GetStatus())),
			LinkStatus:        types.StringValue(p.GetLinkStatus()),
			VirtualCircuitIDs: vcIDs,
		}

		// Sort the ports by role, asserting the API always returns primary for len of 1 responses
		ret[order[p.GetRole()]] = connPort
	}

	return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, ret), nil
}

func parseConnectionVlans(ctx context.Context, conn *metalv1.Interconnection) (*basetypes.ListValue, diag.Diagnostics) {
	nPorts := len(conn.Ports)
	ret := make([]int, 0, nPorts)

	isVLANBasedConn := nPorts != 0 && conn.Ports[0].GetVirtualCircuits()[0].VlanVirtualCircuit != nil

	if isVLANBasedConn {
		order := map[metalv1.InterconnectionPortRole]int{
			metalv1.INTERCONNECTIONPORTROLE_PRIMARY:   0,
			metalv1.INTERCONNECTIONPORTROLE_SECONDARY: 1,
		}

		rawVlans := make([]int, len(conn.Ports))
		for _, p := range conn.Ports {
			rawVlans[order[p.GetRole()]] = int(p.VirtualCircuits[0].VlanVirtualCircuit.GetVnid())
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

func parseConnectionVrfs(ctx context.Context, conn *metalv1.Interconnection) (*basetypes.ListValue, diag.Diagnostics) {
	nPorts := len(conn.Ports)
	ret := make([]string, 0, nPorts)

	isVRFBasedConn := nPorts != 0 && conn.Ports[0].GetVirtualCircuits()[0].VrfVirtualCircuit != nil

	if isVRFBasedConn {
		order := map[metalv1.InterconnectionPortRole]int{
			metalv1.INTERCONNECTIONPORTROLE_PRIMARY:   0,
			metalv1.INTERCONNECTIONPORTROLE_SECONDARY: 1,
		}

		rawVrfs := make([]string, len(conn.Ports))
		for _, p := range conn.Ports {
			vrf := p.VirtualCircuits[0].VrfVirtualCircuit.GetVrf()

			// NB: The VC object on a in Interconnection does not include the
			// full VRF, it's an href instead. No way to remedy this with a
			// 'includes' query param so we need to grab the ID from this
			// instead.
			rawVrfs[order[p.GetRole()]] = path.Base(vrf.GetHref())
		}

		for _, v := range rawVrfs {
			if v != "" {
				ret = append(ret, v)
			}
		}
	}

	vrfs, diags := types.ListValueFrom(ctx, types.StringType, ret)
	if diags.HasError() {
		return nil, diags
	}
	return &vrfs, nil
}
