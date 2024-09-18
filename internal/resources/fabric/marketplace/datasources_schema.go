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
			Computed:    true,
			Description: "Marketplace Offer Type like; PUBLIC, PRIVATE_OFFER",
		},
		"is_auto_renew": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Information about subscription auto renewal",
		},
		"trial": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Subscription Trial",
			Elem:        marketplaceSubscriptionTrialSch(),
		},
		"entitlements": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Subscription entitlements",
			Elem:        marketplaceSubscriptionEntitlementsSch(),
		},
	}
}

func marketplaceSubscriptionTrialSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Marketplace Subscription Trial Enabled",
			},
		},
	}
}

func marketplaceSubscriptionEntitlementsSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subscription Entitlement Id",
			},
			"quantity_entitled": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Entitled Quantity",
			},
			"quantity_consumed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Consumed Quantity",
			},
			"quantity_available": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Available Quantity",
			},
			"asset": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Asset information",
				Elem:        marketplaceSubscriptionAssetSch(),
			},
		},
	}
}

func marketplaceSubscriptionAssetSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the FCR type like; XF_ROUTER",
			},
			"package": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Fabric Cloud Router Package Type",
				Elem:        marketplaceSubscriptionPackageSch(),
			},
		},
	}
}
func marketplaceSubscriptionPackageSch() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud Router package code",
			},
		},
	}
}
