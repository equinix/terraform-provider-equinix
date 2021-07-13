package metal

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func resourceMetalConnection() *schema.Resource {
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
				Description:   "Metro where to the connection will be created",
				ConflictsWith: []string{"facility"},
				ForceNew:      true,
				StateFunc:     toLower,
			},
			"redundancy": {
				// TODO: remove ForceNew and do Update, https://github.com/packethost/packngo/issues/270
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connection redundancy - redundant or primary",
				ForceNew:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connection type - dedicated or shared",
				ForceNew:    true,
			},
			"mode": {
				Type:        schema.TypeString,
				Description: "Mode for connections in IBX facilities with the dedicated type - standard or tunnel",
				Optional:    true,
				Default:     "standard",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization responsible for the connection",
				ForceNew:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the project where the connection is scoped to, only used for type == \"shared\"",
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the connection resource",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the connection resource",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fabric Token from the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard)",
			},
			"speed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port speed in bits per second",
			},
			"ports": {
				Type:        schema.TypeList,
				Elem:        connectionPortSchema(),
				Computed:    true,
				Description: "List of connection ports - primary (`ports[0]`) and secondary (`ports[1]`)",
			},
		},
	}
}

func resourceMetalConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	facility, facOk := d.GetOk("facility")
	metro, metOk := d.GetOk("metro")

	if !(metOk || facOk) {
		return fmt.Errorf("You must set either metro or facility")
	}

	project, projectOk := d.GetOk("project_id")
	connType := packngo.ConnectionType(d.Get("type").(string))
	connMode := packngo.ConnectionMode(d.Get("mode").(string))

	if connType == packngo.ConnectionShared && !projectOk {
		return fmt.Errorf("When you create a \"shared\" connection, you must set project_id")
	}
	if connType == packngo.ConnectionDedicated && projectOk {
		return fmt.Errorf("When you create a \"dedicated\" connection, you mustn't set project_id")
	}

	connReq := packngo.ConnectionCreateRequest{
		Name:       d.Get("name").(string),
		Redundancy: packngo.ConnectionRedundancy(d.Get("redundancy").(string)),
		Type:       connType,
	}

	if connType == packngo.ConnectionShared {
		connReq.Project = project.(string)
	}

	connReq.Mode = connMode

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

	orgId := d.Get("organization_id").(string)

	conn, _, err := client.Connections.OrganizationCreate(orgId, &connReq)
	if err != nil {
		return err
	}

	d.SetId(conn.ID)

	return resourceMetalConnectionRead(d, meta)
}

func resourceMetalConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
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
	client := meta.(*packngo.Client)
	connId := d.Id()

	conn, _, err := client.Connections.Get(
		connId,
		&packngo.GetOptions{Includes: []string{"organization", "facility", "metro", "project"}})
	if err != nil {
		return err
	}

	d.SetId(conn.ID)

	projectId := ""
	if conn.Type == packngo.ConnectionShared {
		projectId = conn.Ports[0].VirtualCircuits[0].Project.ID
	}

	return setMap(d, map[string]interface{}{
		"organization_id": conn.Organization.ID,
		"project_id":      projectId,
		"name":            conn.Name,
		"description":     conn.Description,
		"status":          conn.Status,
		"redundancy":      conn.Redundancy,
		"facility":        conn.Facility.Code,
		"metro":           conn.Metro.Code,
		"token":           conn.Token,
		"type":            conn.Type,
		"speed":           conn.Speed,
		"ports":           getConnectionPorts(conn.Ports),
		"mode":            conn.Mode,
	})
}

func resourceMetalConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	resp, err := client.Connections.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}
	return nil
}
