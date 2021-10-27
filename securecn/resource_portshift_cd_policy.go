package securecn

import (
	"context"
	"log"
	"terraform-provider-securecn/internal/client"
	"terraform-provider-securecn/internal/escher_api/escherClient"
	model2 "terraform-provider-securecn/internal/escher_api/model"
	utils2 "terraform-provider-securecn/internal/utils"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/spf13/cast"
)

var enforcementOptionSchema = &schema.Schema{
	Description:  "The enforcement type for this policy",
	Required:     true,
	Type:         schema.TypeString,
	ValidateFunc: validation.StringInSlice([]string{"FAIL", "IGNORE"}, false),
}

var cdPolicyElementSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"permissible_vulnerability_level": {
			Description:  "The level of risk accepted in this policy",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"NO_RISK", "MEDIUM", "HIGH"}, false),
		},
		"enforcement_option": enforcementOptionSchema,
	},
}

var secretPolicyElementSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"permissible_vulnerability_level": {
			Description:  "The level of risk accepted in this policy",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"NO_KNOWN_RISK", "RISK_IDENTIFIED"}, false),
		},
		"enforcement_option": enforcementOptionSchema,
	},
}

func ResourceCdPolicy() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceCdPolicyCreate,
		ReadContext:   resourceCdPolicyRead,
		UpdateContext: resourceCdPolicyUpdate,
		DeleteContext: resourceCdPolicyDelete,
		Description:   "A SecureCN CD policy",
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
			"deployers": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},
			"api_security_policy": {
				Description: "Specify the cd policy's api security profile",
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_security_profile": {
							Description:  "The id of the api security profile to use for this api policy",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"enforcement_option": enforcementOptionSchema,
					},
				},
			},
			"permission_policy": {
				Description: "Specify the cd policy's permission check profile",
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        cdPolicyElementSchema,
			},
			"secret_policy": {
				Description: "Specify the cd policy's secret check profile",
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        secretPolicyElementSchema,
			},
			"security_context_policy": {
				Description: "Specify the cd policy's security context check profile",
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        cdPolicyElementSchema,
			},
		},
	}
}

func resourceCdPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] creating cd policy")

	err := validateCdPolicyConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	cdPolicy, err := getCdPolicyFromConfig(d, serviceApi)
	if err != nil {
		return diag.FromErr(err)
	}

	params := model2.PostCdPolicyParams{
		Body:    cdPolicy,
		Context: ctx,
	}

	cdPolicyCreated, err := serviceApi.PostCdPolicy(&params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(cdPolicyCreated.Payload.ID))

	return nil
}

func resourceCdPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] reading cd policy")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	params := model2.GetCdPolicyParams{
		Context: ctx,
	}

	cdPolicies, err := serviceApi.GetCdPolicy(&params)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, cdPolicy := range cdPolicies.Payload {
		if string(cdPolicy.ID) == d.Id() {
			return diag.FromErr(updateCdPolicyMutableFields(d, cdPolicy))
		}
	}

	// Tell terraform the cd policy doesn't exist
	d.SetId("")

	return nil
}

func updateCdPolicyMutableFields(d *schema.ResourceData, policy *model2.CdPolicy) error {
	return nil
}

func resourceCdPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] updating cd policy")

	err := validateCdPolicyConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	cdPolicy, err := getCdPolicyFromConfig(d, serviceApi)
	if err != nil {
		return diag.FromErr(err)
	}

	params := model2.PutCdPolicyPolicyIDParams{
		PolicyID: strfmt.UUID(d.Id()),
		Body:     cdPolicy,
		Context:  ctx,
	}

	cdPolicyUpdated, err := serviceApi.PutCdPolicyPolicyID(&params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(cdPolicyUpdated.Payload.ID))

	return resourceCdPolicyRead(ctx, d, m)
}

func resourceCdPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] deleting cd policy")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	params := model2.DeleteCdPolicyPolicyIDParams{
		PolicyID: strfmt.UUID(d.Id()),
		Context:  ctx,
	}

	_, err := serviceApi.DeleteCdPolicyPolicyID(&params)
	if err != nil {
		return diag.FromErr(err)
	}

	// Tell terraform the cd policy doesn't exist
	d.SetId("")

	return nil
}

func validateCdPolicyConfig(d *schema.ResourceData) error {
	log.Printf("[DEBUG] validating cd policy config")

	return nil
}

func getCdPolicyFromConfig(d *schema.ResourceData, api *escherClient.MgmtServiceApiCtx) (*model2.CdPolicy, error) {
	log.Print("[DEBUG] getting cd policy from config")

	name := d.Get(nameFieldName).(string)
	description := d.Get(descriptionFieldName).(string)
	deployers := cast.ToStringSlice(d.Get("deployers"))

	var deployerUUIDs []strfmt.UUID
	for _, deployerId := range deployers {
		deployerUUIDs = append(deployerUUIDs, strfmt.UUID(deployerId))
	}

	cdPolicy := &model2.CdPolicy{
		Name:        &name,
		Deployers:   deployerUUIDs,
		Description: description,
	}

	apiSecurityProfile := utils2.ReadNestedStringFromTF(d, "api_security_policy", "api_security_profile", 0)
	enforcementOption := utils2.ReadNestedStringFromTF(d, "api_security_policy", "enforcement_option", 0)

	if apiSecurityProfile != "" && enforcementOption != "" {
		apiSecurityProfileUUID := strfmt.UUID(apiSecurityProfile)
		cdPolicy.APISecurityCdPolicy = &model2.APISecurityCdPolicyElement{
			APISecurityProfile: &apiSecurityProfileUUID,
			EnforcementOption:  model2.EnforcementOption(enforcementOption),
		}
	}

	permissibleVulnerabilityLevel := utils2.ReadNestedStringFromTF(d, "permission_policy", "permissible_vulnerability_level", 0)
	enforcementOption = utils2.ReadNestedStringFromTF(d, "permission_policy", "enforcement_option", 0)

	if permissibleVulnerabilityLevel != "" && enforcementOption != "" {
		cdPolicy.PermissionCDPolicy = &model2.CdPolicyElement{
			PermissibleVulnerabilityLevel: model2.Risk(permissibleVulnerabilityLevel),
			EnforcementOption:             model2.EnforcementOption(enforcementOption),
		}
	}

	permissibleVulnerabilityLevel = utils2.ReadNestedStringFromTF(d, "secret_policy", "permissible_vulnerability_level", 0)
	enforcementOption = utils2.ReadNestedStringFromTF(d, "secret_policy", "enforcement_option", 0)

	if permissibleVulnerabilityLevel != "" && enforcementOption != "" {
		cdPolicy.SecretCDPolicy = &model2.SecretsCdPolicyElement{
			PermissibleVulnerabilityLevel: model2.CDPipelineSecretsFindingRisk(permissibleVulnerabilityLevel),
			EnforcementOption:             model2.EnforcementOption(enforcementOption),
		}
	}

	permissibleVulnerabilityLevel = utils2.ReadNestedStringFromTF(d, "security_context_policy", "permissible_vulnerability_level", 0)
	enforcementOption = utils2.ReadNestedStringFromTF(d, "security_context_policy", "enforcement_option", 0)

	if permissibleVulnerabilityLevel != "" && enforcementOption != "" {
		cdPolicy.SecurityContextCDPolicy = &model2.CdPolicyElement{
			PermissibleVulnerabilityLevel: model2.Risk(permissibleVulnerabilityLevel),
			EnforcementOption:             model2.EnforcementOption(enforcementOption),
		}
	}

	return cdPolicy, nil
}
