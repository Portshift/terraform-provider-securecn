package securecn

import (
	"context"
	"log"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"terraform-provider-securecn/client"
	"terraform-provider-securecn/escher_api/escherClient"
	"terraform-provider-securecn/escher_api/model"
	"terraform-provider-securecn/utils"
)

func ResourceCiPolicy() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceCiPolicyCreate,
		ReadContext:   resourceCiPolicyRead,
		UpdateContext: resourceCiPolicyUpdate,
		DeleteContext: resourceCiPolicyDelete,
		Description:   "A SecureCN CI policy",
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
			"vulnerability_policy": {
				Description: "Specify the ci policy's vulnerability policy part",
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permissible_vulnerability_level": {
							Description:  "The level of risk accepted in this policy",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL"}, false),
						},
						"enforcement_option": enforcementOptionSchema,
					},
				},
			},
			"dockerfile_scan_policy": {
				Description: "Specify the ci policy's dockerfile scan policy part",
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permissible_dockerfile_scan_severity": {
							Description:  "The scan result severity accepted in this policy",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"INFO", "WARN", "FATAL"}, false),
						},
						"enforcement_option": enforcementOptionSchema,
					},
				},
			},
		},
	}
}

func resourceCiPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] creating ci policy")

	err := validateCiPolicyConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils.GetServiceApi(&httpClientWrapper)

	ciPolicy, err := getCiPolicyFromConfig(d, serviceApi)
	if err != nil {
		return diag.FromErr(err)
	}

	params := model.PostCiPolicyParams{
		Body:    ciPolicy,
		Context: ctx,
	}

	ciPolicyCreated, err := serviceApi.PostCiPolicy(&params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(ciPolicyCreated.Payload.ID))

	return nil
}

func resourceCiPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] reading ci policy")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils.GetServiceApi(&httpClientWrapper)

	params := model.GetCiPolicyParams{
		Context: ctx,
	}

	ciPolicy, err := serviceApi.GetCiPolicy(&params)
	if err != nil {
		return diag.FromErr(err)
	}

	if string(ciPolicy.Payload.ID) == d.Id() {
		return diag.FromErr(updateCiPolicyMutableFields(d, ciPolicy.Payload))
	}

	// Tell terraform the ci policy doesn't exist
	d.SetId("")

	return nil
}

func updateCiPolicyMutableFields(d *schema.ResourceData, policy *model.CiPolicy) error {
	return nil
}

func resourceCiPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] updating ci policy")

	err := validateCiPolicyConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils.GetServiceApi(&httpClientWrapper)

	ciPolicy, err := getCiPolicyFromConfig(d, serviceApi)
	if err != nil {
		return diag.FromErr(err)
	}

	params := model.PutCiPolicyPolicyIDParams{
		PolicyID: strfmt.UUID(d.Id()),
		Body:     ciPolicy,
		Context:  ctx,
	}

	ciPolicyUpdated, err := serviceApi.PutCiPolicyPolicyID(&params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(ciPolicyUpdated.Payload.ID))

	return resourceCiPolicyRead(ctx, d, m)
}

func resourceCiPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] deleting ci policy")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils.GetServiceApi(&httpClientWrapper)

	params := model.DeleteCiPolicyPolicyIDParams{
		PolicyID: strfmt.UUID(d.Id()),
		Context:  ctx,
	}

	_, err := serviceApi.DeleteCiPolicyPolicyID(&params)
	if err != nil {
		return diag.FromErr(err)
	}

	// Tell terraform the ci policy doesn't exist
	d.SetId("")

	return nil
}

func validateCiPolicyConfig(d *schema.ResourceData) error {
	log.Printf("[DEBUG] validating ci policy config")

	return nil
}

func getCiPolicyFromConfig(d *schema.ResourceData, api *escherClient.MgmtServiceApiCtx) (*model.CiPolicy, error) {
	log.Print("[DEBUG] getting ci policy from config")

	name := d.Get(nameFieldName).(string)
	description := d.Get(descriptionFieldName).(string)

	ciPolicy := &model.CiPolicy{
		Name:        &name,
		Description: description,
	}

	permissibleVulnerabilityLevel := utils.ReadNestedStringFromTF(d, "vulnerability_policy", "permissible_vulnerability_level", 0)
	enforcementOption := utils.ReadNestedStringFromTF(d, "vulnerability_policy", "enforcement_option", 0)

	if permissibleVulnerabilityLevel != "" && enforcementOption != "" {
		ciPolicy.VulnerabilityCiPolicy = &model.CiVulnerabilityPolicy{
			PermissibleVulnerabilityLevel: model.VulnerabilitySeverity(permissibleVulnerabilityLevel),
			EnforcementOption:             model.EnforcementOption(enforcementOption),
		}
	}

	permissibleVulnerabilityLevel = utils.ReadNestedStringFromTF(d, "dockerfile_scan_policy", "permissible_dockerfile_scan_severity", 0)
	enforcementOption = utils.ReadNestedStringFromTF(d, "dockerfile_scan_policy", "enforcement_option", 0)

	if permissibleVulnerabilityLevel != "" && enforcementOption != "" {
		ciPolicy.DockerfileScanCiPolicy = &model.CiDockerfileScanPolicy{
			PermissibleDockerfileScanSeverity: model.DockerfileScanSeverity(permissibleVulnerabilityLevel),
			EnforcementOption:                 model.EnforcementOption(enforcementOption),
		}
	}

	return ciPolicy, nil
}
