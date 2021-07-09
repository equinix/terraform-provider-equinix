package metal

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("metal_project_api_key", &resource.Sweeper{
		Name: "metal_project_api_key",
		F:    testSweepProjectAPIKeys,
	})
}

func testSweepProjectAPIKeys(region string) error {
	log.Printf("[DEBUG] Sweeping project API keys")
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("Error getting client for sweeping project API keys: %s", err)
	}
	client := meta.(*packngo.Client)

	_, _, err = client.APIKeys.ProjectList("HELP", nil)

	if err != nil {
		return fmt.Errorf("Error getting list for sweeping project API keys: %s", err)
	}

	return nil
}
