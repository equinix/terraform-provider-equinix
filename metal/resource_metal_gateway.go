package metal

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

var subnetSizes = []int{8, 16, 32, 64, 128}

func intInSlice(valid []int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		for _, val := range valid {
			if v == val {
				return
			}
		}

		es = append(es, fmt.Errorf("expected %s to be one of %v, got %d", k, valid, v))
		return
	}
}

func resourceMetalGateway() *schema.Resource {

	return &schema.Resource{
		Read:   resourceMetalGatewayRead,
		Create: resourceMetalGatewayCreate,
		Delete: resourceMetalGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the Project where the Gateway is scoped to",
				ForceNew:    true,
			},
			"vlan_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the VLAN to associate",
				ForceNew:    true,
			},
			"ip_reservation_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "UUID of the IP Reservation to associate, must be in the same metro as the VLAN",
				ConflictsWith: []string{"private_ipv4_subnet_size"},
				ForceNew:      true,
			},
			"private_ipv4_subnet_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   fmt.Sprintf("Size of the private IPv4 subnet to create for this gateway, one of %v", subnetSizes),
				ConflictsWith: []string{"ip_reservation_id"},
				ValidateFunc:  intInSlice(subnetSizes),
				ForceNew:      true,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the virtual circuit resource",
			},
		},
	}
}

func resourceMetalGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	_, hasIPReservation := d.GetOk("ip_reservation_id")
	_, hasSubnetSize := d.GetOk("private_ipv4_subnet_size")
	if !(hasIPReservation || hasSubnetSize) {
		return fmt.Errorf("You must set either ip_reservation_id or private_ipv4_subnet_size")
	}

	mgcr := packngo.MetalGatewayCreateRequest{
		VirtualNetworkID:      d.Get("vlan_id").(string),
		IPReservationID:       d.Get("ip_reservation_id").(string),
		PrivateIPv4SubnetSize: d.Get("private_ipv4_subnet_size").(int),
	}
	projectId := d.Get("project_id").(string)

	mg, _, err := client.MetalGateways.Create(projectId, &mgcr)
	if err != nil {
		return err
	}
	createWaiter := getMetalGatewayStateWaiter(
		client,
		mg.ID,
		d.Timeout(schema.TimeoutCreate),
		[]string{string(packngo.MetalGatewayActivating)},
		[]string{string(packngo.MetalGatewayActive)},
	)

	_, err = createWaiter.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for metal gateway %s to be created: %s", mg.ID, err.Error())
	}

	d.SetId(mg.ID)

	return resourceMetalGatewayRead(d, meta)
}

func resourceMetalGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	mgId := d.Id()

	includes := &packngo.GetOptions{Includes: []string{"project", "ip_reservation", "virtual_network"}}
	mg, _, err := client.MetalGateways.Get(mgId, includes)
	if err != nil {
		return err
	}

	privateIPv4SubnetSize := uint(0)
	if !mg.IPReservation.Public {
		//privateIPv4SubnetSize = bits.RotateLeft(1, 32-mg.IPReservation.CIDR)
		privateIPv4SubnetSize = 1 << (32 - mg.IPReservation.CIDR)
	}

	return setMap(d, map[string]interface{}{
		"project_id":               mg.Project.ID,
		"vlan_id":                  mg.VirtualNetwork.ID,
		"ip_reservation_id":        mg.IPReservation.ID,
		"private_ipv4_subnet_size": int(privateIPv4SubnetSize),
		"state":                    mg.State,
	})
}

func getMetalGatewayStateWaiter(client *packngo.Client, id string, timeout time.Duration, pending, target []string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			mg, _, err := client.MetalGateways.Get(
				id,
				&packngo.GetOptions{Includes: []string{"project", "port", "virtual_network"}},
			)
			if err != nil {
				return 0, "", err
			}
			return mg, string(mg.State), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

}

func resourceMetalGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	resp, err := client.MetalGateways.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}
	return nil
}
