package metal

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func resourceMetalVirtualCircuit() *schema.Resource {
	return &schema.Resource{
		Read:   resourceMetalVirtualCircuitRead,
		Create: resourceMetalVirtualCircuitCreate,
		Delete: resourceMetalVirtualCircuitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of Connection where the VC is scoped to",
				ForceNew:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the Project where the VC is scoped to",
				ForceNew:    true,
			},
			"port_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the Connection Port where the VC is scoped to",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the Virtual Circuit resource",
				ForceNew:    true,
			},
			"nni_vlan": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Equinix Metal network-to-network VLAN ID",
				ForceNew:    true,
			},
			"vlan_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the VLAN resource",
				ForceNew:    true,
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
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the virtual circuit resource",
			},
		},
	}
}

func resourceMetalVirtualCircuitCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	vncr := packngo.VCCreateRequest{
		VirtualNetworkID: d.Get("vnid").(string),
		NniVLAN:          d.Get("nni_vlan").(int),
		Name:             d.Get("name").(string),
	}
	connId := d.Get("connection_id").(string)
	portId := d.Get("port_id").(string)
	projectId := d.Get("project_id").(string)

	vc, _, err := client.VirtualCircuits.Create(projectId, connId, portId, &vncr, nil)
	if err != nil {
		return err
	}

	d.SetId(vc.ID)

	return resourceMetalVirtualCircuitRead(d, meta)
}

func resourceMetalVirtualCircuitRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	vcId := d.Id()

	vc, _, err := client.VirtualCircuits.Get(
		vcId,
		&packngo.GetOptions{Includes: []string{"project", "port", "virtual_network"}},
	)
	if err != nil {
		return err
	}

	return setMap(d, map[string]interface{}{
		//"connection_id": vc.Connection.ID,
		"project_id": vc.Project.ID,
		//"port_id":       vc.Port.ID,
		"vlan_id":  vc.VirtualNetwork.ID,
		"status":   vc.Status,
		"nni_vlan": vc.NniVLAN,
		"vnid":     vc.VNID,
		"nni_vnid": vc.NniVNID,
		"name":     vc.Name,
	})
}

func resourceMetalVirtualCircuitDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	resp, err := client.VirtualCircuits.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}
	return nil
}
