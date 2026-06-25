package marketplace

import (
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setFabricMap(d *schema.ResourceData, subs *fabricv4.SubscriptionResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}
	marketplaceSubscription := subscriptionMap(subs)
	err := equinix_schema.SetMap(d, marketplaceSubscription)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func subscriptionMap(subs *fabricv4.SubscriptionResponse) map[string]any {
	subscription := make(map[string]any)
	subscription["href"] = subs.GetHref()
	subscription["uuid"] = subs.GetUuid()
	subscription["status"] = subs.GetState()
	subscription["marketplace"] = subs.GetMarketplace()
	subscription["offer_type"] = subs.GetOfferType()
	subscription["is_auto_renew"] = subs.GetIsAutoRenew()
	if subs.Trial != nil {
		trial := subs.GetTrial()
		subscription["trial"] = subscriptionTrialGoToTerraform(&trial)
	}
	if subs.Entitlements != nil {
		entitlements := subs.GetEntitlements()
		subscription["entitlements"] = subscriptionEntitlementsGoToTerraform(entitlements)
	}
	return subscription
}

func subscriptionTrialGoToTerraform(trial *fabricv4.SubscriptionTrial) *schema.Set {
	if trial == nil {
		return nil
	}
	mappedTrial := make(map[string]any)
	mappedTrial["enabled"] = trial.GetEnabled()
	trialSet := schema.NewSet(
		schema.HashResource(marketplaceSubscriptionTrialSch()),
		[]any{mappedTrial},
	)
	return trialSet
}

func subscriptionEntitlementsGoToTerraform(entitlementsList []fabricv4.SubscriptionEntitlementResponse) []map[string]any {
	if entitlementsList == nil {
		return nil
	}
	mappedEntitlements := make([]map[string]any, len(entitlementsList))
	for index, entitlements := range entitlementsList {
		asset := entitlements.GetAsset()
		mappedEntitlements[index] = map[string]any{
			"uuid":              entitlements.GetUuid(),
			"quantity_entitled": entitlements.GetQuantityEntitled(),
			"quantity_consumed": entitlements.GetQuantityConsumed(),
			"asset":             subscriptionAssetGoToTerraform(&asset),
		}
	}
	return mappedEntitlements
}

func subscriptionAssetGoToTerraform(asset *fabricv4.SubscriptionAsset) *schema.Set {
	if asset == nil {
		return nil
	}
	mappedAsset := make(map[string]any)
	mappedAsset["type"] = string(asset.GetType())
	package_ := asset.GetPackage()
	mappedAsset["package"] = subscriptionPackageGoToTerraform(&package_)
	assetSet := schema.NewSet(
		schema.HashResource(marketplaceSubscriptionAssetSch()),
		[]any{mappedAsset},
	)
	return assetSet
}

func subscriptionPackageGoToTerraform(package_ *fabricv4.SubscriptionRouterPackageType) *schema.Set {
	if package_ == nil {
		return nil
	}
	mappedPackage := make(map[string]any)
	mappedPackage["code"] = string(package_.GetCode())
	packageSet := schema.NewSet(
		schema.HashResource(marketplaceSubscriptionPackageSch()),
		[]any{mappedPackage},
	)
	return packageSet
}
