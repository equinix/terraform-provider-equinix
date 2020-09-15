package equinix

import (
	"fmt"
	"strings"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var neAccountSchemaNames = map[string]string{
	"Name":      "name",
	"Number":    "number",
	"Status":    "status",
	"UCMID":     "ucm_id",
	"MetroCode": "metro_code",
}

func dataSourceNeAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNeAccountRead,
		Schema: map[string]*schema.Schema{
			neAccountSchemaNames["Name"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			neAccountSchemaNames["Number"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			neAccountSchemaNames["Status"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Processing", "Submitted", "Staged"}, true),
			},
			neAccountSchemaNames["UCMID"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			neAccountSchemaNames["MetroCode"]: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: stringIsMetroCode(),
			},
		},
	}
}

func dataSourceNeAccountRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	metro := d.Get(neAccountSchemaNames["MetroCode"]).(string)
	name := d.Get(neAccountSchemaNames["Name"]).(string)
	status := d.Get(neAccountSchemaNames["Status"]).(string)
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
		return fmt.Errorf("account query returned no results, please change your search criteria")
	}
	if len(filtered) > 1 {
		return fmt.Errorf("account query returned more than one result, please try more specific search criteria")
	}
	return updateNeAccountResource(filtered[0], metro, d)
}

func updateNeAccountResource(account ne.Account, metroCode string, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%s-%s", metroCode, account.Name))
	if err := d.Set(neAccountSchemaNames["Name"], account.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(neAccountSchemaNames["Number"], account.Number); err != nil {
		return fmt.Errorf("error reading Number: %s", err)
	}
	if err := d.Set(neAccountSchemaNames["Status"], account.Status); err != nil {
		return fmt.Errorf("error reading Status: %s", err)
	}
	if err := d.Set(neAccountSchemaNames["UCMID"], account.UCMID); err != nil {
		return fmt.Errorf("error reading UCMID: %s", err)
	}
	if err := d.Set(neAccountSchemaNames["MetroCode"], metroCode); err != nil {
		return fmt.Errorf("error reading MetroCode: %s", err)
	}
	return nil
}
