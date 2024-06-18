package virtual_circuit

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMetalVirtualCircuitRead,

		Schema: map[string]*schema.Schema{
			"virtual_circuit_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the virtual circuit to lookup",
			},
			"connection_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of Connection where the VC is scoped to",
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
			"port_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the Connection Port where the VC is scoped to",
			},
			"speed": {
				Type:        schema.TypeString,
				Description: "Description of the Virtual Circuit speed. This is for information purposes and is computed when the connection type is shared.",
				Computed:    true,
			},
			"vlan_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the associated VLAN",
			},
			"vrf_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the associated VRF",
			},
			"peer_asn": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The BGP ASN of the peer. The same ASN may be the used across several VCs, but it cannot be the same as the local_asn of the VRF.",
			},
			"subnet": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `A subnet from one of the IP blocks associated with the VRF that we will help create an IP reservation for. Can only be either a /30 or /31.
				 * For a /31 block, it will only have two IP addresses, which will be used for the metal_ip and customer_ip.
				 * For a /30 block, it will have four IP addresses, but the first and last IP addresses are not usable. We will default to the first usable IP address for the metal_ip.`,
			},
			"metal_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Metal IP address for the SVI (Switch Virtual Interface) of the VirtualCircuit. Will default to the first usable IP in the subnet.",
			},
			"customer_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Customer IP address which the CSR switch will peer with. Will default to the other usable IP in the subnet.",
			},
			"md5": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The password that can be set for the VRF BGP peer",
			},
		},
	}
}

func dataSourceMetalVirtualCircuitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vcId := d.Get("virtual_circuit_id").(string)
	d.SetId(vcId)
	return resourceMetalVirtualCircuitRead(ctx, d, meta)
}
