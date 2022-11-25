/*
Copyright 2021 Upbound Inc.
*/

package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/upbound/upjet/pkg/terraform"

	"github.com/bretagne-peiqi/provider-alicloud/apis/v1beta1"
)

const (
	accessKeyID     = "accessKeyID"
	accessKeySecret = "accessKeySecret"
	securityToken   = "securityToken"
	region          = "region"

	// Alibaba Cloud credentials environment variable names
	envAliCloudAcessKey  = "ALICLOUD_ACCESS_KEY"
	envAliCloudSecretKey = "ALICLOUD_SECRET_KEY"
	envAliCloudRegion    = "ALICLOUD_REGION"
	envAliCloudStsToken  = "ALICLOUD_SECURITY_TOKEN"

	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal alicloud credentials as JSON"
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return ps, errors.New(errNoProviderConfig)
		}
		pc := &v1beta1.ProviderConfig{}
		if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
			return ps, errors.Wrap(err, errGetProviderConfig)
		}

		t := resource.NewProviderConfigUsageTracker(client, &v1beta1.ProviderConfigUsage{})
		if err := t.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}

		data, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, client, pc.Spec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		aliCloudCreds := map[string]string{}
		if err := json.Unmarshal(data, &aliCloudCreds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		// set provider configuration
		ps.Configuration = map[string]interface{}{
			"region": "cn-shanghai",
		}
		// set environment variables for sensitive provider configuration
		ps.Env = []string{
			fmt.Sprintf(fmtEnvVar, envAliCloudAcessKey, aliCloudCreds[accessKeyID]),
			fmt.Sprintf(fmtEnvVar, envAliCloudSecretKey, aliCloudCreds[accessKeySecret]),
			fmt.Sprintf(fmtEnvVar, envAliCloudRegion, aliCloudCreds[region]),
			fmt.Sprintf(fmtEnvVar, envAliCloudStsToken, aliCloudCreds[securityToken]),
		}
		// Set credentials in Terraform provider configuration.
		/*ps.Configuration = map[string]any{
			"username": creds["username"],
			"password": creds["password"],
		}*/
		return ps, nil
	}
}
