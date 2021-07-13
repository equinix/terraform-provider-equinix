package metal

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Description: "Equinix Metal network-to-network VLAN ID (optional when the connection has mode=tunnel)",
				Optional:    true,
				ForceNew:    true,
			},
			"vlan_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the VLAN to associate",
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
		VirtualNetworkID: d.Get("vlan_id").(string),
		Name:             d.Get("name").(string),
	}

	connId := d.Get("connection_id").(string)
	portId := d.Get("port_id").(string)
	projectId := d.Get("project_id").(string)

	if nniVlan, ok := d.GetOk("nni_vlan"); ok {
		vncr.NniVLAN = nniVlan.(int)
	}

	conn, _, err := client.Connections.Get(connId, nil)
	if err != nil {
		return err
	}
	if conn.Status == packngo.VCStatusPending {
		return fmt.Errorf("Connection request with name %s and ID %s wasn't approved yet", conn.Name, conn.ID)
	}

	vc, _, err := client.VirtualCircuits.Create(projectId, connId, portId, &vncr, nil)
	if err != nil {
		return err
	}
	createWaiter := getVCStateWaiter(
		client,
		vc.ID,
		d.Timeout(schema.TimeoutCreate),
		[]string{packngo.VCStatusActivating},
		[]string{packngo.VCStatusActive},
	)

	_, err = createWaiter.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for virtual circuit %s to be created: %s", vc.ID, err.Error())
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
		"port_id":    vc.Port.ID,
		"vlan_id": func(d *schema.ResourceData, k string) error {
			if vc.VirtualNetwork != nil {
				return d.Set(k, vc.VirtualNetwork.ID)
			}
			return nil
		},
		"status":   vc.Status,
		"nni_vlan": vc.NniVLAN,
		"vnid":     vc.VNID,
		"nni_vnid": vc.NniVNID,
		"name":     vc.Name,
	})
}

func getVCStateWaiter(client *packngo.Client, id string, timeout time.Duration, pending, target []string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			vc, _, err := client.VirtualCircuits.Get(
				id,
				&packngo.GetOptions{Includes: []string{"project", "port", "virtual_network"}},
			)
			if err != nil {
				return 0, "", err
			}
			return vc, vc.Status, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

}

func resourceMetalVirtualCircuitDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	// we first need to disconnect VLAN from the VC
	_, _, err := client.VirtualCircuits.Update(
		d.Id(),
		&packngo.VCUpdateRequest{VirtualNetworkID: nil},
		nil,
	)
	if err != nil {
		return err
	}

	detachWaiter := getVCStateWaiter(
		client,
		d.Id(),
		d.Timeout(schema.TimeoutDelete),
		[]string{packngo.VCStatusDeactivating},
		[]string{packngo.VCStatusWaiting},
	)

	_, err = detachWaiter.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting virtual circuit %s: %s", d.Id(), err)
	}

	resp, err := client.VirtualCircuits.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}
	return nil
}
