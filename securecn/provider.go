package securecn

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"terraform-provider-securecn/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
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

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func configureProviderClient(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	httpClient := client.NewHttpClient(
		d.Get(AccessKeyFieldName).(string),
		d.Get(SecretKeyFieldName).(string),
		d.Get(ServerUrlFieldName).(string))

	err := installKubectlOnDemand()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	log.Print("[DEBUG] httpClient created successfully")
	return httpClient, nil
}

func installKubectlOnDemand() error {
	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		kubectlDir := "/tmp/terraformbin/"
		kubectlPath := kubectlDir + "kubectl"

		_, err := os.Stat(kubectlPath)
		if err != nil {
			log.Print("[DEBUG] kubectl not available, grabbing it")

			err = os.MkdirAll(kubectlDir, 0755)
			if err != nil {
				return err
			}

			kubectlURL := fmt.Sprintf("https://dl.k8s.io/release/v1.23.0/bin/%s/%s/kubectl", runtime.GOOS, runtime.GOARCH)
			err = downloadFile(kubectlPath, kubectlURL)
			if err != nil {
				return err
			}

			err = os.Chmod(kubectlPath, 0755)
			if err != nil {
				return err
			}

			log.Printf("[DEBUG] kubectl got downloaded to %s", kubectlPath)
		} else {
			log.Printf("[DEBUG] kubectl is available at %s", kubectlPath)
		}

		_ = os.Setenv("PATH", kubectlDir+":"+os.Getenv("PATH"))

	} else {
		log.Printf("[DEBUG] kubectl is available at %s", kubectlPath)
	}

	return nil
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
					Description: "Appsecurity service account secret key to authenticate with",
				},
				ServerUrlFieldName: {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("SECURECN_SERVER_URL", "appsecurity.cisco.com"),
					Description: "Appsecurity server URL",
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
