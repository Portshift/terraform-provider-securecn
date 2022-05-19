package securecn

import (
	"context"
	"errors"
	"log"
	"strings"
	"terraform-provider-securecn/internal/client"
	model2 "terraform-provider-securecn/internal/escher_api/model"
	utils2 "terraform-provider-securecn/internal/utils"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const serverlessRuleGroupName = "Terraform automated rules"

const serverlessRuleNameFieldName = "rule_name"
const serverlessRuleActionFieldName = "action"
const serverlessRuleStatusFieldName = "status"
const serverlessRuleScopeFieldName = "scope"
const serverlessFunctionValidationFieldName = "serverless_function_validation"

const matchByFunctionNameFieldName = "match_by_function_name"
const matchByFunctionArnFieldName = "match_by_function_arn"
const matchByFunctionAnyFieldName = "match_by_function_any"

const serverlessRuleNamesFieldName = "names"
const serverlessRuleArnsFieldName = "arns"
const validationFieldRisk = "risk"
const validationFieldVulnerability = "vulnerability"
const validationFieldSecretsRisk = "secrets_risk"
const validationFieldFunctionPermissionRisk = "function_permission_risk"
const validationFieldPubliclyAccessibleRisk = "publicly_accessible_risk"
const validationFieldDataAccessRisk = "data_access_risk"
const validationFieldIsUnusedFunction = "is_unused_function"
const scopeFieldCloudAccount = "cloud_account"
const scopeFieldRegions = "regions"

func ResourceServerlessRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerlessRuleCreate,
		ReadContext:   resourceServerlessRuleRead,
		UpdateContext: resourceServerlessRuleUpdate,
		DeleteContext: resourceServerlessRuleDelete,
		Description:   "A SecureCN serverless rule",
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			serverlessRuleNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			serverlessRuleStatusFieldName: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ENABLED",
				ValidateFunc: validation.StringInSlice([]string{"ENABLED"}, false),
			},
			serverlessRuleActionFieldName: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ALLOW",
				ValidateFunc: validation.StringInSlice([]string{"ALLOW"}, false),
			},
			serverlessRuleScopeFieldName: {
				Description: "Scope defines the scope of this rule",
				Optional:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						scopeFieldCloudAccount: {
							Type:     schema.TypeString,
							Optional: true,
						},
						scopeFieldRegions: {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			matchByFunctionNameFieldName: {
				Description:  "The rule will match using function names",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{matchByFunctionArnFieldName, matchByFunctionAnyFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						serverlessRuleNamesFieldName: {
							Type:     schema.TypeList,
							MinItems: 1,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			matchByFunctionArnFieldName: {
				Description:  "The rule will match using function arns",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{matchByFunctionNameFieldName, matchByFunctionAnyFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						serverlessRuleArnsFieldName: {
							Required: true,
							Type:     schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			matchByFunctionAnyFieldName: {
				Description:  "The rule will match on any function",
				Type:         schema.TypeBool,
				Optional:     true,
				ExactlyOneOf: []string{matchByFunctionNameFieldName, matchByFunctionArnFieldName},
			},
			serverlessFunctionValidationFieldName: {
				Description: "Define function security validations",
				Type:        schema.TypeList,
				MaxItems:    1,
				MinItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						validationFieldRisk: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"NO_RISK", "LOW", "MEDIUM", "HIGH", "CRITICAL",
							}, true),
						},
						validationFieldVulnerability: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL",
							}, true),
						},
						validationFieldSecretsRisk: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"NO_KNOWN_RISK", "RISK_IDENTIFIED",
							}, true),
						},
						validationFieldFunctionPermissionRisk: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"NO_RISK", "LOW", "MEDIUM", "HIGH", "CRITICAL",
							}, true),
						},
						validationFieldPubliclyAccessibleRisk: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"NO_RISK", "LOW", "MEDIUM",
							}, true),
						},
						validationFieldDataAccessRisk: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"NO_RISK", "LOW",
							}, true),
						},
						validationFieldIsUnusedFunction: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  nil,
						},
					},
				},
			},
		},
	}
}

func resourceServerlessRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] creating serverless rule")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	serverlessRuleFromConfig, err := getServerlessRuleFromConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := serviceApi.CreateServerlessRule(ctx, httpClientWrapper.HttpClient, serverlessRuleFromConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := rule.Payload.ID

	d.SetId(string(ruleId))

	return nil
}

func resourceServerlessRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] reading serverless rule")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	ruleId := d.Id()
	currentRuleInSecureCN, err := serviceApi.GetServerlessRule(ctx, httpClientWrapper.HttpClient, strfmt.UUID(ruleId))
	if err != nil {
		return diag.FromErr(err)
	}

	if currentRuleInSecureCN.Payload.ID == "" || *currentRuleInSecureCN.Payload.Status == "DELETED" {
		// Tell terraform the rule doesn't exist
		d.SetId("")
	} else {
		err = updateServerlessRuleMutableFields(d, currentRuleInSecureCN.Payload)
		return diag.FromErr(err)
	}

	return nil
}

func resourceServerlessRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] updating serverless rule")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	rule, err := getServerlessRuleFromConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	rule.ID = strfmt.UUID(d.Id())

	updatedRule, err := serviceApi.UpdateServerlessRule(ctx, httpClientWrapper.HttpClient, rule, rule.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(updatedRule.Payload.ID))

	return resourceServerlessRuleRead(ctx, d, m)
}

func resourceServerlessRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] deleting serverless rule")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	err := serviceApi.DeleteServerlessRule(ctx, httpClientWrapper.HttpClient, strfmt.UUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	// Tell terraform the rule doesn't exist
	d.SetId("")

	return nil
}

func getServerlessRuleFromConfig(d *schema.ResourceData) (*model2.CdServerlessRule, error) {
	log.Print("[DEBUG] getting serverless rule from config")

	name := d.Get(serverlessRuleNameFieldName).(string)
	action := getServerlessRuleActionFromString(d.Get(serverlessRuleActionFieldName).(string))
	status := getServerlessStatusFromString(d.Get(serverlessRuleStatusFieldName).(string))
	scope, err := getRuleScope(d)
	if err != nil {
		return nil, err
	}

	rule := &model2.CdServerlessRule{
		Action:    &action,
		GroupName: serverlessRuleGroupName,
		ID:        "",
		Name:      &name,
		Scope:     scope,
		Status:    &status,
	}

	ruleType, err := getRuleType(d)
	if err != nil {
		return nil, err
	}

	rule.SetRule(ruleType.(model2.ServerlessRuleType))

	return rule, err
}

func getRuleType(d *schema.ResourceData) (model2.ServerlessRuleType, error) {
	matchByFunctionName := d.Get(matchByFunctionNameFieldName).([]interface{})
	funcValidation, err := getFuncValidationFromConfig(d, serverlessFunctionValidationFieldName)
	if err != nil {
		return nil, err
	}
	if len(matchByFunctionName) != 0 {
		funcNames := utils2.ReadNestedListStringFromTF(d, matchByFunctionNameFieldName, serverlessRuleNamesFieldName, 0)
		ruleType := &model2.FunctionNameServerlessRuleType{
			Names: funcNames,
		}
		ruleType.SetServerlessFunctionValidation(funcValidation)
		return ruleType, nil
	}

	matchByFunctionArns := d.Get(matchByFunctionArnFieldName).([]interface{})
	if len(matchByFunctionArns) != 0 {
		funcArns := utils2.ReadNestedListStringFromTF(d, matchByFunctionArnFieldName, serverlessRuleNamesFieldName, 0)
		ruleType := &model2.FunctionArnServerlessRuleType{
			Arns: funcArns,
		}
		ruleType.SetServerlessFunctionValidation(funcValidation)
		return ruleType, nil
	}

	matchByFunctionAny := d.Get(matchByFunctionAnyFieldName).([]interface{})
	if len(matchByFunctionAny) != 0 {
		ruleType := &model2.FunctionAnyServerlessRuleType{}
		ruleType.SetServerlessFunctionValidation(funcValidation)

		return ruleType, nil
	}

	return nil, errors.New("can't get serverless rule type field from configuration")
}

func getRuleScope(d *schema.ResourceData) ([]*model2.ServerlessRuleScope, error) {
	scope := d.Get(serverlessRuleScopeFieldName).([]interface{})
	if len(scope) == 0 {
		return nil, nil
	}
	scopes := make([]*model2.ServerlessRuleScope, len(scope))
	if len(scope) != 0 {
		for i := 0; i < len(scope); i++ {
			cloudAccount := utils2.ReadNestedStringFromTF(d, serverlessRuleScopeFieldName, scopeFieldCloudAccount, i)
			regions := utils2.ReadNestedListStringFromTF(d, serverlessRuleScopeFieldName, scopeFieldRegions, i)
			ruleScope := &model2.ServerlessRuleScope{
				CloudAccount: &cloudAccount,
				Regions:      regions,
			}
			scopes = append(scopes, ruleScope)
		}
		return scopes, nil
	}

	return nil, errors.New("can't get serverless rule scope field from configuration")
}

func getFuncValidationFromConfig(d *schema.ResourceData, mainField string) (*model2.ServerlessFunctionValidation, error) {
	vulnerability := getVulSeverityFromString(utils2.ReadNestedStringFromTF(d, mainField, validationFieldVulnerability, 0))
	dataAccessRisk := getDataAccessFromString(utils2.ReadNestedStringFromTF(d, mainField, validationFieldDataAccessRisk, 0))
	functionPermissionRisk := getFunctionPermissionRiskFromString(utils2.ReadNestedStringFromTF(d, mainField, validationFieldFunctionPermissionRisk, 0))
	publiclyAccessibleRisk := getPubliclyAccessibleRiskFromString(utils2.ReadNestedStringFromTF(d, mainField, validationFieldPubliclyAccessibleRisk, 0))
	secretsRisk := getSecretsRiskFromString(utils2.ReadNestedStringFromTF(d, mainField, validationFieldSecretsRisk, 0))
	risk := getServerlessFunctionRiskFromString(utils2.ReadNestedStringFromTF(d, mainField, validationFieldRisk, 0))
	isUnusedFunction := utils2.ReadNestedBoolFromTF(d, mainField, validationFieldIsUnusedFunction, 0)

	funcValidation := &model2.ServerlessFunctionValidation{
		Vulnerability:          vulnerability,
		DataAccessRisk:         dataAccessRisk,
		FunctionPermissionRisk: functionPermissionRisk,
		PubliclyAccessibleRisk: publiclyAccessibleRisk,
		SecretsRisk:            secretsRisk,
		Risk:                   risk,
		IsUnusedFunction:       isUnusedFunction,
	}

	return funcValidation, nil
}

func getDataAccessFromString(dataAccessRisk string) model2.ServerlessDataAccessRisk {
	if dataAccessRisk == "" {
		return ""
	}
	dataAccessRiskUpper := strings.ToUpper(dataAccessRisk)

	if dataAccessRiskUpper == "LOW" {
		return model2.ServerlessDataAccessRiskLOW
	}

	return model2.ServerlessDataAccessRiskNORISK
}

func getFunctionPermissionRiskFromString(functionPermissionRisk string) model2.ServerlessPolicyRisk {
	if functionPermissionRisk == "" {
		return ""
	}
	functionPermissionRiskUpper := strings.ToUpper(functionPermissionRisk)

	if functionPermissionRiskUpper == "LOW" {
		return model2.ServerlessPolicyRiskLOW
	} else if functionPermissionRiskUpper == "MEDIUM" {
		return model2.ServerlessPolicyRiskMEDIUM
	} else if functionPermissionRiskUpper == "HIGH" {
		return model2.ServerlessPolicyRiskHIGH
	} else if functionPermissionRiskUpper == "CRITICAL" {
		return model2.ServerlessPolicyRiskCRITICAL
	}

	return model2.ServerlessPolicyRiskNORISK
}

func getPubliclyAccessibleRiskFromString(serverlessPubliclyAccessibleRisk string) model2.ServerlessPubliclyAccessibleRisk {
	if serverlessPubliclyAccessibleRisk == "" {
		return ""
	}
	serverlessPubliclyAccessibleRiskUpper := strings.ToUpper(serverlessPubliclyAccessibleRisk)

	if serverlessPubliclyAccessibleRiskUpper == "LOW" {
		return model2.ServerlessPubliclyAccessibleRiskLOW
	} else if serverlessPubliclyAccessibleRiskUpper == "MEDIUM" {
		return model2.ServerlessPubliclyAccessibleRiskMEDIUM
	}

	return model2.ServerlessPubliclyAccessibleRiskNORISK
}

func getSecretsRiskFromString(serverlessSecretsRisk string) model2.ServerlessSecretsRisk {
	if serverlessSecretsRisk == "" {
		return ""
	}
	serverlessSecretsRiskUpper := strings.ToUpper(serverlessSecretsRisk)

	if serverlessSecretsRiskUpper == "RISK_IDENTIFIED" {
		return model2.ServerlessSecretsRiskRISKIDENTIFIED
	}
	return model2.ServerlessSecretsRiskNOKNOWNRISK
}

func getServerlessFunctionRiskFromString(serverlessFunctionRiskLevel string) model2.ServerlessFunctionRiskLevel {
	if serverlessFunctionRiskLevel == "" {
		return ""
	}
	serverlessFunctionRiskLevelUpper := strings.ToUpper(serverlessFunctionRiskLevel)

	if serverlessFunctionRiskLevelUpper == "LOW" {
		return model2.ServerlessFunctionRiskLevelLOW
	} else if serverlessFunctionRiskLevelUpper == "MEDIUM" {
		return model2.ServerlessFunctionRiskLevelMEDIUM
	} else if serverlessFunctionRiskLevelUpper == "HIGH" {
		return model2.ServerlessFunctionRiskLevelHIGH
	} else if serverlessFunctionRiskLevelUpper == "CRITICAL" {
		return model2.ServerlessFunctionRiskLevelCRITICAL
	}
	return model2.ServerlessFunctionRiskLevelNORISK
}

func getServerlessStatusFromString(status string) model2.ServerlessRuleStatus {
	/*
		for now we support only ENABLED
	*/

	//if status == "ENABLED" {
	return model2.ServerlessRuleStatusENABLED
	//}

}

func getServerlessRuleActionFromString(actionString string) model2.ServerlessRuleAction {
	/*
		for now we support only ALLOW
	*/

	//if action == "ALLOW" {
	return model2.ServerlessRuleActionALLOW
	//}
}

func updateServerlessRuleMutableFields(d *schema.ResourceData, currentRuleInSecureCN *model2.CdServerlessRule) error {
	log.Print("[DEBUG] updating serverless rule mutable fields")

	err := d.Set(serverlessRuleNameFieldName, currentRuleInSecureCN.Name)
	if err != nil {
		return err
	}
	err = d.Set(serverlessRuleActionFieldName, currentRuleInSecureCN.Action)
	if err != nil {
		return err
	}
	err = d.Set(serverlessRuleStatusFieldName, currentRuleInSecureCN.Status)
	if err != nil {
		return err
	}

	err = updateServerlessRuleMutableFieldsValidation(d, currentRuleInSecureCN, err)
	if err != nil {
		return err
	}

	err = updateServerlessRuleMutableFieldsScope(d, currentRuleInSecureCN, err)
	if err != nil {
		return err
	}

	ruleInSecureCN := currentRuleInSecureCN.Rule()
	partTypeInSecureCN := ruleInSecureCN.ServerlessRuleType()
	if partTypeInSecureCN == "FunctionNameServerlessRuleType" {
		_ = d.Set(matchByFunctionArnFieldName, nil)
		_ = d.Set(matchByFunctionAnyFieldName, nil)
		functionInSecureCNNames := ruleInSecureCN.(*model2.FunctionNameServerlessRuleType)
		functionsInTf := make(map[string]interface{})
		functionsInTf[serverlessRuleNamesFieldName] = functionInSecureCNNames.Names
		values := make([]map[string]interface{}, 0, 1)
		values = append(values, functionsInTf)
		err = d.Set(matchByFunctionNameFieldName, values)

	} else if partTypeInSecureCN == "FunctionArnServerlessRuleType" {
		err = d.Set(matchByFunctionNameFieldName, nil)
		err = d.Set(matchByFunctionAnyFieldName, nil)
		functionInSecureCNNames := ruleInSecureCN.(*model2.FunctionArnServerlessRuleType)
		functionsInTf := make(map[string]interface{})
		functionsInTf[serverlessRuleArnsFieldName] = functionInSecureCNNames.Arns
		values := make([]map[string]interface{}, 0, 1)
		values = append(values, functionsInTf)
		err = d.Set(matchByFunctionArnFieldName, values)
	} else if partTypeInSecureCN == "FunctionAnyServerlessRuleType" {
		err = d.Set(matchByFunctionNameFieldName, nil)
		err = d.Set(matchByFunctionArnFieldName, nil)
		functionsInTf := make(map[string]interface{})
		values := make([]map[string]interface{}, 0, 1)
		values = append(values, functionsInTf)
		err = d.Set(matchByFunctionAnyFieldName, values)
	}

	return err
}

func updateServerlessRuleMutableFieldsValidation(d *schema.ResourceData, currentRuleInSecureCN *model2.CdServerlessRule, err error) error {
	funcValidations := make([]map[string]interface{}, 0, 1)
	funcValidation := make(map[string]interface{})
	funcValidation[validationFieldRisk] = currentRuleInSecureCN.Rule().ServerlessFunctionValidation().Risk
	funcValidation[validationFieldVulnerability] = currentRuleInSecureCN.Rule().ServerlessFunctionValidation().Vulnerability
	funcValidation[validationFieldSecretsRisk] = currentRuleInSecureCN.Rule().ServerlessFunctionValidation().SecretsRisk
	funcValidation[validationFieldFunctionPermissionRisk] = currentRuleInSecureCN.Rule().ServerlessFunctionValidation().FunctionPermissionRisk
	funcValidation[validationFieldPubliclyAccessibleRisk] = currentRuleInSecureCN.Rule().ServerlessFunctionValidation().PubliclyAccessibleRisk
	funcValidation[validationFieldDataAccessRisk] = currentRuleInSecureCN.Rule().ServerlessFunctionValidation().DataAccessRisk
	funcValidation[validationFieldIsUnusedFunction] = currentRuleInSecureCN.Rule().ServerlessFunctionValidation().IsUnusedFunction
	funcValidations = append(funcValidations, funcValidation)
	err = d.Set(serverlessFunctionValidationFieldName, funcValidations)
	if err != nil {
		return err
	}
	return nil
}

func updateServerlessRuleMutableFieldsScope(d *schema.ResourceData, currentRuleInSecureCN *model2.CdServerlessRule, err error) error {
	funcScope := make([]map[string]interface{}, 0, 1)
	scopeInSecureCN := currentRuleInSecureCN.Scope
	for _, singleScopeInSecureCN := range scopeInSecureCN {
		singleScopeInTf := make(map[string]interface{})
		cloudAccountInSecureCN := *singleScopeInSecureCN.CloudAccount
		regionsInSecureCn := singleScopeInSecureCN.Regions
		regionsInTf := make([]string, 0, len(regionsInSecureCn))
		for _, singleRegionInSingleScopeSecureCN := range regionsInTf {
			regionsInTf = append(regionsInTf, singleRegionInSingleScopeSecureCN)
		}
		singleScopeInTf[scopeFieldCloudAccount] = cloudAccountInSecureCN
		singleScopeInTf[scopeFieldRegions] = regionsInTf
		funcScope = append(funcScope, singleScopeInTf)
	}

	err = d.Set(serverlessRuleScopeFieldName, funcScope)
	if err != nil {
		return err
	}
	return nil
}
