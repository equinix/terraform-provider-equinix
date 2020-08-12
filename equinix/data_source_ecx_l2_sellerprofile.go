package equinix

import (
	"bytes"
	"fmt"

	"github.com/equinix/ecx-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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

var ecxL2SellerProfileMetrosSchemaNames = map[string]string{
	"Code":    "code",
	"Name":    "name",
	"IBXes":   "ibxes",
	"Regions": "regions",
}

var ecxL2SellerProfileAdditionalInfosSchemaNames = map[string]string{
	"Name":             "name",
	"Description":      "description",
	"DataType":         "data_type",
	"IsMandatory":      "mandatory",
	"IsCaptureInEmail": "captured_in_email",
}

func dataSourceECXL2SellerProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceECXL2SellerProfileRead,
		Schema: map[string]*schema.Schema{
			ecxL2SellerProfileSchemaNames["UUID"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			ecxL2SellerProfileSchemaNames["Name"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			ecxL2SellerProfileSchemaNames["Description"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			ecxL2SellerProfileSchemaNames["SpeedFromAPI"]: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			ecxL2SellerProfileSchemaNames["AllowCustomSpeed"]: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			ecxL2SellerProfileSchemaNames["RequiredRedundancy"]: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			ecxL2SellerProfileSchemaNames["Encapsulation"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			ecxL2SellerProfileSchemaNames["GlobalOrganization"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			ecxL2SellerProfileSchemaNames["OrganizationName"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			ecxL2SellerProfileSchemaNames["SpeedBand"]: {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      ecxL2ServiceProfileSpeedBandHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ecxL2ServiceProfileSpeedBandSchemaNames["Speed"]: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						ecxL2ServiceProfileSpeedBandSchemaNames["SpeedUnit"]: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			ecxL2SellerProfileSchemaNames["Metros"]: {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      ecxL2SellerProfileMetroHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ecxL2SellerProfileMetrosSchemaNames["Code"]: {
							Type:     schema.TypeString,
							Computed: true,
						},
						ecxL2SellerProfileMetrosSchemaNames["Name"]: {
							Type:     schema.TypeString,
							Computed: true,
						},
						ecxL2SellerProfileMetrosSchemaNames["IBXes"]: {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						ecxL2SellerProfileMetrosSchemaNames["Regions"]: {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			ecxL2SellerProfileSchemaNames["AdditionalInfos"]: {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      ecxL2SellerProfileAdditionalInfoHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ecxL2SellerProfileAdditionalInfosSchemaNames["Name"]: {
							Type:     schema.TypeString,
							Computed: true,
						},
						ecxL2SellerProfileAdditionalInfosSchemaNames["Description"]: {
							Type:     schema.TypeString,
							Computed: true,
						},
						ecxL2SellerProfileAdditionalInfosSchemaNames["DataType"]: {
							Type:     schema.TypeString,
							Computed: true,
						},
						ecxL2SellerProfileAdditionalInfosSchemaNames["IsMandatory"]: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						ecxL2SellerProfileAdditionalInfosSchemaNames["IsCaptureInEmail"]: {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceECXL2SellerProfileRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	name := d.Get(ecxL2SellerProfileSchemaNames["Name"]).(string)
	orgName := d.Get(ecxL2SellerProfileSchemaNames["OrganizationName"]).(string)
	orgGlobalName := d.Get(ecxL2SellerProfileSchemaNames["GlobalOrganization"]).(string)
	profiles, err := conf.ecx.GetL2SellerProfiles()
	if err != nil {
		return err
	}
	var filteredProfiles []ecx.L2ServiceProfile
	for _, profile := range profiles {
		if name != "" && profile.Name != name {
			continue
		}
		if orgName != "" && profile.OrganizationName != orgName {
			continue
		}
		if orgGlobalName != "" && profile.GlobalOrganization != orgGlobalName {
			continue
		}
		filteredProfiles = append(filteredProfiles, profile)
	}
	if len(filteredProfiles) < 1 {
		return fmt.Errorf("profile query returned no results, please change your search criteria")
	}
	if len(filteredProfiles) > 1 {
		return fmt.Errorf("query returned more than one result, please try more specific search criteria")
	}
	return updateECXL2SellerProfileResource(filteredProfiles[0], d)
}

func updateECXL2SellerProfileResource(profile ecx.L2ServiceProfile, d *schema.ResourceData) error {
	d.SetId(profile.UUID)
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
