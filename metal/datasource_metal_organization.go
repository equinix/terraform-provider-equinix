package metal

import (
	"fmt"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalOrganizationRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Description:   "The organization name",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"organization_id"},
			},
			"organization_id": {
				Type:          schema.TypeString,
				Description:   "The UUID of the organization resource",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description string",
				Computed:    true,
			},

			"website": {
				Type:        schema.TypeString,
				Description: "Website link",
				Computed:    true,
			},

			"twitter": {
				Type:        schema.TypeString,
				Description: "Twitter handle",
				Computed:    true,
			},
			"logo": {
				Type:        schema.TypeString,
				Description: "Logo URL",
				Computed:    true,
			},
			"project_ids": {
				Type:        schema.TypeList,
				Description: "UUIDs of project resources which belong to this organization",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func findOrgByName(os []packngo.Organization, name string) (*packngo.Organization, error) {
	results := make([]packngo.Organization, 0)
	for _, o := range os {
		if o.Name == name {
			results = append(results, o)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no organization found with name %s", name)
	}
	return nil, fmt.Errorf("too many organizations found with name %s (found %d, expected 1)", name, len(results))
}

func dataSourceMetalOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	nameRaw, nameOK := d.GetOk("name")
	orgIdRaw, orgIdOK := d.GetOk("organization_id")

	if !orgIdOK && !nameOK {
		return fmt.Errorf("You must supply organization_id or name")
	}
	var org *packngo.Organization

	if nameOK {
		name := nameRaw.(string)

		os, _, err := client.Organizations.List(nil)
		if err != nil {
			return err
		}

		org, err = findOrgByName(os, name)
		if err != nil {
			return err
		}
	} else {
		orgId := orgIdRaw.(string)
		var err error
		org, _, err = client.Organizations.Get(orgId, nil)
		if err != nil {
			return err
		}
	}
	projectIds := []string{}

	for _, p := range org.Projects {
		projectIds = append(projectIds, path.Base(p.URL))
	}

	d.Set("organization_id", org.ID)
	d.Set("name", org.Name)
	d.Set("description", org.Description)
	d.Set("website", org.Website)
	d.Set("twitter", org.Twitter)
	d.Set("logo", org.Logo)
	d.Set("project_ids", projectIds)
	d.SetId(org.ID)

	return nil
}
