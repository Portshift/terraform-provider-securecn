package securecn

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"log"
	"os"

	"terraform-provider-securecn/client"
)

const ClusterResourceName = "securecn_k8s_cluster"
const ConnectionRuleResourceName = "securecn_connection_rule"
const EnvironmentResourceName = "securecn_environment"
const DeploymentRuleResourceName = "securecn_deployment_rule"
const DeployerResourceName = "securecn_deployer"
const CiPolicyResourceName = "securecn_ci_policy"
const CdPolicyResourceName = "securecn_cd_policy"
const AccessKeyFieldName = "access_key"
const SecretKeyFieldName = "secret_key"
const ServerUrlFieldName = "server_url"

func configureProviderClient(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	httpClient := client.NewHttpClient(
		d.Get(AccessKeyFieldName).(string),
		d.Get(SecretKeyFieldName).(string),
		d.Get(ServerUrlFieldName).(string))

	log.Print("[DEBUG] httpClient created successfully")
	return httpClient, nil
}

func Provider() plugin.ProviderFunc {
	return func() *schema.Provider {
		return &schema.Provider{
			Schema: map[string]*schema.Schema{
				AccessKeyFieldName: {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("SECURECN_ACCESS_KEY", os.Getenv("SECURECN_ACCESS_KEY")),
					Description: "SecureCN service account access key to authenticate with",
				},
				SecretKeyFieldName: {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("SECURECN_SECRET_KEY", os.Getenv("SECURECN_SECRET_KEY")),
					Description: "SecureCN service account secret key to authenticate with",
				},
				ServerUrlFieldName: {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("SECURECN_SERVER_URL", "securecn.cisco.com"),
					Description: "SecureCN server URL",
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				ClusterResourceName:        ResourceCluster(),
				ConnectionRuleResourceName: ResourceConnectionRule(),
				EnvironmentResourceName:    ResourceEnvironment(),
				DeploymentRuleResourceName: ResourceDeploymentRule(),
				DeployerResourceName:       ResourceDeployer(),
				CiPolicyResourceName:       ResourceCiPolicy(),
				CdPolicyResourceName:       ResourceCdPolicy(),
			},
			ConfigureContextFunc: configureProviderClient,
		}
	}
}
