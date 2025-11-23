package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/melbahja/goph"
)

func TestAccDokkuNetworkCreate(t *testing.T) {
	networkName := fmt.Sprintf("test-network-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDokkuNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "dokku_network" "test" {
	name = "%s"
}				
`, networkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDokkuNetworkExists("dokku_network.test", networkName),
				),
			},
		},
	})
}

func TestAccDokkuNetworkUpdate(t *testing.T) {
	networkNameA := fmt.Sprintf("test-network-%s", acctest.RandString(10))
	networkNameB := fmt.Sprintf("test-network-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDokkuNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "dokku_network" "test" {
	name = "%s"
}				
`, networkNameA),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDokkuNetworkExists("dokku_network.test", networkNameA),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "dokku_network" "test" {
	name = "%s"
}				
`, networkNameB),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDokkuNetworkExists("dokku_network.test", networkNameB),
				),
			},
		},
	})
}

func testAccCheckDokkuNetworkExists(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Network ID not present")
		}

		sshClient := testAccProvider.Meta().(*goph.Client)

		network := NewDokkuNetwork(rs.Primary.ID)
		err := dokkuNetworkRead(network, sshClient)
		if err != nil {
			return fmt.Errorf("Error retrieving network info")
		}

		if network.Name != name {
			return fmt.Errorf("Network names missmatch")
		}
		return nil
	}
}

func testAccCheckDokkuNetworkDestroy(s *terraform.State) error {
	sshClient := testAccProvider.Meta().(*goph.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dokku_network" {
			continue
		}

		exists := dokkuNetworkExists(rs.Primary.ID, sshClient)

		if exists {
			return fmt.Errorf("Dokku network %s should not exist", rs.Primary.ID)
		}
	}

	return nil
}
