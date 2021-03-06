package yandex

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceYandexIAMServiceAccount(t *testing.T) {
	accountName := "sa" + acctest.RandString(10)
	accountDesc := "Service Account desc"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIAMServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataServiceAccount(accountName, accountDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.yandex_iam_service_account.bar",
						"name", accountName),
					resource.TestCheckResourceAttr("data.yandex_iam_service_account.bar",
						"description", accountDesc),
					resource.TestCheckResourceAttr("data.yandex_iam_service_account.bar",
						"folder_id", getExampleFolderID()),
				),
			},
		},
	})
}

func testAccDataServiceAccount(name, desc string) string {
	return fmt.Sprintf(`
data "yandex_iam_service_account" "bar" {
    service_account_id = "${yandex_iam_service_account.foo.id}"
}

resource "yandex_iam_service_account" "foo" {
    name        = "%s"
	description = "%s"
}`, name, desc)
}
