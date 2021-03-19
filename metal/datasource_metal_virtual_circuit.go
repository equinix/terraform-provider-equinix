package metal

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalVirtualCircuit() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalVirtualCircuitRead,

		Schema: map[string]*schema.Schema{
			"virtual_circuit_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the VirtualCircuit to lookup",
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vnid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nni_vnid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nni_vlan": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
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
	d.SetId(vc.ID)

	return nil
}
