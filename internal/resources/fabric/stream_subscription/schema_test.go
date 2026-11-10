package streamsubscription

import (
	"context"
	"testing"

	datasource "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestStreamSubscriptionCredentialsAreSensitive(t *testing.T) {
	ctx := context.Background()

	sink := resourceSchema(ctx).Attributes["sink"].(resource.SingleNestedAttribute)
	credential := sink.Attributes["credential"].(resource.SingleNestedAttribute)
	settings := sink.Attributes["settings"].(resource.SingleNestedAttribute)
	for _, name := range []string{"access_token", "integration_key", "api_key", "password"} {
		if !credential.Attributes[name].IsSensitive() {
			t.Errorf("resource sink.credential.%s is not sensitive", name)
		}
	}
	if !settings.Attributes["application_key"].IsSensitive() {
		t.Error("resource sink.settings.application_key is not sensitive")
	}

	dsSink := dataSourceStreamSubscriptionByID(ctx).Attributes["sink"].(datasource.SingleNestedAttribute)
	dsCredential := dsSink.Attributes["credential"].(datasource.SingleNestedAttribute)
	dsSettings := dsSink.Attributes["settings"].(datasource.SingleNestedAttribute)
	for _, name := range []string{"access_token", "integration_key", "api_key", "password"} {
		if !dsCredential.Attributes[name].IsSensitive() {
			t.Errorf("data source sink.credential.%s is not sensitive", name)
		}
	}
	if !dsSettings.Attributes["application_key"].IsSensitive() {
		t.Error("data source sink.settings.application_key is not sensitive")
	}
}
