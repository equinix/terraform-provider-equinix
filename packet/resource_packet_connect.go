// WARNING: Packet Connect has been deprecated, and will be removed in release 2.7.0.

package packet

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketConnect() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketConnectCreate,
		Read:   resourcePacketConnectRead,
		Delete: resourcePacketConnectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"facility": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_speed": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"provider_payload": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"vxlan": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		DeprecationMessage: "packet_connect has now been deprecated, and will be removed in release 2.7.0.",
	}
}

func waitForConnectStatus(d *schema.ResourceData, target string, pending string, meta interface{}) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{pending},
		Target:     []string{target},
		Refresh:    connectRefreshFunc(d, meta),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func connectRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*packngo.Client)

	return func() (interface{}, string, error) {
		if err := resourcePacketConnectRead(d, meta); err != nil {
			return nil, "", err
		}

		if status, ok := d.GetOk("status"); ok {
			projectId := d.Get("project_id").(string)
			c, _, err := client.Connects.Get(d.Id(), projectId, nil)
			if err != nil {
				return nil, "", friendlyError(err)
			}
			return c, status.(string), nil
		}

		return nil, "", nil
	}
}

func resourcePacketConnectCreate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("packet_connect is deprecated in provider version 2.7.0")
}

func resourcePacketConnectRead(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("packet_connect is deprecated in provider version 2.7.0")
}

func resourcePacketConnectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	pc, _, err := client.Connects.Deprovision(d.Id(), d.Get("project_id").(string), false)
	if err != nil {
		return friendlyError(err)
	}
	_, err = waitForConnectStatus(d, "DEPROVISIONED", "DEPROVISIONING", meta)
	if err != nil {
		return friendlyError(err)
	}

	_, err = client.Connects.Delete(d.Id(), pc.ProjectID)
	if err != nil {
		return friendlyError(err)
	}

	return nil
}
