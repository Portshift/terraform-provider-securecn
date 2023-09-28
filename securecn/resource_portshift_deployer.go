package securecn

import (
	"context"
	"fmt"
	"log"
	"terraform-provider-securecn/internal/client"
	"terraform-provider-securecn/internal/escher_api/escherClient"
	model2 "terraform-provider-securecn/internal/escher_api/model"
	utils2 "terraform-provider-securecn/internal/utils"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDeployer() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceDeployerCreate,
		ReadContext:   resourceDeployerRead,
		UpdateContext: resourceDeployerUpdate,
		DeleteContext: resourceDeployerDelete,
		Description:   "A SecureCN deployer",
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			nameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator_deployer": {
				Description: "Create and modify an operator deployer's properties",
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Description:  "The id of the kubernetes cluster in SecureCN of this deployer",
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsUUID,
						},
						"namespace": {
							Description: "The namespace of the ServiceAccount of this deployer",
							Required:    true,
							ForceNew:    true,
							Type:        schema.TypeString,
						},
						"service_account": {
							Description: "The Kubernetes ServiceAccount name of the deployer",
							Required:    true,
							Type:        schema.TypeString,
						},
						"security_check": {
							Description: "Enable security checks for this deployer",
							Optional:    true,
							Type:        schema.TypeBool,
							Default:     false,
						},
					},
				},
			},
		},
	}
}

func resourceDeployerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] creating deployer")

	err := validateDeployerConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	deployerFromConfig, err := getDeployerFromConfig(ctx, d, serviceApi)
	if err != nil {
		return diag.FromErr(err)
	}

	deployer, err := serviceApi.CreateDeployer(ctx, deployerFromConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(deployer.Payload.ID()))

	return nil
}

func resourceDeployerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] reading deployer")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	deployerId := d.Id()

	deployer, err := serviceApi.GetDeployerById(ctx, strfmt.UUID(deployerId))
	if err != nil {
		return diag.FromErr(err)
	}

	if deployer == nil {
		// Tell terraform the deployer doesn't exist
		d.SetId("")
	} else {
		return diag.FromErr(updateDeployerMutableFields(d, deployer))
	}

	return nil
}

func updateDeployerMutableFields(d *schema.ResourceData, deployer model2.Deployer) error {
	return nil
}

func resourceDeployerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] updating deployer")

	err := validateDeployerConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	deployer, err := getDeployerFromConfig(ctx, d, serviceApi)
	if err != nil {
		return diag.FromErr(err)
	}

	deployer.SetID(strfmt.UUID(d.Id()))

	updatedDeployer, err := serviceApi.UpdateDeployer(ctx, deployer)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(updatedDeployer.Payload.ID()))

	return resourceDeployerRead(ctx, d, m)
}

func resourceDeployerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] deleting deployer")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	_, err := serviceApi.DeleteDeployer(ctx, strfmt.UUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	// Tell terraform the deployer doesn't exist
	d.SetId("")

	return nil
}

func validateDeployerConfig(d *schema.ResourceData) error {
	log.Printf("[DEBUG] validating deployer config")

	if d.Get("operator_deployer") == nil {
		return fmt.Errorf("currently only operator_deployer is supported and operator_deployer is mandatory")
	}

	return nil
}

func getDeployerFromConfig(ctx context.Context, d *schema.ResourceData, api *escherClient.MgmtServiceApiCtx) (model2.Deployer, error) {
	log.Print("[DEBUG] getting deployer from config")

	name := d.Get(nameFieldName).(string)
	clusterId := utils2.ReadNestedStringFromTF(d, "operator_deployer", "cluster_id", 0)
	namespaceName := utils2.ReadNestedStringFromTF(d, "operator_deployer", "namespace", 0)
	securityCheck := utils2.ReadNestedBoolFromTF(d, "operator_deployer", "security_check", 0)
	ruleCreation := utils2.ReadNestedBoolFromTF(d, "operator_deployer", "rule_creation", 0)
	serviceAccountName := utils2.ReadNestedStringFromTF(d, "operator_deployer", "service_account", 0)

	serviceAccounts, err := api.GetDeployersServiceAccountsByNamespace(ctx, strfmt.UUID(clusterId), namespaceName)
	if err != nil {
		return nil, err
	}

	var serviceAccount *model2.ServiceAccountInfo
	for _, sai := range serviceAccounts.Payload {
		if sai.Name == serviceAccountName {
			serviceAccount = sai
		}
	}
	if serviceAccount == nil {
		return nil, fmt.Errorf("failed to find %s service account on cluster: %+v", serviceAccountName, serviceAccounts.Payload)
	}

	namespaces, err := api.GetKubernetesClustersKubernetesClusterIDNamespaces(ctx, strfmt.UUID(clusterId))
	if err != nil {
		return nil, err
	}

	var namespace *model2.KubernetesNamespaceResponse
	for _, ns := range namespaces.Payload {
		if ns.Name == namespaceName {
			namespace = ns
		}
	}
	if namespace == nil {
		return nil, fmt.Errorf("failed to find %s namespace on cluster", namespaceName)
	}

	deployer := &model2.OperatorDeployer{
		ClusterID:     strfmt.UUID(clusterId),
		SecurityCheck: &securityCheck,
		RuleCreation:  &ruleCreation,
		NamespaceID:   namespace.ID,
	}

	deployer.SetDeployer(name)
	deployer.SetDeployerID(&serviceAccount.ID)

	return deployer, nil
}
