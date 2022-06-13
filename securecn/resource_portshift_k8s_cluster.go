package securecn

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"terraform-provider-securecn/internal/client"
	"terraform-provider-securecn/internal/escher_api/escherClient"
	"terraform-provider-securecn/internal/escher_api/model"
	utils2 "terraform-provider-securecn/internal/utils"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/spf13/cast"
)

const secureCNBundleFilePath = "securecn_bundle.tar.gz"
const scriptFilePath = "install_bundle.sh"
const yamlFilePath = "securecn_bundle.yml"
const patchDnsFilePath = "patch_dns.sh"
const certsGenFilePath = "certs_gen.sh"
const vaultCertsGenFilePath = "certs_gen_vault.sh"
const vaultCertsFolder = "vault_certs"
const uninstallAgentCommandFormat = "KUBECONFIG=%s kubectl get cm -n portshift portshift-uninstaller -o jsonpath='{.data.config}' | KUBECONFIG=%s bash"
const useK8sContextCommandFormat = "kubectl config use-context"
const viewK8sConfigCommand = "kubectl config view --raw"

const KubernetesClusterContextFieldName = "kubernetes_cluster_context"
const NameFieldName = "name"
const CiImageValidationFieldName = "ci_image_validation"
const RestrictRegistries = "restrict_registries"
const CdPodTemplateFieldName = "cd_pod_template"
const ConnectionsControlFieldName = "connections_control"
const IstioAlreadyInstalledFieldName = "istio_already_installed"
const IstioVersionFieldName = "istio_version"
const IstioIngressEnabledFieldName = "istio_ingress_enabled"
const IstioIngressAnnotationsFieldName = "istio_ingress_annotations"
const EnableApiIntelligenceDASTFieldName = "api_intelligence_dast"
const EnableAutoLabelFieldName = "auto_labeling"
const HoldApplicationUntilProxyStartsFieldName = "hold_application_until_proxy_starts"
const ExternalCAFieldName = "external_ca"
const InternalRegistryFieldName = "internal_registry"
const ServiceDiscoveryIsolationFieldName = "service_discovery_isolation"
const TLSInspectionFieldName = "tls_inspection"
const TokenInjectionFieldName = "token_injection"
const SkipReadyCheckFieldName = "skip_ready_check"
const TracingSupportFieldName = "tracing_support"
const TraceAnalyzerFieldName = "trace_analyzer"
const SpecReconstructionFieldName = "spec_reconstruction"
const SidecarResourcesFieldName = "sidecar_resources"
const MultiClusterCommunicationSupportFieldName = "multi_cluster_communication_support"
const MultiClusterCommunicationSupportCertsPathFieldName = MultiClusterCommunicationSupportFieldName + "_certs_path"
const InspectIncomingClusterConnectionsFieldName = "inspect_incoming_cluster_connections"
const FailCloseFieldName = "fail_close"
const PersistentStorageFieldName = "persistent_storage"
const ExternalHttpsProxyFieldName = "external_https_proxy"
const OrchestrationTypeFieldName = "orchestration_type"
const MinimumReplicasFieldName = "minimum_replicas"
const CiImageSignerValidationEnabledFieldName = "ci_image_signer_validation_enabled"

func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,
		Description:   "A SecureCN k8s cluster",
		Schema: map[string]*schema.Schema{
			KubernetesClusterContextFieldName:         {Type: schema.TypeString, Required: true, ForceNew: true, Description: "The k8s context name of the cluster", ValidateFunc: validation.StringIsNotEmpty},
			NameFieldName:                             {Type: schema.TypeString, Required: true, Description: "The name of cluster in SecureCN"},
			CiImageValidationFieldName:                {Type: schema.TypeBool, Optional: true, Default: false, Description: "Identify pods only if the image hash matches the value generated by the CI plugin or entered manually in the UI"},
			CdPodTemplateFieldName:                    {Type: schema.TypeBool, Optional: true, Default: false, Description: "Identify pod templates only originating from SecureCN CD plugin"},
			RestrictRegistries:                        {Type: schema.TypeBool, Optional: true, Default: false, Description: "Workload from untrusted registries will be marked as 'unknown'"},
			ConnectionsControlFieldName:               {Type: schema.TypeBool, Optional: true, Default: true, Description: "Enable connections control"},
			IstioAlreadyInstalledFieldName:			   {Type: schema.TypeBool, Optional: true, Default: false, Description: "if false, istio will be installed"},
			IstioVersionFieldName:                     {Type: schema.TypeString, Optional: true, Default: nil, Computed: true, Description: "if IstioAlreadyInstalled, this specify its version"},
			IstioIngressEnabledFieldName:              {Type: schema.TypeBool, Optional: true, Computed: true, Description: "If installing Istio, use Istio ingress"},
			IstioIngressAnnotationsFieldName:          {Type: schema.TypeMap, Elem: schema.TypeString, Optional: true, Default: map[string]string{}, Description: "If enabling Istio ingress, use Istio these ingress annotation"},
			MultiClusterCommunicationSupportFieldName: {Type: schema.TypeBool, Optional: true, Default: false, Description: "Enable multi cluster communication"},
			MultiClusterCommunicationSupportCertsPathFieldName: {Type: schema.TypeString, Optional: true, Default: "", Description: "Multi cluster certs path. Only valid if " + MultiClusterCommunicationSupportFieldName + " is true",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					path := val.(string)
					if _, err := os.Stat(path); path != "" && os.IsNotExist(err) {
						errs = append(errs, fmt.Errorf("folder %s does not exist", path))
					}
					return
				},
			},
			InspectIncomingClusterConnectionsFieldName: {Type: schema.TypeBool, Optional: true, Default: false, Description: "Enable enforcement and visibility of connections from external IP sources"},
			FailCloseFieldName:                         {Type: schema.TypeBool, Optional: true, Default: false, Description: "When enabled, workloads and connections will be blocked in case SecureCN agent is not responding"},
			PersistentStorageFieldName:                 {Type: schema.TypeBool, Optional: true, Default: false, Description: "Allow SecureCN agent to save the policy persistently, so it will be available after a restart of the pod. This will Require 128MB of storage for the agent pod."},
			ExternalHttpsProxyFieldName:                {Type: schema.TypeString, Optional: true, Default: "", Description: "Proxy definitions for outgoing HTTPS traffic from the cluster, if needed"},
			OrchestrationTypeFieldName: {Type: schema.TypeString, Optional: true, Default: "KUBERNETES", Description: "Orchestration type of the kubernetes cluster",
				ValidateFunc: validation.StringInSlice([]string{
					"GKE", "OPENSHIFT", "RANCHER", "AKS", "EKS", "KUBERNETES", "IKS",
				}, true),
			},
			EnableApiIntelligenceDASTFieldName:       {Type: schema.TypeBool, Optional: true, Default: false, Description: "Enable API Intelligence DAST integration"},
			EnableAutoLabelFieldName:                 {Type: schema.TypeBool, Optional: true, Default: false, Description: "Enable auto labeling of Kubernetes namespaces"},
			HoldApplicationUntilProxyStartsFieldName: {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether the controller should hold the application until the proxy starts"},
			ServiceDiscoveryIsolationFieldName:       {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether the service discovery isolation is enabled"},
			TLSInspectionFieldName:                   {Type: schema.TypeBool, Optional: true, Computed: true, Description: "Indicates whether the TLS inspection is enabled"},
			TokenInjectionFieldName:                  {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether the token injection is enabled"},
			SkipReadyCheckFieldName:                  {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether the cluster installation should be async"},
			TracingSupportFieldName:                  {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether to install tracing support, enable for apiSecurity accounts."},
			TraceAnalyzerFieldName:                   {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether the trace analyzer is enabled"},
			SpecReconstructionFieldName:              {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether the OpenAPI specification reconstruction is enabled"},
			ExternalCAFieldName: {
				Description: "Use an external CA for this cluster",
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description:  "The id of the external CA",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"name": {
							Description: "The name of the external CA",
							Optional:    true,
							Type:        schema.TypeString,
						},
					},
				},
			},
			InternalRegistryFieldName: {
				Description: "Use an internal container registry for this cluster",
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Description: "The url of the internal registry",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			SidecarResourcesFieldName: {
				Description: "Define resource limits for Istio sidecars",
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_init_limits_cpu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"proxy_init_limits_memory": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"proxy_init_requests_cpu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"proxy_init_requests_memory": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"proxy_limits_cpu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"proxy_limits_memory": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"proxy_requests_cpu": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"proxy_requests_memory": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			MinimumReplicasFieldName: {Type: schema.TypeInt, Optional: true, Default: 1, Description: "minimum number of controller replicas",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					minReplicas := val.(int)
					if minReplicas < 1 || minReplicas > 5 {
						errs = append(errs, fmt.Errorf("%s should be between 1 and 5 (inclusive)", MinimumReplicasFieldName))
					}
					return
				},
			},
			CiImageSignerValidationEnabledFieldName: {Type: schema.TypeBool, Optional: true, Default: false, Description: "indicates whether ci image signer validation is Enabled"},
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] creating cluster")

	err := validateConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	kubernetesCluster, err := getClusterFromConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	secureCNCluster, err := serviceApi.CreateKubernetesCluster(ctx, httpClientWrapper.HttpClient, kubernetesCluster)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId := secureCNCluster.Payload.ID
	k8sContext := d.Get(KubernetesClusterContextFieldName).(string)
	multiClusterFolder := d.Get(MultiClusterCommunicationSupportCertsPathFieldName).(string)
	tokenInjection := d.Get(TokenInjectionFieldName).(bool)
	skipReadyCheck := d.Get(SkipReadyCheckFieldName).(bool)

	err = installAgent(ctx, serviceApi, httpClientWrapper, clusterId, k8sContext, multiClusterFolder, tokenInjection, skipReadyCheck)
	if err != nil {
		_ = serviceApi.DeleteKubernetesCluster(ctx, httpClientWrapper.HttpClient, clusterId)
		return diag.FromErr(err)
	}

	d.SetId(string(clusterId))
	return resourceClusterRead(ctx, d, m)
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] reading cluster")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	clusterId := d.Id()

	secureCNCluster, err := serviceApi.GetKubernetesClusterById(ctx, httpClientWrapper.HttpClient, strfmt.UUID(clusterId))
	if err != nil {
		return diag.FromErr(err)
	}

	if secureCNCluster.Payload.ID == "" {
		// Tell terraform the cluster doesn't exist
		d.SetId("")
	} else {
		kubernetesClusterFromConfig, err := getClusterFromConfig(d)
		if err != nil {
			return diag.FromErr(err)
		}
		k8sContext := d.Get(KubernetesClusterContextFieldName).(string)
		certsFolder := d.Get(MultiClusterCommunicationSupportCertsPathFieldName).(string)
		tokenInjection := d.Get(TokenInjectionFieldName).(bool)
		skipReadyCheck := d.Get(SkipReadyCheckFieldName).(bool)
		err = updateAgent(k8sContext, certsFolder, kubernetesClusterFromConfig, tokenInjection, *secureCNCluster.Payload.IsMultiCluster, *secureCNCluster.Payload.EnableConnectionsControl, *secureCNCluster.Payload.AgentFailClose, *secureCNCluster.Payload.IsPersistent, *secureCNCluster.Payload.ProxyConfiguration, serviceApi, httpClientWrapper, strfmt.UUID(clusterId), skipReadyCheck)
		if err != nil {
			return diag.FromErr(err)
		}
		updateMutableFields(d, secureCNCluster.Payload)
	}

	return nil
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] updating cluster")

	err := validateConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)

	kubernetesClusterFromConfig, err := getClusterFromConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId := d.Id()
	updatedCluster, err := serviceApi.UpdateKubernetesCluster(ctx, httpClientWrapper.HttpClient, kubernetesClusterFromConfig, strfmt.UUID(clusterId))
	if err != nil {
		return diag.FromErr(err)
	}

	kubernetesClusterFromConfig, err = getClusterFromConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	k8sContext := d.Get(KubernetesClusterContextFieldName).(string)
	certsFolder := d.Get(MultiClusterCommunicationSupportCertsPathFieldName).(string)
	tokenInjection := d.Get(TokenInjectionFieldName).(bool)

	err = updateAgent(k8sContext, certsFolder, kubernetesClusterFromConfig, tokenInjection, *updatedCluster.Payload.IsMultiCluster, *updatedCluster.Payload.EnableConnectionsControl, *updatedCluster.Payload.AgentFailClose, *updatedCluster.Payload.IsPersistent, *updatedCluster.Payload.ProxyConfiguration, serviceApi, httpClientWrapper, strfmt.UUID(clusterId), false) //TODO update other fields, tests
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(updatedCluster.Payload.ID))
	return resourceClusterRead(ctx, d, m)
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] deleting cluster")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	clusterId := strfmt.UUID(d.Id())

	k8sContext := d.Get(KubernetesClusterContextFieldName).(string)
	err := deleteAgent(k8sContext)
	if err != nil {
		return diag.FromErr(err)
	}

	err = serviceApi.DeleteKubernetesCluster(ctx, httpClientWrapper.HttpClient, clusterId)
	if err != nil {
		return diag.FromErr(err)
	}

	// Tell terraform the cluster doesn't exist
	d.SetId("")

	return nil
}

func installAgent(ctx context.Context, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper, clusterId strfmt.UUID, context string, multiClusterFolder string, tokenInjection bool, skipReadyCheck bool) error {
	log.Print("[DEBUG] installing agent")
	err := downloadAndExtractBundle(ctx, serviceApi, httpClientWrapper, clusterId)
	if err != nil {
		return err
	}

	if tokenInjection {
		err = utils2.MakeExecutable(vaultCertsGenFilePath)
		if err != nil {
			return err
		}
	}

	kubeconfig, err := createTempKubeconfig(context)
	if err != nil {
		return err
	}

	defer os.Remove(kubeconfig)

	output, err := utils2.ExecuteScript(scriptFilePath, multiClusterFolder, skipReadyCheck, kubeconfig)
	if err != nil {
		return fmt.Errorf("%s:\n%s", err, output)
	}

	log.Print("[DEBUG] agent installed successfully")

	defer os.Remove(yamlFilePath)
	defer os.Remove(scriptFilePath)

	if multiClusterFolder != "" {
		defer os.Remove(patchDnsFilePath)
		defer os.Remove(certsGenFilePath)
	}

	if tokenInjection {
		defer os.Remove(vaultCertsGenFilePath)
		defer os.RemoveAll(vaultCertsFolder)
	}

	return nil
}

func deleteAgent(context string) error {
	log.Printf("[DEBUG] deleting agent from context: " + context)

	kubeconfig, err := createTempKubeconfig(context)
	if err != nil {
		return err
	}

	defer os.Remove(kubeconfig)

	_, err = utils2.ExecBashCommand(fmt.Sprintf(uninstallAgentCommandFormat, kubeconfig, kubeconfig))
	if err != nil {
		return err
	}

	return nil
}

func createTempKubeconfig(context string) (string, error) {
	log.Print("[DEBUG] changing k8s context to " + context)

	kubeconfig, err := utils2.ExecBashCommand(viewK8sConfigCommand)
	if err != nil {
		log.Print("[DEBUG] failed to print k8s config: " + err.Error())
		return "", err
	}

	kubeconfigfile, err := ioutil.TempFile(".", "kubeconfig")
	if err != nil {
		log.Print("[DEBUG] failed to create temporary k8s config: " + err.Error())
		return "", err
	}

	_, err = kubeconfigfile.WriteString(kubeconfig)
	if err != nil {
		log.Print("[DEBUG] failed to write temporary k8s config: " + err.Error())
		return "", err
	}

	err = kubeconfigfile.Close()
	if err != nil {
		log.Print("[DEBUG] failed to close temporary k8s config: " + err.Error())
		return "", err
	}

	changeContextCommand := fmt.Sprintf("KUBECONFIG=%s %s %s", kubeconfigfile.Name(), useK8sContextCommandFormat, context)
	_, err = utils2.ExecBashCommand(changeContextCommand)
	if err != nil {
		log.Print("[DEBUG] failed to change k8s context: " + err.Error())
		return "", err
	}

	return kubeconfigfile.Name(), nil
}

func downloadAndExtractBundle(ctx context.Context, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper, clusterId strfmt.UUID) error {
	log.Print("[DEBUG] downloading and extracting bundle")
	err := downloadInstallBundle(ctx, serviceApi, httpClientWrapper.HttpClient, clusterId)
	if err != nil {
		return err
	}
	open, err := os.Open(secureCNBundleFilePath)
	if err != nil {
		return err
	}
	defer os.Remove(secureCNBundleFilePath)

	err = utils2.ExtractTarGz(open)
	if err != nil {
		return err
	}
	return nil
}

func getClusterFromConfig(d *schema.ResourceData) (*model.KubernetesCluster, error) {
	log.Print("[DEBUG] getting cluster from config")

	clusterName := d.Get(NameFieldName).(string)
	ciImageValidation := d.Get(CiImageValidationFieldName).(bool)
	cdPodTemplate := d.Get(CdPodTemplateFieldName).(bool)
	restrictRegistries := d.Get(RestrictRegistries).(bool)
	connectionsControl := d.Get(ConnectionsControlFieldName).(bool)
	istioAlredyInstalled := d.Get(IstioAlreadyInstalledFieldName).(bool)
	istioVersion := d.Get(IstioVersionFieldName).(string)
	istioIngressEnabled := d.Get(IstioIngressEnabledFieldName).(bool)
	istioIngressAnnotationsRaw := cast.ToStringMapString(d.Get(IstioIngressAnnotationsFieldName))
	var istioIngressAnnotations []*model.KubernetesAnnotation
	for k, v := range istioIngressAnnotationsRaw {
		istioIngressAnnotations = append(istioIngressAnnotations, &model.KubernetesAnnotation{
			Key:   &k,
			Value: &v,
		})
	}
	supportsMultiClusterCommunication := d.Get(MultiClusterCommunicationSupportFieldName).(bool)
	inspectIncomingClusterConnections := d.Get(InspectIncomingClusterConnectionsFieldName).(bool)
	failClose := d.Get(FailCloseFieldName).(bool)
	ciImageSignerValidationEnabled := d.Get(CiImageSignerValidationEnabledFieldName).(bool)
	persistentStorage := d.Get(PersistentStorageFieldName).(bool)
	externalHttpsProxy := d.Get(ExternalHttpsProxyFieldName).(string)
	orchestrationType := d.Get(OrchestrationTypeFieldName).(string)
	minimumReplicas := d.Get(MinimumReplicasFieldName).(int)

	enableProxy := externalHttpsProxy != ""
	clusterPodDefinitionSource := model.ClusterPodDefinitionSourceKUBERNETES
	if cdPodTemplate {
		clusterPodDefinitionSource = model.ClusterPodDefinitionSourceCD
	}

	proxyConfig := &model.ProxyConfiguration{
		EnableProxy: &enableProxy,
		HTTPSProxy:  externalHttpsProxy,
	}
	istioParams := &model.IstioInstallationParameters{
		IsIstioAlreadyInstalled: &istioAlredyInstalled,
		IstioVersion:            istioVersion,
	}
	enableAPIIntelligenceDAST := d.Get(EnableApiIntelligenceDASTFieldName).(bool)
	enableAutoLabel := d.Get(EnableAutoLabelFieldName).(bool)
	holdApplicationUntilProxyStarts := d.Get(HoldApplicationUntilProxyStartsFieldName).(bool)
	enableServiceDiscoveryIsolation := d.Get(ServiceDiscoveryIsolationFieldName).(bool)
	enableTLSInspection := d.Get(TLSInspectionFieldName).(bool)
	enableTokenInjection := d.Get(TokenInjectionFieldName).(bool)
	tracingSupportEnabled := d.Get(TracingSupportFieldName).(bool)
	traceAnalyzerEnabled := d.Get(TraceAnalyzerFieldName).(bool)
	specReconstructionEnabled := d.Get(SpecReconstructionFieldName).(bool)
	installTracingSupport := tracingSupportEnabled || traceAnalyzerEnabled || specReconstructionEnabled

	cluster := &model.KubernetesCluster{
		AgentFailClose:                    &failClose,
		APIIntelligenceDAST:               &enableAPIIntelligenceDAST,
		AutoLabelEnabled:                  &enableAutoLabel,
		CiImageValidation:                 &ciImageValidation,
		ClusterPodDefinitionSource:        clusterPodDefinitionSource,
		EnableConnectionsControl:          &connectionsControl,
		ID:                                "",
		IsHoldApplicationUntilProxyStarts: &holdApplicationUntilProxyStarts,
		IsIstioIngressEnabled:             &istioIngressEnabled,
		IsMultiCluster:                    &supportsMultiClusterCommunication,
		IsPersistent:                      &persistentStorage,
		IstioIngressAnnotations:           istioIngressAnnotations,
		IstioInstallationParameters:       istioParams,
		Name:                              &clusterName,
		OrchestrationType:                 &orchestrationType,
		PreserveOriginalSourceIP:          &inspectIncomingClusterConnections,
		ProxyConfiguration:                proxyConfig,
		RestrictRegistires:                &restrictRegistries,
		ServiceDiscoveryIsolationEnabled:  &enableServiceDiscoveryIsolation,
		TLSInspectionEnabled:              &enableTLSInspection,
		TokenInjectionEnabled:             &enableTokenInjection,
		MinimalNumberOfControllerReplicas: minimumReplicas,
		CiImageSignerValidationEnabled:    &ciImageSignerValidationEnabled,
	}

	if installTracingSupport {
		cluster.TracingSupportSettings = &model.TracingSupportSettings{
			InstallTracingSupport:    &installTracingSupport,
			TraceAnalyzerEnabled:     &traceAnalyzerEnabled,
			SpecReconstructorEnabled: &specReconstructionEnabled,
		}
	}

	externalCaId := utils2.ReadNestedStringFromTF(d, ExternalCAFieldName, "id", 0)
	externalCaName := utils2.ReadNestedStringFromTF(d, ExternalCAFieldName, "name", 0)
	if externalCaId != "" {
		cluster.ExternalCa = &model.ExternalCaDetails{
			ID:   strfmt.UUID(externalCaId),
			Name: externalCaName,
		}
	}

	internalRegistryUrl := utils2.ReadNestedStringFromTF(d, InternalRegistryFieldName, "url", 0)
	if internalRegistryUrl != "" {
		internalRegistryEnabled := true
		cluster.InternalRegistryParameters = &model.InternalRegistryParameters{
			InternalRegistryEnabled: &internalRegistryEnabled,
			InternalRegistry:        internalRegistryUrl,
		}
	}

	proxyInitLimitsCpu := utils2.ReadNestedIntFromTF(d, SidecarResourcesFieldName, "proxy_init_limits_cpu", 0)
	proxyInitLimitsMemory := utils2.ReadNestedIntFromTF(d, SidecarResourcesFieldName, "proxy_init_limits_memory", 0)
	proxyInitRequestsCpu := utils2.ReadNestedIntFromTF(d, SidecarResourcesFieldName, "proxy_init_requests_cpu", 0)
	proxyInitRequestsMemory := utils2.ReadNestedIntFromTF(d, SidecarResourcesFieldName, "proxy_init_requests_memory", 0)
	proxyLimitsCpu := utils2.ReadNestedIntFromTF(d, SidecarResourcesFieldName, "proxy_limits_cpu", 0)
	proxyLimitsMemory := utils2.ReadNestedIntFromTF(d, SidecarResourcesFieldName, "proxy_limits_memory", 0)
	proxyRequestsCpu := utils2.ReadNestedIntFromTF(d, SidecarResourcesFieldName, "proxy_requests_cpu", 0)
	proxyRequestsMemory := utils2.ReadNestedIntFromTF(d, SidecarResourcesFieldName, "proxy_requests_memory", 0)
	cluster.SidecarsResources = &model.SidecarsResource{
		ProxyInitLimitsCPU:      int64(proxyInitLimitsCpu),
		ProxyInitLimitsMemory:   int64(proxyInitLimitsMemory),
		ProxyInitRequestsCPU:    int64(proxyInitRequestsCpu),
		ProxyInitRequestsMemory: int64(proxyInitRequestsMemory),
		ProxyLimitsCPU:          int64(proxyLimitsCpu),
		ProxyLimitsMemory:       int64(proxyLimitsMemory),
		ProxyRequestCPU:         int64(proxyRequestsCpu),
		ProxyRequestMemory:      int64(proxyRequestsMemory),
	}

	return cluster, nil
}

func downloadInstallBundle(ctx context.Context, serviceApi *escherClient.MgmtServiceApiCtx, client *http.Client, clusterId strfmt.UUID) error {
	log.Print("[DEBUG] downloading file")

	file, err := os.Create(secureCNBundleFilePath)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	err = serviceApi.DownloadKubernetesSecureCNBundle(ctx, client, buffer, clusterId)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, buffer)
	if err != nil {
		return err
	}

	return nil
}

func updateMutableFields(d *schema.ResourceData, secureCNCluster *model.KubernetesCluster) {
	log.Print("[DEBUG] updating mutable fields agent")

	_ = d.Set(NameFieldName, secureCNCluster.Name)
	_ = d.Set(CiImageValidationFieldName, secureCNCluster.CiImageValidation)
	_ = d.Set(CdPodTemplateFieldName, secureCNCluster.ClusterPodDefinitionSource == "CD")
	_ = d.Set(ConnectionsControlFieldName, secureCNCluster.EnableConnectionsControl)
	_ = d.Set(IstioAlreadyInstalledFieldName, secureCNCluster.IstioInstallationParameters.IsIstioAlreadyInstalled)
	_ = d.Set(IstioVersionFieldName, secureCNCluster.IstioInstallationParameters.IstioVersion)
	_ = d.Set(MultiClusterCommunicationSupportFieldName, secureCNCluster.IsMultiCluster)
	_ = d.Set(InspectIncomingClusterConnectionsFieldName, secureCNCluster.PreserveOriginalSourceIP)
	_ = d.Set(FailCloseFieldName, secureCNCluster.AgentFailClose)
	_ = d.Set(CiImageSignerValidationEnabledFieldName, secureCNCluster.CiImageSignerValidationEnabled)
	_ = d.Set(PersistentStorageFieldName, secureCNCluster.IsPersistent)
	_ = d.Set(ExternalHttpsProxyFieldName, secureCNCluster.ProxyConfiguration.HTTPSProxy)
	_ = d.Set(OrchestrationTypeFieldName, secureCNCluster.OrchestrationType)
	_ = d.Set(TLSInspectionFieldName, secureCNCluster.TLSInspectionEnabled)
	_ = d.Set(TokenInjectionFieldName, secureCNCluster.TokenInjectionEnabled)
	_ = d.Set(ServiceDiscoveryIsolationFieldName, secureCNCluster.ServiceDiscoveryIsolationEnabled)
}

func validateConfig(d *schema.ResourceData) error {
	log.Printf("[DEBUG] validating config")

	isMultiCluster := d.Get(MultiClusterCommunicationSupportFieldName).(bool)
	multiClusterFolder := d.Get(MultiClusterCommunicationSupportCertsPathFieldName).(string)
	connectionsControl := d.Get(ConnectionsControlFieldName).(bool)
	inspectIncomingClusterConnections := d.Get(InspectIncomingClusterConnectionsFieldName).(bool)
	if isMultiCluster && multiClusterFolder == "" {
		return errors.New(fmt.Sprintf("invalid configuration. %s can't be empty when %s is true", MultiClusterCommunicationSupportCertsPathFieldName, MultiClusterCommunicationSupportFieldName))
	}

	if !connectionsControl && isMultiCluster {
		return errors.New(fmt.Sprintf("invalid configuration. %s is off but %s is on ", ConnectionsControlFieldName, MultiClusterCommunicationSupportFieldName))
	}

	if !connectionsControl && inspectIncomingClusterConnections {
		return errors.New(fmt.Sprintf("invalid configuration. %s is off but %s is on", MultiClusterCommunicationSupportCertsPathFieldName, InspectIncomingClusterConnectionsFieldName))
	}

	return nil
}

func updateAgent(context string, multiClusterFolder string, clusterInTerraformConfig *model.KubernetesCluster, tokenInjection bool, prevIsMultiCluster bool, prevConnectionControl bool, prevAgentFailClose bool, prevIsPersistent bool, prevProxyConfiguration model.ProxyConfiguration, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper, clusterId strfmt.UUID, skipReadyCheck bool) error {
	needsUpdate := *clusterInTerraformConfig.IsMultiCluster != prevIsMultiCluster
	needsUpdate = needsUpdate || *clusterInTerraformConfig.EnableConnectionsControl != prevConnectionControl
	needsUpdate = needsUpdate || *clusterInTerraformConfig.AgentFailClose != prevAgentFailClose
	needsUpdate = needsUpdate || *clusterInTerraformConfig.IsPersistent != prevIsPersistent
	needsUpdate = needsUpdate || *clusterInTerraformConfig.ProxyConfiguration.EnableProxy != *prevProxyConfiguration.EnableProxy
	needsUpdate = needsUpdate || clusterInTerraformConfig.ProxyConfiguration.HTTPSProxy != prevProxyConfiguration.HTTPSProxy

	if needsUpdate {
		log.Print("[DEBUG] updating agent")
		err := deleteAgent(context)
		if err != nil {
			return err
		}
		err = installAgent(nil, serviceApi, httpClientWrapper, clusterId, context, multiClusterFolder, tokenInjection, skipReadyCheck)
		if err != nil {
			return err
		}
	}
	return nil
}
