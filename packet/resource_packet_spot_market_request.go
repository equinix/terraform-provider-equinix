package packet

import (
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketSpotMarketRequest() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketSpotMarketRequestCreate,
		Read:   resourcePacketSpotMarketRequestRead,
		Delete: resourcePacketSpotMarketRequestDelete,

		Schema: map[string]*schema.Schema{
			"devices_min": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"devices_max": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"max_bid_price": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
				ForceNew: true,
			},
			"facilities": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"instance_parameters": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"billing_cycle": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"plan": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"operating_system": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"hostname": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"termintation_time": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"always_pxe": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"features": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"locked": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"project_ssh_keys": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"user_ssh_keys": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"userdata": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"wait_for_devices": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
		Timeouts: resourceDefaultTimeouts,
	}
}

func resourcePacketSpotMarketRequestCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	var waitForDevices bool

	facilitiesRaw := d.Get("facilities").([]interface{})
	facilities := []string{}

	for _, f := range facilitiesRaw {
		facilities = append(facilities, f.(string))
	}

	params := packngo.SpotMarketRequestInstanceParameters{
		Hostname:        d.Get("instance_parameters.0.hostname").(string),
		BillingCycle:    d.Get("instance_parameters.0.billing_cycle").(string),
		Plan:            d.Get("instance_parameters.0.plan").(string),
		OperatingSystem: d.Get("instance_parameters.0.operating_system").(string),
	}

	if val, ok := d.GetOk("instance_parameters.0.userdata"); ok {
		params.UserData = val.(string)
	}

	if val, ok := d.GetOk("instance_parameters.0.always_pxe"); ok {
		params.AlwaysPXE = val.(bool)
	}

	if val, ok := d.GetOk("instance_parameters.0.description"); ok {
		params.Description = val.(string)
	}

	if val, ok := d.GetOk("instance_parameters.0.features"); ok {
		temp := val.([]interface{})
		for _, i := range temp {
			if i != nil {
				params.Features = append(params.Features, i.(string))
			}
		}
	}

	if val, ok := d.GetOk("wait_for_devices"); ok {
		waitForDevices = val.(bool)
	}

	if val, ok := d.GetOk("instance_parameters.0.locked"); ok {
		params.Locked = val.(bool)
	}

	if val, ok := d.GetOk("instance_parameters.0.project_ssh_keys"); ok {
		temp := val.([]interface{})
		for _, i := range temp {
			if i != nil {
				params.ProjectSSHKeys = append(params.ProjectSSHKeys, i.(string))
			}
		}
	}

	if val, ok := d.GetOk("instance_parameters.0.user_ssh_keys"); ok {
		temp := val.([]interface{})
		for _, i := range temp {
			if i != nil {
				params.UserSSHKeys = append(params.UserSSHKeys, i.(string))
			}
		}
	}

	smrc := &packngo.SpotMarketRequestCreateRequest{
		DevicesMax:  d.Get("devices_max").(int),
		DevicesMin:  d.Get("devices_min").(int),
		MaxBidPrice: d.Get("max_bid_price").(float64),
		FacilityIDs: facilities,
		Parameters:  params,
	}

	smr, _, err := client.SpotMarketRequests.Create(smrc, d.Get("project_id").(string))
	if err != nil {
		return err
	}

	d.SetId(smr.ID)

	if waitForDevices {
		stateConf := &resource.StateChangeConf{
			Pending:        []string{"not_done"},
			Target:         []string{"done"},
			Refresh:        resourceStateRefreshFunc(d, meta),
			Timeout:        d.Timeout(schema.TimeoutCreate),
			MinTimeout:     5 * time.Second,
			Delay:          3 * time.Second, // Wait 10 secs before starting
			NotFoundChecks: 600,             //Setting high number, to support long timeouts
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return err
		}
	}

	return resourcePacketSpotMarketRequestRead(d, meta)
}

func resourcePacketSpotMarketRequestRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	smr, _, err := client.SpotMarketRequests.Get(d.Id(), &packngo.GetOptions{Includes: []string{"project", "devices", "facilities"}})
	if err != nil {
		return err
	}

	deviceIDs := make([]string, len(smr.Devices))
	for i, d := range smr.Devices {
		deviceIDs[i] = d.ID
	}

	facilityIDs := make([]string, len(smr.Facilities))
	if len(smr.Facilities) > 0 {
		for i, f := range smr.Facilities {
			facilityIDs[i] = f.ID
		}
	}
	d.Set("project_id", smr.Project.ID)

	return nil
}

func resourcePacketSpotMarketRequestDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	var waitForDevices bool

	if val, ok := d.GetOk("wait_for_devices"); ok {
		waitForDevices = val.(bool)
	}
	if waitForDevices {
		smr, _, err := client.SpotMarketRequests.Get(d.Id(), &packngo.GetOptions{Includes: []string{"project", "devices", "facilities"}})
		if err != nil {
			return nil
		}

		stateConf := &resource.StateChangeConf{
			Pending:        []string{"not_done"},
			Target:         []string{"done"},
			Refresh:        resourceStateRefreshFunc(d, meta),
			Timeout:        d.Timeout(schema.TimeoutDelete),
			MinTimeout:     5 * time.Second,
			Delay:          3 * time.Second, // Wait 10 secs before starting
			NotFoundChecks: 600,             //Setting high number, to support long timeouts
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return err
		}

		for _, d := range smr.Devices {
			_, err := client.Devices.Delete(d.ID)
			if err != nil {
				return err
			}
		}
	}
	_, err := client.SpotMarketRequests.Delete(d.Id(), true)
	if err != nil {
		return nil
	}
	return nil
}

func resourceStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*packngo.Client)
		smr, _, err := client.SpotMarketRequests.Get(d.Id(), &packngo.GetOptions{Includes: []string{"project", "devices", "facilities"}})

		if err != nil {
			return nil, "", err

		}
		var finished bool

		for _, d := range smr.Devices {

			dev, _, _ := client.Devices.Get(d.ID, nil)
			if dev.State != "active" {
				break
			} else {
				finished = true
			}
		}
		if finished {
			return smr, "done", nil
		}
		return nil, "not_done", nil
	}
}
