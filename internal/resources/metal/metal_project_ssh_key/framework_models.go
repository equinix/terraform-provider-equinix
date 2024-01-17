package metal_project_ssh_key

import (
	"path"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
)

type ResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	PublicKey   types.String `tfsdk:"public_key"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	Created     types.String `tfsdk:"created"`
	Updated     types.String `tfsdk:"updated"`
	OwnerID     types.String `tfsdk:"owner_id"`
	ProjectID   types.String `tfsdk:"project_id"`
}

func (m *ResourceModel) parse(key *packngo.SSHKey) diag.Diagnostics {
	m.ID = types.StringValue(key.ID)
	m.Name = types.StringValue(key.Label)
	m.PublicKey = types.StringValue(key.Key)
	m.Fingerprint = types.StringValue(key.FingerPrint)
	m.Created = types.StringValue(key.Created)
	m.Updated = types.StringValue(key.Updated)
	m.OwnerID = types.StringValue(path.Base(key.Owner.Href))
	m.ProjectID = m.OwnerID
	return nil
}

// TODO (ocobles) ideally we would embed ResourceModel instead of
// explicitly define all the ResourceModel fields again in DataSourceModel
// https://github.com/hashicorp/terraform-plugin-framework/issues/242
type DataSourceModel struct {
	Search      types.String `tfsdk:"search"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	PublicKey   types.String `tfsdk:"public_key"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	Created     types.String `tfsdk:"created"`
	Updated     types.String `tfsdk:"updated"`
	OwnerID     types.String `tfsdk:"owner_id"`
	ProjectID   types.String `tfsdk:"project_id"`
}

func (m *DataSourceModel) parse(key *packngo.SSHKey) diag.Diagnostics {
	m.ID = types.StringValue(key.ID)
	m.Name = types.StringValue(key.Label)
	m.PublicKey = types.StringValue(key.Key)
	m.Fingerprint = types.StringValue(key.FingerPrint)
	m.Created = types.StringValue(key.Created)
	m.Updated = types.StringValue(key.Updated)
	m.OwnerID = types.StringValue(path.Base(key.Owner.Href))
	m.ProjectID = m.OwnerID
	return nil
}