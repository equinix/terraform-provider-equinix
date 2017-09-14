package packet

import (
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketDeviceCreate,
		Read:   resourcePacketDeviceRead,
		Update: resourcePacketDeviceUpdate,
		Delete: resourcePacketDeviceDelete,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"operating_system": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"facility": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"plan": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"billing_cycle": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"root_password": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"locked": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"access_public_ipv6": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"access_public_ipv4": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"access_private_ipv4": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"network": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"gateway": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"family": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},

						"cidr": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},

						"public": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_data": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"public_ipv4_subnet_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"ipxe_script_url": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_data"},
			},

			"always_pxe": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"hardware_reservation_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "next-available" && len(old) > 0 {
						return true
					}
					return false
				},
			},

			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePacketDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	createRequest := &packngo.DeviceCreateRequest{
		Hostname:             d.Get("hostname").(string),
		Plan:                 d.Get("plan").(string),
		Facility:             d.Get("facility").(string),
		OS:                   d.Get("operating_system").(string),
		BillingCycle:         d.Get("billing_cycle").(string),
		ProjectID:            d.Get("project_id").(string),
		PublicIPv4SubnetSize: d.Get("public_ipv4_subnet_size").(int),
	}

	if attr, ok := d.GetOk("user_data"); ok {
		createRequest.UserData = attr.(string)
	}

	if attr, ok := d.GetOk("ipxe_script_url"); ok {
		createRequest.IPXEScriptURL = attr.(string)
	}

	if attr, ok := d.GetOk("hardware_reservation_id"); ok {
		createRequest.HardwareReservationID = attr.(string)
	}

	if createRequest.OS == "custom_ipxe" {
		if createRequest.IPXEScriptURL == "" && createRequest.UserData == "" {
			return friendlyError(errors.New("\"ipxe_script_url\" or \"user_data\"" +
				" must be provided when \"custom_ipxe\" OS is selected."))
		}
	}

	if createRequest.OS != "custom_ipxe" && createRequest.IPXEScriptURL != "" {
		return friendlyError(errors.New("\"ipxe_script_url\" argument provided, but" +
			" OS is not \"custom_ipxe\". Please verify and fix device arguments."))
	}

	if attr, ok := d.GetOk("always_pxe"); ok {
		createRequest.AlwaysPXE = attr.(bool)
	}

	tags := d.Get("tags.#").(int)
	if tags > 0 {
		createRequest.Tags = make([]string, 0, tags)
		for i := 0; i < tags; i++ {
			key := fmt.Sprintf("tags.%d", i)
			createRequest.Tags = append(createRequest.Tags, d.Get(key).(string))
		}
	}

	newDevice, _, err := client.Devices.Create(createRequest)
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(newDevice.ID)

	// Wait for the device so we can get the networking attributes that show up after a while.
	_, err = waitForDeviceAttribute(d, "active", []string{"queued", "provisioning"}, "state", meta)
	if err != nil {
		if isForbidden(err) {
			// If the device doesn't get to the active state, we can't recover it from here.
			d.SetId("")

			return errors.New("provisioning time limit exceeded; the Packet team will investigate")
		}
		return err
	}

	return resourcePacketDeviceRead(d, meta)
}

func resourcePacketDeviceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	device, _, err := client.Devices.Get(d.Id())
	if err != nil {
		err = friendlyError(err)

		// If the device somehow already destroyed, mark as succesfully gone.
		if isNotFound(err) {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", device.Hostname)
	d.Set("plan", device.Plan.Slug)
	d.Set("facility", device.Facility.Code)
	d.Set("operating_system", device.OS.Slug)
	d.Set("state", device.State)
	d.Set("billing_cycle", device.BillingCycle)
	d.Set("locked", device.Locked)
	d.Set("created", device.Created)
	d.Set("updated", device.Updated)
	d.Set("ipxe_script_url", device.IPXEScriptURL)
	d.Set("always_pxe", device.AlwaysPXE)
	d.Set("root_password", device.RootPassword)

	if len(device.HardwareReservation.Href) > 0 {
		d.Set("hardware_reservation_id", path.Base(device.HardwareReservation.Href))
	}

	tags := make([]string, 0, len(device.Tags))
	for _, tag := range device.Tags {
		tags = append(tags, tag)
	}
	d.Set("tags", tags)

	var (
		ipv4SubnetSize int
		host           string
		networks       = make([]map[string]interface{}, 0, 1)
	)
	for _, ip := range device.Network {
		network := map[string]interface{}{
			"address": ip.Address,
			"gateway": ip.Gateway,
			"family":  ip.AddressFamily,
			"cidr":    ip.CIDR,
			"public":  ip.Public,
		}
		networks = append(networks, network)

		// Initial device IPs are fixed and marked as "Management"
		if ip.Management {
			if ip.AddressFamily == 4 {
				if ip.Public {
					host = ip.Address
					ipv4SubnetSize = ip.CIDR
					d.Set("access_public_ipv4", ip.Address)
				} else {
					d.Set("access_private_ipv4", ip.Address)
				}
			} else {
				d.Set("access_public_ipv6", ip.Address)
			}
		}
	}
	d.Set("network", networks)
	d.Set("public_ipv4_subnet_size", ipv4SubnetSize)

	if host != "" {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": host,
		})
	}

	return nil
}

func resourcePacketDeviceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	if d.HasChange("locked") {
		var action func(string) (*packngo.Response, error)
		if d.Get("locked").(bool) {
			action = client.Devices.Lock
		} else {
			action = client.Devices.Unlock
		}
		if _, err := action(d.Id()); err != nil {
			return friendlyError(err)
		}
	}

	return resourcePacketDeviceRead(d, meta)
}

func resourcePacketDeviceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	if _, err := client.Devices.Delete(d.Id()); err != nil {
		return friendlyError(err)
	}

	return nil
}

func waitForDeviceAttribute(d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newDeviceStateRefreshFunc(d, attribute, meta),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func newDeviceStateRefreshFunc(d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*packngo.Client)

	return func() (interface{}, string, error) {
		if err := resourcePacketDeviceRead(d, meta); err != nil {
			return nil, "", err
		}

		if attr, ok := d.GetOk(attribute); ok {
			device, _, err := client.Devices.Get(d.Id())
			if err != nil {
				return nil, "", friendlyError(err)
			}
			return &device, attr.(string), nil
		}

		return nil, "", nil
	}
}

// powerOnAndWait Powers on the device and waits for it to be active.
func powerOnAndWait(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	_, err := client.Devices.PowerOn(d.Id())
	if err != nil {
		return friendlyError(err)
	}

	_, err = waitForDeviceAttribute(d, "active", []string{"off"}, "state", client)
	return err
}
