package equinix

import (
	"fmt"
	"strings"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var networkAccountSchemaNames = map[string]string{
	"Name":      "name",
	"Number":    "number",
	"Status":    "status",
	"UCMID":     "ucm_id",
	"MetroCode": "metro_code",
}

func dataSourceNetworkAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkAccountRead,
		Schema: map[string]*schema.Schema{
			networkAccountSchemaNames["Name"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			networkAccountSchemaNames["Number"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			networkAccountSchemaNames["Status"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Processing", "Submitted", "Staged"}, true),
			},
			networkAccountSchemaNames["UCMID"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			networkAccountSchemaNames["MetroCode"]: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: stringIsMetroCode(),
			},
		},
	}
}

func dataSourceNetworkAccountRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	metro := d.Get(networkAccountSchemaNames["MetroCode"]).(string)
	name := d.Get(networkAccountSchemaNames["Name"]).(string)
	status := d.Get(networkAccountSchemaNames["Status"]).(string)
	accounts, err := conf.ne.GetAccounts(metro)
	if err != nil {
		return err
	}
	var filtered []ne.Account
	for _, account := range accounts {
		if name != "" && account.Name != name {
			continue
		}
		if status != "" && !strings.EqualFold(account.Status, status) {
			continue
		}
		filtered = append(filtered, account)
	}
	if len(filtered) < 1 {
		return fmt.Errorf("network account query returned no results, please change your search criteria")
	}
	if len(filtered) > 1 {
		return fmt.Errorf("network account query returned more than one result, please try more specific search criteria")
	}
	return updateNetworkAccountResource(filtered[0], metro, d)
}

func updateNetworkAccountResource(account ne.Account, metroCode string, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%s-%s", metroCode, account.Name))
	if err := d.Set(networkAccountSchemaNames["Name"], account.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(networkAccountSchemaNames["Number"], account.Number); err != nil {
		return fmt.Errorf("error reading Number: %s", err)
	}
	if err := d.Set(networkAccountSchemaNames["Status"], account.Status); err != nil {
		return fmt.Errorf("error reading Status: %s", err)
	}
	if err := d.Set(networkAccountSchemaNames["UCMID"], account.UCMID); err != nil {
		return fmt.Errorf("error reading UCMID: %s", err)
	}
	if err := d.Set(networkAccountSchemaNames["MetroCode"], metroCode); err != nil {
		return fmt.Errorf("error reading MetroCode: %s", err)
	}
	return nil
}
