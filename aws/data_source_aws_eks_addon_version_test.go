package aws

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAWSEksAddonVersionDataSource_basic(t *testing.T) {
	var addonVersion := eks.AddonVersionInfo
	dataSourceResourceName := "data.aws_eks_addon_version.test"
	addonName := "vpc-cni"
	kubernetesVersion := "1.20"
	ctx := context.TODO()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWSEks(t); testAccPreCheckAWSEksAddon(t) },
		ErrorCheck:        testAccErrorCheck(t, eks.EndpointsID),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAWSEksAddonDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksAddonDataSourceConfig_Basic(rName, addonName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksAddonExists(ctx, dataSourceResourceName, &addonVersion),
					testAccMatchResourceAttrRegionalARN(dataSourceResourceName, "arn", "eks", regexp.MustCompile(fmt.Sprintf("addon/%s/%s/.+$", rName, addonName))),
					resource.TestCheckResourceAttrPair(resourceName, "arn", dataSourceResourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "addon_version", dataSourceResourceName, "addon_version"),
					resource.TestCheckResourceAttrPair(resourceName, "service_account_role_arn", dataSourceResourceName, "service_account_role_arn"),
					resource.TestCheckResourceAttrPair(resourceName, "created_at", dataSourceResourceName, "created_at"),
					resource.TestCheckResourceAttrPair(resourceName, "modified_at", dataSourceResourceName, "modified_at"),
					resource.TestCheckResourceAttrPair(resourceName, "tags.%", dataSourceResourceName, "tags.%"),
				),
			},
		},
	})
}

func testAccAWSEksAddonVersionDataSourceConfig_Basic(rName, addonName string) string {
	return composeConfig(testAccAWSEksAddonConfig_Base(rName), fmt.Sprintf(`
data "aws_eks_addon_version" "test" {
  addon_name   = %q
  kubernetes_version = %q
}
`, addonName, kubernetesVersion))
}
