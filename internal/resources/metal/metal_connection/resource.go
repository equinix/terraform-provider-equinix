package metal_connection

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/equinix/terraform-provider-equinix/internal/converters"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
	"golang.org/x/exp/slices"
)

func Resource() *schema.Resource {
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
				Deprecated:    "Use metro instead of facility.  For more information, read the migration guide: https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices",
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
				StateFunc:     converters.ToLowerIf,
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
			"contact_email": {
				Type:        schema.TypeString,
				Description: "The preferred email used for communication and notifications about the Equinix Fabric interconnection. Required when using a Project API key. Optional and defaults to the primary user email address when using a User API key",
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // TODO(displague) packngo needs updating
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
				Optional:    true,
				Computed:    true,
				Description: fmt.Sprintf("Port speed. Required for a_side connections. Allowed values are %s", strings.Join(speeds, ", ")),
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
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

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
		vlans = converters.IfArrToIntArr(d.Get("vlans").([]interface{}))
	}
	connRedundancy := packngo.ConnectionRedundancy(d.Get("redundancy").(string))

	connReq := packngo.ConnectionCreateRequest{
		Name:       d.Get("name").(string),
		Redundancy: connRedundancy,
		Type:       connType,
	}

	// missing email is tolerated for user keys (can't be reasonably detected)
	if contactEmail, ok := d.GetOk("contact_email"); ok {
		connReq.ContactEmail = contactEmail.(string)
	}

	speedRaw, speedOk := d.GetOk("speed")

	// missing speed is tolerated only for shared connections of type z_side
	// https://github.com/equinix/terraform-provider-equinix/issues/276
	if (connType == packngo.ConnectionDedicated) || (tokenType == "a_side") {
		if !speedOk {
			return fmt.Errorf("you must set speed, it's optional only for shared connections of type z_side")
		}
		speed, err := speedStrToUint(speedRaw.(string))
		if err != nil {
			return err
		}
		connReq.Speed = speed
	}

	// this could be generalized, see $ grep "d.Get(\"tags" *
	tags := d.Get("tags.#").(int)
	if tags > 0 {
		connReq.Tags = converters.IfArrToStringArr(d.Get("tags").([]interface{}))
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
				return equinix_errors.FriendlyError(err)
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
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	if d.HasChange("locked") {
		var action func(string) (*packngo.Response, error)
		if d.Get("locked").(bool) {
			action = client.Devices.Lock
		} else {
			action = client.Devices.Unlock
		}
		if _, err := action(d.Id()); err != nil {
			return equinix_errors.FriendlyError(err)
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

	// TODO(displague) packngo does not implement ContactEmail for update
	// if d.HasChange("contact_email" {}

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
			return equinix_errors.FriendlyError(fmt.Errorf("garbage in tags: %s", ts))
		}
	}

	if !reflect.DeepEqual(ur, packngo.ConnectionUpdateRequest{}) {
		if _, _, err := client.Connections.Update(d.Id(), &ur, nil); err != nil {
			return equinix_errors.FriendlyError(err)
		}
	}

	// Don't update VLANs until _after_ the main ConnectionUpdateRequest has succeeded
	if d.HasChange("vlans") {
		connType := packngo.ConnectionType(d.Get("type").(string))

		if connType == packngo.ConnectionShared {
			old, new := d.GetChange("vlans")
			oldVlans := converters.IfArrToIntArr(old.([]interface{}))
			newVlans := converters.IfArrToIntArr(new.([]interface{}))
			maxVlans := int(math.Max(float64(len(oldVlans)), float64(len(newVlans))))

			ports := d.Get("ports").([]interface{})

			for i := 0; i < maxVlans; i++ {
				if d.HasChange(fmt.Sprintf("vlans.%d", i)) {
					if i+1 > len(newVlans) {
						// The VNID was removed; unassign the old VNID
						if _, _, err := updateHiddenVirtualCircuitVNID(client, ports[i].(map[string]interface{}), ""); err != nil {
							return equinix_errors.FriendlyError(err)
						}
					} else {
						j := slices.Index(oldVlans, newVlans[i])
						if j > i {
							// The VNID was moved to a different list index; unassign the VNID for the old index so that it is available for reassignment
							if _, _, err := updateHiddenVirtualCircuitVNID(client, ports[j].(map[string]interface{}), ""); err != nil {
								return equinix_errors.FriendlyError(err)
							}
						}
						// Assign the VNID (whether it is new or moved) to the correct port
						if _, _, err := updateHiddenVirtualCircuitVNID(client, ports[i].(map[string]interface{}), strconv.Itoa(newVlans[i])); err != nil {
							return equinix_errors.FriendlyError(err)
						}
					}
				}
			}
		} else {
			return fmt.Errorf("when you update a \"dedicated\" connection, you cannot set vlans")
		}
	}

	return resourceMetalConnectionRead(d, meta)
}

func updateHiddenVirtualCircuitVNID(client *packngo.Client, port map[string]interface{}, newVNID string) (*packngo.VirtualCircuit, *packngo.Response, error) {
	// This function is used to update the implicit virtual circuits attached to a shared `metal_connection` resource
	// Do not use this function for a non-shared `metal_connection`
	vcids := (port["virtual_circuit_ids"]).([]interface{})
	vcid := vcids[0].(string)
	ucr := packngo.VCUpdateRequest{}
	ucr.VirtualNetworkID = &newVNID
	return client.VirtualCircuits.Update(vcid, &ucr, nil)
}

func resourceMetalConnectionRead(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

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

	return equinix_schema.SetMap(d, map[string]interface{}{
		"organization_id":    conn.Organization.ID,
		"project_id":         projectId,
		"contact_email":      conn.ContactEmail,
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
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal
	resp, err := client.Connections.Delete(d.Id(), true)
	if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
		return equinix_errors.FriendlyError(err)
	}
	return nil
}
