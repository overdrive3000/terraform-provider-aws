package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAWSEksAddonVersionDataSource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	dataSourceResourceName := "data.aws_eks_addon_version.test"
	addonName := "vpc-cni"
	// ctx := context.TODO()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:        testAccErrorCheck(t, eks.EndpointsID),
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksAddonVersionDataSourceConfig_Basic(rName, addonName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceResourceName, "addon_name", "vpc-cni"),
				),
			},
		},
	})
}

func testAccAWSEksAddonVersionDataSourceConfig_Basic(rName, addonName string) string {
	return composeConfig(testAccAWSEksAddonConfig_Base(rName), fmt.Sprintf(`
data "aws_eks_addon_version" "test" {
  addon_name   = %[2]q
  kubernetes_version = aws_eks_cluster.test.version

	depends_on = [
		aws_eks_cluster.test
	]
}
`, rName, addonName))
}
