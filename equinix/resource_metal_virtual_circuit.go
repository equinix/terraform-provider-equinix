package equinix

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func resourceMetalVirtualCircuit() *schema.Resource {
	return &schema.Resource{
		Read:   resourceMetalVirtualCircuitRead,
		Create: resourceMetalVirtualCircuitCreate,
		Update: resourceMetalVirtualCircuitUpdate,
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
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Virtual Circuit resource",
			},
			"speed": {
				Type:        schema.TypeString,
				Description: "Description of the Virtual Circuit speed. This is for information purposes and is computed when the connection type is shared.",
				Optional:    true,
				Computed:    true,
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "Tags attached to the virtual circuit",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
			"vrf_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "UUID of the VLAN to associate",
				ConflictsWith: []string{"vnid"},
				ForceNew:      true,
			},
			"peer_as": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"vrf_id"},
				Description: `A subnet from one of the IP blocks associated with the VRF that we will help create an IP reservation for. Can only be either a /30 or /31.
				 * For a /31 block, it will only have two IP addresses, which will be used for the metal_ip and customer_ip.
				 * For a /30 block, it will have four IP addresses, but the first and last IP addresses are not usable. We will default to the first usable IP address for the metal_ip`,
				ForceNew: true,
			},
			"subnet": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"vrf_id"},
				Description: `A subnet from one of the IP blocks associated with the VRF that we will help create an IP reservation for. Can only be either a /30 or /31.
				 * For a /31 block, it will only have two IP addresses, which will be used for the metal_ip and customer_ip.
				 * For a /30 block, it will have four IP addresses, but the first and last IP addresses are not usable. We will default to the first usable IP address for the metal_ip.`,
			},
			"metal_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"vrf_id"},
				Description:  "The IP address that’s set as “our” IP that is configured on the rack_local_vlan SVI. Will default to the first usable IP in the subnet.",
			},
			"customer_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"vrf_id"},
				Description:  "The IP address set as the customer IP which the CSR switch will peer with. Will default to the other usable IP in the subnet.",
			},
			"md5": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The password that can be set for the VRF BGP peer",
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
	client := meta.(*Config).metal
	vncr := packngo.VCCreateRequest{
		VirtualNetworkID: d.Get("vlan_id").(string),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Speed:            d.Get("speed").(string),
		VRFID:            d.Get("vrf_id").(string),
		PeerAs:           d.Get("peer_as").(int),
		Subnet:           d.Get("subnet").(string),
		MetalIP:          d.Get("metal_ip").(string),
		CustomerIP:       d.Get("customer_ip").(string),
		MD5:              d.Get("md5").(string),
	}

	connId := d.Get("connection_id").(string)
	portId := d.Get("port_id").(string)
	projectId := d.Get("project_id").(string)

	tags := d.Get("tags.#").(int)
	if tags > 0 {
		vncr.Tags = convertStringArr(d.Get("tags").([]interface{}))
	}

	if nniVlan, ok := d.GetOk("nni_vlan"); ok {
		vncr.NniVLAN = nniVlan.(int)
	}
	conn, _, err := client.Connections.Get(connId, nil)
	if err != nil {
		return err
	}
	if conn.Status == string(packngo.VCStatusPending) {
		return fmt.Errorf("Connection request with name %s and ID %s wasn't approved yet", conn.Name, conn.ID)
	}

	vc, _, err := client.VirtualCircuits.Create(projectId, connId, portId, &vncr, nil)
	if err != nil {
		log.Printf("[DEBUG] Error creating virtual circuit: %s", err)
		return err
	}
	// TODO: offer to wait while VCStatusPending
	createWaiter := getVCStateWaiter(
		client,
		vc.ID,
		d.Timeout(schema.TimeoutCreate),
		[]string{string(packngo.VCStatusActivating)},
		[]string{string(packngo.VCStatusActive)},
	)

	_, err = createWaiter.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for virtual circuit %s to be created: %s", vc.ID, err.Error())
	}

	d.SetId(vc.ID)

	return resourceMetalVirtualCircuitRead(d, meta)
}

func resourceMetalVirtualCircuitRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metal
	vcId := d.Id()

	vc, _, err := client.VirtualCircuits.Get(
		vcId,
		&packngo.GetOptions{Includes: []string{"project", "port", "virtual_network", "vrf"}},
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
		"vrf_id": func(d *schema.ResourceData, k string) error {
			if vc.VRF != nil {
				return d.Set(k, vc.VRF.ID)
			}
			return nil
		},
		"status":      vc.Status,
		"nni_vlan":    vc.NniVLAN,
		"vnid":        vc.VNID,
		"nni_vnid":    vc.NniVNID,
		"name":        vc.Name,
		"speed":       vc.Speed,
		"description": vc.Description,
		"tags":        vc.Tags,
		"peer_as":     vc.PeerAs,
		"subnet":      vc.Subnet,
		"metal_ip":    vc.MetalIP,
		"customer_ip": vc.CustomerIP,
		"md5":         vc.MD5,
	})
}

func getVCStateWaiter(client *packngo.Client, id string, timeout time.Duration, pending, target []string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			vc, _, err := client.VirtualCircuits.Get(
				id,
				&packngo.GetOptions{Includes: []string{
					"project", "port", "virtual_network",
					"vrf",
				}}, // TODO: we are not using the returned VC. Remove the includes?
			)
			if err != nil {
				return 0, "", err
			}
			return vc, string(vc.Status), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}

func resourceMetalVirtualCircuitUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metal

	ur := packngo.VCUpdateRequest{}
	if d.HasChange("vnid") {
		vnid := d.Get("vnid").(string)
		ur.VirtualNetworkID = &vnid
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		ur.Name = &name
	}

	if d.HasChange("description") {
		desc := d.Get("description").(string)
		ur.Description = &desc
	}

	if d.HasChange("speed") {
		speed := d.Get("speed").(string)
		ur.Speed = speed
	}

	if d.HasChange("tags") {
		ts := d.Get("tags")
		sts := []string{}

		switch ts.(type) {
		case []interface{}:
			for _, v := range ts.([]interface{}) {
				sts = append(sts, v.(string))
			}
			ur.Tags = &sts
		default:
			return friendlyError(fmt.Errorf("garbage in tags: %s", ts))
		}
	}

	if !reflect.DeepEqual(ur, packngo.VCUpdateRequest{}) {
		if _, _, err := client.VirtualCircuits.Update(d.Id(), &ur, nil); err != nil {
			return friendlyError(err)
		}
	}
	return resourceMetalVirtualCircuitRead(d, meta)
}

func resourceMetalVirtualCircuitDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metal
	// we first need to disconnect VLAN from the VC
	empty := ""
	_, _, err := client.VirtualCircuits.Update(
		d.Id(),
		&packngo.VCUpdateRequest{VirtualNetworkID: &empty},
		nil,
	)
	if err != nil {
		return err
	}

	detachWaiter := getVCStateWaiter(
		client,
		d.Id(),
		d.Timeout(schema.TimeoutDelete),
		[]string{string(packngo.VCStatusDeactivating)},
		[]string{string(packngo.VCStatusWaiting)},
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
