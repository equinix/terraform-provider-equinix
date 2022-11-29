package equinix

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

var (
	mega          uint64 = 1000 * 1000
	giga          uint64 = 1000 * mega
	allowedSpeeds        = []struct {
		Int uint64
		Str string
	}{
		{50 * mega, "50Mbps"},
		{200 * mega, "200Mbps"},
		{500 * mega, "500Mbps"},
		{1 * giga, "1Gbps"},
		{2 * giga, "2Gbps"},
		{5 * giga, "5Gbps"},
		{10 * giga, "10Gbps"},
	}
)

func speedStrToUint(speed string) (uint64, error) {
	allowedStrings := []string{}
	for _, allowedSpeed := range allowedSpeeds {
		if allowedSpeed.Str == speed {
			return allowedSpeed.Int, nil
		}
		allowedStrings = append(allowedStrings, allowedSpeed.Str)
	}
	return 0, fmt.Errorf("invalid speed string: %s. Allowed strings: %s", speed, strings.Join(allowedStrings, ", "))
}

func speedUintToStr(speed uint64) (string, error) {
	allowedUints := []uint64{}
	for _, allowedSpeed := range allowedSpeeds {
		if speed == allowedSpeed.Int {
			return allowedSpeed.Str, nil
		}
		allowedUints = append(allowedUints, allowedSpeed.Int)
	}
	return "", fmt.Errorf("%d is not allowed speed value. Allowed values: %v", speed, allowedUints)
}

func resourceMetalConnection() *schema.Resource {
	speeds := []string{}
	for _, allowedSpeed := range allowedSpeeds {
		speeds = append(speeds, allowedSpeed.Str)
	}
	return &schema.Resource{
		Read:   resourceMetalConnectionRead,
		Create: resourceMetalConnectionCreate,
		Delete: resourceMetalConnectionDelete,
		Update: resourceMetalConnectionUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the connection resource",
				ForceNew:    true,
			},
			"facility": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "Facility where the connection will be created",
				ConflictsWith: []string{"metro"},
				ForceNew:      true,
			},
			"metro": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "Metro where the connection will be created",
				ConflictsWith: []string{"facility"},
				ForceNew:      true,
				StateFunc:     toLower,
			},
			"redundancy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connection redundancy - redundant or primary",
				ValidateFunc: validation.StringInSlice([]string{
					string(packngo.ConnectionRedundant),
					string(packngo.ConnectionPrimary),
				}, false),
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connection type - dedicated or shared",
				ForceNew:    true,
				ValidateFunc: validation.StringInSlice([]string{
					string(packngo.ConnectionDedicated),
					string(packngo.ConnectionShared),
				}, false),
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the project where the connection is scoped to. Required with type \"shared\"",
				ForceNew:    true,
			},
			"speed": {
				Type:        schema.TypeString,
				Required:    true,
				Description: fmt.Sprintf("Port speed. Allowed values are %s", strings.Join(speeds, ", ")),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the connection resource",
			},
			"mode": {
				Type:        schema.TypeString,
				Description: "Mode for connections in IBX facilities with the dedicated type - standard or tunnel",
				Optional:    true,
				Default:     "standard",
				ValidateFunc: validation.StringInSlice([]string{
					string(packngo.ConnectionModeStandard),
					string(packngo.ConnectionModeTunnel),
				}, false),
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "Tags attached to the connection",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"vlans": {
				Type:        schema.TypeList,
				Description: "Only used with shared connection. VLANs to attach. Pass one vlan for Primary/Single connection and two vlans for Redundant connection",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				MaxItems:    2,
			},
			"service_token_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Only used with shared connection. Type of service token to use for the connection, a_side or z_side",
			},
			"organization_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "ID of the organization responsible for the connection. Applicable with type \"dedicated\"",
				ForceNew:     true,
				AtLeastOneOf: []string{"organization_id", "project_id"},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the connection resource",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Only used with shared connection. Fabric Token required to continue the setup process with [equinix_ecx_l2_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_ecx_l2_connection) or from the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard)",
				Deprecated:  "If your organization already has connection service tokens enabled, use `service_tokens` instead",
			},
			"ports": {
				Type:        schema.TypeList,
				Elem:        connectionPortSchema(),
				Computed:    true,
				Description: "List of connection ports - primary (`ports[0]`) and secondary (`ports[1]`)",
			},
			"service_tokens": {
				Type:        schema.TypeList,
				Description: "Only used with shared connection. List of service tokens required to continue the setup process with [equinix_ecx_l2_connection](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_ecx_l2_connection) or from the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard)",
				Computed:    true,
				Elem:        serviceTokenSchema(),
			},
		},
	}
}

func resourceMetalConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	meta.(*Config).addModuleToMetalUserAgent(d)
	client := meta.(*Config).metal

	facility, facOk := d.GetOk("facility")
	metro, metOk := d.GetOk("metro")

	if !(metOk || facOk) {
		return fmt.Errorf("you must set either metro or facility")
	}

	connType := packngo.ConnectionType(d.Get("type").(string))

	connMode := d.Get("mode").(string)

	tokenTypeRaw, tokenTypeOk := d.GetOk("service_token_type")
	tokenType := packngo.FabricServiceTokenType(tokenTypeRaw.(string))

	vlans := []int{}
	vlansNum := d.Get("vlans.#").(int)
	if vlansNum > 0 {
		vlans = convertIntArr2(d.Get("vlans").([]interface{}))
	}
	connRedundancy := packngo.ConnectionRedundancy(d.Get("redundancy").(string))

	speed, err := speedStrToUint(d.Get("speed").(string))
	if err != nil {
		return err
	}

	connReq := packngo.ConnectionCreateRequest{
		Name:       d.Get("name").(string),
		Redundancy: connRedundancy,
		Type:       connType,
		Speed:      speed,
	}

	// this could be generalized, see $ grep "d.Get(\"tags" *
	tags := d.Get("tags.#").(int)
	if tags > 0 {
		connReq.Tags = convertStringArr(d.Get("tags").([]interface{}))
	}

	if metOk {
		connReq.Metro = metro.(string)
	}
	if facOk {
		connReq.Facility = facility.(string)
	}

	desc, descOk := d.GetOk("description")
	if descOk {
		description := desc.(string)
		connReq.Description = &description
	}

	projectId, projectIdOk := d.GetOk("project_id")
	if connType == packngo.ConnectionShared {
		if !projectIdOk {
			return fmt.Errorf("you must set project_id for \"shared\" connection")
		}
		if connMode == string(packngo.ConnectionModeTunnel) {
			return fmt.Errorf("tunnel mode is not supported for \"shared\" connections")
		}
		if connRedundancy == packngo.ConnectionPrimary && vlansNum == 2 {
			return fmt.Errorf("when you create a \"shared\" connection without redundancy, you must only set max 1 vlan")
		}
		connReq.VLANs = vlans
		connReq.ServiceTokenType = tokenType
		conn, _, err := client.Connections.ProjectCreate(projectId.(string), &connReq)
		if err != nil {
			return err
		}
		d.SetId(conn.ID)
	} else {
		organizationId, organizationIdOk := d.GetOk("organization_id")
		if !organizationIdOk {
			if !projectIdOk {
				return fmt.Errorf("you must set one of organization_id or project_id for \"dedicated\" connection")
			}
			proj, _, err := client.Projects.Get(projectId.(string), &packngo.GetOptions{Includes: []string{"organization"}})
			if err != nil {
				return friendlyError(err)
			}
			organizationId = proj.Organization.ID
		}
		if tokenTypeOk {
			return fmt.Errorf("when you create a \"dedicated\" connection, you must not set service_token_type")
		}
		if vlansNum > 0 {
			return fmt.Errorf("when you create a \"dedicated\" connection, you must not set vlans")
		}
		connReq.Mode = packngo.ConnectionMode(connMode)
		conn, _, err := client.Connections.OrganizationCreate(organizationId.(string), &connReq)
		if err != nil {
			return err
		}
		d.SetId(conn.ID)
	}

	return resourceMetalConnectionRead(d, meta)
}

func resourceMetalConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	meta.(*Config).addModuleToMetalUserAgent(d)
	client := meta.(*Config).metal

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
	ur := packngo.ConnectionUpdateRequest{}

	if d.HasChange("description") {
		desc := d.Get("description").(string)
		ur.Description = &desc
	}

	if d.HasChange("mode") {
		mode := packngo.ConnectionMode(d.Get("mode").(string))
		ur.Mode = &mode
	}

	if d.HasChange("redundancy") {
		redundancy := packngo.ConnectionRedundancy(d.Get("redundancy").(string))
		ur.Redundancy = redundancy
	}

	if d.HasChange("tags") {
		ts := d.Get("tags")
		sts := []string{}

		switch ts.(type) {
		case []interface{}:
			for _, v := range ts.([]interface{}) {
				sts = append(sts, v.(string))
			}
			ur.Tags = sts
		default:
			return friendlyError(fmt.Errorf("garbage in tags: %s", ts))
		}
	}

	if !reflect.DeepEqual(ur, packngo.ConnectionUpdateRequest{}) {
		if _, _, err := client.Connections.Update(d.Id(), &ur, nil); err != nil {
			return friendlyError(err)
		}
	}
	return resourceMetalConnectionRead(d, meta)
}

func resourceMetalConnectionRead(d *schema.ResourceData, meta interface{}) error {
	meta.(*Config).addModuleToMetalUserAgent(d)
	client := meta.(*Config).metal

	connId := d.Id()
	conn, _, err := client.Connections.Get(
		connId,
		&packngo.GetOptions{Includes: []string{"service_tokens", "organization", "facility", "metro", "project"}})
	if err != nil {
		return err
	}

	d.SetId(conn.ID)

	projectId := d.Get("project_id").(string)
	// fix the project id get when it's added straight to the Connection API resource
	// https://github.com/packethost/packngo/issues/317
	if conn.Type == packngo.ConnectionShared {
		projectId = conn.Ports[0].VirtualCircuits[0].Project.ID
	}
	mode := "standard"
	if conn.Mode != nil {
		mode = string(*conn.Mode)
	}
	side := ""
	if len(conn.Tokens) > 0 {
		side = string(conn.Tokens[0].ServiceTokenType)
	}
	speed := "0"
	if conn.Speed > 0 {
		speed, err = speedUintToStr(conn.Speed)
		if err != nil {
			return err
		}
	}
	serviceTokens, err := getServiceTokens(conn.Tokens)
	if err != nil {
		return err
	}

	vlans := getConnectionVlans(conn)
	if vlans != nil {
		d.Set("vlans", vlans)
	}

	return setMap(d, map[string]interface{}{
		"organization_id":    conn.Organization.ID,
		"project_id":         projectId,
		"name":               conn.Name,
		"description":        conn.Description,
		"status":             conn.Status,
		"redundancy":         conn.Redundancy,
		"facility":           conn.Facility.Code,
		"metro":              conn.Metro.Code,
		"token":              conn.Token,
		"type":               conn.Type,
		"speed":              speed,
		"ports":              getConnectionPorts(conn.Ports),
		"mode":               mode,
		"tags":               conn.Tags,
		"service_tokens":     serviceTokens,
		"service_token_type": side,
	})
}

func resourceMetalConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	meta.(*Config).addModuleToMetalUserAgent(d)
	client := meta.(*Config).metal
	resp, err := client.Connections.Delete(d.Id(), true)
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}
	return nil
}
