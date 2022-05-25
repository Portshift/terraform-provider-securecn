package securecn

import (
	"context"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"terraform-provider-securecn/internal/client"
	model2 "terraform-provider-securecn/internal/escher_api/model"
	utils2 "terraform-provider-securecn/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const trustedSignerNameFieldName = "name"
const trustedSignerKeysFieldName = "keys"
const trustedSignerClustersFieldName = "clusters"

func ResourceTrustedSigner() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceTrustedSignerCreate,
		ReadContext:   resourceTrustedSignerRead,
		UpdateContext: resourceTrustedSignerUpdate,
		DeleteContext: resourceTrustedSignerDelete,
		Description:   "A SecureCN TrustedSigner",
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			trustedSignerNameFieldName: {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			trustedSignerKeysFieldName: {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			trustedSignerClustersFieldName: {
				Optional: true,
				Type:     schema.TypeList,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},
		},
	}
}

func resourceTrustedSignerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] creating trustedSigner")

	err := validateTrustedSignerConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	trustedSignerFromConfig, err := getTrustedSignerFromConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	trustedSigner, err := serviceApi.CreateTrustedSigner(ctx, httpClientWrapper.HttpClient, trustedSignerFromConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(trustedSigner.Payload.ID))

	return nil
}

func resourceTrustedSignerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] reading trustedSigner")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	trustedSignerId := d.Id()

	trustedSigner, err := serviceApi.GetTrustedSignerById(ctx, httpClientWrapper.HttpClient, strfmt.UUID(trustedSignerId))

	if err != nil {
		return diag.FromErr(err)
	}

	if trustedSigner == nil {
		// Tell terraform the trustedSigner doesn't exist
		d.SetId("")
	} else {
		return diag.FromErr(updateTrustedSignerMutableFields(d, trustedSigner.GetPayload()))
	}

	return nil
}

func resourceTrustedSignerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] updating trustedSigner")

	err := validateTrustedSignerConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	trustedSigner, err := getTrustedSignerFromConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = serviceApi.UpdateTrustedSigner(ctx, httpClientWrapper.HttpClient, trustedSigner, trustedSigner.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTrustedSignerRead(ctx, d, m)
}

func resourceTrustedSignerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] deleting trustedSigner")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	_, err := serviceApi.DeleteTrustedSigner(ctx, httpClientWrapper.HttpClient, strfmt.UUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	// Tell terraform the trustedSigner doesn't exist
	d.SetId("")

	return nil
}

func validateTrustedSignerConfig(d *schema.ResourceData) error {
	log.Printf("[DEBUG] validating trustedSigner config")

	return nil
}

func getTrustedSignerFromConfig(d *schema.ResourceData) (*model2.TrustedSigner, error) {
	log.Print("[DEBUG] getting trustedSigner from config")

	name := d.Get(trustedSignerNameFieldName).(string)
	trustedSignerKeys := getTrustedSignerKeysFromConfig(d)
	trustedSignerClusters := getTrustedSignerClustersKeysFromConfig(d)

	trustedSigner := &model2.TrustedSigner{
		Name:                  &name,
		Keys:                  trustedSignerKeys,
		TrustedSignerClusters: trustedSignerClusters,
	}

	return trustedSigner, nil
}

func getTrustedSignerClustersKeysFromConfig(d *schema.ResourceData) []*model2.TrustedSignerCluster {
	clustersInConfig := d.Get(trustedSignerClustersFieldName).([]interface{})
	clusters := make([]*model2.TrustedSignerCluster, 0, len(clustersInConfig))
	for _, clusterId := range clustersInConfig {
		cluster := &model2.TrustedSignerCluster{
			ID: strfmt.UUID(clusterId.(string)),
		}
		clusters = append(clusters, cluster)
	}

	return clusters
}

func getTrustedSignerKeysFromConfig(d *schema.ResourceData) []*model2.TrustedSignerKey {
	keysMap := d.Get(trustedSignerKeysFieldName).(map[string]interface{})
	keys := make([]*model2.TrustedSignerKey, 0, len(keysMap))
	for k, v := range keysMap {
		value := v.(string)
		key := &model2.TrustedSignerKey{
			Key:  &k,
			Name: &value,
		}
		keys = append(keys, key)
	}

	return filterEmptyTrustedSignerKeys(keys)
}

func filterEmptyTrustedSignerKeys(labels []*model2.TrustedSignerKey) []*model2.TrustedSignerKey {
	var ans []*model2.TrustedSignerKey
	for _, label := range labels {
		if label != nil {
			ans = append(ans, label)
		}
	}
	return ans
}

func updateTrustedSignerMutableFields(d *schema.ResourceData, currentTrustedSignerInSecureCn *model2.TrustedSigner) error {
	err := d.Set(trustedSignerNameFieldName, currentTrustedSignerInSecureCn.Name)
	if err != nil {
		return err
	}

	err = updateTrustedSignerMutableFieldsKeys(d, currentTrustedSignerInSecureCn)
	if err != nil {
		return err
	}

	err = updateTrustedSignerMutableFieldsClusters(d, currentTrustedSignerInSecureCn)
	if err != nil {
		return err
	}
	return nil
}

func updateTrustedSignerMutableFieldsKeys(d *schema.ResourceData, currentSignerInSecureCN *model2.TrustedSigner) error {
	keysInSecureCN := currentSignerInSecureCN.Keys
	signerKeys := make(map[string]string, len(keysInSecureCN))
	for _, singleKeyInSecureCN := range keysInSecureCN {
		nameInSecureCn := singleKeyInSecureCN.Name
		keyInSecureCn := singleKeyInSecureCN.Key
		signerKeys[*nameInSecureCn] = *keyInSecureCn
	}

	err := d.Set(trustedSignerKeysFieldName, signerKeys)
	if err != nil {
		return err
	}
	return nil
}

func updateTrustedSignerMutableFieldsClusters(d *schema.ResourceData, currentRuleInSecureCN *model2.TrustedSigner) error {
	clustersInSecureCN := currentRuleInSecureCN.TrustedSignerClusters
	signerClusters := make([]string, 0, len(clustersInSecureCN))
	for _, singleKeyInSecureCN := range clustersInSecureCN {
		signerClusters = append(signerClusters, singleKeyInSecureCN.ID.String())
	}

	err := d.Set(trustedSignerClustersFieldName, signerClusters)
	if err != nil {
		return err
	}
	return nil
}
