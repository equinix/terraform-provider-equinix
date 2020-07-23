package equinix

import (
	"ecx-go/v3"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var ecxPortSchemaNames = map[string]string{
	"UUID":          "uuid",
	"Name":          "name",
	"Region":        "region",
	"IBX":           "ibx",
	"MetroCode":     "metro_code",
	"Priority":      "priority",
	"Encapsulation": "encapsulation",
	"Buyout":        "buyout",
	"Bandwidth":     "bandwidth",
	"Status":        "status",
}

func dataSourceECXPort() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceECXPortRead,
		Schema: map[string]*schema.Schema{
			ecxPortSchemaNames["UUID"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			ecxPortSchemaNames["Name"]: {
				Type:     schema.TypeString,
				Required: true,
			},
			ecxPortSchemaNames["Region"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			ecxPortSchemaNames["IBX"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			ecxPortSchemaNames["MetroCode"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			ecxPortSchemaNames["Priority"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			ecxPortSchemaNames["Encapsulation"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			ecxPortSchemaNames["Buyout"]: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			ecxPortSchemaNames["Bandwidth"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			ecxPortSchemaNames["Status"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceECXPortRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	name := d.Get(ecxPortSchemaNames["Name"]).(string)
	ports, err := conf.ecx.GetUserPorts()
	if err != nil {
		return err
	}
	var filteredPorts []ecx.Port
	for _, port := range ports {
		if port.Name == name {
			filteredPorts = append(filteredPorts, port)
		}
	}
	if len(filteredPorts) < 1 {
		return fmt.Errorf("profile query returned no results, please change your search criteria")
	}
	if len(filteredPorts) > 1 {
		return fmt.Errorf("query returned more than one result, please try more specific search criteria")
	}
	return updateECXPortResource(ports[0], d)
}

func updateECXPortResource(port ecx.Port, d *schema.ResourceData) error {
	d.SetId(port.UUID)
	if err := d.Set(ecxPortSchemaNames["UUID"], port.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if err := d.Set(ecxPortSchemaNames["Region"], port.Region); err != nil {
		return fmt.Errorf("error reading Region: %s", err)
	}
	if err := d.Set(ecxPortSchemaNames["IBX"], port.IBX); err != nil {
		return fmt.Errorf("error reading IBX: %s", err)
	}
	if err := d.Set(ecxPortSchemaNames["MetroCode"], port.MetroCode); err != nil {
		return fmt.Errorf("error reading MetroCode: %s", err)
	}
	if err := d.Set(ecxPortSchemaNames["Priority"], port.Priority); err != nil {
		return fmt.Errorf("error reading Priority: %s", err)
	}
	if err := d.Set(ecxPortSchemaNames["Encapsulation"], port.Encapsulation); err != nil {
		return fmt.Errorf("error reading Encapsulation: %s", err)
	}
	if err := d.Set(ecxPortSchemaNames["Buyout"], port.Buyout); err != nil {
		return fmt.Errorf("error reading Buyout: %s", err)
	}
	if err := d.Set(ecxPortSchemaNames["Bandwidth"], port.Bandwidth); err != nil {
		return fmt.Errorf("error reading Bandwidth: %s", err)
	}
	if err := d.Set(ecxPortSchemaNames["Status"], port.Status); err != nil {
		return fmt.Errorf("error reading Status: %s", err)
	}
	return nil
}
