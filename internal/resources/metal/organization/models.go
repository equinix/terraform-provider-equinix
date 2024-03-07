package organization

import (
	"context"

	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/packethost/packngo"
)

type AddressResourceModel struct {
	Address types.String `tfsdk:"address"`
	City    types.String `tfsdk:"city"`
	Country types.String `tfsdk:"country"`
	State   types.String `tfsdk:"state"`
	ZipCode types.String `tfsdk:"zip_code"`
}

type ResourceModel struct {
	ID          types.String                                          `tfsdk:"id"`
	Name        types.String                                          `tfsdk:"name"`
	Description types.String                                          `tfsdk:"description"`
	Website     types.String                                          `tfsdk:"website"`
	Twitter     types.String                                          `tfsdk:"twitter"`
	Logo        types.String                                          `tfsdk:"logo"`
	Created     types.String                                          `tfsdk:"created"`
	Updated     types.String                                          `tfsdk:"updated"`
	Address     fwtypes.ListNestedObjectValueOf[AddressResourceModel] `tfsdk:"address"`
}

func (m *ResourceModel) parse(ctx context.Context, org *packngo.Organization) diag.Diagnostics {
	m.ID = types.StringValue(org.ID)
	m.Name = types.StringValue(org.Name)
	m.Description = types.StringValue(org.Description)
	m.Website = types.StringValue(org.Website)
	m.Twitter = types.StringValue(org.Twitter)
	m.Logo = types.StringValue(org.Logo)
	m.Created = types.StringValue(org.Created)
	m.Updated = types.StringValue(org.Updated)

	m.Address = parseAddress(ctx, org.Address)

	return nil
}

type DataSourceModel struct {
	ID             types.String                                          `tfsdk:"id"`
	Name           types.String                                          `tfsdk:"name"`
	OrganizationID types.String                                          `tfsdk:"organization_id"`
	Description    types.String                                          `tfsdk:"description"`
	Website        types.String                                          `tfsdk:"website"`
	Twitter        types.String                                          `tfsdk:"twitter"`
	Logo           types.String                                          `tfsdk:"logo"`
	ProjectIDs     []types.List                                          `tfsdk:"project_ids"`
	Address        fwtypes.ListNestedObjectValueOf[AddressResourceModel] `tfsdk:"address"` // List of Address
}

func (m *DataSourceModel) parse(ctx context.Context, org *packngo.Organization) diag.Diagnostics {
	var diags diag.Diagnostics
	// Convert Metal Organization data to the Terraform state
	m.ID = types.StringValue(org.ID)
	m.Name = types.StringValue(org.Name)
	m.OrganizationID = types.StringValue(org.ID)
	m.Description = types.StringValue(org.Description)
	m.Website = types.StringValue(org.Website)
	m.Twitter = types.StringValue(org.Twitter)
	m.Logo = types.StringValue(org.Logo)
	m.Address = parseAddress(ctx, org.Address)

	projects := make([]string, len(org.Projects))
	pList := make([]basetypes.ListValue, len(org.Projects))
	for i, p := range org.Projects {
		projects[i] = p.ID
		projList, _ := types.ListValueFrom(ctx, types.StringType, p.ID)
		pList = append(pList, projList)
	}
	m.ProjectIDs = pList

	return diags
}

func parseAddress(ctx context.Context, addr packngo.Address) fwtypes.ListNestedObjectValueOf[AddressResourceModel] {
	addressresourcemodel := make([]AddressResourceModel, 1)
	addressresourcemodel[0] = AddressResourceModel{
		Address: types.StringValue(addr.Address),
		City:    types.StringPointerValue(addr.City),
		Country: types.StringValue(addr.Country),
		State:   types.StringPointerValue(addr.State),
		ZipCode: types.StringValue(addr.ZipCode),
	}
	return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, addressresourcemodel)
}
