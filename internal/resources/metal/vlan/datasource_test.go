package vlan_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vlan"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceMetalVlan_byVxlanMetro(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDatasourceVlanCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalVlanConfig_byVxlanMetro(rs, metro, "tfacc-vlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "vxlan",
						"data.equinix_metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "id",
						"data.equinix_metal_vlan.dsvlan", "id",
					),
					resource.TestCheckResourceAttr(
						"equinix_metal_vlan.barvlan", "vxlan", "6",
					),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_vlan.bardsvlan", "vxlan", "6",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.barvlan", "id",
						"data.equinix_metal_vlan.bardsvlan", "id",
					),
				),
			},
		},
	})
}

func testAccDataSourceMetalVlanConfig_byVxlanMetro(projSuffix, metro, desc string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "equinix_metal_vlan" "foovlan" {
    project_id = equinix_metal_project.foobar.id
    metro = "%s"
    description = "%s"
    vxlan = 5
}

data "equinix_metal_vlan" "dsvlan" {
    metro = equinix_metal_vlan.foovlan.metro
    project_id = equinix_metal_vlan.foovlan.project_id
    vxlan = equinix_metal_vlan.foovlan.vxlan
}

resource "equinix_metal_vlan" "barvlan" {
    project_id = equinix_metal_project.foobar.id
    metro = equinix_metal_vlan.foovlan.metro
    description = "%s"
    vxlan = 6
}

data "equinix_metal_vlan" "bardsvlan" {
    metro = equinix_metal_vlan.barvlan.metro
    project_id = equinix_metal_vlan.barvlan.project_id
    vxlan = equinix_metal_vlan.barvlan.vxlan
}
`, projSuffix, metro, desc, desc)
}

func TestAccDataSourceMetalVlan_byVlanId(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDatasourceVlanCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalVlanConfig_byVlanId(rs, metro, "tfacc-vlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "vxlan",
						"data.equinix_metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "project_id",
						"data.equinix_metal_vlan.dsvlan", "project_id",
					),
				),
			},
		},
	})
}

func testAccDataSourceMetalVlanConfig_byVlanId(projSuffix, metro, desc string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "equinix_metal_vlan" "foovlan" {
    project_id = equinix_metal_project.foobar.id
    metro = "%s"
    description = "%s"
    vxlan = 5
}

data "equinix_metal_vlan" "dsvlan" {
    vlan_id = equinix_metal_vlan.foovlan.id
}
`, projSuffix, metro, desc)
}

func TestAccDataSourceMetalVlan_byProjectId(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDatasourceVlanCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalVlanConfig_byProjectId(rs, metro, "tfacc-vlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "vxlan",
						"data.equinix_metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "project_id",
						"data.equinix_metal_vlan.dsvlan", "project_id",
					),
				),
			},
		},
	})
}

func testAccDataSourceMetalVlanConfig_byProjectId(projSuffix, metro, desc string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "equinix_metal_vlan" "foovlan" {
    project_id = equinix_metal_project.foobar.id
    metro = "%s"
    description = "%s"
    vxlan = 5
}

data "equinix_metal_vlan" "dsvlan" {
    project_id = equinix_metal_vlan.foovlan.project_id
}
`, projSuffix, metro, desc)
}

func TestMetalVlan_matchingVlan(t *testing.T) {
	type args struct {
		vlans     []metalv1.VirtualNetwork
		vxlan     int
		projectID string
		facility  string
		metro     string
	}
	tests := []struct {
		name    string
		args    args
		want    *metalv1.VirtualNetwork
		wantErr bool
	}{
		{
			name: "MatchingVLAN",
			args: args{
				vlans:     []metalv1.VirtualNetwork{{Vxlan: metalv1.PtrInt32(123)}},
				vxlan:     123,
				projectID: "",
				facility:  "",
				metro:     "",
			},
			want:    &metalv1.VirtualNetwork{Vxlan: metalv1.PtrInt32(123)},
			wantErr: false,
		},
		{
			name: "MatchingFac",
			args: args{
				vlans:    []metalv1.VirtualNetwork{{AdditionalProperties: map[string]interface{}{"facility_code": "fac"}}},
				facility: "fac",
			},
			want:    &metalv1.VirtualNetwork{AdditionalProperties: map[string]interface{}{"facility_code": "fac"}},
			wantErr: false,
		},
		{
			name: "MatchingMet",
			args: args{
				vlans: []metalv1.VirtualNetwork{{MetroCode: metalv1.PtrString("met")}},
				metro: "met",
			},
			want:    &metalv1.VirtualNetwork{MetroCode: metalv1.PtrString("met")},
			wantErr: false,
		},
		{
			name: "SecondMatch",
			args: args{
				vlans: []metalv1.VirtualNetwork{{AdditionalProperties: map[string]interface{}{"facility_code": "fac"}}, {MetroCode: metalv1.PtrString("met")}},
				metro: "met",
			},
			want:    &metalv1.VirtualNetwork{MetroCode: metalv1.PtrString("met")},
			wantErr: false,
		},
		{
			name: "TwoMatches",
			args: args{
				vlans: []metalv1.VirtualNetwork{{MetroCode: metalv1.PtrString("met")}, {MetroCode: metalv1.PtrString("met")}},
				metro: "met",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ComplexMatch",
			args: args{
				vlans: []metalv1.VirtualNetwork{{Vxlan: metalv1.PtrInt32(987), AdditionalProperties: map[string]interface{}{"facility_code": "fac"}, MetroCode: metalv1.PtrString("skip")}, {Vxlan: metalv1.PtrInt32(123), AdditionalProperties: map[string]interface{}{"facility_code": "fac"}, MetroCode: metalv1.PtrString("met")}, {Vxlan: metalv1.PtrInt32(456), AdditionalProperties: map[string]interface{}{"facility_code": "fac"}, MetroCode: metalv1.PtrString("nope")}},
				metro: "met",
			},
			want:    &metalv1.VirtualNetwork{Vxlan: metalv1.PtrInt32(123), AdditionalProperties: map[string]interface{}{"facility_code": "fac"}, MetroCode: metalv1.PtrString("met")},
			wantErr: false,
		},
		{
			name: "NoMatch",
			args: args{
				vlans:     nil,
				vxlan:     123,
				projectID: "pid",
				facility:  "fac",
				metro:     "met",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vlan.MatchingVlan(tt.args.vlans, tt.args.vxlan, tt.args.projectID, tt.args.facility, tt.args.metro)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchingVlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchingVlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testAccMetalDatasourceVlanCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).Metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_vlan" {
			continue
		}
		if _, _, err := client.ProjectVirtualNetworks.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Data source VLAN still exists")
		}
	}

	return nil
}

func TestAccDataSourceMetalVlan_byVxlanMetro_upgradeFromVersion(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckMetal(t) },
		CheckDestroy: testAccMetalDatasourceVlanCheckDestroyed,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"equinix": {
						VersionConstraint: "1.29.0", // latest version with resource defined on SDKv2
						Source:            "equinix/equinix",
					},
				},
				Config: testAccDataSourceMetalVlanConfig_byVxlanMetro(rs, metro, "tfacc-vlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "vxlan",
						"data.equinix_metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "id",
						"data.equinix_metal_vlan.dsvlan", "id",
					),
					resource.TestCheckResourceAttr(
						"equinix_metal_vlan.barvlan", "vxlan", "6",
					),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_vlan.bardsvlan", "vxlan", "6",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.barvlan", "id",
						"data.equinix_metal_vlan.bardsvlan", "id",
					),
				),
			},
			{
				ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
				Config:                   testAccDataSourceMetalVlanConfig_byVxlanMetro(rs, metro, "tfacc-vlan"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
