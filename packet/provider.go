package packet

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a schema.Provider for managing Packet infrastructure.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PACKET_AUTH_TOKEN", nil),
				Description: "The API auth key for API operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"packet_device":            resourcePacketDevice(),
			"packet_ssh_key":           resourcePacketSSHKey(),
			"packet_project":           resourcePacketProject(),
			"packet_volume":            resourcePacketVolume(),
			"packet_volume_attachment": resourcePacketVolumeAttachment(),
			"packet_reserved_ip_block": resourcePacketReservedIPBlock(),
			"packet_ip_attachment":     resourcePacketIPAttachment(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func normalizeJSON(jsonValue interface{}) (string, error) {
	var j map[string]interface{}
	if jsonValue == nil {
		return "", nil
	}
	switch jsonValue.(type) {
	case string:
		s := jsonValue.(string)
		if s == "" {
			return "", nil
		}
		err := json.Unmarshal([]byte(s), &j)
		if err != nil {
			return s, err
		}
	case map[string]interface{}:
		j = jsonValue.(map[string]interface{})
	default:
		return "", fmt.Errorf("%v is not recognized as a JSON value", jsonValue)
	}
	bytes, err := json.Marshal(j)
	if err != nil {
		return "", fmt.Errorf("error marshaling intermediate map %v", j)
	}
	return string(bytes[:]), nil
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AuthToken: d.Get("auth_token").(string),
	}
	return config.Client(), nil
}
