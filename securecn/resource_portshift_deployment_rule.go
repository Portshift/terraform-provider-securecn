package securecn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"terraform-provider-securecn/internal/client"
	"terraform-provider-securecn/internal/escher_api/escherClient"
	model2 "terraform-provider-securecn/internal/escher_api/model"
	utils2 "terraform-provider-securecn/internal/utils"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const deploymentRuleGroupName = "Terraform automated rules"

const deploymentRuleNameFieldName = "rule_name"
const deploymentRuleActionFieldName = "action"
const deploymentRuleStatusFieldName = "status"
const deploymentRuleScopeFieldName = "scope"

const matchByPodNameFieldName = "match_by_pod_name"
const matchByPodLabelFieldName = "match_by_pod_label"
const matchByPodAnyFieldName = "match_by_pod_any"

const deploymentRuleNamesFieldName = "names"
const deploymentRuleLabelsFieldName = "labels"
const deploymentRuleVulnerabilitySeverityFieldName = "vulnerability_severity_level"
const deploymentRuleVulnerabilityOnViolationActionFieldName = "vulnerability_on_violation_action"
const deploymentRulePSPProfileFieldName = "psp_profile"
const deploymentRulePSPOnViolationActionFieldName = "psp_on_violation_action"

func ResourceDeploymentRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeploymentRuleCreate,
		ReadContext:   resourceDeploymentRuleRead,
		UpdateContext: resourceDeploymentRuleUpdate,
		DeleteContext: resourceDeploymentRuleDelete,
		Description:   "A SecureCN deployment rule",
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			deploymentRuleNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			deploymentRuleActionFieldName: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ALLOW",
				ValidateFunc: validation.StringInSlice([]string{"ALLOW"}, false),
			},
			deploymentRuleStatusFieldName: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ENABLED",
				ValidateFunc: validation.StringInSlice([]string{"ENABLED"}, false),
			},
			deploymentRuleScopeFieldName: {
				Description:  "Scope defines the scope of this rule",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ANY",
				ValidateFunc: validation.StringInSlice([]string{"ANY", "CLUSTER", "ENVIRONMENT"}, true),
			},
			matchByPodNameFieldName: {
				Description:  "The rule will match using pod names",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{matchByPodLabelFieldName, matchByPodAnyFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						deploymentRuleNamesFieldName: {
							Type:     schema.TypeList,
							MinItems: 1,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						deploymentRuleVulnerabilitySeverityFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL",
							}, true),
						},
						deploymentRuleVulnerabilityOnViolationActionFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"BLOCK", "DETECT",
							}, true),
						},
						deploymentRulePSPProfileFieldName: {
							Optional: true,
							Type:     schema.TypeString,
						},
						deploymentRulePSPOnViolationActionFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"BLOCK", "DETECT", "ENFORCE",
							}, true),
						},
					},
				},
			},
			matchByPodLabelFieldName: {
				Description:  "The rule will match using pod labels",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{matchByPodNameFieldName, matchByPodAnyFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						deploymentRuleLabelsFieldName: {
							Required: true,
							Type:     schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						deploymentRuleVulnerabilitySeverityFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL",
							}, true),
						},
						deploymentRuleVulnerabilityOnViolationActionFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"BLOCK", "DETECT",
							}, true),
						},
						deploymentRulePSPProfileFieldName: {
							Optional: true,
							Type:     schema.TypeString,
						},
						deploymentRulePSPOnViolationActionFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"BLOCK", "DETECT", "ENFORCE",
							}, true),
						},
					},
				},
			},
			matchByPodAnyFieldName: {
				Description:  "The rule will match on any pod",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{matchByPodNameFieldName, matchByPodLabelFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						deploymentRuleVulnerabilitySeverityFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL",
							}, true),
						},
						deploymentRuleVulnerabilityOnViolationActionFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"BLOCK", "DETECT",
							}, true),
						},
						deploymentRulePSPProfileFieldName: {
							Optional: true,
							Type:     schema.TypeString,
						},
						deploymentRulePSPOnViolationActionFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"BLOCK", "DETECT", "ENFORCE",
							}, true),
						},
					},
				},
			},
		},
	}
}

func resourceDeploymentRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] creating deployment rule")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	deploymentRuleFromConfig, err := getDeploymentRuleFromConfig(ctx, d, serviceApi, httpClientWrapper)
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := serviceApi.CreateDeploymentRule(ctx, httpClientWrapper.HttpClient, deploymentRuleFromConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := rule.Payload.ID

	d.SetId(string(ruleId))

	return nil
}

func resourceDeploymentRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] reading deployment rule")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	ruleId := d.Id()

	currentRuleInSecureCN, err := serviceApi.GetDeploymentRule(ctx, httpClientWrapper.HttpClient, strfmt.UUID(ruleId))
	if err != nil {
		return diag.FromErr(err)
	}

	if currentRuleInSecureCN.Payload.ID == "" || currentRuleInSecureCN.Payload.Status == "DELETED" {
		// Tell terraform the rule doesn't exist
		d.SetId("")
	} else {
		err = updateDeploymentRuleMutableFields(d, currentRuleInSecureCN.Payload)
		return diag.FromErr(err)
	}

	return nil
}

func resourceDeploymentRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] updating deployment rule")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	rule, err := getDeploymentRuleFromConfig(ctx, d, serviceApi, httpClientWrapper)
	if err != nil {
		return diag.FromErr(err)
	}

	rule.ID = strfmt.UUID(d.Id())

	updatedRule, err := serviceApi.UpdateDeploymentRule(ctx, httpClientWrapper.HttpClient, rule, rule.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(updatedRule.Payload.ID))

	return resourceDeploymentRuleRead(ctx, d, m)
}

func resourceDeploymentRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] deleting deployment rule")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	err := serviceApi.DeleteDeploymentRule(ctx, httpClientWrapper.HttpClient, strfmt.UUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	// Tell terraform the rule doesn't exist
	d.SetId("")

	return nil
}

func validateDeploymentRuleFromConfig(deploymentRuleFromConfig *model2.CdAppRule) error {
	log.Printf("[DEBUG] validating deployment rule config")
	if deploymentRuleFromConfig.App().WorkloadRuleType() == "PodNameWorkloadRuleType" {
		app := deploymentRuleFromConfig.App().(*model2.PodNameWorkloadRuleType)
		pspPolicy := app.PodValidation.PodSecurityPolicy
		vulPolicy := app.PodValidation.Vulnerability
		err := validatePodValidation(pspPolicy, vulPolicy)
		if err != nil {
			return err
		}
	}

	if deploymentRuleFromConfig.App().WorkloadRuleType() == "PodLabelWorkloadRuleType" {
		app := deploymentRuleFromConfig.App().(*model2.PodLabelWorkloadRuleType)
		pspPolicy := app.PodValidation.PodSecurityPolicy
		vulPolicy := app.PodValidation.Vulnerability
		err := validatePodValidation(pspPolicy, vulPolicy)
		if err != nil {
			return err
		}
	}

	if deploymentRuleFromConfig.App().WorkloadRuleType() == "PodAnyWorkloadRuleType" {
		app := deploymentRuleFromConfig.App().(*model2.PodAnyWorkloadRuleType)
		pspPolicy := app.PodValidation.PodSecurityPolicy
		vulPolicy := app.PodValidation.Vulnerability
		err := validatePodValidation(pspPolicy, vulPolicy)
		if err != nil {
			return err
		}
	}

	return nil
}

func validatePodValidation(pspPolicy *model2.PodSecurityPolicyValidation, vulPolicy *model2.VulnerabilityValidation) error {
	if pspPolicy != nil && ((pspPolicy.OnViolationAction != "" && pspPolicy.PodSecurityPolicyID == nil) || (pspPolicy.OnViolationAction == "" && pspPolicy.PodSecurityPolicyID != nil)) {
		return fmt.Errorf("invalid psp policy configuration. if 1 is set, the other must also be set. %s, %s", deploymentRulePSPProfileFieldName, deploymentRulePSPOnViolationActionFieldName)
	}
	if vulPolicy != nil && ((vulPolicy.OnViolationAction != "" && vulPolicy.HighestVulnerabilityAllowed == "") || (vulPolicy.OnViolationAction == "" && vulPolicy.HighestVulnerabilityAllowed != "")) {
		return fmt.Errorf("invalid vulnerability policy configuration. if one of the the fields is set, the other must also be set. %s, %s", deploymentRuleVulnerabilitySeverityFieldName, deploymentRuleVulnerabilityOnViolationActionFieldName)
	}
	return nil
}

func getDeploymentRuleFromConfig(ctx context.Context, d *schema.ResourceData, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper) (*model2.CdAppRule, error) {
	log.Print("[DEBUG] getting deployment rule from config")

	name := d.Get(deploymentRuleNameFieldName).(string)
	action := getRuleActionFromString(d.Get(deploymentRuleActionFieldName).(string))
	status := getStatusFromString(d.Get(deploymentRuleStatusFieldName).(string))
	scope := getScopeFromString(d.Get(deploymentRuleScopeFieldName).(string))

	rule := &model2.CdAppRule{
		ID:        "",
		Action:    action,
		GroupName: deploymentRuleGroupName,
		Name:      &name,
		Status:    status,
		Scope:     scope,
	}

	app, err := getApp(ctx, d, serviceApi, httpClientWrapper)
	if err != nil {
		return nil, err
	}

	rule.SetApp(app.(model2.WorkloadRuleType))

	err = validateDeploymentRuleFromConfig(rule)

	return rule, err
}

func getScopeFromString(scope string) model2.WorkloadRuleScopeType {
	switch strings.ToLower(scope) {
	case "cluster":
		return model2.WorkloadRuleScopeTypeClusterNameRuleType
	case "environment":
		return model2.WorkloadRuleScopeTypeEnvironmentNameRuleType
	default:
		return model2.WorkloadRuleScopeTypeAnyRuleType
	}
}

func getApp(ctx context.Context, d *schema.ResourceData, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper) (model2.WorkloadRuleType, error) {
	matchByPodName := d.Get(matchByPodNameFieldName).([]interface{})
	if len(matchByPodName) != 0 {
		podNames := utils2.ReadNestedListStringFromTF(d, matchByPodNameFieldName, deploymentRuleNamesFieldName, 0)
		podValidation, err := getPodValidationFromConfig(ctx, d, matchByPodNameFieldName, serviceApi, httpClientWrapper)
		if err != nil {
			return nil, err
		}
		app := &model2.PodNameWorkloadRuleType{
			Names:         podNames,
			PodValidation: podValidation,
		}

		return app, nil
	}

	matchByPodLabel := d.Get(matchByPodLabelFieldName).([]interface{})
	if len(matchByPodLabel) != 0 {
		podLabels := utils2.GetLabelsFromMap(utils2.ReadNestedMapStringFromTF(d, matchByPodLabelFieldName, deploymentRuleLabelsFieldName, 0))
		podValidation, err := getPodValidationFromConfig(ctx, d, matchByPodLabelFieldName, serviceApi, httpClientWrapper)
		if err != nil {
			return nil, err
		}
		app := &model2.PodLabelWorkloadRuleType{
			Labels:        podLabels,
			PodValidation: podValidation,
		}

		return app, nil
	}

	matchByPodAny := d.Get(matchByPodAnyFieldName).([]interface{})
	if len(matchByPodAny) != 0 {
		podValidation, err := getPodValidationFromConfig(ctx, d, matchByPodAnyFieldName, serviceApi, httpClientWrapper)
		if err != nil {
			return nil, err
		}
		app := &model2.PodAnyWorkloadRuleType{
			PodValidation: podValidation,
		}

		return app, nil
	}

	return nil, errors.New("can't get deployment rule app field from configuration")
}

func getPodValidationFromConfig(ctx context.Context, d *schema.ResourceData, mainField string, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper) (*model2.PodValidation, error) {
	vulnerabilityValidation := getVulnerabilityValidationFromConfig(d, mainField)
	podSecurityPolicyValidation, err := getPspValidationFromConfig(ctx, d, mainField, serviceApi, httpClientWrapper)
	if err != nil {
		return nil, err
	}

	podValidation := &model2.PodValidation{
		PodSecurityPolicy: podSecurityPolicyValidation,
		Vulnerability:     vulnerabilityValidation,
	}
	return podValidation, nil
}

func getPspValidationFromConfig(ctx context.Context, d *schema.ResourceData, mainField string, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper) (*model2.PodSecurityPolicyValidation, error) {
	pspProfileName := utils2.ReadNestedStringFromTF(d, mainField, deploymentRulePSPProfileFieldName, 0)
	if pspProfileName == "" {
		return nil, nil
	}
	pspProfileId, err := getPspProfileIdFromName(ctx, serviceApi, httpClientWrapper, pspProfileName)
	if err != nil {
		return nil, fmt.Errorf("%v\nmake sure a profile with that name exists: %s", err, pspProfileName)
	}

	actionString := utils2.ReadNestedStringFromTF(d, mainField, deploymentRulePSPOnViolationActionFieldName, 0)
	pspAction := getOnViolationActionFromString(actionString)
	shouldMutate := isEnforce(actionString)
	podSecurityPolicyValidation := &model2.PodSecurityPolicyValidation{
		OnViolationAction:   pspAction,
		PodSecurityPolicyID: pspProfileId,
		ShouldMutate:        &shouldMutate,
	}
	return podSecurityPolicyValidation, nil
}

func getPspProfileIdFromName(ctx context.Context, serviceApi *escherClient.MgmtServiceApiCtx, wrapper client.HttpClientWrapper, podSecurityPolicyProfileName string) (*strfmt.UUID, error) {
	if podSecurityPolicyProfileName == "" {
		return nil, nil
	}
	pspId, err := serviceApi.GetPspIdByName(ctx, wrapper.HttpClient, podSecurityPolicyProfileName)

	if err != nil {
		return nil, fmt.Errorf("failed to get psp profile by name: %v. %v", podSecurityPolicyProfileName, err)
	}

	return &pspId.Payload, nil

}

func getVulnerabilityValidationFromConfig(d *schema.ResourceData, mainField string) *model2.VulnerabilityValidation {
	vulSeverity := getVulSeverityFromString(utils2.ReadNestedStringFromTF(d, mainField, deploymentRuleVulnerabilitySeverityFieldName, 0))
	vulAction := getOnViolationActionFromString(utils2.ReadNestedStringFromTF(d, mainField, deploymentRuleVulnerabilityOnViolationActionFieldName, 0))
	if vulSeverity == "" || vulAction == "" {
		return nil
	}
	vulnerabilityValidation := &model2.VulnerabilityValidation{
		HighestVulnerabilityAllowed: vulSeverity,
		OnViolationAction:           vulAction,
	}
	return vulnerabilityValidation
}

func getOnViolationActionFromString(action string) model2.OnViolationAction {
	if action == "" {
		return ""
	}
	actionUpper := strings.ToUpper(action)

	if actionUpper == "DETECT" {
		return model2.OnViolationActionDETECT
	}
	if actionUpper == "BLOCK" {
		return model2.OnViolationActionBLOCK
	}

	// ENFORCE can only be ENFORCE_AND_DETECT for now
	return model2.OnViolationActionDETECT
}

func getVulSeverityFromString(severity string) model2.VulnerabilitySeverity {
	if severity == "" {
		return ""
	}
	severityUpper := strings.ToUpper(severity)

	if severityUpper == "LOW" {
		return model2.VulnerabilitySeverityLOW
	} else if severityUpper == "MEDIUM" {
		return model2.VulnerabilitySeverityMEDIUM
	} else if severityUpper == "HIGH" {
		return model2.VulnerabilitySeverityHIGH
	} else if severityUpper == "CRITICAL" {
		return model2.VulnerabilitySeverityCRITICAL
	}

	return model2.VulnerabilitySeverityUNKNOWN
}

func isEnforce(pspAction string) bool {
	return pspAction == "ENFORCE"
}

func getStatusFromString(status string) model2.AppRuleStatus {
	/*
		for now we support only ENABLED
	*/

	//if status == "ENABLED" {
	return model2.AppRuleStatusENABLED
	//}

}

func getRuleActionFromString(actionString string) model2.AppRuleType {
	/*
		for now we support only ALLOW
	*/

	//if action == "ALLOW" {
	return model2.AppRuleTypeALLOW
	//}
}

func updateDeploymentRuleMutableFields(d *schema.ResourceData, currentRuleInSecureCN *model2.CdAppRule) error {
	log.Print("[DEBUG] updating deployment rule mutable fields")

	err := d.Set(deploymentRuleNameFieldName, currentRuleInSecureCN.Name)
	if err != nil {
		return err
	}
	err = d.Set(deploymentRuleActionFieldName, currentRuleInSecureCN.Action)
	if err != nil {
		return err
	}
	err = d.Set(deploymentRuleStatusFieldName, currentRuleInSecureCN.Status)
	if err != nil {
		return err
	}

	appInSecureCN := currentRuleInSecureCN.App()

	partTypeInSecureCN := appInSecureCN.WorkloadRuleType()
	if partTypeInSecureCN == "PodNameWorkloadRuleType" {
		_ = d.Set(matchByPodLabelFieldName, nil)
		_ = d.Set(matchByPodAnyFieldName, nil)
		appInSecureCNNames := appInSecureCN.(*model2.PodNameWorkloadRuleType)
		appsInTf := make(map[string]interface{})
		appsInTf[deploymentRuleVulnerabilitySeverityFieldName] = appInSecureCNNames.PodValidation.Vulnerability.HighestVulnerabilityAllowed
		appsInTf[deploymentRuleVulnerabilityOnViolationActionFieldName] = appInSecureCNNames.PodValidation.Vulnerability.OnViolationAction
		appsInTf[deploymentRulePSPProfileFieldName] = appInSecureCNNames.PodValidation.PodSecurityPolicy.PodSecurityPolicyName
		mutate := appInSecureCNNames.PodValidation.PodSecurityPolicy.ShouldMutate
		pspAction := appInSecureCNNames.PodValidation.PodSecurityPolicy.OnViolationAction
		if *mutate {
			pspAction = "ENFORCE"
		}
		appsInTf[deploymentRulePSPOnViolationActionFieldName] = pspAction
		appsInTf[deploymentRuleNamesFieldName] = appInSecureCNNames.Names
		values := make([]map[string]interface{}, 0, 1)
		values = append(values, appsInTf)
		err = d.Set(matchByPodNameFieldName, values)
	} else if partTypeInSecureCN == "PodLabelWorkloadRuleType" {
		err = d.Set(matchByPodNameFieldName, nil)
		err = d.Set(matchByPodAnyFieldName, nil)
		appInSecureCNLabels := appInSecureCN.(*model2.PodLabelWorkloadRuleType)
		appsInTf := make(map[string]interface{})
		appsInTf[deploymentRuleVulnerabilitySeverityFieldName] = appInSecureCNLabels.PodValidation.Vulnerability.HighestVulnerabilityAllowed
		appsInTf[deploymentRuleVulnerabilityOnViolationActionFieldName] = appInSecureCNLabels.PodValidation.Vulnerability.OnViolationAction
		appsInTf[deploymentRulePSPProfileFieldName] = appInSecureCNLabels.PodValidation.PodSecurityPolicy.PodSecurityPolicyName
		mutate := appInSecureCNLabels.PodValidation.PodSecurityPolicy.ShouldMutate
		pspAction := appInSecureCNLabels.PodValidation.PodSecurityPolicy.OnViolationAction
		if *mutate {
			pspAction = "ENFORCE"
		}
		appsInTf[deploymentRulePSPOnViolationActionFieldName] = pspAction
		appsInTf[deploymentRuleLabelsFieldName] = utils2.GetListStringFromLabels(appInSecureCNLabels.Labels)
		values := make([]map[string]interface{}, 0, 1)
		values = append(values, appsInTf)
		err = d.Set(matchByPodLabelFieldName, values)
	} else if partTypeInSecureCN == "PodAnyWorkloadRuleType" {
		err = d.Set(matchByPodNameFieldName, nil)
		err = d.Set(matchByPodLabelFieldName, nil)
		appInSecureCNAny := appInSecureCN.(*model2.PodAnyWorkloadRuleType)
		appsInTf := make(map[string]interface{})
		appsInTf[deploymentRuleVulnerabilitySeverityFieldName] = appInSecureCNAny.PodValidation.Vulnerability.HighestVulnerabilityAllowed
		appsInTf[deploymentRuleVulnerabilityOnViolationActionFieldName] = appInSecureCNAny.PodValidation.Vulnerability.OnViolationAction
		appsInTf[deploymentRulePSPProfileFieldName] = appInSecureCNAny.PodValidation.PodSecurityPolicy.PodSecurityPolicyName
		mutate := appInSecureCNAny.PodValidation.PodSecurityPolicy.ShouldMutate
		pspAction := appInSecureCNAny.PodValidation.PodSecurityPolicy.OnViolationAction
		if *mutate {
			pspAction = "ENFORCE"
		}
		appsInTf[deploymentRulePSPOnViolationActionFieldName] = pspAction
		values := make([]map[string]interface{}, 0, 1)
		values = append(values, appsInTf)
		err = d.Set(matchByPodAnyFieldName, values)
	}

	return err
}
