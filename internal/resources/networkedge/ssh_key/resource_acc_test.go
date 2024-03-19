package ssh_key

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/nprintf"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("equinix_network_ssh_key", &resource.Sweeper{
		Name: "equinix_network_ssh_key",
		F:    testSweepNetworkSSHKey,
	})
}

func testSweepNetworkSSHKey(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping Network SSH keys: %s", err)
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	keys, err := config.Ne.GetSSHPublicKeys()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching NetworkSSHKey list: %s", err)
		return err
	}
	nonSweepableCount := 0
	for _, key := range keys {
		if !isSweepableTestResource(ne.StringValue(key.Name)) {
			nonSweepableCount++
			continue
		}
		if err := config.Ne.DeleteSSHPublicKey(ne.StringValue(key.UUID)); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkSSHKey resource %s (%s): %s", ne.StringValue(key.UUID), ne.StringValue(key.Name), err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkSSHKey resource %s (%s)", ne.StringValue(key.UUID), ne.StringValue(key.Name))
		}
	}
	if nonSweepableCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonSweepableCount)
	}
	return nil
}

func TestAccNetworkSSHKey(t *testing.T) {
	context := map[string]interface{}{
		"resourceName": "test",
		"name":         fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"public_key":   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCXdzXBHaVpKpdO0udnB+4JOgUq7APO2rPXfrevvlZrps98AtlwXXVWZ5duRH5NFNfU4G9HCSiAPsebgjY0fG85tcShpXfHfACLt0tBW8XhfLQP2T6S50FQ1brBdURMDCMsD7duOXqvc0dlbs2/KcswHvuUmqVzob3bz7n1bQ48wIHsPg4ARqYhy5LN3OkllJH/6GEfqi8lKZx01/P/gmJMORcJujuOyXRB+F2iXBVYdhjML3Qg4+tEekBcVZOxUbERRZ0pvQ52Y6wUhn2VsjljixyqeOdmD0m6DayDQgSWms6bKPpBqN7zhXXk4qe8bXT4tQQba65b2CQ2A91jw2KgM/YZNmjyUJ+Rf1cQosJf9twqbAZDZ6rAEmj9zzvQ5vD/CGuzxdVMkePLlUK4VGjPu7cVzhXrnq4318WqZ5/lNiCST8NQ0fssChN8ANUzr/p/wwv3faFMVNmjxXTZMsbMFT/fbb2MVVuqNFN65drntlg6/xEao8gZROuRYiakBx8= user@host",
		"type":         "RSA",
	}
	resourceName := fmt.Sprintf("equinix_network_ssh_key.%s", context["resourceName"].(string))
	var key ne.SSHPublicKey
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSSHKey(context),
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkSSHKeyExists(resourceName, &key),
					testAccNetworkSSHKeyAttributes(&key, context),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkSSHKey(ctx map[string]interface{}) string {
	return nprintf.NPrintf(`
resource "equinix_network_ssh_key" "%{resourceName}" {
  name       = "%{name}"
  public_key = "%{public_key}"
  type       = "%{type}"
}
`, ctx)
}

func testAccNetworkSSHKeyExists(resourceName string, key *ne.SSHPublicKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		client := testAccProvider.Meta().(*config.Config).Ne
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		resp, err := client.GetSSHPublicKey(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error when fetching SSH public key '%s': %s", rs.Primary.ID, err)
		}
		*key = *resp
		return nil
	}
}

func testAccNetworkSSHKeyAttributes(key *ne.SSHPublicKey, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["name"]; ok && ne.StringValue(key.Name) != v.(string) {
			return fmt.Errorf("name does not match %v - %v", ne.StringValue(key.Name), v)
		}
		if v, ok := ctx["public_key"]; ok && ne.StringValue(key.Value) != v.(string) {
			return fmt.Errorf("public_key does not match %v - %v", ne.StringValue(key.Value), v)
		}
		if v, ok := ctx["type"]; ok && ne.StringValue(key.Type) != v.(string) {
			return fmt.Errorf("type does not match %v - %v", ne.StringValue(key.Type), v)
		}
		return nil
	}
}
