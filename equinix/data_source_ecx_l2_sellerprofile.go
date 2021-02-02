package equinix

import (
	"bytes"
	"context"
	"fmt"

	"github.com/equinix/ecx-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform/helper/hashcode"
)

var ecxL2SellerProfileSchemaNames = map[string]string{
	"UUID":               "uuid",
	"Name":               "name",
	"Description":        "description",
	"SpeedFromAPI":       "speed_from_api",
	"AllowCustomSpeed":   "speed_customization_allowed",
	"RequiredRedundancy": "redundancy_required",
	"Encapsulation":      "encapsulation",
	"GlobalOrganization": "organization_global_name",
	"OrganizationName":   "organization_name",
	"SpeedBand":          "speed_band",
	"Metros":             "metro",
	"AdditionalInfos":    "additional_info",
}

var ecxL2SellerProfileDescriptions = map[string]string{
	"UUID":               "Unique identifier of the seller profile",
	"Name":               "Name of the seller profile",
	"Description":        "Seller Profile text description",
	"SpeedFromAPI":       "Boolean that indicates if seller is deriving connection speed from an API call",
	"AllowCustomSpeed":   "Boolean that indicates if seller allows customer to enter a custom connection speed",
	"RequiredRedundancy": "Boolean that indicate if seller requires connections to be redundant",
	"Encapsulation":      "Seller profile's encapsulation (either Dot1q or QinQ)",
	"GlobalOrganization": "Name of seller's global organization",
	"OrganizationName":   "Name of seller's organization",
	"SpeedBand":          "One or more specifications of speed/bandwidth supported by given seller profile",
	"Metros":             "One or more specifications of metro locations supported by seller profile",
	"AdditionalInfos":    "One or more specifications of additional buyer information attributes that can be provided in connection definition that uses given seller profile",
}

var ecxL2SellerProfileMetrosSchemaNames = map[string]string{
	"Code":    "code",
	"Name":    "name",
	"IBXes":   "ibxes",
	"Regions": "regions",
}

var ecxL2SellerProfileMetrosDescriptions = map[string]string{
	"Code":    "Location metro code",
	"Name":    "Location metro nam",
	"IBXes":   "List of IBXes supported within given metro",
	"Regions": "List of regions supported within given metro",
}

var ecxL2SellerProfileAdditionalInfosSchemaNames = map[string]string{
	"Name":             "name",
	"Description":      "description",
	"DataType":         "data_type",
	"IsMandatory":      "mandatory",
	"IsCaptureInEmail": "captured_in_email",
}

var ecxL2SellerProfileAdditionalInfosDescriptions = map[string]string{
	"Name":             "Name of additional information attribute",
	"Description":      "Textual description of additional information attribute",
	"DataType":         "Data type of additional information attribute. Either BOOLEAN, INTEGER or STRING",
	"IsMandatory":      "Specifies if additional information attribute is mandatory to create connection",
	"IsCaptureInEmail": "Specified if additional information attribute can be captured in email",
}

func dataSourceECXL2SellerProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceECXL2SellerProfileRead,
		Description: "Use this data source to get details of Equinix Fabric layer 2	seller profile with a given name and / or organization",
		Schema: map[string]*schema.Schema{
			ecxL2SellerProfileSchemaNames["UUID"]: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: ecxL2SellerProfileDescriptions["UUID"],
			},
			ecxL2SellerProfileSchemaNames["Name"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  ecxL2SellerProfileDescriptions["Name"],
			},
			ecxL2SellerProfileSchemaNames["Description"]: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: ecxL2SellerProfileDescriptions["Description"],
			},
			ecxL2SellerProfileSchemaNames["SpeedFromAPI"]: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: ecxL2SellerProfileDescriptions["SpeedFromAPI"],
			},
			ecxL2SellerProfileSchemaNames["AllowCustomSpeed"]: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: ecxL2SellerProfileDescriptions["AllowCustomSpeed"],
			},
			ecxL2SellerProfileSchemaNames["RequiredRedundancy"]: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: ecxL2SellerProfileDescriptions["RequiredRedundancy"],
			},
			ecxL2SellerProfileSchemaNames["Encapsulation"]: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: ecxL2SellerProfileDescriptions["Encapsulation"],
			},
			ecxL2SellerProfileSchemaNames["GlobalOrganization"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  ecxL2SellerProfileDescriptions["GlobalOrganization"],
			},
			ecxL2SellerProfileSchemaNames["OrganizationName"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  ecxL2SellerProfileDescriptions["OrganizationName"],
			},
			ecxL2SellerProfileSchemaNames["SpeedBand"]: {
				Type:        schema.TypeSet,
				Computed:    true,
				Set:         ecxL2ServiceProfileSpeedBandHash,
				Description: ecxL2SellerProfileDescriptions["SpeedBand"],
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ecxL2ServiceProfileSpeedBandSchemaNames["Speed"]: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: ecxL2ServiceProfileSpeedBandDescriptions["Speed"],
						},
						ecxL2ServiceProfileSpeedBandSchemaNames["SpeedUnit"]: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: ecxL2ServiceProfileSpeedBandDescriptions["SpeedUnit"],
						},
					},
				},
			},
			ecxL2SellerProfileSchemaNames["Metros"]: {
				Type:        schema.TypeSet,
				Computed:    true,
				Set:         ecxL2SellerProfileMetroHash,
				Description: ecxL2SellerProfileDescriptions["Metros"],
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ecxL2SellerProfileMetrosSchemaNames["Code"]: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: ecxL2SellerProfileMetrosDescriptions["Code"],
						},
						ecxL2SellerProfileMetrosSchemaNames["Name"]: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: ecxL2SellerProfileMetrosDescriptions["Name"],
						},
						ecxL2SellerProfileMetrosSchemaNames["IBXes"]: {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: ecxL2SellerProfileMetrosDescriptions["IBXes"],
						},
						ecxL2SellerProfileMetrosSchemaNames["Regions"]: {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: ecxL2SellerProfileMetrosDescriptions["Regions"],
						},
					},
				},
			},
			ecxL2SellerProfileSchemaNames["AdditionalInfos"]: {
				Type:        schema.TypeSet,
				Computed:    true,
				Set:         ecxL2SellerProfileAdditionalInfoHash,
				Description: ecxL2SellerProfileDescriptions["AdditionalInfos"],
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ecxL2SellerProfileAdditionalInfosSchemaNames["Name"]: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: ecxL2SellerProfileAdditionalInfosDescriptions["Name"],
						},
						ecxL2SellerProfileAdditionalInfosSchemaNames["Description"]: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: ecxL2SellerProfileAdditionalInfosDescriptions["Description"],
						},
						ecxL2SellerProfileAdditionalInfosSchemaNames["DataType"]: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: ecxL2SellerProfileAdditionalInfosDescriptions["DataType"],
						},
						ecxL2SellerProfileAdditionalInfosSchemaNames["IsMandatory"]: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: ecxL2SellerProfileAdditionalInfosDescriptions["IsMandatory"],
						},
						ecxL2SellerProfileAdditionalInfosSchemaNames["IsCaptureInEmail"]: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: ecxL2SellerProfileAdditionalInfosDescriptions["IsCaptureInEmail"],
						},
					},
				},
			},
		},
	}
}

func dataSourceECXL2SellerProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*Config)
	var diags diag.Diagnostics
	name := d.Get(ecxL2SellerProfileSchemaNames["Name"]).(string)
	orgName := d.Get(ecxL2SellerProfileSchemaNames["OrganizationName"]).(string)
	orgGlobalName := d.Get(ecxL2SellerProfileSchemaNames["GlobalOrganization"]).(string)
	profiles, err := conf.ecx.GetL2SellerProfiles()
	if err != nil {
		return diag.FromErr(err)
	}
	var filteredProfiles []ecx.L2ServiceProfile
	for _, profile := range profiles {
		if name != "" && ecx.StringValue(profile.Name) != name {
			continue
		}
		if orgName != "" && ecx.StringValue(profile.OrganizationName) != orgName {
			continue
		}
		if orgGlobalName != "" && ecx.StringValue(profile.GlobalOrganization) != orgGlobalName {
			continue
		}
		filteredProfiles = append(filteredProfiles, profile)
	}
	if len(filteredProfiles) < 1 {
		return diag.Errorf("profile query returned no results, please change your search criteria")
	}
	if len(filteredProfiles) > 1 {
		return diag.Errorf("query returned more than one result, please try more specific search criteria")
	}
	if err := updateECXL2SellerProfileResource(filteredProfiles[0], d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateECXL2SellerProfileResource(profile ecx.L2ServiceProfile, d *schema.ResourceData) error {
	d.SetId(ecx.StringValue(profile.UUID))
	if err := d.Set(ecxL2SellerProfileSchemaNames["UUID"], profile.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if err := d.Set(ecxL2SellerProfileSchemaNames["Name"], profile.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(ecxL2SellerProfileSchemaNames["Description"], profile.Description); err != nil {
		return fmt.Errorf("error reading Description: %s", err)
	}
	if err := d.Set(ecxL2SellerProfileSchemaNames["SpeedFromAPI"], profile.SpeedFromAPI); err != nil {
		return fmt.Errorf("error reading SpeedFromAPI: %s", err)
	}
	if err := d.Set(ecxL2SellerProfileSchemaNames["AllowCustomSpeed"], profile.AllowCustomSpeed); err != nil {
		return fmt.Errorf("error reading AllowCustomSpeed: %s", err)
	}
	if err := d.Set(ecxL2SellerProfileSchemaNames["RequiredRedundancy"], profile.RequiredRedundancy); err != nil {
		return fmt.Errorf("error reading RequiredRedundancy: %s", err)
	}
	if err := d.Set(ecxL2SellerProfileSchemaNames["Encapsulation"], profile.Encapsulation); err != nil {
		return fmt.Errorf("error reading Encapsulation: %s", err)
	}
	if err := d.Set(ecxL2SellerProfileSchemaNames["GlobalOrganization"], profile.GlobalOrganization); err != nil {
		return fmt.Errorf("error reading GlobalOrganization: %s", err)
	}
	if err := d.Set(ecxL2SellerProfileSchemaNames["OrganizationName"], profile.OrganizationName); err != nil {
		return fmt.Errorf("error reading OrganizationName: %s", err)
	}
	if err := d.Set(ecxL2SellerProfileSchemaNames["SpeedBand"], flattenECXL2ServiceProfileSpeedBands(profile.SpeedBands)); err != nil {
		return fmt.Errorf("error reading SpeedBand: %s", err)
	}
	if err := d.Set(ecxL2SellerProfileSchemaNames["Metros"], flattenECXL2SellerProfileMetros(profile.Metros)); err != nil {
		return fmt.Errorf("error reading Metros: %s", err)
	}
	if err := d.Set(ecxL2SellerProfileSchemaNames["AdditionalInfos"], flattenECXL2SellerProfileAdditionalInfos(profile.AdditionalInfos)); err != nil {
		return fmt.Errorf("error reading AdditionalInfos: %s", err)
	}
	return nil
}

func flattenECXL2SellerProfileMetros(metros []ecx.L2SellerProfileMetro) interface{} {
	transformed := make([]interface{}, len(metros))
	for i := range metros {
		transformed[i] = map[string]interface{}{
			ecxL2SellerProfileMetrosSchemaNames["Code"]:    metros[i].Code,
			ecxL2SellerProfileMetrosSchemaNames["Name"]:    metros[i].Name,
			ecxL2SellerProfileMetrosSchemaNames["IBXes"]:   metros[i].IBXes,
			ecxL2SellerProfileMetrosSchemaNames["Regions"]: metros[i].Regions,
		}
	}
	return transformed
}

func flattenECXL2SellerProfileAdditionalInfos(infos []ecx.L2SellerProfileAdditionalInfo) interface{} {
	transformed := make([]interface{}, len(infos))
	for i := range infos {
		transformed[i] = map[string]interface{}{
			ecxL2SellerProfileAdditionalInfosSchemaNames["Name"]:             infos[i].Name,
			ecxL2SellerProfileAdditionalInfosSchemaNames["Description"]:      infos[i].Description,
			ecxL2SellerProfileAdditionalInfosSchemaNames["DataType"]:         infos[i].DataType,
			ecxL2SellerProfileAdditionalInfosSchemaNames["IsMandatory"]:      infos[i].IsMandatory,
			ecxL2SellerProfileAdditionalInfosSchemaNames["IsCaptureInEmail"]: infos[i].IsCaptureInEmail,
		}
	}
	return transformed
}

func ecxL2ServiceProfileSpeedBandHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m[ecxL2ServiceProfileSpeedBandSchemaNames["Speed"]].(int)))
	buf.WriteString(fmt.Sprintf("%s-", m[ecxL2ServiceProfileSpeedBandSchemaNames["SpeedUnit"]].(string)))
	return hashcode.String(buf.String())
}

func ecxL2SellerProfileMetroHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m[ecxL2SellerProfileMetrosSchemaNames["Code"]].(string)))
	return hashcode.String(buf.String())
}

func ecxL2SellerProfileAdditionalInfoHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m[ecxL2SellerProfileAdditionalInfosSchemaNames["Name"]].(string)))
	return hashcode.String(buf.String())
}
