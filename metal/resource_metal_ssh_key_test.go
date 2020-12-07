package metal

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("metal_ssh_key", &resource.Sweeper{
		Name: "metal_ssh_key",
		F:    testSweepSSHKeys,
	})
}

func testSweepSSHKeys(region string) error {
	log.Printf("[DEBUG] Sweeping ssh keys")
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("Error getting client for sweeping ssh keys: %s", err)
	}
	client := meta.(*packngo.Client)

	sshkeys, _, err := client.SSHKeys.List()
	if err != nil {
		return fmt.Errorf("Error getting list for sweeping ssh keys: %s", err)
	}
	ids := []string{}
	for _, k := range sshkeys {
		if strings.HasPrefix(k.Label, "tfacc-") {
			ids = append(ids, k.ID)
		}
	}
	for _, id := range ids {
		log.Printf("Removing ssh key %s", id)
		resp, err := client.SSHKeys.Delete(id)
		if err != nil && resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Error deleting ssh key %s", err)
		}
	}
	return nil
}

func TestAccMetalSSHKey_Basic(t *testing.T) {
	var key packngo.SSHKey
	rInt := acctest.RandInt()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalSSHKeyConfig_basic(rInt, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalSSHKeyExists("metal_ssh_key.foobar", &key),
					resource.TestCheckResourceAttr(
						"metal_ssh_key.foobar", "name", fmt.Sprintf("tfacc-user-key-%d", rInt)),
					resource.TestCheckResourceAttr(
						"metal_ssh_key.foobar", "public_key", publicKeyMaterial),
					resource.TestCheckResourceAttrSet(
						"metal_ssh_key.foobar", "owner_id"),
				),
			},
		},
	})
}

func TestAccMetalSSHKey_ProjectBasic(t *testing.T) {
	rInt := acctest.RandInt()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalSSHKeyConfig_projectBasic(rInt, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_project.test", "id",
						"metal_project_ssh_key.foobar", "project_id",
					),
				),
			},
		},
	})
}

func TestAccMetalSSHKey_Update(t *testing.T) {
	var key packngo.SSHKey
	rInt := acctest.RandInt()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalSSHKeyConfig_basic(rInt, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalSSHKeyExists("metal_ssh_key.foobar", &key),
					resource.TestCheckResourceAttr(
						"metal_ssh_key.foobar", "name", fmt.Sprintf("tfacc-user-key-%d", rInt)),
					resource.TestCheckResourceAttr(
						"metal_ssh_key.foobar", "public_key", publicKeyMaterial),
				),
			},
			{
				Config: testAccCheckMetalSSHKeyConfig_basic(rInt+1, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalSSHKeyExists("metal_ssh_key.foobar", &key),
					resource.TestCheckResourceAttr(
						"metal_ssh_key.foobar", "name", fmt.Sprintf("tfacc-user-key-%d", rInt+1)),
					resource.TestCheckResourceAttr(
						"metal_ssh_key.foobar", "public_key", publicKeyMaterial),
				),
			},
		},
	})
}

func TestAccMetalSSHKey_projectImportBasic(t *testing.T) {
	sshKey, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalSSHKeyConfig_projectBasic(acctest.RandInt(), sshKey),
			},
			{
				ResourceName:      "metal_project_ssh_key.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMetalSSHKey_importBasic(t *testing.T) {
	sshKey, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalSSHKeyConfig_basic(acctest.RandInt(), sshKey),
			},
			{
				ResourceName:      "metal_ssh_key.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMetalSSHKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_ssh_key" {
			continue
		}
		if _, _, err := client.SSHKeys.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("SSH key still exists")
		}
	}

	return nil
}

func testAccCheckMetalSSHKeyExists(n string, key *packngo.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*packngo.Client)

		foundKey, _, err := client.SSHKeys.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if foundKey.ID != rs.Primary.ID {
			return fmt.Errorf("SSh Key not found: %v - %v", rs.Primary.ID, foundKey)
		}

		*key = *foundKey

		fmt.Printf("key: %v", key)
		return nil
	}
}

func testAccCheckMetalSSHKeyConfig_basic(rInt int, publicSshKey string) string {
	return fmt.Sprintf(`
resource "metal_ssh_key" "foobar" {
    name = "tfacc-user-key-%d"
    public_key = "%s"
}`, rInt, publicSshKey)
}

func testAccCheckMetalSSHKeyConfig_projectBasic(rInt int, publicSshKey string) string {
	return fmt.Sprintf(`

resource "metal_project" "test" {
    name = "tfacc-project-key-test-%d"
}

resource "metal_project_ssh_key" "foobar" {
    name = "tfacc-project-key-%d"
    public_key = "%s"
	project_id = metal_project.test.id
}`, rInt, rInt, publicSshKey)
}
