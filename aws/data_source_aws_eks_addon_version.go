package aws

import (
	"context"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
	gversion "github.com/hashicorp/go-version"
)

func dataSourceAwsEksAddonVersion() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAwsEksAddonVersionRead,
		Schema: map[string]*schema.Schema{
			"addon_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"kubernetes_version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"default_version": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"addon_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAwsEksAddonVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*AWSClient).eksconn
	ignoreTagsConfig := meta.(*AWSClient).IgnoreTagsConfig

	addonName := d.Get("addon_name").(string)
	kubernetesVersion := d.Get("kubernetes_version").(string)
	defaultVersion := d.Get("default_version").(bool)

	input := &eks.DescribeAddonVersionsInput{
		AddonName:         aws.String(addonName),
		KubernetesVersion: aws.String(kubernetesVersion),
	}


	output, err := conn.DescribeAddonVersionsWithContext(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading EKS Addon Version for (%s): %w", addonName, err))
	}

	if len(output.Addons) == 0 || len(output.Addons[0].AddonVersions) == 0 {
		return diag.FromErr(fmt.Errorf("EKS Addon Version (%s) not found", addonName))
	}

	var versions []*gversion.Version

	for _, addonVersionInfo := range output.Addons[0].AddonVersions {
		v, err := parseVersion(*addonVersionInfo.AddonVersion)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error reading EKS Addon Version for (%s): %w"), addonName, err)
		}
		// If defaultVersion is true look for default addon version for given
		// Kubernetes version. Default version is the Addon Verison used for EKS when
		// creating a cluster in specific Kubernetes version.
		if defaultVersion {
			for _, comp := range addonVersionInfo.Compatibilities {
				if *comp.DefaultVersion {
					versions = append(versions, v)
				}
			}
		} else {
			versions = append(versions, v)
		}
	}

	if len(versions) == 0 {
		return diag.FromErr(fmt.Errorf("no addon version found for Addon (%q) using Kubernetes version (%q)", addonName, kubernetesVersion))
	}

	// Sort versions to get latest one
	sort.SliceStable(versions, func(i, j int) bool {
		return versions[j].LessThan(versions[i])
	})
	

	d.SetId(fmt.Sprintf("%s:%s", addonName, kubernetesVersion))
	d.Set("addon_version", versions[0].Original())

	return nil
}

func parseVersion(v string) (*gversion.Version, error) {
	version, err := gversion.NewVersion(v)
	if err != nil {
		return nil, err
	}
	return version, nil
}