package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNutanixVPCListDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCListDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_vpc_list.test", "entities.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_vpc_list.test", "api_version"),
				),
			},
		},
	})
}

func TestAccNutanixVPCListDataSource_Name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCListDataSourceConfigWithName(randIntBetween(25, 45)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_vpc_list.test", "entities.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_vpc_list.test", "api_version"),
					resource.TestCheckResourceAttr(
						"data.nutanix_vpc_list.test", "entities.0.spec.0.resources.0.externally_routable_prefix_list.0.prefix_length", "16"),
					resource.TestCheckResourceAttr(
						"data.nutanix_vpc_list.test", "entities.0.spec.0.resources.0.externally_routable_prefix_list.0.ip", "172.31.0.0"),
					resource.TestCheckResourceAttr(
						"data.nutanix_vpc_list.test", "entities.0.spec.0.resources.0.common_domain_name_server_ip_list.0.ip", "8.8.8.9"),
				),
			},
		},
	})
}

func testAccVPCListDataSourceConfig() string {
	return (`
	data "nutanix_vpc_list" "test" {
	}
`)
}

func testAccVPCListDataSourceConfigWithName(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "acctest-managed" {
	cluster_uuid = local.cluster1
	name        = "acctest-managed-%[1]d"
	description = "Description of my unit test VLAN"
	vlan_id     = %[1]d
	subnet_type = "VLAN"
	subnet_ip          = "10.250.140.0"
	default_gateway_ip = "10.250.140.1"
	prefix_length = 24
	is_external = true
	ip_config_pool_list_ranges = ["10.250.140.10 10.250.140.20"]
}

resource "nutanix_vpc" "test" {
	name = "acctest-managed-%[1]d"
  
  
	external_subnet_reference_uuid = [
	  resource.nutanix_subnet.acctest-managed.id
	]
  
	common_domain_name_server_ip_list{
			ip = "8.8.8.9"
	}
  
	externally_routable_prefix_list{
	  ip=  "172.31.0.0"
	  prefix_length= 16
	}
  }
	data "nutanix_vpc_list" "test" {
		depends_on = [
			resource.nutanix_vpc.test
		]
	}
`, r)
}