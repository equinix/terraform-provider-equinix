package organization

import (
	"context"
	"fmt"
	"regexp"

	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/packethost/packngo"
)

type AddressResourceModel struct {
	address types.String  `tfsdk:"address"`
	city    *types.String `tfsdk:"city"`
	country types.String  `tfsdk:"country"`
	state   *types.String `tfsdk:"state"`
	zipCode types.String  `tfsdk:"zip_code"`
}

type ResourceModel struct {
	id          types.String                                          `tfsdk:"id"`
	name        types.String                                          `tfsdk:"name"`
	description types.String                                          `tfsdk:"description"`
	website     types.String                                          `tfsdk:"website"`
	twitter     types.String                                          `tfsdk:"twitter"`
	logo        types.String                                          `tfsdk:"logo"`
	created     types.String                                          `tfsdk:"created"`
	updated     types.String                                          `tfsdk:"updated"`
	address     fwtypes.ListNestedObjectValueOf[AddressResourceModel] `tfsdk:"address"` // List of Address
}

func (m *ResourceModel) parse(ctx context.Context, org *packngo.Organization) diag.Diagnostics {
	m.id = types.StringValue(org.ID)
	m.name = types.StringValue(org.Name)
	m.description = types.StringValue(org.Description)
	m.website = types.StringValue(org.Website)
	m.twitter = types.StringValue(org.Twitter)
	m.logo = types.StringValue(org.Logo)
	m.created = types.StringValue(org.Created)
	m.updated = types.StringValue(org.Updated)

	var addressResModels []AddressResourceModel
	city := types.StringValue(*org.Address.City)
	state := types.StringValue(*org.Address.State)
	addressResModels = append(addressResModels, AddressResourceModel{
		address: types.StringValue(org.Address.Address),
		city:    &city,
		country: types.StringValue(org.Address.Country),
		state:   &state,
		zipCode: types.StringValue(org.Address.ZipCode),
	})
	m.address = fwtypes.NewListNestedObjectValueOfValueSlice(ctx, addressResModels)
	return nil
}

type DataSourceModel struct {
	id              types.String                                          `tfsdk:"id"`
	name            types.String                                          `tfsdk:"name"`
	organization_id types.String                                          `tfsdk:"organization_id"`
	description     types.String                                          `tfsdk:"description"`
	website         types.String                                          `tfsdk:"website"`
	twitter         types.String                                          `tfsdk:"twitter"`
	logo            types.String                                          `tfsdk:"logo"`
	project_ids     []types.List                                          `tfsdk:"project_ids"`
	address         fwtypes.ListNestedObjectValueOf[AddressResourceModel] `tfsdk:"address"` // List of Address
}

func (m *DataSourceModel) parse(ctx context.Context, org *packngo.Organization) diag.Diagnostics {
	var diags diag.Diagnostics
	m.id = types.StringValue(org.ID)
	m.name = types.StringValue(org.Name)
	m.organization_id = types.StringValue(org.ID)
	m.description = types.StringValue(org.Description)
	m.website = types.StringValue(org.Website)
	m.twitter = types.StringValue(org.Twitter)

	projects := make([]string, len(org.Projects))
	pList := make([]basetypes.ListValue, len(org.Projects))
	for i, p := range org.Projects {
		projects[i] = p.ID
		projList, _ := types.ListValueFrom(ctx, types.StringType, p.ID)
		pList = append(pList, projList)

	}

	m.project_ids = pList

	m.logo = types.StringValue(org.Logo)

	addressresourcemodel := make([]AddressResourceModel, 1)

	cityValue := types.StringValue(*org.Address.City)
	stateValue := types.StringValue(*org.Address.State)
	arm := AddressResourceModel{
		address: types.StringValue(org.Address.Address),
		city:    &cityValue,
		country: types.StringValue(org.Address.Country),
		state:   &stateValue,
		zipCode: types.StringValue(org.Address.ZipCode),
	}
	addressresourcemodel[0] = arm
	m.address = fwtypes.NewListNestedObjectValueOfValueSlice(ctx, addressresourcemodel)

	return diags
}

func StringToRegex(pattern string) (regExp *regexp.Regexp) {
	regExp, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}

	return regExp
}
