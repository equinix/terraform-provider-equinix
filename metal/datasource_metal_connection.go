package metal

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func connectionPortSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"speed": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"link_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_circuit_ids": {
				Computed: true,
				Type:     schema.TypeList,
				Elem:     schema.TypeString,
			},
		},
	}
}

func dataSourceMetalConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalConnectionRead,

		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Connection to lookup",
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"redundancy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"facility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"speed": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ports": {
				Type:     schema.TypeList,
				Elem:     connectionPortSchema(),
				Computed: true,
			},
		},
	}
}

func getConnectionPorts(cps []packngo.ConnectionPort) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0, 1)

	for _, p := range cps {
		vcIDs := []string{}
		for _, vc := range p.VirtualCircuits {
			vcIDs = append(vcIDs, vc.ID)
		}
		connPort := map[string]interface{}{
			"name":                p.Name,
			"role":                string(p.Role),
			"speed":               p.Speed,
			"status":              p.Status,
			"link_status":         p.LinkStatus,
			"virtual_circuit_ids": vcIDs,
		}
		ret = append(ret, connPort)
	}
	return ret
}

func dataSourceMetalConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	connId := d.Get("connection_id").(string)

	conn, _, err := client.Connections.Get(
		connId,
		&packngo.GetOptions{Includes: []string{"organization","facility"}})
	if err != nil {
		return err
	}

	d.Set("connection_id", conn.ID)
	d.Set("organization_id", conn.Organization.ID)
	d.Set("name", conn.Name)
	d.Set("description", conn.Description)
	d.Set("status", conn.Status)
	d.Set("redundancy", conn.Redundancy)
	d.Set("facility", conn.Facility.Code)
	d.Set("token", conn.Token)
	d.Set("type", conn.Type)
	d.Set("speed", conn.Speed)
	d.Set("ports", getConnectionPorts(conn.Ports))
	d.SetId(conn.ID)

	return nil
}
