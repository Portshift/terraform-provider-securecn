package securecn

import (
	"context"
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-securecn/client"
	"terraform-provider-securecn/escher_api/escherClient"
	"terraform-provider-securecn/escher_api/model"
	"terraform-provider-securecn/utils"
)

const nameFieldName = "name"
const descriptionFieldName = "description"
const kubernetesEnvironmentFieldName = "kubernetes_environment"
const clusterNameFieldName = "cluster_name"
const namespacesNamesFieldName = "namespaces_by_names"
const namespacesLabelsFieldName = "namespaces_by_labels"

func ResourceEnvironment() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,
		Description:   "A SecureCN environment",
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			nameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			descriptionFieldName: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			kubernetesEnvironmentFieldName: {
				Description: "The kubernetes environments to include in the SecureCN env",
				Required:    true,
				Type:        schema.TypeList,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						clusterNameFieldName: {
							Description: "The name of the kubernetes cluster in SecureCN",
							Type:        schema.TypeString,
							Required:    true,
						},
						namespacesNamesFieldName: {
							Description: "The env will match using namespace name",
							Optional:    true,
							Type:        schema.TypeList,
							MinItems:    1,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						namespacesLabelsFieldName: {
							Description: "The source will match using namespace labels",
							Optional:    true,
							Type:        schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] creating environment")

	err := validateEnvironmentConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils.GetServiceApi(&httpClientWrapper)

	environmentFromConfig, err := getEnvironmentFromConfig(ctx, d, serviceApi, httpClientWrapper)
	if err != nil {
		return diag.FromErr(err)
	}

	environment, err := serviceApi.CreateEnvironment(ctx, httpClientWrapper.HttpClient, environmentFromConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	envId := environment.Payload.ID

	d.SetId(string(envId))

	return nil
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] reading environment")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils.GetServiceApi(&httpClientWrapper)
	envId := d.Id()

	currentEnvInSecureCN, err := serviceApi.GetEnvironment(ctx, httpClientWrapper.HttpClient, strfmt.UUID(envId))
	if err != nil {
		return diag.FromErr(err)
	}

	if currentEnvInSecureCN.Payload.ID == "" {
		// Tell terraform the env doesn't exist
		d.SetId("")
	} else {
		err = updateEnvironmentMutableFields(d, currentEnvInSecureCN.Payload)
		return diag.FromErr(err)
	}

	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] updating environment")

	err := validateEnvironmentConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils.GetServiceApi(&httpClientWrapper)

	environment, err := getEnvironmentFromConfig(ctx, d, serviceApi, httpClientWrapper)
	if err != nil {
		return diag.FromErr(err)
	}
	environment.ID = strfmt.UUID(d.Id())

	updatedEnv, err := serviceApi.UpdateEnvironment(ctx, httpClientWrapper.HttpClient, environment, environment.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(updatedEnv.Payload.ID))

	return resourceEnvironmentRead(ctx, d, m)
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] deleting environment")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils.GetServiceApi(&httpClientWrapper)
	envId := strfmt.UUID(d.Id())
	err := serviceApi.DeleteEnvironment(ctx, httpClientWrapper.HttpClient, envId)
	if err != nil {
		return diag.FromErr(err)
	}

	// Tell terraform the env doesn't exist
	d.SetId("")

	return nil
}

func validateEnvironmentConfig(d *schema.ResourceData) error {
	log.Printf("[DEBUG] validating config")
	allKubernetesEnvironments := d.Get(kubernetesEnvironmentFieldName).([]interface{})
	for index, _ := range allKubernetesEnvironments {
		namespaceNames := utils.ReadNestedListStringFromTF(d, kubernetesEnvironmentFieldName, namespacesNamesFieldName, index)
		namespaceLabels := utils.GetLabelsFromMap(utils.ReadNestedMapStringFromTF(d, kubernetesEnvironmentFieldName, namespacesLabelsFieldName, index))

		if namespaceNames != nil && namespaceLabels != nil {
			return fmt.Errorf("kubernetes_environment.%d.namespaces_by_names\": one of `kubernetes_environment.%d.namespaces_by_names,kubernetes_environment.[%d].namespaces_by_labels` must be specified", index, index, index)
		}
	}
	return nil
}

func getEnvironmentFromConfig(ctx context.Context, d *schema.ResourceData, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper) (*model.Environment, error) {
	log.Print("[DEBUG] getting environment from config")

	name := d.Get(nameFieldName).(string)
	desc := d.Get(descriptionFieldName).(string)

	kubernetesEnvs := make([]*model.KubernetesEnvironment, 0)

	allKubernetesEnvironments := d.Get(kubernetesEnvironmentFieldName).([]interface{})
	for index, _ := range allKubernetesEnvironments {
		clusterName := utils.ReadNestedStringFromTF(d, kubernetesEnvironmentFieldName, clusterNameFieldName, index)
		log.Printf("[DEBUG] %v clusterName: %v", index, clusterName)
		namespaceNames := utils.ReadNestedListStringFromTF(d, kubernetesEnvironmentFieldName, namespacesNamesFieldName, index)
		log.Printf("[DEBUG] %v namespaceNames: %v", index, namespaceNames)
		namespaceLabels := utils.GetLabelsFromMap(utils.ReadNestedMapStringFromTF(d, kubernetesEnvironmentFieldName, namespacesLabelsFieldName, index))
		log.Printf("[DEBUG] %v namespaceLabels: %v", index, namespaceLabels)

		clusterId, err := serviceApi.GetKubernetesClusterIdByName(ctx, httpClientWrapper.HttpClient, clusterName)
		if err != nil {
			return nil, err
		}

		kubeEnv := createKubernetesEnvFromConfig(namespaceNames, namespaceLabels, &clusterId.Payload)

		kubernetesEnvs = append(kubernetesEnvs, kubeEnv)
	}

	env := &model.Environment{
		ID:                     "",
		Description:            desc,
		KubernetesEnvironments: kubernetesEnvs,
		Name:                   &name,
	}

	return env, nil
}

func createKubernetesEnvFromConfig(namespaceNames []string, namespaceLabels []*model.Label, clusterID *strfmt.UUID) *model.KubernetesEnvironment {

	return &model.KubernetesEnvironment{
		ID:                "",
		KubernetesCluster: clusterID,
		NamespaceLabels:   namespaceLabels,
		Namespaces:        namespaceNames,
	}

}

func updateEnvironmentMutableFields(d *schema.ResourceData, currentEnvInSecureCN *model.Environment) error {
	log.Print("[DEBUG] updating environment mutable fields")

	err := d.Set(nameFieldName, currentEnvInSecureCN.Name)
	if err != nil {
		return err
	}

	err = d.Set(descriptionFieldName, currentEnvInSecureCN.Description)
	if err != nil {
		return err
	}

	err = mutateKubernetesEnvs(d, currentEnvInSecureCN)
	return err
}

func mutateKubernetesEnvs(d *schema.ResourceData, currentEnvInSecureCN *model.Environment) error {

	envsInSecureCN := currentEnvInSecureCN.KubernetesEnvironments

	envsInTf := d.Get(kubernetesEnvironmentFieldName).([]interface{})

	for index, envInSecureCN := range envsInSecureCN {
		if index < len(envsInTf) {
			envInTf := envsInTf[index].(map[string]interface{})
			envInTf[clusterNameFieldName] = envInSecureCN.KubernetesClusterName
			envInTf[namespacesNamesFieldName] = envInSecureCN.NamespaceLabels
			envInTf[namespacesNamesFieldName] = envInSecureCN.Namespaces
		} else {
			newEnv := make(map[string]interface{}, 0)
			newEnv[clusterNameFieldName] = envInSecureCN.KubernetesClusterName
			newEnv[namespacesNamesFieldName] = envInSecureCN.NamespaceLabels
			newEnv[namespacesNamesFieldName] = envInSecureCN.Namespaces
			envsInTf = append(envsInTf, newEnv)
		}
	}

	err := d.Set(kubernetesEnvironmentFieldName, envsInTf)
	return err
}
