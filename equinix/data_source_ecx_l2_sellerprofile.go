package equinix

import (
	"fmt"

	"github.com/equinix/ecx-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var ecxL2SellerProfileSchemaNames = map[string]string{
	"UUID": "uuid",
	"Name": "name",
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
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceECXL2SellerProfileRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	name := d.Get(ecxPortSchemaNames["Name"]).(string)
	profiles, err := conf.ecx.GetL2SellerProfiles()
	if err != nil {
		return err
	}
	var filteredProfiles []ecx.L2ServiceProfile
	for _, profile := range profiles {
		if profile.Name == name {
			filteredProfiles = append(filteredProfiles, profile)
		}
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
	return nil
}
