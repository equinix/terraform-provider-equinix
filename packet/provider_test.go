package packet

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"packet": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PACKET_AUTH_TOKEN"); v == "" {
		t.Fatal("PACKET_AUTH_TOKEN must be set for acceptance tests")
	}
}

func TestNormalizeJSON(t *testing.T) {

	s1 := `
    { "disks": [{"partitions": [{"label": "BIOS",
      "size": 4096,
      "number": 1},
     {"label": "SWAP", "size": "3993600", "number": 2},
     {"label": "ROOT", "size": 0, "number": 3}],
    "device": "/dev/sda",
    "wipeTable": true}],
  "filesystems": [{"mount": {"format": "ext4",
     "device": "/dev/sda3",
     "create": {"options": ["-L", "ROOT"]},
     "point": "/"}},
   {"mount": {"format": "swap",
     "device": "/dev/sda2",
     "create": {"options": ["-L", "SWAP"]},
     "point": "none"}}]
    }`

	s2 := `
    {
  "disks": [
    {
      "device": "/dev/sda",
      "wipeTable": true,
      "partitions": [
        {
          "label": "BIOS",
          "number": 1,
          "size": 4096
        },
        {
          "label": "SWAP",
          "number": 2,
          "size": "3993600"
        },
        {
          "label": "ROOT",
          "number": 3,
          "size": 0
        }
      ]
}
 ],
  "filesystems": [
    {
      "mount": {
        "device": "/dev/sda3",
        "format": "ext4",
        "point": "/",
        "create": {
          "options": [
            "-L",
            "ROOT"
          ]
        }
      }
    },
    {
      "mount": {
        "device": "/dev/sda2",
        "format": "swap",
        "point": "none",
        "create": {
          "options": [
            "-L",
            "SWAP"
          ]
        }
      }
    }
  ]
}`

	m := map[string]interface{}{
		"disks": []map[string]interface{}{
			map[string]interface{}{
				"partitions": []map[string]interface{}{
					map[string]interface{}{
						"label":  "BIOS",
						"size":   4096,
						"number": 1,
					},
					map[string]interface{}{
						"label":  "SWAP",
						"size":   "3993600",
						"number": 2,
					},
					map[string]interface{}{
						"label":  "ROOT",
						"size":   0,
						"number": 3,
					},
				},
				"device":    "/dev/sda",
				"wipeTable": true,
			},
		},
		"filesystems": []map[string]interface{}{
			map[string]interface{}{
				"mount": map[string]interface{}{
					"device": "/dev/sda3",
					"format": "ext4",
					"point":  "/",
					"create": map[string][]string{
						"options": []string{"-L", "ROOT"},
					},
				},
			},
			map[string]interface{}{
				"mount": map[string]interface{}{
					"device": "/dev/sda2",
					"format": "swap",
					"point":  "none",
					"create": map[string][]string{
						"options": []string{"-L", "SWAP"},
					},
				},
			},
		},
	}

	r1, err := normalizeJSON(s1)
	if err != nil {
		t.Errorf("While parsing first test json string: %s", err)
	}
	r2, err := normalizeJSON(s2)
	if err != nil {
		t.Errorf("While parsing second test json string: %s", err)
	}
	rm, err := normalizeJSON(m)
	if err != nil {
		t.Errorf("While parsing second test map: %s", err)
	}
	if r1 != r2 {
		t.Error("Canonical JSON from first and second test string differ.")
	}
	if r1 != rm {
		t.Error("Canonical JSON from first test string and test map differ.")
	}
	if r2 != rm {
		t.Error("Canonical JSON from second test string and test map differ.")
	}
}
