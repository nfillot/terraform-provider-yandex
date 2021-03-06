package yandex

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/vpc/v1"
)

func TestAccDataSourceVPCSubnet(t *testing.T) {
	t.Parallel()

	subnetName1 := acctest.RandomWithPrefix("tf-subnet-1")
	subnetName2 := acctest.RandomWithPrefix("tf-subnet-2")
	subnetDesc1 := "Description for test subnet #1"
	subnetDesc2 := "Description for test subnet #2"

	folderID := getExampleFolderID()
	var network vpc.Network

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckVPCNetworkDestroy,
			testAccCheckVPCSubnetDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVPCSubnetConfig(subnetName1, subnetDesc1, subnetName2, subnetDesc2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCNetworkExists("yandex_vpc_network.foo", &network),

					testAccDataSourceVPCSubnetExists("data.yandex_vpc_subnet.bar1"),
					testAccDataSourceVPCSubnetExists("data.yandex_vpc_subnet.bar2"),

					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "name", subnetName1),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "description", subnetDesc1),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "folder_id", folderID),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "zone", "ru-central1-b"),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "v4_cidr_blocks.#", "1"),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "v4_cidr_blocks.0", "172.16.1.0/24"),
					resource.TestCheckResourceAttrSet("data.yandex_vpc_subnet.bar1", "network_id"),

					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "name", subnetName2),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "description", subnetDesc2),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "folder_id", folderID),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "zone", "ru-central1-c"),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "v4_cidr_blocks.#", "1"),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "v4_cidr_blocks.0", "172.16.2.0/24"),
					resource.TestCheckResourceAttrSet("data.yandex_vpc_subnet.bar2", "network_id"),
				),
			},
		},
	})
}

func TestAccDataSourceVPCSubnetV6(t *testing.T) {
	t.Skip("waiting ipv6 support in subnets")
	t.Parallel()

	subnetName1 := acctest.RandomWithPrefix("tf-subnet-1")
	subnetName2 := acctest.RandomWithPrefix("tf-subnet-2")
	subnetDesc1 := "Description for test subnet #1 with IPv6"
	subnetDesc2 := "Description for test subnet #2 with IPv6"

	folderID := getExampleFolderID()
	var network vpc.Network

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckVPCNetworkDestroy,
			testAccCheckVPCSubnetDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVPCSubnetConfigV6(subnetName1, subnetDesc1, subnetName2, subnetDesc2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCNetworkExists("yandex_vpc_network.foo", &network),

					testAccDataSourceVPCSubnetExists("data.yandex_vpc_subnet.bar1"),
					testAccDataSourceVPCSubnetExists("data.yandex_vpc_subnet.bar2"),

					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "name", subnetName1),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "description", subnetDesc1),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "folder_id", folderID),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "zone", "ru-central1-b"),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "v4_cidr_blocks.#", "1"),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "v4_cidr_blocks.0", "172.16.1.0/24"),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "v6_cidr_blocks.#", "1"),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar1", "v6_cidr_blocks.0", "fd01::/64"),
					resource.TestCheckResourceAttrSet("data.yandex_vpc_subnet.bar1", "network_id"),

					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "name", subnetName2),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "description", subnetDesc2),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "folder_id", folderID),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "zone", "ru-central1-c"),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "v4_cidr_blocks.#", "1"),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "v4_cidr_blocks.0", "172.16.2.0/24"),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "v6_cidr_blocks.#", "1"),
					resource.TestCheckResourceAttr("data.yandex_vpc_subnet.bar2", "v6_cidr_blocks.0", "fd02::/64"),
					resource.TestCheckResourceAttrSet("data.yandex_vpc_subnet.bar2", "network_id"),
				),
			},
		},
	})
}

func testAccDataSourceVPCSubnetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if ds.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.sdk.VPC().Subnet().Get(context.Background(), &vpc.GetSubnetRequest{
			SubnetId: ds.Primary.ID,
		})

		if err != nil {
			return err
		}

		if found.Id != ds.Primary.ID {
			return fmt.Errorf("Subnet not found")
		}

		return nil
	}
}

//revive:disable:var-naming
func testAccDataSourceVPCSubnetConfig(name1, desc1, name2, desc2 string) string {
	return fmt.Sprintf(`
data "yandex_vpc_subnet" "bar1" {
	subnet_id = "${yandex_vpc_subnet.foo1.id}"
}

data "yandex_vpc_subnet" "bar2" {
	subnet_id = "${yandex_vpc_subnet.foo2.id}"
}

resource "yandex_vpc_network" "foo" {
	name        = "%s"
	description = "description for test"
}

resource "yandex_vpc_subnet" "foo1" {
	name           = "%s"
	network_id     = "${yandex_vpc_network.foo.id}"
	description    = "%s"
	v4_cidr_blocks = ["172.16.1.0/24"]
	zone           = "ru-central1-b"
}

resource "yandex_vpc_subnet" "foo2" {
	name           = "%s"
	network_id     = "${yandex_vpc_network.foo.id}"
	description    = "%s"
	v4_cidr_blocks = ["172.16.2.0/24"]
	zone           = "ru-central1-c"
}
`, acctest.RandomWithPrefix("tf-network"), name1, desc1, name2, desc2)
}

func testAccDataSourceVPCSubnetConfigV6(name1, desc1, name2, desc2 string) string {
	return fmt.Sprintf(`
data "yandex_vpc_subnet" "bar1" {
	subnet_id = "${yandex_vpc_subnet.foo1.id}"
}

data "yandex_vpc_subnet" "bar2" {
	subnet_id = "${yandex_vpc_subnet.foo2.id}"
}

resource "yandex_vpc_network" "foo" {
	name        = "%s"
	description = "description for test"
}

resource "yandex_vpc_subnet" "foo1" {
	name           = "%s"
	network_id     = "${yandex_vpc_network.foo.id}"
	description    = "%s"
	v4_cidr_blocks = ["172.16.1.0/24"]
	v6_cidr_blocks = ["fd01::/64"]
	zone           = "ru-central1-b"
}

resource "yandex_vpc_subnet" "foo2" {
	name           = "%s"
	network_id     = "${yandex_vpc_network.foo.id}"
	description    = "%s"
	v4_cidr_blocks = ["172.16.2.0/24"]
    v6_cidr_blocks = ["fd02::/64"]
    zone           = "ru-central1-c"
}
`, acctest.RandomWithPrefix("tf-network"), name1, desc1, name2, desc2)
}
