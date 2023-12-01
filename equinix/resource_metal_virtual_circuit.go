package equinix

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func resourceMetalVirtualCircuit() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout:   diagnosticsWrapper(resourceMetalVirtualCircuitRead),
		CreateContext:        diagnosticsWrapper(resourceMetalVirtualCircuitCreate),
		UpdateWithoutTimeout: diagnosticsWrapper(resourceMetalVirtualCircuitUpdate),
		DeleteContext:        diagnosticsWrapper(resourceMetalVirtualCircuitDelete),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				// TODO: implement SuppressDiffFunc for input with units to bps without units
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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "UUID of the VLAN to associate",
				ExactlyOneOf: []string{"vlan_id", "vrf_id"},
			},
			"vrf_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "UUID of the VRF to associate",
				ExactlyOneOf: []string{"vlan_id", "vrf_id"},
				ForceNew:     true,
			},
			"peer_asn": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"vrf_id"},
				Description:  "The BGP ASN of the peer. The same ASN may be the used across several VCs, but it cannot be the same as the local_asn of the VRF.",
				ForceNew:     true,
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
				Description:  "The Metal IP address for the SVI (Switch Virtual Interface) of the VirtualCircuit. Will default to the first usable IP in the subnet.",
			},
			"customer_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"vrf_id"},
				Description:  "The Customer IP address which the CSR switch will peer with. Will default to the other usable IP in the subnet.",
			},
			"md5": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
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

func resourceMetalVirtualCircuitCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal
	vncr := packngo.VCCreateRequest{
		VirtualNetworkID: d.Get("vlan_id").(string),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Speed:            d.Get("speed").(string),
		VRFID:            d.Get("vrf_id").(string),
		PeerASN:          d.Get("peer_asn").(int),
		Subnet:           d.Get("subnet").(string),
		MetalIP:          d.Get("metal_ip").(string),
		CustomerIP:       d.Get("customer_ip").(string),
		MD5:              d.Get("md5").(string),
	}

	connId := d.Get("connection_id").(string)
	portId := d.Get("port_id").(string)
	projectId := d.Get("project_id").(string)

	


	vc, _, err := client.VirtualCircuits.Create(projectId, connId, portId, &vncr, nil)
	if err != nil {
		log.Printf("[DEBUG] Error creating virtual circuit: %s", err)
		return err
	}
	// TODO: offer to wait while VCStatusPending
	createWaiter := getVCStateWaiter(
		client,
		vc.ID,
		d.Timeout(schema.TimeoutCreate)-30*time.Second,
		[]string{string(packngo.VCStatusActivating)},
		[]string{string(packngo.VCStatusActive)},
	)

	_, err = createWaiter.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("Error waiting for virtual circuit %s to be created: %s", vc.ID, err.Error())
	}

	d.SetId(vc.ID)

	return resourceMetalVirtualCircuitRead(ctx, d, meta)
}

func resourceMetalVirtualCircuitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal
	vcId := d.Id()

	vc, _, err := client.VirtualCircuits.Get(
		vcId,
		&packngo.GetOptions{Includes: []string{"project", "virtual_network", "vrf"}},
	)
	if err != nil {
		return err
	}

	// TODO: use API field from VC responses when available The regexp is
	// optimistic, not guaranteed. This affects resource imports. "port" is not
	// in the Includes above to assure the Href needed below.
	connectionID := "" // vc.Connection.ID is not available yet
	portID := ""       // vc.Port.ID would be available with ?include=port
	connectionRe := regexp.MustCompile("/connections/([0-9a-z-]+)/ports/([0-9a-z-]+)")
	matches := connectionRe.FindStringSubmatch(vc.Port.Href.Href)
	if len(matches) == 3 {
		connectionID = matches[1]
		portID = matches[2]
	} else {
		log.Printf("[DEBUG] Could not parse connection and port ID from port href %s", vc.Port.Href.Href)
	}

	return setMap(d, map[string]interface{}{
		"project_id": vc.Project.ID,
		"port_id":    portID,
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
		"speed":       strconv.Itoa(vc.Speed),
		"description": vc.Description,
		"tags":        vc.Tags,
		"peer_asn":    vc.PeerASN,
		"subnet":      vc.Subnet,
		"metal_ip":    vc.MetalIP,
		"customer_ip": vc.CustomerIP,
		"md5":         vc.MD5,
		"connection_id": func(d *schema.ResourceData, k string) error {
			if connectionID != "" {
				return d.Set(k, connectionID)
			}
			return nil
		},
	})
}

func getVCStateWaiter(client *packngo.Client, id string, timeout time.Duration, pending, target []string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			vc, _, err := client.VirtualCircuits.Get(id, nil)
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

func resourceMetalVirtualCircuitUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	ur := packngo.VCUpdateRequest{}
	if d.HasChange("vlan_id") {
		vnid := d.Get("vlan_id").(string)
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
	return resourceMetalVirtualCircuitRead(ctx, d, meta)
}

func resourceMetalVirtualCircuitDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	resp, err := client.VirtualCircuits.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}

	deleteWaiter := getVCStateWaiter(
		client,
		d.Id(),
		d.Timeout(schema.TimeoutDelete)-30*time.Second,
		[]string{string(packngo.VCStatusDeleting)},
		[]string{},
	)

	_, err = deleteWaiter.WaitForStateContext(ctx)
	if ignoreResponseErrors(httpForbidden, httpNotFound)(nil, err) != nil {
		return fmt.Errorf("Error deleting virtual circuit %s: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
