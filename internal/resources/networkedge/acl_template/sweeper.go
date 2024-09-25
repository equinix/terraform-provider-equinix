package acl_template

import (
	"context"
	"fmt"
	"log"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/equinix/terraform-provider-equinix/internal/sweep"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_network_acl_template", &resource.Sweeper{
		Name: "equinix_network_acl_template",
		F:    testSweepNetworkACLTemplate,
	})
}

func testSweepNetworkACLTemplate(region string) error {
	config, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping Network ACL Templates: %s", err)
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	templates, err := config.Ne.GetACLTemplates()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching Network ACL Templates list: %s", err)
		return err
	}
	nonSweepableCount := 0
	for _, template := range templates {
		if !sweep.IsSweepableTestResource(ne.StringValue(template.Name)) {
			nonSweepableCount++
			continue
		}
		if err := config.Ne.DeleteACLTemplate(ne.StringValue(template.UUID)); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkACLTemplate resource %s (%s): %s", ne.StringValue(template.UUID), ne.StringValue(template.Name), err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkACLTemplate resource %s (%s)", ne.StringValue(template.UUID), ne.StringValue(template.Name))
		}
	}
	if nonSweepableCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonSweepableCount)
	}
	return nil
}
