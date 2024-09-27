package virtual_circuit

import (
	"context"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/converters"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Description: `Use this resource to associate VLAN with a Dedicated Port from [Equinix Fabric - software-defined interconnections](https://deploy.equinix.com/developers/docs/metal/interconnections/introduction/#associating-a-vlan-with-a-dedicated-port).

See the [Virtual Routing and Forwarding documentation](https://deploy.equinix.com/developers/docs/metal/layer2-networking/vrf/) for product details and API reference material.`,
		ReadWithoutTimeout:   resourceMetalVirtualCircuitRead,
		CreateContext:        resourceMetalVirtualCircuitCreate,
		UpdateWithoutTimeout: resourceMetalVirtualCircuitUpdate,
		DeleteContext:        resourceMetalVirtualCircuitDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"connection_id", "virtual_circuit_id"},
				Description:  "UUID of Connection where the VC is scoped to.  Only used for dedicated connections",
				ForceNew:     true,
			},
			"virtual_circuit_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "UUID of an existing VC to configure. Used in the case of shared interconnections where the VC has already been created.",
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
				Type:          schema.TypeInt,
				Description:   "Equinix Metal network-to-network VLAN ID (optional when the connection has mode=tunnel)",
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"virtual_circuit_id"},
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Metal IP address for the SVI (Switch Virtual Interface) of the VirtualCircuit. Will default to the first usable IP in the subnet.",
				Computed:    true,
			},
			"customer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Customer IP address which the CSR switch will peer with. Will default to the other usable IP in the subnet.",
				Computed:    true,
			},
			"subnet_ipv6": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `A subnet from one of the IPv6 blocks associated with the VRF that we will help create an IP reservation for. Can only be either a /126 or /127.
				 * For a /127 block, it will only have two IP addresses, which will be used for the metal_ip and customer_ip.
				 * For a /126 block, it will have four IP addresses, but the first and last IP addresses are not usable. We will default to the first usable IP address for the metal_ip.`,
			},
			"metal_ipv6": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Metal IPv6 address for the SVI (Switch Virtual Interface) of the VirtualCircuit. Will default to the first usable IP in the IPv6 subnet.",
				Computed:    true,
			},
			"customer_ipv6": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Customer IPv6 address which the CSR switch will peer with. Will default to the other usable IP in the IPv6 subnet.",
				Computed:    true,
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
				Description: "VNID VLAN parameter, see https://deploy.equinix.com/developers/docs/metal/interconnections/introduction/",
			},
			"nni_vnid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Nni VLAN ID parameter, see https://deploy.equinix.com/developers/docs/metal/interconnections/introduction/",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the virtual circuit resource",
			},
		},
	}
}

func resourceMetalVirtualCircuitCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if _, ok := d.GetOk("virtual_circuit_id"); ok {
		vcId := d.Get("virtual_circuit_id").(string)
		d.SetId(vcId)
		return resourceMetalVirtualCircuitUpdate(ctx, d, meta)
	} else {
		client := meta.(*config.Config).NewMetalClientForSDK(d)
		vncr := metalv1.VirtualCircuitCreateInput{}

		connId := d.Get("connection_id").(string)
		portId := d.Get("port_id").(string)
		projectId := d.Get("project_id").(string)
		name := d.Get("name").(string)

		if _, ok := d.GetOk("vlan_id"); ok {
			vncr.VlanVirtualCircuitCreateInput = &metalv1.VlanVirtualCircuitCreateInput{
				ProjectId:   projectId,
				Name:        &name,
				Description: metalv1.PtrString(d.Get("description").(string)),
				Speed:       metalv1.PtrString(d.Get("speed").(string)),
				Vnid:        metalv1.PtrString(d.Get("vlan_id").(string)),
			}
		} else {
			vncr.VrfVirtualCircuitCreateInput = &metalv1.VrfVirtualCircuitCreateInput{
				ProjectId:   projectId,
				Name:        &name,
				Description: metalv1.PtrString(d.Get("description").(string)),
				Speed:       metalv1.PtrString(d.Get("speed").(string)),
				Vrf:         d.Get("vrf_id").(string),
				// TODO: woof
				Md5:          *metalv1.NewNullableString(metalv1.PtrString(d.Get("md5").(string))),
				PeerAsn:      int64(d.Get("peer_asn").(int)),
				Subnet:       d.Get("subnet").(string),
				CustomerIp:   metalv1.PtrString(d.Get("customer_ip").(string)),
				MetalIp:      metalv1.PtrString(d.Get("metal_ip").(string)),
				SubnetIpv6:   metalv1.PtrString(d.Get("subnet_ipv6").(string)),
				CustomerIpv6: metalv1.PtrString(d.Get("customer_ipv6").(string)),
				MetalIpv6:    metalv1.PtrString(d.Get("metal_ipv6").(string)),
			}
		}

		tags := d.Get("tags.#").(int)
		if tags > 0 {
			if _, ok := d.GetOk("vlan_id"); ok {
				vncr.VlanVirtualCircuitCreateInput.Tags = converters.IfArrToStringArr(d.Get("tags").([]interface{}))
			} else {
				vncr.VrfVirtualCircuitCreateInput.Tags = converters.IfArrToStringArr(d.Get("tags").([]interface{}))
			}
		}

		if nniVlan, ok := d.GetOk("nni_vlan"); ok {
			if _, ok := d.GetOk("vlan_id"); ok {
				vncr.VlanVirtualCircuitCreateInput.NniVlan = metalv1.PtrInt32(int32(nniVlan.(int)))
			} else {
				vncr.VrfVirtualCircuitCreateInput.NniVlan = int32(nniVlan.(int))
			}
		}

		conn, _, err := client.InterconnectionsApi.GetInterconnection(ctx, connId).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		if conn.GetStatus() == string(metalv1.VLANVIRTUALCIRCUITSTATUS_PENDING) {
			return diag.Errorf("Connection request with name %s and ID %s wasn't approved yet", conn.GetName(), conn.GetId())
		}

		vc, _, err := client.InterconnectionsApi.CreateInterconnectionPortVirtualCircuit(ctx, connId, portId).VirtualCircuitCreateInput(vncr).Execute()
		if err != nil {
			log.Printf("[DEBUG] Error creating virtual circuit: %s", err)
			return diag.FromErr(err)
		}

		var vcId string

		if vc.VlanVirtualCircuit != nil {
			vcId = vc.VlanVirtualCircuit.GetId()
		} else {
			vcId = vc.VrfVirtualCircuit.GetId()
		}

		// TODO: offer to wait while VCStatusPending
		createWaiter := getVCStateWaiter(
			ctx,
			client,
			vcId,
			d.Timeout(schema.TimeoutCreate)-30*time.Second,
			[]string{string(metalv1.VLANVIRTUALCIRCUITSTATUS_ACTIVATING)},
			[]string{string(metalv1.VLANVIRTUALCIRCUITSTATUS_ACTIVE)},
		)

		_, err = createWaiter.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for virtual circuit %s to be created: %s", vcId, err.Error())
		}

		d.SetId(vcId)

		return resourceMetalVirtualCircuitRead(ctx, d, meta)
	}
}

func resourceMetalVirtualCircuitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)
	vcId := d.Id()

	vc, _, err := client.InterconnectionsApi.GetVirtualCircuit(ctx, vcId).
		Include([]string{"project", "virtual_network", "vrf"}).
		Execute()

	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: use API field from VC responses when available The regexp is
	// optimistic, not guaranteed. This affects resource imports. "port" is not
	// in the Includes above to assure the Href needed below.
	var portHref string

	if vc.VlanVirtualCircuit != nil {
		portHref = vc.VlanVirtualCircuit.Port.GetHref()
	} else {
		portHref = vc.VrfVirtualCircuit.Port.GetHref()
	}
	connectionID := "" // vc.Connection.ID is not available yet
	portID := ""       // vc.Port.ID would be available with ?include=port
	connectionRe := regexp.MustCompile("/connections/([0-9a-z-]+)/ports/([0-9a-z-]+)")
	matches := connectionRe.FindStringSubmatch(portHref)
	if len(matches) == 3 {
		connectionID = matches[1]
		portID = matches[2]
	} else {
		log.Printf("[DEBUG] Could not parse connection and port ID from port href %s", portHref)
	}
	var errs []error

	if connectionID != "" {
		errs = append(errs, d.Set("connection_id", connectionID))
	}
	d.Set("port_id", portID)

	if vc.VlanVirtualCircuit != nil {
		errs = append(errs, d.Set("project_id", vc.VlanVirtualCircuit.Project.GetId()))
		// TODO: blarg, spec has virtual network as Href, so these attrs arent directly available
		errs = append(errs, d.Set("vlan_id", vc.VlanVirtualCircuit.VirtualNetwork.AdditionalProperties["id"]))
		errs = append(errs, d.Set("status", vc.VlanVirtualCircuit.GetStatus()))
		errs = append(errs, d.Set("nni_vlan", vc.VlanVirtualCircuit.GetNniVlan()))
		errs = append(errs, d.Set("vnid", vc.VlanVirtualCircuit.GetVnid()))
		// TODO: this attribute isn't mentioned in the spec
		errs = append(errs, d.Set("nni_vnid", vc.VlanVirtualCircuit.AdditionalProperties["nni_vnid"]))
		errs = append(errs, d.Set("name", vc.VlanVirtualCircuit.GetName()))
		errs = append(errs, d.Set("speed", strconv.Itoa(int(vc.VlanVirtualCircuit.GetSpeed()))))
		errs = append(errs, d.Set("description", vc.VlanVirtualCircuit.GetDescription()))
		errs = append(errs, d.Set("tags", vc.VlanVirtualCircuit.GetTags()))
	} else {
		errs = append(errs, d.Set("project_id", vc.VrfVirtualCircuit.Project.GetId()))
		errs = append(errs, d.Set("vrf_id", vc.VrfVirtualCircuit.Vrf.GetId()))
		errs = append(errs, d.Set("status", vc.VrfVirtualCircuit.GetStatus()))
		errs = append(errs, d.Set("nni_vlan", vc.VrfVirtualCircuit.GetNniVlan()))
		// TODO: this attribute isn't mentioned in the spec
		errs = append(errs, d.Set("nni_vnid", vc.VrfVirtualCircuit.AdditionalProperties["nni_vnid"]))
		errs = append(errs, d.Set("name", vc.VrfVirtualCircuit.GetName()))
		errs = append(errs, d.Set("speed", strconv.Itoa(int(vc.VrfVirtualCircuit.GetSpeed()))))
		errs = append(errs, d.Set("description", vc.VrfVirtualCircuit.GetDescription()))
		errs = append(errs, d.Set("tags", vc.VrfVirtualCircuit.GetTags()))
		errs = append(errs, d.Set("peer_asn", vc.VrfVirtualCircuit.GetPeerAsn()))
		errs = append(errs, d.Set("subnet", vc.VrfVirtualCircuit.GetSubnet()))
		errs = append(errs, d.Set("metal_ip", vc.VrfVirtualCircuit.GetMetalIp()))
		errs = append(errs, d.Set("customer_ip", vc.VrfVirtualCircuit.GetCustomerIp()))
		errs = append(errs, d.Set("subnet_ipv6", vc.VrfVirtualCircuit.GetSubnetIpv6()))
		errs = append(errs, d.Set("metal_ipv6", vc.VrfVirtualCircuit.GetMetalIpv6()))
		errs = append(errs, d.Set("customer_ipv6", vc.VrfVirtualCircuit.GetCustomerIpv6()))

		errs = append(errs, d.Set("md5", vc.VrfVirtualCircuit.GetMd5()))
	}

	return diag.FromErr(errors.Join(errs...))
}

func getVCStateWaiter(ctx context.Context, client *metalv1.APIClient, id string, timeout time.Duration, pending, target []string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			vc, resp, err := client.InterconnectionsApi.GetVirtualCircuit(ctx, id).Execute()
			if err != nil {
				if resp != nil {
					// The resource delete function uses this waiter and relies
					// on it to return an ErrorResponse error so it can treat
					// a 404 as success.  This conversion is done here for now
					// to avoid a larger refactoring.
					err = equinix_errors.FriendlyErrorForMetalGo(err, resp)
				}
				return 0, "", err
			}
			vcStatus := ""
			if vc.VlanVirtualCircuit != nil {
				vcStatus = string(vc.VlanVirtualCircuit.GetStatus())
			} else {
				vcStatus = string(vc.VrfVirtualCircuit.GetStatus())
			}
			return vc, vcStatus, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}

func resourceMetalVirtualCircuitUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)
	needsUpdate := false

	ur := metalv1.VirtualCircuitUpdateInput{}

	if _, ok := d.GetOk("vlan_id"); ok {
		ur.VlanVirtualCircuitUpdateInput = &metalv1.VlanVirtualCircuitUpdateInput{}
		if d.HasChange("vlan_id") {
			needsUpdate = true
			vnid := d.Get("vlan_id").(string)
			ur.VlanVirtualCircuitUpdateInput.Vnid = &vnid
		}

		if d.HasChange("name") {
			needsUpdate = true
			name := d.Get("name").(string)
			ur.VlanVirtualCircuitUpdateInput.Name = &name
		}

		if d.HasChange("description") {
			needsUpdate = true
			desc := d.Get("description").(string)
			ur.VlanVirtualCircuitUpdateInput.Description = &desc
		}

		if d.HasChange("speed") {
			needsUpdate = true
			speed := d.Get("speed").(string)
			ur.VlanVirtualCircuitUpdateInput.Speed = &speed
		}

		if d.HasChange("tags") {
			needsUpdate = true
			ts := d.Get("tags")
			sts := []string{}

			switch ts.(type) {
			case []interface{}:
				for _, v := range ts.([]interface{}) {
					sts = append(sts, v.(string))
				}
				ur.VlanVirtualCircuitUpdateInput.Tags = sts
			default:
				return diag.Errorf("garbage in tags: %s", ts)
			}
		}
	} else {
		ur.VrfVirtualCircuitUpdateInput = &metalv1.VrfVirtualCircuitUpdateInput{}

		if d.HasChange("name") {
			needsUpdate = true
			name := d.Get("name").(string)
			ur.VrfVirtualCircuitUpdateInput.Name = &name
		}

		if d.HasChange("description") {
			needsUpdate = true
			desc := d.Get("description").(string)
			ur.VrfVirtualCircuitUpdateInput.Description = &desc
		}

		if d.HasChange("speed") {
			needsUpdate = true
			speed := d.Get("speed").(string)
			ur.VrfVirtualCircuitUpdateInput.Speed = &speed
		}

		if d.HasChange("tags") {
			needsUpdate = true
			ts := d.Get("tags")
			sts := []string{}

			switch ts.(type) {
			case []interface{}:
				for _, v := range ts.([]interface{}) {
					sts = append(sts, v.(string))
				}
				ur.VrfVirtualCircuitUpdateInput.Tags = sts
			default:
				return diag.Errorf("garbage in tags: %s", ts)
			}
		}

		if d.HasChange("subnet") {
			needsUpdate = true
			subnet := metalv1.PtrString(d.Get("subnet").(string))
			ur.VrfVirtualCircuitUpdateInput.Subnet = subnet
		}

		if d.HasChange("customer_ip") {
			needsUpdate = true
			customer_ip := metalv1.PtrString(d.Get("customer_ip").(string))
			ur.VrfVirtualCircuitUpdateInput.CustomerIp = customer_ip
		}

		if d.HasChange("metal_ip") {
			needsUpdate = true
			metal_ip := metalv1.PtrString(d.Get("metal_ip").(string))
			ur.VrfVirtualCircuitUpdateInput.MetalIp = metal_ip
		}

		if d.HasChange("subnet_ipv6") {
			needsUpdate = true
			subnet_ipv6 := metalv1.PtrString(d.Get("subnet_ipv6").(string))
			ur.VrfVirtualCircuitUpdateInput.SubnetIpv6 = subnet_ipv6
		}

		if d.HasChange("customer_ipv6") {
			needsUpdate = true
			customer_ipv6 := metalv1.PtrString(d.Get("customer_ipv6").(string))
			ur.VrfVirtualCircuitUpdateInput.CustomerIpv6 = customer_ipv6
		}

		if d.HasChange("metal_ipv6") {
			needsUpdate = true
			metal_ipv6 := metalv1.PtrString(d.Get("metal_ipv6").(string))
			ur.VrfVirtualCircuitUpdateInput.MetalIpv6 = metal_ipv6
		}
	}

	if needsUpdate {
		if _, _, err := client.InterconnectionsApi.UpdateVirtualCircuit(ctx, d.Id()).VirtualCircuitUpdateInput(ur).Execute(); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceMetalVirtualCircuitRead(ctx, d, meta)
}

func resourceMetalVirtualCircuitDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if _, ok := d.GetOk("virtual_circuit_id"); ok {
		return nil
	}

	client := meta.(*config.Config).NewMetalClientForSDK(d)

	_, resp, err := client.InterconnectionsApi.DeleteVirtualCircuit(ctx, d.Id()).Execute()
	if err != nil {
		if resp != nil {
			// equinix_error.HttpNotFound and similar do not short-circuit
			// based on response code, so we have to convert to a FriendlyError
			// in order to use existing checks for equinix_errors.IgnoreHttpResponseErrors
			err = equinix_errors.FriendlyErrorForMetalGo(err, resp)
		}
		if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err) != nil {
			return diag.FromErr(err)
		}
	}

	deleteWaiter := getVCStateWaiter(
		ctx,
		client,
		d.Id(),
		d.Timeout(schema.TimeoutDelete)-30*time.Second,
		[]string{string(metalv1.VLANVIRTUALCIRCUITSTATUS_DELETING)},
		[]string{},
	)

	_, err = deleteWaiter.WaitForStateContext(ctx)
	if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(nil, err) != nil {
		return diag.Errorf("Error deleting virtual circuit %s: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
