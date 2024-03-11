package project_ssh_key

import (
	"path"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func (m *ResourceModel) parse(key *metalv1.SSHKey) diag.Diagnostics {
	m.ID = types.StringValue(key.GetId())
	m.Name = types.StringValue(key.GetLabel())
	m.PublicKey = types.StringValue(key.GetKey())
	m.Fingerprint = types.StringValue(key.GetFingerprint())
	m.Created = types.StringValue(key.CreatedAt.GoString())
	m.Updated = types.StringValue(key.UpdatedAt.GoString())

	ownerID := key.AdditionalProperties["owner"].(map[string]interface{})
	m.OwnerID = types.StringValue(path.Base(ownerID["href"].(string)))
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

func (m *DataSourceModel) parse(key *metalv1.SSHKey) diag.Diagnostics {
	m.ID = types.StringValue(key.GetId())
	m.Name = types.StringValue(key.GetLabel())
	m.PublicKey = types.StringValue(key.GetKey())
	m.Fingerprint = types.StringValue(key.GetFingerprint())
	m.Created = types.StringValue(key.CreatedAt.GoString())
	m.Updated = types.StringValue(key.UpdatedAt.GoString())

	ownerID := key.AdditionalProperties["owner"].(map[string]interface{})
	m.OwnerID = types.StringValue(path.Base(ownerID["href"].(string)))
	m.ProjectID = m.OwnerID
	return nil
}
