package securecn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"regexp"
	"strings"
	"terraform-provider-securecn/internal/client"
	model2 "terraform-provider-securecn/internal/escher_api/model"
	utils2 "terraform-provider-securecn/internal/utils"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const connectionRuleGroupName = "Terraform automated rules"

const connectionRuleNameFieldName = "rule_name"

const connectionRuleActionNameFieldName = "action"
const connectionRuleStatusNameFieldName = "status"
const sourceIpRangeFieldName = "source_by_ip_range"
const sourceExternalFieldName = "source_by_external"
const sourcePodNameFieldName = "source_by_pod_name"
const sourcePodLabelFieldName = "source_by_pod_label"
const sourcePodAnyFieldName = "source_by_pod_any"
const destinationAddressIpRangeFieldName = "destination_by_address_ip_range"
const destinationAddressDomainFieldName = "destination_by_address_domain"
const destinationExternalFieldName = "destination_by_external"
const destinationPodNameFieldName = "destination_by_pod_name"
const destinationPodLabelFieldName = "destination_by_pod_label"
const destinationPodAnyFieldName = "destination_by_pod_any"

const ipsFieldName = "ips"
const domainsFieldName = "domains"
const connectionRuleNamesFieldName = "names"
const connectionRuleNamesLabelsFieldName = "labels"
const connectionRuleVulnerabilitySeverityFieldName = "vulnerability_severity_level"

func ResourceConnectionRule() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceConnectionRuleCreate,
		ReadContext:   resourceConnectionRuleRead,
		UpdateContext: resourceConnectionRuleUpdate,
		DeleteContext: resourceConnectionRuleDelete,
		Description:   "A SecureCN k8s connection rule",
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			connectionRuleNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			connectionRuleActionNameFieldName: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ALLOW",
				ValidateFunc: func(value interface{}, key string) (warns []string, errs []error) {
					action := value.(string)
					if action != "ALLOW" {
						errs = append(errs, fmt.Errorf("only ALLOW action is supported"))
					}
					return
				},
			},
			connectionRuleStatusNameFieldName: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ENABLED",
				ValidateFunc: func(value interface{}, key string) (warns []string, errs []error) {
					status := value.(string)
					if status != "ENABLED" {
						errs = append(errs, fmt.Errorf("only ENABLED status is supported"))
					}
					return
				},
			},
			sourceIpRangeFieldName: {
				Description:  "The source will match using ip ranges",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{sourceExternalFieldName, sourcePodNameFieldName, sourcePodLabelFieldName, sourcePodAnyFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ipsFieldName: {
							Type:     schema.TypeList,
							MinItems: 1,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			sourceExternalFieldName: {
				Description:  "The source will match on external connections",
				Optional:     true,
				ExactlyOneOf: []string{sourceIpRangeFieldName, sourcePodNameFieldName, sourcePodLabelFieldName, sourcePodAnyFieldName},
				Type:         schema.TypeBool,
				Default:      false,
				ValidateFunc: func(value interface{}, key string) (warns []string, errs []error) {
					isExternal := value.(bool)
					if !isExternal {
						errs = append(errs, fmt.Errorf("if %s is set, it must be true", sourceExternalFieldName))
					}
					return
				},
			},
			sourcePodNameFieldName: {
				Description:  "The source will match using pod names",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{sourceIpRangeFieldName, sourceExternalFieldName, sourcePodLabelFieldName, sourcePodAnyFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						connectionRuleNamesFieldName: {
							Type:     schema.TypeList,
							MinItems: 1,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						connectionRuleVulnerabilitySeverityFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL",
							}, true),
						},
					},
				},
			},
			sourcePodLabelFieldName: {
				Description:  "The source will match using pod labels",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{sourceIpRangeFieldName, sourceExternalFieldName, sourcePodNameFieldName, sourcePodAnyFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						connectionRuleNamesLabelsFieldName: {
							Required: true,
							Type:     schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						connectionRuleVulnerabilitySeverityFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL",
							}, true),
						},
					},
				},
			},
			sourcePodAnyFieldName: {
				Description:  "The source will match on any pod (with given vulnerability severity (or higher) if configured)",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{sourceIpRangeFieldName, sourceExternalFieldName, sourcePodNameFieldName, sourcePodLabelFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						connectionRuleVulnerabilitySeverityFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL",
							}, true),
						},
					},
				},
			},
			destinationAddressIpRangeFieldName: {
				Description:  "The destination match will match using ip ranges",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{destinationAddressDomainFieldName, destinationExternalFieldName, destinationPodNameFieldName, destinationPodLabelFieldName, destinationPodAnyFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ipsFieldName: {
							Type:     schema.TypeList,
							MinItems: 1,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			destinationAddressDomainFieldName: {
				Description:  "The destination match will match using domain names",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{destinationAddressIpRangeFieldName, destinationExternalFieldName, destinationPodNameFieldName, destinationPodLabelFieldName, destinationPodAnyFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						domainsFieldName: {
							Type:     schema.TypeList,
							MinItems: 1,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			destinationExternalFieldName: {
				Description:  "The destination will match on external connections",
				Optional:     true,
				ExactlyOneOf: []string{destinationAddressIpRangeFieldName, destinationAddressDomainFieldName, destinationPodNameFieldName, destinationPodLabelFieldName, destinationPodAnyFieldName},
				Type:         schema.TypeBool,
				Default:      false,
				ValidateFunc: func(value interface{}, key string) (warns []string, errs []error) {
					isExternal := value.(bool)
					if isExternal == false {
						errs = append(errs, fmt.Errorf("if %s is set, it must be true", destinationExternalFieldName))
					}
					return
				},
			},
			destinationPodNameFieldName: {
				Description:  "The destination will match using pod names",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{destinationAddressIpRangeFieldName, destinationAddressDomainFieldName, destinationExternalFieldName, destinationPodLabelFieldName, destinationPodAnyFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						connectionRuleNamesFieldName: {
							Type:     schema.TypeList,
							MinItems: 1,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						connectionRuleVulnerabilitySeverityFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL",
							}, true),
						},
					},
				},
			},
			destinationPodLabelFieldName: {
				Description:  "The destination will match using pod labels",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{destinationAddressIpRangeFieldName, destinationAddressDomainFieldName, destinationExternalFieldName, destinationPodNameFieldName, destinationPodAnyFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						connectionRuleNamesLabelsFieldName: {
							Required: true,
							Type:     schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						connectionRuleVulnerabilitySeverityFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL",
							}, true),
						},
					},
				},
			},
			destinationPodAnyFieldName: {
				Description:  "The destination will match on any pod (with given vulnerability severity (or higher) if configured)",
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{destinationAddressIpRangeFieldName, destinationAddressDomainFieldName, destinationExternalFieldName, destinationPodNameFieldName, destinationPodLabelFieldName},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						connectionRuleVulnerabilitySeverityFieldName: {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"UNKNOWN", "LOW", "MEDIUM", "HIGH", "CRITICAL",
							}, true),
						},
					},
				},
			},
		},
	}
}

func resourceConnectionRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] creating cd connection policy")

	err := validateConnectionRuleConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	ruleConfig, err := getConnectionRuleFromConfig(d)

	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := serviceApi.CreateConnectionRule(ctx, httpClientWrapper.HttpClient, ruleConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := rule.Payload.ID

	d.SetId(string(ruleId))

	return nil
}

func resourceConnectionRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] reading cd connection policy")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	ruleId := d.Id()

	currentRuleInSecureCN, err := serviceApi.GetCdConnectionsRule(ctx, httpClientWrapper.HttpClient, strfmt.UUID(ruleId))
	if err != nil {
		return diag.FromErr(err)
	}

	if currentRuleInSecureCN.Payload.ID == "" || currentRuleInSecureCN.Payload.Status == "DELETED" {
		// Tell terraform the rule doesn't exist
		d.SetId("")
	} else {
		updateConnectionRuleMutableFields(d, currentRuleInSecureCN.Payload)
	}

	return nil
}

func resourceConnectionRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] updating cd connection policy")

	err := validateConnectionRuleConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	rule, err := getConnectionRuleFromConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	rule.ID = strfmt.UUID(d.Id())

	//err = updateNonMutableFields(err, serviceApi, httpClientWrapper, rule)
	//if err != nil {
	//	return err
	//}

	updatedRule, err := serviceApi.UpdateCdConnectionsRule(ctx, httpClientWrapper.HttpClient, rule, rule.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(updatedRule.Payload.ID))

	return resourceConnectionRuleRead(ctx, d, m)
}

//func updateNonMutableFields(err error, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper, rule *model.CdConnectionRule) error {
//	currentRuleInSecureCN, err := serviceApi.GetCdConnectionsRule(httpClientWrapper.HttpClient, rule.ID)
//	if err != nil {
//		return err
//	}
//	rule.Name = currentRuleInSecureCN.Payload.Name
//	return nil
//}

func resourceConnectionRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] deleting cd connection policy")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	ruleId := strfmt.UUID(d.Id())

	err := serviceApi.DeleteCdConnectionsRule(ctx, httpClientWrapper.HttpClient, ruleId)
	if err != nil {
		return diag.FromErr(err)
	}

	// Tell terraform the rule doesn't exist
	d.SetId("")

	return nil
}

func validateConnectionRuleConfig(d *schema.ResourceData) error {
	log.Printf("[DEBUG] validating config")

	ips := utils2.ReadNestedListStringFromTF(d, sourceIpRangeFieldName, ipsFieldName, 0)
	for _, ip := range ips {
		_, _, err := net.ParseCIDR(ip)
		if err != nil {
			return errors.New(fmt.Sprintf("invalid configuration. all ips in %s must be a valid CIDR. example: 1.1.1.1/24(IPv4), 002::1234:abcd:ffff:c0a8:101/64(IPv6). ip = %s", sourceIpRangeFieldName, ip))
		}
	}

	ips = utils2.ReadNestedListStringFromTF(d, destinationAddressIpRangeFieldName, ipsFieldName, 0)
	for _, ip := range ips {
		_, _, err := net.ParseCIDR(ip)
		if err != nil {
			return errors.New(fmt.Sprintf("invalid configuration. all ips in %s must be a valid CIDR. example: 1.1.1.1/24(IPv4), 002::1234:abcd:ffff:c0a8:101/64(IPv6). ip = %s", destinationAddressIpRangeFieldName, ip))
		}
	}

	domains := utils2.ReadNestedListStringFromTF(d, destinationAddressDomainFieldName, domainsFieldName, 0)
	for _, domain := range domains {
		match, _ := regexp.MatchString("^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$", domain)
		if match != true {
			return errors.New(fmt.Sprintf("invalid configuration. all domains in %s must be a valid domain name. example: www.domain.com. domain = %s", destinationAddressDomainFieldName, domain))
		}
	}

	return nil
}

func getConnectionRuleFromConfig(d *schema.ResourceData) (*model2.CdConnectionRule, error) {
	log.Print("[DEBUG] getting connection rule from config")

	ruleName := d.Get(connectionRuleNameFieldName).(string)
	//groupName := d.Get(groupNameFieldName).(string)
	status := d.Get(connectionRuleStatusNameFieldName).(string)
	action := d.Get(connectionRuleActionNameFieldName).(string)

	rule := &model2.CdConnectionRule{
		Action:    getConnectionRuleActionFromString(action),
		ID:        "",
		Name:      ruleName,
		GroupName: connectionRuleGroupName,
		Status:    status,
	}

	source, err := getSource(d)
	if err != nil {
		return nil, err
	}

	destination, err := getDestination(d)
	if err != nil {
		return nil, err
	}

	rule.SetSource(source.(model2.ConnectionRulePart))
	rule.SetDestination(destination.(model2.ConnectionRulePart))

	return rule, err
}

func getConnectionRuleActionFromString(action string) model2.ConnectionRuleAction {
	/*
		for now we support only ALLOW
	*/

	//if action == "ALLOW" {
	return model2.ConnectionRuleActionALLOW
	//}
}

func getSource(d *schema.ResourceData) (interface{}, error) {
	sourceIps := utils2.ReadNestedListStringFromTF(d, sourceIpRangeFieldName, ipsFieldName, 0)
	sourceExternal := d.Get(sourceExternalFieldName).(bool)
	sourcePodNames := utils2.ReadNestedListStringFromTF(d, sourcePodNameFieldName, connectionRuleNamesFieldName, 0)
	sourcePodLabels := utils2.ReadNestedMapStringFromTF(d, sourcePodLabelFieldName, connectionRuleNamesLabelsFieldName, 0)
	sourcePodAny := isPodAny(d, sourcePodAnyFieldName)

	if len(sourceIps) != 0 {
		source := &model2.IPRangeConnectionRulePart{
			Networks: sourceIps,
		}
		return source, nil
	} else if sourceExternal {
		source := &model2.ExternalConnectionRulePart{}

		return source, nil
	} else if len(sourcePodNames) != 0 {
		vulString := utils2.ReadNestedStringFromTF(d, sourcePodNameFieldName, connectionRuleVulnerabilitySeverityFieldName, 0)
		source := &model2.PodNameConnectionRulePart{
			Environments:               nil,
			Names:                      sourcePodNames,
			VulnerabilitySeverityLevel: strings.ToUpper(vulString),
		}

		return source, nil

	} else if len(sourcePodLabels) != 0 {
		vulString := utils2.ReadNestedStringFromTF(d, sourcePodLabelFieldName, connectionRuleVulnerabilitySeverityFieldName, 0)
		labels := utils2.GetLabelsFromMap(sourcePodLabels)

		source := &model2.PodLablesConnectionRulePart{
			Environments:               nil,
			Labels:                     labels,
			VulnerabilitySeverityLevel: strings.ToUpper(vulString),
		}

		return source, nil

	} else if sourcePodAny {
		vulString := utils2.ReadNestedStringFromTF(d, sourcePodAnyFieldName, connectionRuleVulnerabilitySeverityFieldName, 0)
		source := &model2.PodAnyConnectionRulePart{
			Environments:               nil,
			VulnerabilitySeverityLevel: strings.ToUpper(vulString),
		}

		return source, nil
	} else {
		return nil, errors.New(fmt.Sprintf("failed to get source. all source fields were empty"))
	}
}

func getDestination(d *schema.ResourceData) (interface{}, error) {
	destinationIps := utils2.ReadNestedListStringFromTF(d, destinationAddressIpRangeFieldName, ipsFieldName, 0)
	destinationDomains := utils2.ReadNestedListStringFromTF(d, destinationAddressDomainFieldName, domainsFieldName, 0)
	destinationExternal := d.Get(destinationExternalFieldName).(bool)
	destinationPodNames := utils2.ReadNestedListStringFromTF(d, destinationPodNameFieldName, connectionRuleNamesFieldName, 0)
	destinationPodLabels := utils2.ReadNestedMapStringFromTF(d, destinationPodLabelFieldName, connectionRuleNamesLabelsFieldName, 0)
	destinationPodAny := isPodAny(d, destinationPodAnyFieldName)

	if len(destinationIps) != 0 {
		destination := &model2.IPRangeConnectionRulePart{
			Networks: destinationIps,
		}
		return destination, nil
	} else if len(destinationDomains) != 0 {
		destination := &model2.FqdnConnectionRulePart{
			FqdnAddresses: destinationDomains,
		}
		return destination, nil
	} else if destinationExternal {
		destination := &model2.ExternalConnectionRulePart{}

		return destination, nil
	} else if len(destinationPodNames) != 0 {
		vulString := utils2.ReadNestedStringFromTF(d, destinationPodNameFieldName, connectionRuleVulnerabilitySeverityFieldName, 0)
		destination := &model2.PodNameConnectionRulePart{
			Environments:               nil,
			Names:                      destinationPodNames,
			VulnerabilitySeverityLevel: strings.ToUpper(vulString),
		}

		return destination, nil

	} else if len(destinationPodLabels) != 0 {
		vulString := utils2.ReadNestedStringFromTF(d, destinationPodLabelFieldName, connectionRuleVulnerabilitySeverityFieldName, 0)
		labels := utils2.GetLabelsFromMap(destinationPodLabels)
		destination := &model2.PodLablesConnectionRulePart{
			Environments:               nil,
			Labels:                     labels,
			VulnerabilitySeverityLevel: strings.ToUpper(vulString),
		}

		return destination, nil

	} else if destinationPodAny {
		vulString := utils2.ReadNestedStringFromTF(d, destinationPodAnyFieldName, connectionRuleVulnerabilitySeverityFieldName, 0)
		destination := &model2.PodAnyConnectionRulePart{
			Environments:               nil,
			VulnerabilitySeverityLevel: strings.ToUpper(vulString),
		}

		return destination, nil
	} else {
		return nil, errors.New(fmt.Sprintf("failed to get destination. all destination fields were empty"))
	}
}

func isPodAny(d *schema.ResourceData, mainField string) bool {
	ipsData, exists := d.GetOk(mainField)

	if exists == true {
		interfaces := ipsData.([]interface{})
		if interfaces != nil {
			return true
		}
	}

	return false
}

func updateConnectionRuleMutableFields(d *schema.ResourceData, currentRuleInSecureCN *model2.CdConnectionRule) {
	log.Print("[DEBUG] updating cd connection rule mutable fields")

	_ = d.Set(connectionRuleNameFieldName, currentRuleInSecureCN.Name)
	_ = d.Set(connectionRuleActionNameFieldName, currentRuleInSecureCN.Action)
	_ = d.Set(connectionRuleStatusNameFieldName, currentRuleInSecureCN.Status)

	mutateSource(d, currentRuleInSecureCN)
	mutateDestination(d, currentRuleInSecureCN)
}

func mutateDestination(d *schema.ResourceData, currentRule *model2.CdConnectionRule) {
	destination := currentRule.Destination()
	destinationPartType := destination.ConnectionRulePartType()
	if destinationPartType == "PodNameConnectionRulePart" {
		mainField := destinationPodNameFieldName
		updateByPodNames(d, destination, mainField)
		_ = d.Set(destinationAddressIpRangeFieldName, nil)
		_ = d.Set(destinationAddressDomainFieldName, nil)
		_ = d.Set(destinationExternalFieldName, nil)
		_ = d.Set(destinationPodLabelFieldName, nil)
		_ = d.Set(destinationPodAnyFieldName, nil)
	} else if destinationPartType == "PodLablesConnectionRulePart" {
		mainField := destinationPodLabelFieldName
		updateByLabels(d, destination, mainField)
		_ = d.Set(destinationAddressIpRangeFieldName, nil)
		_ = d.Set(destinationAddressDomainFieldName, nil)
		_ = d.Set(destinationExternalFieldName, nil)
		_ = d.Set(destinationPodNameFieldName, nil)
		_ = d.Set(destinationPodAnyFieldName, nil)
	} else if destinationPartType == "PodAnyConnectionRulePart" {
		_ = d.Set(destinationPodAnyFieldName, destination)
		_ = d.Set(destinationAddressIpRangeFieldName, nil)
		_ = d.Set(destinationAddressDomainFieldName, nil)
		_ = d.Set(destinationExternalFieldName, nil)
		_ = d.Set(destinationPodNameFieldName, nil)
		_ = d.Set(destinationPodLabelFieldName, nil)
	} else if destinationPartType == "IpRangeConnectionRulePart" {
		mainField := destinationAddressIpRangeFieldName
		updateByIps(d, destination, mainField)
		_ = d.Set(destinationAddressDomainFieldName, nil)
		_ = d.Set(destinationExternalFieldName, nil)
		_ = d.Set(destinationPodNameFieldName, nil)
		_ = d.Set(destinationPodLabelFieldName, nil)
		_ = d.Set(destinationPodAnyFieldName, nil)
	} else if destinationPartType == "FqdnConnectionRulePart" {
		mainField := destinationAddressDomainFieldName
		updateByDomains(d, destination, mainField)
		_ = d.Set(destinationAddressIpRangeFieldName, nil)
		_ = d.Set(destinationExternalFieldName, nil)
		_ = d.Set(destinationPodNameFieldName, nil)
		_ = d.Set(destinationPodLabelFieldName, nil)
		_ = d.Set(destinationPodAnyFieldName, nil)
	} else if destinationPartType == "ExternalConnectionRulePart" {
		_ = d.Set(destinationExternalFieldName, destination)
		_ = d.Set(destinationAddressIpRangeFieldName, nil)
		_ = d.Set(destinationAddressDomainFieldName, nil)
		_ = d.Set(destinationPodNameFieldName, nil)
		_ = d.Set(destinationPodLabelFieldName, nil)
		_ = d.Set(destinationPodAnyFieldName, nil)
	}
}

func mutateSource(d *schema.ResourceData, currentRuleInSecureCN *model2.CdConnectionRule) {
	source := currentRuleInSecureCN.Source()
	currentSourcePartTypeInSecureCN := source.ConnectionRulePartType()
	if currentSourcePartTypeInSecureCN == "PodNameConnectionRulePart" {
		mainField := sourcePodNameFieldName
		updateByPodNames(d, source, mainField)
		_ = d.Set(sourceIpRangeFieldName, nil)
		_ = d.Set(sourceExternalFieldName, nil)
		_ = d.Set(sourcePodLabelFieldName, nil)
		_ = d.Set(sourcePodAnyFieldName, nil)

	} else if currentSourcePartTypeInSecureCN == "PodLablesConnectionRulePart" {
		mainField := sourcePodLabelFieldName
		updateByLabels(d, source, mainField)
		_ = d.Set(sourcePodNameFieldName, nil)
		_ = d.Set(sourceIpRangeFieldName, nil)
		_ = d.Set(sourceExternalFieldName, nil)
		_ = d.Set(sourcePodAnyFieldName, nil)
	} else if currentSourcePartTypeInSecureCN == "PodAnyConnectionRulePart" {
		currentSourceInSecureCN := source.(*model2.PodAnyConnectionRulePart)
		_ = d.Set(sourcePodAnyFieldName, []map[string]string{{
			connectionRuleVulnerabilitySeverityFieldName: currentSourceInSecureCN.VulnerabilitySeverityLevel,
		}})
		_ = d.Set(sourcePodLabelFieldName, nil)
		_ = d.Set(sourcePodNameFieldName, nil)
		_ = d.Set(sourceIpRangeFieldName, nil)
		_ = d.Set(sourceExternalFieldName, nil)
	} else if currentSourcePartTypeInSecureCN == "IpRangeConnectionRulePart" {
		mainField := sourceIpRangeFieldName
		updateByIps(d, source, mainField)
		_ = d.Set(sourcePodLabelFieldName, nil)
		_ = d.Set(sourcePodNameFieldName, nil)
		_ = d.Set(sourceExternalFieldName, nil)
		_ = d.Set(sourcePodAnyFieldName, nil)
	} else if currentSourcePartTypeInSecureCN == "ExternalConnectionRulePart" {
		currentSourceInSecureCN := source.(*model2.ExternalConnectionRulePart)
		_ = d.Set(sourceExternalFieldName, currentSourceInSecureCN)
		_ = d.Set(sourceIpRangeFieldName, nil)
		_ = d.Set(sourcePodLabelFieldName, nil)
		_ = d.Set(sourcePodNameFieldName, nil)
		_ = d.Set(sourcePodAnyFieldName, nil)
	}
}

func updateByPodNames(d *schema.ResourceData, part model2.ConnectionRulePart, mainField string) {
	currentPartInSecureCN := part.(*model2.PodNameConnectionRulePart)
	currentPartInTerraform := d.Get(mainField)
	if currentPartInTerraform == nil || len(currentPartInTerraform.([]interface{})) == 0 {
		_ = d.Set(mainField, currentPartInSecureCN)
	} else {
		terraformPart := currentPartInTerraform.([]interface{})[0]
		for key, value := range terraformPart.(map[string]interface{}) {
			if key == connectionRuleNamesFieldName {
				updateStringSliceSubField(d, mainField, connectionRuleNamesFieldName, terraformPart, currentPartInSecureCN.Names)
			}
			if key == connectionRuleVulnerabilitySeverityFieldName {
				updateStringSubField(d, mainField, connectionRuleVulnerabilitySeverityFieldName, terraformPart, value, currentPartInSecureCN.VulnerabilitySeverityLevel)
			}

		}
	}
}

func updateByLabels(d *schema.ResourceData, part model2.ConnectionRulePart, mainField string) {
	currentPartInSecureCN := part.(*model2.PodLablesConnectionRulePart)
	currentPartInTerraform := d.Get(mainField)
	if currentPartInTerraform == nil || len(currentPartInTerraform.([]interface{})) == 0 {
		_ = d.Set(mainField, currentPartInSecureCN)
	} else {
		terraformPart := currentPartInTerraform.([]interface{})[0]

		for key, value := range terraformPart.(map[string]interface{}) {
			if key == connectionRuleNamesLabelsFieldName {
				updateLabelMapSubField(d, mainField, connectionRuleNamesLabelsFieldName, terraformPart, currentPartInSecureCN.Labels)
			}
			if key == connectionRuleVulnerabilitySeverityFieldName {
				updateStringSubField(d, mainField, connectionRuleVulnerabilitySeverityFieldName, terraformPart, value, currentPartInSecureCN.VulnerabilitySeverityLevel)
			}

		}
	}
}

func updateByIps(d *schema.ResourceData, part model2.ConnectionRulePart, mainField string) {
	currentPartInSecureCN := part.(*model2.IPRangeConnectionRulePart)
	currentPartInTerraform := d.Get(mainField)
	if currentPartInTerraform == nil || len(currentPartInTerraform.([]interface{})) == 0 {
		_ = d.Set(mainField, currentPartInSecureCN)
	} else {
		terraformPart := currentPartInTerraform.([]interface{})[0]
		for key := range terraformPart.(map[string]interface{}) {
			if key == ipsFieldName {
				updateStringSliceSubField(d, mainField, ipsFieldName, terraformPart, currentPartInSecureCN.Networks)
			}
		}
	}
}

func updateByDomains(d *schema.ResourceData, part model2.ConnectionRulePart, mainField string) {
	currentPartInSecureCN := part.(*model2.FqdnConnectionRulePart)
	currentPartInTerraform := d.Get(mainField)
	if currentPartInTerraform == nil || len(currentPartInTerraform.([]interface{})) == 0 {
		_ = d.Set(mainField, currentPartInSecureCN)
	} else {
		terraformPart := currentPartInTerraform.([]interface{})[0]
		for key := range terraformPart.(map[string]interface{}) {
			if key == domainsFieldName {
				updateStringSliceSubField(d, mainField, domainsFieldName, terraformPart, currentPartInSecureCN.FqdnAddresses)
			}
		}
	}
}

func updateStringSubField(d *schema.ResourceData, mainField string, subField string, terraformPart interface{}, valueInTerraform interface{}, valueInSecureCN string) {
	if valueInTerraform != valueInSecureCN {
		fieldInTerraform := terraformPart.(map[string]interface{})
		fieldInTerraform[subField] = valueInSecureCN
		newValues := make([]interface{}, 0, len(fieldInTerraform))
		newValues = append(newValues, fieldInTerraform)
		_ = d.Set(mainField, newValues)
	}
}

func updateStringSliceSubField(d *schema.ResourceData, mainField string, subField string, terraformPart interface{}, secureCNPart []string) {
	valuesInTerraform := getDataInTerraformAsStringSlice(terraformPart, subField)
	if !utils2.IsStringSlicesIdentical(valuesInTerraform, secureCNPart) {
		fieldInTerraform := terraformPart.(map[string]interface{})
		fieldInTerraform[subField] = secureCNPart
		newValues := make([]interface{}, 0, len(fieldInTerraform))
		newValues = append(newValues, fieldInTerraform)
		_ = d.Set(mainField, newValues)
	}
}

func updateLabelMapSubField(d *schema.ResourceData, mainField string, subField string, terraformPart interface{}, secureCNPart []*model2.Label) {

	labelsInTerraform := getDataInTerraformAsLabelsSlice(terraformPart, subField)
	if !reflect.DeepEqual(labelsInTerraform, secureCNPart) {
		_ = d.Set(mainField, nil)
	}
}

func getDataInTerraformAsStringSlice(inter interface{}, subfield string) []string {
	imap := inter.(map[string]interface{})
	sub := imap[subfield].([]interface{})
	values := make([]string, 0, len(sub))
	for _, s := range sub {
		values = append(values, s.(string))
	}
	return values
}

func getDataInTerraformAsLabelsSlice(inter interface{}, subfield string) []*model2.Label {
	imap := inter.(map[string]interface{})
	sub := imap[subfield].(map[string]interface{})
	values := make([]*model2.Label, 0, len(sub))
	for k, v := range sub {
		label := &model2.Label{
			Key:   k,
			Value: v.(string),
		}

		values = append(values, label)
	}
	return values
}
