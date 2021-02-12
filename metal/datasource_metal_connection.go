package metal

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func connectionPortSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
        "
    }
}

func dataSourceMetalConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalConnectionRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"connection_id"},
			},
			"connection_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"redundancy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"facility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"speed": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     connectionPortSchema(),
			},
		},
	}
}

func findConnectionByName(os []packngo.Connection, name string) (*packngo.Connection, error) {
	results := make([]packngo.Connection, 0)
	for _, o := range os {
		if o.Name == name {
			results = append(results, o)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no connection found with name %s", name)
	}
	return nil, fmt.Errorf("too many connections found with name %s (found %d, expected 1)", name, len(results))
}

func dataSourceMetalConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	nameRaw, nameOK := d.GetOk("name")
	orgIdRaw, orgIdOK := d.GetOk("connection_id")

	if !orgIdOK && !nameOK {
		return fmt.Errorf("You must supply connection_id or name")
	}
	var org *packngo.Connection

	if nameOK {
		name := nameRaw.(string)

		os, _, err := client.Connections.List(nil)
		if err != nil {
			return err
		}

		org, err = findOrgByName(os, name)
		if err != nil {
			return err
		}
	} else {
		orgId := orgIdRaw.(string)
		log.Println(orgId)
		var err error
		org, _, err = client.Connections.Get(orgId, nil)
		if err != nil {
			return err
		}
	}
	projectIds := []string{}

	for _, p := range org.Projects {
		projectIds = append(projectIds, filepath.Base(p.URL))
	}

	d.Set("connection_id", org.ID)
	d.Set("name", org.Name)
	d.Set("description", org.Description)
	d.Set("website", org.Website)
	d.Set("twitter", org.Twitter)
	d.Set("logo", org.Logo)
	d.Set("project_ids", projectIds)
	d.SetId(org.ID)

	return nil
}
