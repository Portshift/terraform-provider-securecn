package escherClient

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"io"
	"log"
	"net/http"
	auth2 "terraform-provider-securecn/internal/escher_api/auth"
	model2 "terraform-provider-securecn/internal/escher_api/model"
)

const (
	// DefaultHost is the default Host
	// found in Meta (info) section of spec file
	DefaultHost string = "localhost"
	// DefaultBasePath is the default BasePath
	// found in Meta (info) section of spec file
	DefaultBasePath string = "/api"
)

// DefaultSchemes are the default schemes found in Meta (info) section of spec file
var DefaultSchemes = []string{"https"}

// DefaultTransportConfig creates a TransportConfig with the
// default settings taken from the meta section of the spec file.
func DefaultTransportConfig() *TransportConfig {
	return &TransportConfig{
		Host:     DefaultHost,
		BasePath: DefaultBasePath,
		Schemes:  DefaultSchemes,
	}
}

// TransportConfig contains the transport related info,
// found in the meta section of the spec file.
type TransportConfig struct {
	Host     string
	BasePath string
	Schemes  []string
}

// WithHost overrides the default host,
// provided by the meta section of the spec file.
func (cfg *TransportConfig) WithHost(host string) *TransportConfig {
	cfg.Host = host
	return cfg
}

// WithBasePath overrides the default basePath,
// provided by the meta section of the spec file.
func (cfg *TransportConfig) WithBasePath(basePath string) *TransportConfig {
	cfg.BasePath = basePath
	return cfg
}

// WithSchemes overrides the default schemes,
// provided by the meta section of the spec file.
func (cfg *TransportConfig) WithSchemes(schemes []string) *TransportConfig {
	cfg.Schemes = schemes
	return cfg
}

type MgmtServiceApiCtx struct {
	auth      runtime.ClientAuthInfoWriter
	accessKey string
	runtime   *auth2.Runtime
}

func CreateServiceApi(url, accessKey, secretKey string, client *http.Client) (*MgmtServiceApiCtx, error) {
	secretKeyBytes, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret key: %v", err)
	}

	apiCtx := createMgmtServiceApiCtx(url, client)

	apiCtx.setServiceKeys(accessKey, secretKeyBytes)

	return apiCtx, nil
}

func createMgmtServiceApiCtx(mgmtHost string, httpClient *http.Client) *MgmtServiceApiCtx {
	cfg := DefaultTransportConfig().WithHost(mgmtHost)
	transport := auth2.NewWithClient(cfg.Host, cfg.BasePath, cfg.Schemes, httpClient)

	return &MgmtServiceApiCtx{
		runtime: transport,
	}
}

func (serviceMgmtApi *MgmtServiceApiCtx) DownloadKubernetesSecureCNBundle(ctx context.Context, client *http.Client, writer io.Writer, clusterUUID strfmt.UUID) error {
	params := &model2.DownloadBundleParams{
		KubernetesClusterID: clusterUUID,
		Timeout:             0,
		Context:             ctx,
	}

	_, err := serviceMgmtApi.downloadBundle(params, writer, client)

	if err != nil {
		return fmt.Errorf("failed to get SecureCN bundle: %v", err)
	}

	return nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) CreateKubernetesCluster(ctx context.Context, client *http.Client, cluster *model2.KubernetesCluster) (*model2.PostKubernetesClustersCreated, error) {
	log.Print("[DEBUG] creating cluster")

	params := &model2.PostKubernetesClustersParams{
		Cluster:    cluster,
		Timeout:    0,
		Context:    ctx,
		HTTPClient: client,
	}

	newCluster, err := serviceMgmtApi.postKubernetesClusters(params, client)

	if err != nil {
		return nil, err
	}

	return newCluster, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) CreateConnectionRule(ctx context.Context, client *http.Client, rule *model2.CdConnectionRule) (*model2.PostCdConnectionsRuleCreated, error) {
	log.Print("[DEBUG] creating cd connections rule")

	params := &model2.PostCdConnectionsRuleParams{
		Body:       rule,
		Context:    ctx,
		HTTPClient: client,
	}

	newRule, err := serviceMgmtApi.postCdConnectionsRule(params)

	if err != nil {
		return nil, err
	}

	return newRule, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) CreateEnvironment(ctx context.Context, client *http.Client, env *model2.Environment) (*model2.PostEnvironmentsCreated, error) {
	log.Print("[DEBUG] creating environment")

	params := &model2.PostEnvironmentsParams{
		Body:       env,
		Context:    ctx,
		HTTPClient: client,
	}

	newEnv, err := serviceMgmtApi.postEnvironments(params, client)

	if err != nil {
		log.Printf("[DEBUG] failed creating environment %v", err)
		return nil, err
	}

	return newEnv, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) CreateDeploymentRule(ctx context.Context, client *http.Client, rule *model2.CdAppRule) (*model2.PostCdDeploymentRuleCreated, error) {
	log.Print("[DEBUG] deployment rule")

	params := &model2.PostCdDeploymentRuleParams{
		Body:       rule,
		Context:    ctx,
		HTTPClient: client,
	}

	newRule, err := serviceMgmtApi.postCdDeploymentRule(params)

	if err != nil {
		log.Printf("[DEBUG] failed creating deployment rule %v", err)
		return nil, err
	}

	return newRule, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) GetKubernetesClusterById(ctx context.Context, client *http.Client, clusterId strfmt.UUID) (*model2.GetKubernetesClustersKubernetesClusterIDOK, error) {
	log.Print("[DEBUG] getting cluster")

	params := &model2.GetKubernetesClustersKubernetesClusterIDParams{
		KubernetesClusterID: clusterId,
		Context:             ctx,
		HTTPClient:          client,
	}
	cluster, err := serviceMgmtApi.getKubernetesClustersKubernetesClusterID(params)

	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes cluster: %v", err)
	}

	return cluster, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) GetKubernetesClusterIdByName(ctx context.Context, client *http.Client, kubernetesClusterName string) (*model2.GetKubernetesClustersKubernetesClusterNameOK, error) {
	log.Print("[DEBUG] getting cluster")

	params := &model2.GetKubernetesClustersKubernetesClusterNameParams{
		KubernetesClusterName: kubernetesClusterName,
		Context:               ctx,
		HTTPClient:            client,
	}
	clusterId, err := serviceMgmtApi.getKubernetesClustersKubernetesClusterName(params)

	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes cluster id: %v", err)
	}

	return clusterId, nil
}

/*
GetKubernetesClustersKubernetesClusterIDNamespaces lists namespaces on a specific kubernetes cluster
*/
func (serviceMgmtApi *MgmtServiceApiCtx) GetKubernetesClustersKubernetesClusterIDNamespaces(ctx context.Context, clusterId strfmt.UUID) (*model2.GetKubernetesClustersKubernetesClusterIDNamespacesOK, error) {
	registry := new(strfmt.Registry)

	params := model2.GetKubernetesClustersKubernetesClusterIDNamespacesParams{
		KubernetesClusterID: clusterId,
		Context:             ctx,
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "GetKubernetesClustersKubernetesClusterIDNamespaces",
		Method:             "GET",
		PathPattern:        "/kubernetesClusters/{kubernetesClusterId}/namespaces",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             &params,
		Reader:             &model2.GetKubernetesClustersKubernetesClusterIDNamespacesReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.GetKubernetesClustersKubernetesClusterIDNamespacesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*model2.GetKubernetesClustersKubernetesClusterIDNamespacesDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

func (serviceMgmtApi *MgmtServiceApiCtx) GetCdConnectionsRule(ctx context.Context, client *http.Client, ruleId strfmt.UUID) (*model2.GetCdRuleIDConnectionsRuleOK, error) {
	log.Print("[DEBUG] getting cd connection rule")

	params := &model2.GetCdRuleIDConnectionsRuleParams{
		RuleID:     ruleId,
		Context:    ctx,
		HTTPClient: client,
	}
	rule, err := serviceMgmtApi.getCdRuleIDConnectionsRule(params)

	if err != nil {
		return nil, fmt.Errorf("failed to get cd connection rule: %v", err)
	}

	return rule, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) GetEnvironment(ctx context.Context, client *http.Client, envId strfmt.UUID) (*model2.GetEnvironmentsEnvIDOK, error) {
	log.Print("[DEBUG] getting environment")

	params := &model2.GetEnvironmentsEnvIDParams{
		EnvID:      envId,
		Context:    ctx,
		HTTPClient: client,
	}
	env, err := serviceMgmtApi.getEnvironmentsEnvID(params)

	if err != nil {
		return nil, fmt.Errorf("failed to get environment. id: %v, : %v", envId, err)
	}

	return env, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) GetDeploymentRule(ctx context.Context, client *http.Client, ruleId strfmt.UUID) (*model2.GetCdRuleIDDeploymentRuleOK, error) {
	log.Print("[DEBUG] getting deployment rule")

	params := &model2.GetCdRuleIDDeploymentRuleParams{
		RuleID:     ruleId,
		Context:    ctx,
		HTTPClient: client,
	}
	rule, err := serviceMgmtApi.getCdRuleIDDeploymentRule(params)

	if err != nil {
		return nil, fmt.Errorf("failed to get rule. id: %v, : %v", ruleId, err)
	}

	return rule, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) UpdateKubernetesCluster(ctx context.Context, client *http.Client, cluster *model2.KubernetesCluster, clusterId strfmt.UUID) (*model2.PutKubernetesClustersKubernetesClusterIDOK, error) {
	log.Print("[DEBUG] updating cluster")

	params := &model2.PutKubernetesClustersKubernetesClusterIDParamsWriter{
		Body:                cluster,
		KubernetesClusterID: clusterId,
		Context:             ctx,
		HTTPClient:          client,
	}
	cluster.ID = ""
	updatedCluster, err := serviceMgmtApi.putKubernetesClustersKubernetesClusterID(params)

	if err != nil {
		return nil, fmt.Errorf("failed to update kubernetes cluster: %v", err)
	}

	return updatedCluster, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) UpdateCdConnectionsRule(ctx context.Context, client *http.Client, rule *model2.CdConnectionRule, ruleId strfmt.UUID) (*model2.PutCdRuleIDConnectionsRuleOK, error) {
	log.Print("[DEBUG] updating cd connections rule")

	params := &model2.PutCdRuleIDConnectionsRuleParams{
		Body:       rule,
		RuleID:     ruleId,
		Context:    ctx,
		HTTPClient: client,
	}

	updatedRule, err := serviceMgmtApi.putCdRuleIDConnectionsRule(params)

	if err != nil {
		return nil, fmt.Errorf("failed to update cd connections rule: %v", err)
	}

	return updatedRule, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) UpdateEnvironment(ctx context.Context, client *http.Client, env *model2.Environment, envId strfmt.UUID) (*model2.PutEnvironmentsEnvIDOK, error) {
	log.Print("[DEBUG] updating environment")

	params := &model2.PutEnvironmentsEnvIDParams{
		Body:       env,
		EnvID:      envId,
		Context:    ctx,
		HTTPClient: client,
	}

	updatedEnv, err := serviceMgmtApi.putEnvironmentsEnvID(params)

	if err != nil {
		return nil, fmt.Errorf("failed to update environment: %v", err)
	}

	return updatedEnv, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) UpdateDeploymentRule(ctx context.Context, client *http.Client, rule *model2.CdAppRule, ruleId strfmt.UUID) (*model2.PutCdRuleIDDeploymentRuleOK, error) {
	log.Print("[DEBUG] updating deployment rule")

	params := &model2.PutCdRuleIDDeploymentRuleParams{
		Body:       rule,
		RuleID:     ruleId,
		Context:    ctx,
		HTTPClient: client,
	}

	updatedRule, err := serviceMgmtApi.putCdRuleIDDeploymentRule(params)

	if err != nil {
		return nil, fmt.Errorf("failed to update deployment rule: %v", err)
	}

	return updatedRule, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) DeleteKubernetesCluster(ctx context.Context, client *http.Client, clusterId strfmt.UUID) error {
	log.Print("[DEBUG] deleting cluster")

	params := &model2.DeleteKubernetesClustersKubernetesClusterIDParams{
		KubernetesClusterID: clusterId,
		Context:             ctx,
		HTTPClient:          client,
	}
	_, err := serviceMgmtApi.deleteKubernetesClustersKubernetesClusterID(params)

	if err != nil {
		return fmt.Errorf("failed to delete kubernetes cluster: %v", err)
	}

	return nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) DeleteCdConnectionsRule(ctx context.Context, client *http.Client, ruleId strfmt.UUID) error {
	log.Print("[DEBUG] deleting cd connections rule")
	params := &model2.DeleteCdRuleIDConnectionsRuleParams{
		RuleID:     ruleId,
		Context:    ctx,
		HTTPClient: client,
	}
	_, err := serviceMgmtApi.deleteCdRuleIDConnectionsRule(params)

	if err != nil {
		return fmt.Errorf("failed to delete cd connections rule: %v", err)
	}

	return nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) DeleteEnvironment(ctx context.Context, client *http.Client, envId strfmt.UUID) error {
	log.Print("[DEBUG] deleting environment")

	params := &model2.DeleteEnvironmentsEnvIDParams{
		EnvID:      envId,
		Context:    ctx,
		HTTPClient: client,
	}
	_, err := serviceMgmtApi.deleteEnvironmentEnvID(params)

	if err != nil {
		return fmt.Errorf("failed to delete environment: %v", err)
	}

	return nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) DeleteDeploymentRule(ctx context.Context, client *http.Client, ruleId strfmt.UUID) error {
	log.Print("[DEBUG] deleting deployment rule")

	params := &model2.DeleteCdRuleIDDeploymentRuleParams{
		RuleID:     ruleId,
		Context:    ctx,
		HTTPClient: client,
	}
	_, err := serviceMgmtApi.deleteCdRuleIDDeploymentRule(params)

	if err != nil {
		return fmt.Errorf("failed to delete cd deployment rule: %v", err)
	}

	return nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) GetPspIdByName(ctx context.Context, client *http.Client, podSecurityPolicyProfileName string) (*model2.GetCdPodSecurityPolicyProfilesPodSecurityPolicyProfileNameOK, error) {

	params := &model2.GetCdPodSecurityPolicyProfilesPodSecurityPolicyProfileNameParams{
		PodSecurityPolicyProfileName: podSecurityPolicyProfileName,
		Context:                      ctx,
		HTTPClient:                   client,
	}

	pspId, err := serviceMgmtApi.getCdPodSecurityPolicyProfilesPodSecurityPolicyProfileName(params)

	if err != nil {
		return nil, fmt.Errorf("failed to get psp id by name: %v. %v", podSecurityPolicyProfileName, err)
	}

	return pspId, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) downloadBundle(params *model2.DownloadBundleParams, writer io.Writer, c *http.Client) (*model2.GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK, error) {
	registry := new(strfmt.Registry)
	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGz",
		Method:             "GET",
		PathPattern:        "/kubernetesClusters/{kubernetesClusterId}/download_bundle",
		ProducesMediaTypes: []string{"application/gzip"},
		ConsumesMediaTypes: []string{"application/gzip"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             params,
		Reader:             &model2.GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzReader{Formats: *registry, Writer: writer},
		Context:            params.Context,
		Client:             c,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*model2.GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

func (serviceMgmtApi *MgmtServiceApiCtx) postKubernetesClusters(params *model2.PostKubernetesClustersParams, client *http.Client) (*model2.PostKubernetesClustersCreated, error) {
	registry := new(strfmt.Registry)
	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PostKubernetesClusters",
		Method:             "POST",
		PathPattern:        "/kubernetesClusters",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.PostKubernetesClustersReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             client,
	})

	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PostKubernetesClustersCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*model2.PostKubernetesClustersDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

func (serviceMgmtApi *MgmtServiceApiCtx) getKubernetesClustersKubernetesClusterID(params *model2.GetKubernetesClustersKubernetesClusterIDParams) (*model2.GetKubernetesClustersKubernetesClusterIDOK, error) {
	registry := new(strfmt.Registry)
	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "GetKubernetesClustersKubernetesClusterID",
		Method:             "GET",
		PathPattern:        "/kubernetesClusters/{kubernetesClusterId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.GetKubernetesClustersKubernetesClusterIDReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.GetKubernetesClustersKubernetesClusterIDOK)
	if ok {
		return success, nil
	}

	// unexpected success response
	unexpectedSuccess := result.(*model2.GetKubernetesClustersKubernetesClusterIDDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
getKubernetesClustersKubernetesClusterName gets the kubernetes cluster id with the given name
*/
func (serviceMgmtApi *MgmtServiceApiCtx) getKubernetesClustersKubernetesClusterName(params *model2.GetKubernetesClustersKubernetesClusterNameParams) (*model2.GetKubernetesClustersKubernetesClusterNameOK, error) {
	registry := new(strfmt.Registry)

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "GetKubernetesClustersKubernetesClusterName",
		Method:             "GET",
		PathPattern:        "/cd/kubernetesClusters/{kubernetesClusterName}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.GetKubernetesClustersKubernetesClusterNameReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.GetKubernetesClustersKubernetesClusterNameOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*model2.GetKubernetesClustersKubernetesClusterNameDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

func (serviceMgmtApi *MgmtServiceApiCtx) deleteKubernetesClustersKubernetesClusterID(params *model2.DeleteKubernetesClustersKubernetesClusterIDParams) (*model2.DeleteKubernetesClustersKubernetesClusterIDNoContent, error) {
	registry := new(strfmt.Registry)

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "DeleteKubernetesClustersKubernetesClusterID",
		Method:             "DELETE",
		PathPattern:        "/kubernetesClusters/{kubernetesClusterId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.DeleteKubernetesClustersKubernetesClusterIDReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.DeleteKubernetesClustersKubernetesClusterIDNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*model2.DeleteKubernetesClustersKubernetesClusterIDDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

func (serviceMgmtApi *MgmtServiceApiCtx) putKubernetesClustersKubernetesClusterID(params *model2.PutKubernetesClustersKubernetesClusterIDParamsWriter) (*model2.PutKubernetesClustersKubernetesClusterIDOK, error) {
	registry := new(strfmt.Registry)
	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PutKubernetesClustersKubernetesClusterID",
		Method:             "PUT",
		PathPattern:        "/kubernetesClusters/{kubernetesClusterId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.PutKubernetesClustersKubernetesClusterIDReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PutKubernetesClustersKubernetesClusterIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*model2.PutKubernetesClustersKubernetesClusterIDDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
deleteCdRuleIDConnectionsRule deletes a cd connection rule
*/
func (serviceMgmtApi *MgmtServiceApiCtx) deleteCdRuleIDConnectionsRule(params *model2.DeleteCdRuleIDConnectionsRuleParams) (*model2.DeleteCdRuleIDConnectionsRuleNoContent, error) {
	registry := new(strfmt.Registry)
	if params == nil {
		params = model2.NewDeleteCdRuleIDConnectionsRuleParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "DeleteCdRuleIDConnectionsRule",
		Method:             "DELETE",
		PathPattern:        "/cd/{ruleId}/connectionsRule",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.DeleteCdRuleIDConnectionsRuleReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.DeleteCdRuleIDConnectionsRuleNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	msg := fmt.Sprintf("unexpected success response for DeleteCdRuleIDConnectionsRule: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("delete cd connection rule", msg, 400)
}

/*
getCdRuleIDConnectionsRule gets a cd connection rule
*/
func (serviceMgmtApi *MgmtServiceApiCtx) getCdRuleIDConnectionsRule(params *model2.GetCdRuleIDConnectionsRuleParams) (*model2.GetCdRuleIDConnectionsRuleOK, error) {
	registry := new(strfmt.Registry)
	if params == nil {
		params = model2.NewGetCdRuleIDConnectionsRuleParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "GetCdRuleIDConnectionsRule",
		Method:             "GET",
		PathPattern:        "/cd/{ruleId}/connectionsRule",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.GetCdRuleIDConnectionsRuleReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.GetCdRuleIDConnectionsRuleOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetCdRuleIDConnectionsRule: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("get cd connection rule", msg, 400)
}

/*
postCdConnectionsRule adds cd connection rule
*/
func (serviceMgmtApi *MgmtServiceApiCtx) postCdConnectionsRule(params *model2.PostCdConnectionsRuleParams) (*model2.PostCdConnectionsRuleCreated, error) {
	registry := new(strfmt.Registry)
	if params == nil {
		params = model2.NewPostCdConnectionsRuleParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PostCdConnectionsRule",
		Method:             "POST",
		PathPattern:        "/cd/connectionsRule",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.PostCdConnectionsRuleReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PostCdConnectionsRuleCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostCdConnectionsRule: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("post cd connection rule", msg, 400)

}

/*
putCdRuleIDConnectionsRule updates a cd connection rule
*/
func (serviceMgmtApi *MgmtServiceApiCtx) putCdRuleIDConnectionsRule(params *model2.PutCdRuleIDConnectionsRuleParams) (*model2.PutCdRuleIDConnectionsRuleOK, error) {
	registry := new(strfmt.Registry)
	if params == nil {
		params = model2.NewPutCdRuleIDConnectionsRuleParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PutCdRuleIDConnectionsRule",
		Method:             "PUT",
		PathPattern:        "/cd/{ruleId}/connectionsRule",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.PutCdRuleIDConnectionsRuleReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PutCdRuleIDConnectionsRuleOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PutCdRuleIDConnectionsRule: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("put cd connection rule", msg, 400)

}

/*
postEnvironments adds a new environment
*/
func (serviceMgmtApi *MgmtServiceApiCtx) postEnvironments(params *model2.PostEnvironmentsParams, client *http.Client) (*model2.PostEnvironmentsCreated, error) {
	registry := new(strfmt.Registry)
	if params == nil {
		params = model2.NewPostEnvironmentsParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PostEnvironments",
		Method:             "POST",
		PathPattern:        "/environments",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.PostEnvironmentsReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             client,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PostEnvironmentsCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostEnvironments: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("post env", msg, 400)
}

/*
deleteEnvironmentEnvID deletes a SecureCN environment by ID.
*/
func (serviceMgmtApi *MgmtServiceApiCtx) deleteEnvironmentEnvID(params *model2.DeleteEnvironmentsEnvIDParams) (*model2.DeleteEnvironmentEnvIDNoContent, error) {
	registry := new(strfmt.Registry)
	// TODO: Validate the params before sending
	if params == nil {
		params = model2.NewDeleteEnvironmentEnvIDParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "DeleteEnvironmentEnvID",
		Method:             "DELETE",
		PathPattern:        "/environments/{envId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             params,
		Reader:             &model2.DeleteEnvironmentsEnvIDReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.DeleteEnvironmentEnvIDNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*model2.DeleteEnvironmentEnvIDDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
getEnvironmentsEnvID gets an environment
*/
func (serviceMgmtApi *MgmtServiceApiCtx) getEnvironmentsEnvID(params *model2.GetEnvironmentsEnvIDParams) (*model2.GetEnvironmentsEnvIDOK, error) {
	registry := new(strfmt.Registry)
	if params == nil {
		params = model2.NewGetEnvironmentsEnvIDParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "GetEnvironmentsEnvID",
		Method:             "GET",
		PathPattern:        "/environments/{envId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.GetEnvironmentsEnvIDReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.GetEnvironmentsEnvIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetEnvironmentsEnvID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("get env", msg, 400)
}

/*
putEnvironmentsEnvID edits an existing SecureCN environment

Edit an existing SecureCN environment.

*/
func (serviceMgmtApi *MgmtServiceApiCtx) putEnvironmentsEnvID(params *model2.PutEnvironmentsEnvIDParams) (*model2.PutEnvironmentsEnvIDOK, error) {
	registry := new(strfmt.Registry)
	if params == nil {
		params = model2.NewPutEnvironmentsEnvIDParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PutEnvironmentsEnvID",
		Method:             "PUT",
		PathPattern:        "/environments/{envId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.PutEnvironmentsEnvIDReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PutEnvironmentsEnvIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PutEnvironmentsEnvID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("put env", msg, 400)
}

/*
getCdPodSecurityPolicyProfilesPodSecurityPolicyProfileName gets an id of a psp profile by name
*/
func (serviceMgmtApi *MgmtServiceApiCtx) getCdPodSecurityPolicyProfilesPodSecurityPolicyProfileName(params *model2.GetCdPodSecurityPolicyProfilesPodSecurityPolicyProfileNameParams) (*model2.GetCdPodSecurityPolicyProfilesPodSecurityPolicyProfileNameOK, error) {
	registry := new(strfmt.Registry)

	if params == nil {
		params = model2.NewGetCdPodSecurityPolicyProfilesPodSecurityPolicyProfileNameParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "GetCdPodSecurityPolicyProfilesPodSecurityPolicyProfileName",
		Method:             "GET",
		PathPattern:        "/cd/podSecurityPolicyProfiles/{podSecurityPolicyProfileName}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.GetCdPodSecurityPolicyProfilesPodSecurityPolicyProfileNameReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.GetCdPodSecurityPolicyProfilesPodSecurityPolicyProfileNameOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetCdPodSecurityPolicyProfilesPodSecurityPolicyProfileName: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("get psp id by name", msg, 400)
}

/*
postCdDeploymentRule adds cd deployment rule
*/
func (serviceMgmtApi *MgmtServiceApiCtx) postCdDeploymentRule(params *model2.PostCdDeploymentRuleParams) (*model2.PostCdDeploymentRuleCreated, error) {
	registry := new(strfmt.Registry)

	if params == nil {
		params = model2.NewPostCdDeploymentRuleParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PostCdDeploymentRule",
		Method:             "POST",
		PathPattern:        "/cd/deploymentRule",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.PostCdDeploymentRuleReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PostCdDeploymentRuleCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostCdDeploymentRule: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("post deployment rule", msg, 400)
}

/*
getCdRuleIDDeploymentRule gets a cd deployment rule
*/
func (serviceMgmtApi *MgmtServiceApiCtx) getCdRuleIDDeploymentRule(params *model2.GetCdRuleIDDeploymentRuleParams) (*model2.GetCdRuleIDDeploymentRuleOK, error) {
	registry := new(strfmt.Registry)

	if params == nil {
		params = model2.NewGetCdRuleIDDeploymentRuleParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "GetCdRuleIDDeploymentRule",
		Method:             "GET",
		PathPattern:        "/cd/{ruleId}/deploymentRule",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.GetCdRuleIDDeploymentRuleReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.GetCdRuleIDDeploymentRuleOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetCdRuleIDDeploymentRule: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("get deployment rule", msg, 400)
}

/*
putCdRuleIDDeploymentRule updates a cd deployment rule
*/
func (serviceMgmtApi *MgmtServiceApiCtx) putCdRuleIDDeploymentRule(params *model2.PutCdRuleIDDeploymentRuleParams) (*model2.PutCdRuleIDDeploymentRuleOK, error) {
	registry := new(strfmt.Registry)

	if params == nil {
		params = model2.NewPutCdRuleIDDeploymentRuleParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PutCdRuleIDDeploymentRule",
		Method:             "PUT",
		PathPattern:        "/cd/{ruleId}/deploymentRule",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.PutCdRuleIDDeploymentRuleReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PutCdRuleIDDeploymentRuleOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PutCdRuleIDDeploymentRule: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("put deployment rule", msg, 400)
}

/*
deleteCdRuleIDDeploymentRule deletes a cd deployment rule
*/
func (serviceMgmtApi *MgmtServiceApiCtx) deleteCdRuleIDDeploymentRule(params *model2.DeleteCdRuleIDDeploymentRuleParams) (*model2.DeleteCdRuleIDDeploymentRuleNoContent, error) {
	registry := new(strfmt.Registry)

	if params == nil {
		params = model2.NewDeleteCdRuleIDDeploymentRuleParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "DeleteCdRuleIDDeploymentRule",
		Method:             "DELETE",
		PathPattern:        "/cd/{ruleId}/deploymentRule",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &model2.DeleteCdRuleIDDeploymentRuleReader{Formats: *registry},
		AuthInfo:           serviceMgmtApi.auth,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.DeleteCdRuleIDDeploymentRuleNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DeleteCdRuleIDDeploymentRule: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("delete deployment rule", msg, 400)
}

func (serviceMgmtApi *MgmtServiceApiCtx) setServiceKeys(accessKey string, secretKey []byte) {
	secretKeyStr := base64.StdEncoding.EncodeToString(secretKey)
	serviceMgmtApi.auth = auth2.NewAuth(accessKey, secretKeyStr, "global/services/portshift_request")
	serviceMgmtApi.runtime.DefaultAuthentication = serviceMgmtApi.auth
}

func (serviceMgmtApi *MgmtServiceApiCtx) GetDeployerById(ctx context.Context, deployerId strfmt.UUID) (model2.Deployer, error) {
	log.Print("[DEBUG] getting deployer")

	params := model2.GetDeployersParams{
		Context: ctx,
	}

	deployers, err := serviceMgmtApi.getDeployers(&params)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployers: %v", err)
	}

	for _, deployer := range deployers.Payload {
		if deployer.ID() == deployerId {
			return deployer, nil
		}
	}

	return nil, nil
}

func (serviceMgmtApi *MgmtServiceApiCtx) getDeployers(params *model2.GetDeployersParams) (*model2.GetDeployersOK, error) {
	registry := new(strfmt.Registry)

	// TODO: Validate the params before sending
	if params == nil {
		params = model2.NewGetDeployersParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "GetDeployers",
		Method:             "GET",
		PathPattern:        "/deployers",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             params,
		Reader:             &model2.GetDeployersReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.GetDeployersOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetDeployers: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("read deployer", msg, 400)
}

func (serviceMgmtApi *MgmtServiceApiCtx) CreateDeployer(ctx context.Context, deployer model2.Deployer) (*model2.PostDeployersCreated, error) {
	registry := new(strfmt.Registry)

	params := &model2.PostDeployersParams{
		Body:    deployer,
		Context: ctx,
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PostDeployers",
		Method:             "POST",
		PathPattern:        "/deployers",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             params,
		Reader:             &model2.PostDeployersReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PostDeployersCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostDeployers: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("post deployer", msg, 400)
}

/*
GetDeployersServiceAccountsByNamespace lists all the service account on the system
*/
func (serviceMgmtApi *MgmtServiceApiCtx) GetDeployersServiceAccountsByNamespace(ctx context.Context, clusterId strfmt.UUID, namespace string) (*model2.GetDeployersServiceAccountsOK, error) {
	registry := new(strfmt.Registry)

	params := model2.GetDeployersServiceAccountsParams{
		KubernetesClusterID: clusterId,
		NamespaceName:       &namespace,
		Context:             ctx,
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "GetDeployersServiceAccounts",
		Method:             "GET",
		PathPattern:        "/deployers/serviceAccounts",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             &params,
		Reader:             &model2.GetDeployersServiceAccountsReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.GetDeployersServiceAccountsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetDeployersServiceAccounts: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("list deployer SAs", msg, 400)
}

/*
DeleteDeployer deletes an deployer
*/
func (serviceMgmtApi *MgmtServiceApiCtx) DeleteDeployer(ctx context.Context, uuid strfmt.UUID) (*model2.DeleteDeployersDeployerIDNoContent, error) {
	registry := new(strfmt.Registry)

	params := model2.DeleteDeployersDeployerIDParams{
		DeployerID: uuid,
		Context:    ctx,
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "DeleteDeployersDeployerID",
		Method:             "DELETE",
		PathPattern:        "/deployers/{deployerId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             &params,
		Reader:             &model2.DeleteDeployersDeployerIDReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.DeleteDeployersDeployerIDNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DeleteDeployersDeployerID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("delete deployer", msg, 400)
}

func (serviceMgmtApi *MgmtServiceApiCtx) UpdateDeployer(ctx context.Context, deployer model2.Deployer) (*model2.PutDeployersDeployerIDOK, error) {
	registry := new(strfmt.Registry)

	params := model2.PutDeployersDeployerIDParams{
		DeployerID: deployer.ID(),
		Body:       deployer,
		Context:    ctx,
	}

	deployer.SetID("")

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PutDeployersDeployerID",
		Method:             "PUT",
		PathPattern:        "/deployers/{deployerId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             &params,
		Reader:             &model2.PutDeployersDeployerIDReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PutDeployersDeployerIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PutDeployersDeployerID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("update deployer", msg, 400)
}

/*
PostCdPolicy sets the current c d policy at least one cd policy element should be present
*/
func (serviceMgmtApi *MgmtServiceApiCtx) PostCdPolicy(params *model2.PostCdPolicyParams) (*model2.PostCdPolicyCreated, error) {
	registry := new(strfmt.Registry)

	// TODO: Validate the params before sending
	if params == nil {
		params = model2.NewPostCdPolicyParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PostCdPolicy",
		Method:             "POST",
		PathPattern:        "/cdPolicy",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             params,
		Reader:             &model2.PostCdPolicyReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PostCdPolicyCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostCdPolicy: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("create cd policy", msg, 400)
}

func (serviceMgmtApi *MgmtServiceApiCtx) GetCdPolicy(params *model2.GetCdPolicyParams) (*model2.GetCdPolicyOK, error) {
	registry := new(strfmt.Registry)

	// TODO: Validate the params before sending
	if params == nil {
		params = model2.NewGetCdPolicyParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "GetCdPolicy",
		Method:             "GET",
		PathPattern:        "/cdPolicy",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             params,
		Reader:             &model2.GetCdPolicyReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.GetCdPolicyOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetCdPolicy: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("get cd policy", msg, 400)
}

func (serviceMgmtApi *MgmtServiceApiCtx) PutCdPolicyPolicyID(params *model2.PutCdPolicyPolicyIDParams) (*model2.PutCdPolicyPolicyIDOK, error) {
	registry := new(strfmt.Registry)

	// TODO: Validate the params before sending
	if params == nil {
		params = model2.NewPutCdPolicyPolicyIDParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PutCdPolicyPolicyID",
		Method:             "PUT",
		PathPattern:        "/cdPolicy/{policyId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             params,
		Reader:             &model2.PutCdPolicyPolicyIDReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PutCdPolicyPolicyIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PutCdPolicyPolicyID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("put cd policy", msg, 400)
}

func (serviceMgmtApi *MgmtServiceApiCtx) DeleteCdPolicyPolicyID(params *model2.DeleteCdPolicyPolicyIDParams) (*model2.DeleteCdPolicyPolicyIDNoContent, error) {
	registry := new(strfmt.Registry)

	// TODO: Validate the params before sending
	if params == nil {
		params = model2.NewDeleteCdPolicyPolicyIDParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "DeleteCdPolicyPolicyID",
		Method:             "DELETE",
		PathPattern:        "/cdPolicy/{policyId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             params,
		Reader:             &model2.DeleteCdPolicyPolicyIDReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.DeleteCdPolicyPolicyIDNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DeleteCdPolicyPolicyID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("delete cd policy", msg, 400)
}

func (serviceMgmtApi *MgmtServiceApiCtx) PostCiPolicy(params *model2.PostCiPolicyParams) (*model2.PostCiPolicyCreated, error) {
	registry := new(strfmt.Registry)

	// TODO: Validate the params before sending
	if params == nil {
		params = model2.NewPostCiPolicyParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PostCiPolicy",
		Method:             "POST",
		PathPattern:        "/ciPolicy",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             params,
		Reader:             &model2.PostCiPolicyReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PostCiPolicyCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostCiPolicy: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("post ci policy", msg, 400)
}

func (serviceMgmtApi *MgmtServiceApiCtx) GetCiPolicy(params *model2.GetCiPolicyParams) (*model2.GetCiPolicyOK, error) {
	registry := new(strfmt.Registry)

	// TODO: Validate the params before sending
	if params == nil {
		params = model2.NewGetCiPolicyParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "GetCiPolicy",
		Method:             "GET",
		PathPattern:        "/ciPolicy",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             params,
		Reader:             &model2.GetCiPolicyReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.GetCiPolicyOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetCiPolicy: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("get ci policy", msg, 400)
}

func (serviceMgmtApi *MgmtServiceApiCtx) DeleteCiPolicyPolicyID(params *model2.DeleteCiPolicyPolicyIDParams) (*model2.DeleteCiPolicyPolicyIDNoContent, error) {
	registry := new(strfmt.Registry)

	// TODO: Validate the params before sending
	if params == nil {
		params = model2.NewDeleteCiPolicyPolicyIDParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "DeleteCiPolicyPolicyID",
		Method:             "DELETE",
		PathPattern:        "/ciPolicy/{policyId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             params,
		Reader:             &model2.DeleteCiPolicyPolicyIDReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.DeleteCiPolicyPolicyIDNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DeleteCiPolicyPolicyID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("delete ci policy", msg, 400)
}

func (serviceMgmtApi *MgmtServiceApiCtx) PutCiPolicyPolicyID(params *model2.PutCiPolicyPolicyIDParams) (*model2.PutCiPolicyPolicyIDOK, error) {
	registry := new(strfmt.Registry)

	// TODO: Validate the params before sending
	if params == nil {
		params = model2.NewPutCiPolicyPolicyIDParams()
	}

	result, err := serviceMgmtApi.runtime.Submit(&runtime.ClientOperation{
		ID:                 "PutCiPolicyPolicyID",
		Method:             "PUT",
		PathPattern:        "/ciPolicy/{policyId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		AuthInfo:           serviceMgmtApi.auth,
		Params:             params,
		Reader:             &model2.PutCiPolicyPolicyIDReader{Formats: *registry},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*model2.PutCiPolicyPolicyIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PutCiPolicyPolicyID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	return nil, runtime.NewAPIError("put ci policy", msg, 400)
}
