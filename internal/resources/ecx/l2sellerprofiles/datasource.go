package l2sellerprofiles

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/equinix/ecx-go/v2"
	"github.com/equinix/terraform-provider-equinix/internal/comparisons"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/equinix/terraform-provider-equinix/internal/resources/ecx/l2sellerprofile"
	"github.com/equinix/terraform-provider-equinix/internal/resources/ecx/l2serviceprofile"
	"github.com/equinix/terraform-provider-equinix/internal/validaters"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var ecxL2SellerProfilesSchemaNames = map[string]string{
	"NameRegex":          "name_regex",
	"Metros":             "metro_codes",
	"SpeedBands":         "speed_bands",
	"OrganizationName":   "organization_name",
	"GlobalOrganization": "organization_global_name",
	"Profiles":           "profiles",
}

var ecxL2SellerProfilesDescriptions = map[string]string{
	"NameRegex":          "A regex string to apply on returned seller profile names and filter search results",
	"Metros":             "List of metro codes of locations that should be served by resulting profiles",
	"SpeedBands":         "List of speed bands that should be supported by resulting profiles",
	"OrganizationName":   "Name of seller's organization",
	"GlobalOrganization": "Name of seller's global organization",
	"Profiles":           "Resulting list of profiles that match filtering criteria",
}

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceECXL2SellerProfilesRead,
		Description: "Use this data source to get list of Equinix Fabric layer 2 seller profiles",
		Schema: map[string]*schema.Schema{
			ecxL2SellerProfilesSchemaNames["NameRegex"]: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  ecxL2SellerProfilesDescriptions["NameRegex"],
			},
			ecxL2SellerProfilesSchemaNames["Metros"]: {
				Type:        schema.TypeSet,
				Optional:    true,
				MinItems:    1,
				Description: ecxL2SellerProfilesDescriptions["Metros"],
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validaters.StringIsMetroCode(),
				},
			},
			ecxL2SellerProfilesSchemaNames["SpeedBands"]: {
				Type:        schema.TypeSet,
				Optional:    true,
				MinItems:    1,
				Description: ecxL2SellerProfilesDescriptions["SpeedBands"],
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validaters.StringIsSpeedBand(),
				},
			},
			ecxL2SellerProfilesSchemaNames["OrganizationName"]: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  ecxL2SellerProfilesDescriptions["OrganizationName"],
			},
			ecxL2SellerProfilesSchemaNames["GlobalOrganization"]: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  ecxL2SellerProfilesDescriptions["GlobalOrganization"],
			},
			ecxL2SellerProfilesSchemaNames["Profiles"]: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: ecxL2SellerProfilesDescriptions["Profiles"],
				Elem: &schema.Resource{
					Schema: l2sellerprofile.CreateECXL2SellerProfileSchema(),
				},
			},
		},
	}
}

func dataSourceECXL2SellerProfilesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*config.Config)
	var diags diag.Diagnostics
	profiles, err := conf.ECXClient.GetL2SellerProfiles()
	if err != nil {
		return diag.FromErr(err)
	}
	var filteredProfiles []ecx.L2ServiceProfile
	nameRegex := d.Get(ecxL2SellerProfilesSchemaNames["NameRegex"]).(string)
	metros := converters.ExpandSetToStringList(d.Get(ecxL2SellerProfilesSchemaNames["Metros"]).(*schema.Set))
	speedBands := converters.ExpandSetToStringList(d.Get(ecxL2SellerProfilesSchemaNames["SpeedBands"]).(*schema.Set))
	orgName := d.Get(ecxL2SellerProfilesSchemaNames["OrganizationName"]).(string)
	globalOrgName := d.Get(ecxL2SellerProfilesSchemaNames["GlobalOrganization"]).(string)
	for _, profile := range profiles {
		if nameRegex != "" {
			r := regexp.MustCompile(nameRegex)
			if !r.MatchString(ecx.StringValue(profile.Name)) {
				continue
			}
		}
		if len(metros) > 0 && !comparisons.AtLeastOneStringFound(metros, flattenECXL2SellerProfileMetroCodes(profile.Metros)) {
			continue
		}
		if len(speedBands) > 0 && !comparisons.AtLeastOneStringFound(speedBands, flattenECXL2SellerProfileSpeedBands(profile.SpeedBands)) {
			continue
		}
		if orgName != "" && !strings.EqualFold(ecx.StringValue(profile.OrganizationName), orgName) {
			continue
		}
		if globalOrgName != "" && !strings.EqualFold(ecx.StringValue(profile.GlobalOrganization), globalOrgName) {
			continue
		}
		filteredProfiles = append(filteredProfiles, profile)
	}
	if len(filteredProfiles) < 1 {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "profile query returned no results, please change your search criteria",
			AttributePath: cty.Path{cty.GetAttrStep{}},
		})
		return diags
	}
	if err := updateECXL2SellerProfilesResource(filteredProfiles, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateECXL2SellerProfilesResource(profiles []ecx.L2ServiceProfile, d *schema.ResourceData) error {
	d.SetId("ecxL2SellerProfiles")
	if err := d.Set(ecxL2SellerProfilesSchemaNames["Profiles"], flattenECXL2SellerProfiles(profiles)); err != nil {
		return fmt.Errorf("error reading profiles: %s", err)
	}
	return nil
}

func flattenECXL2SellerProfileMetroCodes(metros []ecx.L2SellerProfileMetro) []string {
	transformed := make([]string, len(metros))
	for i := range metros {
		transformed[i] = ecx.StringValue(metros[i].Code)
	}
	return transformed
}

func flattenECXL2SellerProfileSpeedBands(speedBands []ecx.L2ServiceProfileSpeedBand) []string {
	transformed := make([]string, len(speedBands))
	for i := range speedBands {
		transformed[i] = fmt.Sprintf("%d%s", ecx.IntValue(speedBands[i].Speed), ecx.StringValue(speedBands[i].SpeedUnit))
	}
	return transformed
}

func flattenECXL2SellerProfiles(profiles []ecx.L2ServiceProfile) interface{} {
	transformed := make([]interface{}, len(profiles))
	for i := range profiles {
		transformedProfile := make(map[string]interface{})
		transformedProfile[l2sellerprofile.EcxL2SellerProfileSchemaNames["UUID"]] = profiles[i].UUID
		transformedProfile[l2sellerprofile.EcxL2SellerProfileSchemaNames["Name"]] = profiles[i].Name
		transformedProfile[l2sellerprofile.EcxL2SellerProfileSchemaNames["Description"]] = profiles[i].Description
		transformedProfile[l2sellerprofile.EcxL2SellerProfileSchemaNames["SpeedFromAPI"]] = profiles[i].SpeedFromAPI
		transformedProfile[l2sellerprofile.EcxL2SellerProfileSchemaNames["AllowCustomSpeed"]] = profiles[i].AllowCustomSpeed
		transformedProfile[l2sellerprofile.EcxL2SellerProfileSchemaNames["RequiredRedundancy"]] = profiles[i].RequiredRedundancy
		transformedProfile[l2sellerprofile.EcxL2SellerProfileSchemaNames["Encapsulation"]] = profiles[i].Encapsulation
		transformedProfile[l2sellerprofile.EcxL2SellerProfileSchemaNames["GlobalOrganization"]] = profiles[i].GlobalOrganization
		transformedProfile[l2sellerprofile.EcxL2SellerProfileSchemaNames["OrganizationName"]] = profiles[i].OrganizationName
		transformedProfile[l2sellerprofile.EcxL2SellerProfileSchemaNames["SpeedBand"]] = l2serviceprofile.FlattenECXL2ServiceProfileSpeedBands(profiles[i].SpeedBands)
		transformedProfile[l2sellerprofile.EcxL2SellerProfileSchemaNames["Metros"]] = l2sellerprofile.FlattenECXL2SellerProfileMetros(profiles[i].Metros)
		transformedProfile[l2sellerprofile.EcxL2SellerProfileSchemaNames["AdditionalInfos"]] = l2sellerprofile.FlattenECXL2SellerProfileAdditionalInfos(profiles[i].AdditionalInfos)
		transformed[i] = transformedProfile
	}
	return transformed
}