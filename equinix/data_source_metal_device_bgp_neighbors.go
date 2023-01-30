package equinix

import (
	metalv1 "github.com/equinix-labs/metal-go/metal/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func bgpNeighborSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"address_family": {
				Type:        schema.TypeInt,
				Description: "IP address version, 4 or 6",
				Computed:    true,
			},
			"customer_as": {
				Type:        schema.TypeInt,
				Description: "Local autonomous system number",
				Computed:    true,
			},
			"customer_ip": {
				Type:        schema.TypeString,
				Description: "Local used peer IP address",
				Computed:    true,
			},
			"md5_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether BGP session is password enabled",
				Computed:    true,
			},
			"md5_password": {
				Type:        schema.TypeString,
				Description: "BGP session password in plaintext (not a checksum)",
				Computed:    true,
				Sensitive:   true,
			},
			"multihop": {
				Type:        schema.TypeBool,
				Description: "Whether the neighbor is in EBGP multihop session",
				Computed:    true,
			},
			"peer_as": {
				Type:        schema.TypeInt,
				Description: "Peer AS number (different than customer_as for EBGP)",
				Computed:    true,
			},
			"peer_ips": {
				Type:        schema.TypeList,
				Description: "Array of IP addresses of this neighbor's peers",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"routes_in": {
				Type:        schema.TypeList,
				Description: "Array of incoming routes. Each route has attributes",
				Computed:    true,
				Elem:        bgpRouteSchema(),
			},
			"routes_out": {
				Type:        schema.TypeList,
				Description: "Array of outgoing routes in the same format",
				Computed:    true,
				Elem:        bgpRouteSchema(),
			},
		},
	}
}

func bgpRouteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"route": {
				Type:        schema.TypeString,
				Description: "CIDR expression of route (IP/mask)",
				Computed:    true,
			},
			"exact": {
				Type:        schema.TypeBool,
				Description: "Whether the route is exact",
				Computed:    true,
			},
		},
	}
}

func dataSourceMetalDeviceBGPNeighbors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalDeviceBGPNeighborsRead,
		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:        schema.TypeString,
				Description: "UUID of BGP-enabled device whose neighbors to list",
				Required:    true,
			},
			"bgp_neighbors": {
				Type:        schema.TypeList,
				Description: "Array of BGP neighbor records",
				Computed:    true,
				Elem:        bgpNeighborSchema(),
			},
		},
	}
}

func dataSourceMetalDeviceBGPNeighborsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metalgo
	deviceID := d.Get("device_id").(string)

	bgpNeighborsRaw, _, err := client.DevicesApi.GetBgpNeighborData(nil, deviceID).Execute()
	if err != nil {
		return err
	}

	d.Set("bgp_neighbors", getBgpNeighbors(bgpNeighborsRaw))
	d.SetId(deviceID)
	return nil
}

func getRoutesSlice(routes []metalv1.BgpRoute) []map[string]interface{} {
	ret := []map[string]interface{}{}
	for _, r := range routes {
		ret = append(ret, map[string]interface{}{
			"route": r.Route, "exact": r.Exact,
		})
	}
	return ret
}

func getBgpNeighbors(ns *metalv1.BgpSessionNeighbors) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0, 1)
	for _, n := range ns.BgpNeighbors {
		neighbor := map[string]interface{}{
			"address_family": n.AddressFamily,
			"customer_as":    n.CustomerAs,
			"customer_ip":    n.CustomerIp,
			"md5_enabled":    n.Md5Enabled,
			"md5_password":   n.Md5Password,
			"multihop":       n.Multihop,
			"peer_as":        n.PeerAs,
			"peer_ips":       n.PeerIps,
			"routes_in":      getRoutesSlice(n.RoutesIn),
			"routes_out":     getRoutesSlice(n.RoutesOut),
		}
		ret = append(ret, neighbor)
	}
	return ret
}
