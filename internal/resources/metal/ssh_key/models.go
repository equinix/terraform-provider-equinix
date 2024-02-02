package ssh_key

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
}

func (m *ResourceModel) parse(key *packngo.SSHKey) diag.Diagnostics {
	m.ID = types.StringValue(key.ID)
	m.Name = types.StringValue(key.Label)
	m.PublicKey = types.StringValue(key.Key)
	m.Fingerprint = types.StringValue(key.FingerPrint)
	m.Created = types.StringValue(key.Created)
	m.Updated = types.StringValue(key.Updated)
	m.OwnerID = types.StringValue(path.Base(key.Owner.Href))
	return nil
}
