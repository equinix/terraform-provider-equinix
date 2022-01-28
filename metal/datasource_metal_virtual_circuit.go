package metal

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalVirtualCircuit() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalVirtualCircuitRead,

		Schema: map[string]*schema.Schema{
			"virtual_circuit_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the virtual circuit to lookup",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the virtual circuit",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the virtual circuit",
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "Tags attached to the virtual circuit",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the virtual circuit",
			},
			"vnid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "VNID VLAN parameter, see https://metal.equinix.com/developers/docs/networking/fabric/",
			},
			"nni_vnid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Nni VLAN ID parameter, see https://metal.equinix.com/developers/docs/networking/fabric/",
			},
			"nni_vlan": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Nni VLAN parameter, see https://metal.equinix.com/developers/docs/networking/fabric/",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the projct to which the virtual circuit belongs",
			},
		},
	}
}

func dataSourceMetalVirtualCircuitRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	vcId := d.Get("virtual_circuit_id").(string)

	vc, _, err := client.VirtualCircuits.Get(
		vcId,
		&packngo.GetOptions{Includes: []string{"project"}})
	if err != nil {
		return err
	}

	d.Set("virtual_circuit_id", vc.ID)
	d.Set("name", vc.Name)
	d.Set("status", vc.Status)
	d.Set("vnid", vc.VNID)
	d.Set("nni_vnid", vc.NniVNID)
	d.Set("nni_vlan", vc.NniVLAN)
	d.Set("project_id", vc.Project.ID)
	d.Set("description", vc.Description)
	tags := d.Get("tags.#").(int)
	if tags > 0 {
		vc.Tags = convertStringArr(d.Get("tags").([]interface{}))
	}
	d.SetId(vc.ID)

	return nil
}
