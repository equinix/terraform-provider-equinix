package marketplace

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func fabricMarketplaceSubscriptionDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Subscription URI information",
		},
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Equinix-assigned marketplace identifier",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Subscription Status like; ACTIVE, EXPIRED, CANCELLED, GRACE_PERIOD",
		},
		"marketplace": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Marketplace like; AWS, GCP, AZURE, REDHAT",
		},
		"offer_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Marketplace Offer Type like; PUBLIC, PRIVATE_OFFER",
		},
		"is_auto_renew": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Information about subscription auto renewal",
		},
		"trial": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Subscription Trial",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: marketplaceSubscriptionTrialSch(),
			},
		},
		"entitlements": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Subscription entitlements",
			Elem: &schema.Resource{
				Schema: marketplaceSubscriptionEntitlementsSch(),
			},
		},
	}
}

func marketplaceSubscriptionTrialSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Description: "Marketplace Subscription Trial Enabled",
		},
	}
}

func marketplaceSubscriptionEntitlementsSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Subscription Entitlement Id",
		},
		"quantity_entitled": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Entitled Quantity",
		},
		"quantity_consumed": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Consumed Quantity",
		},
		"quantity_available": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Available Quantity",
		},
		"asset": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Description: "Asset information",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: marketplaceSubscriptionAssetSch(),
			},
		},
	}
}

func marketplaceSubscriptionAssetSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Defines the FCR type like; XF_ROUTER",
		},
		"package": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Fabric Cloud Router Package Type",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: marketplaceSubscriptionPackageSch(),
			},
		},
	}
}
func marketplaceSubscriptionPackageSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"code": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Cloud Router package code",
		},
	}
}
