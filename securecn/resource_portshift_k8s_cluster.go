package securecn

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"
	"terraform-provider-securecn/internal/client"
	"terraform-provider-securecn/internal/escher_api/escherClient"
	"terraform-provider-securecn/internal/escher_api/model"
	utils2 "terraform-provider-securecn/internal/utils"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/spf13/cast"
)

const installationDirPrefix = ".kubernetes_controller_installation_path_"
const secureCNBundleFilePath = "securecn_bundle.tar.gz"
const scriptFilePath = "install_bundle.sh"
const uninstallCmd = "./" + scriptFilePath + " --uninstall"

const vaultCertsGenFilePath = "certs_gen_vault.sh"
const tracingCertsFilePath = "certs_gen_tracing.sh"
const forceRemoveVaultCmd = " FORCE_REMOVE_VAULT=\"TRUE\""
const bash = " bash"
const getPortshiftPodsFormat = "KUBECONFIG=%s kubectl get pods -n portshift"
const describePortshiftPodsFormat = "KUBECONFIG=%s kubectl describe pods -n portshift"
const useK8sContextCommandFormat = "kubectl config use-context"
const viewK8sConfigCommand = "kubectl config view --raw"

const KubernetesClusterContextFieldName = "kubernetes_cluster_context"
const NameFieldName = "name"
const CiImageValidationFieldName = "ci_image_validation"
const RestrictRegistriesFieldName = "restrict_registries"
const CdPodTemplateFieldName = "cd_pod_template"
const ConnectionsControlFieldName = "connections_control"
const KubernetesSecurityFieldName = "kubernetes_security"
const IstioAlreadyInstalledFieldName = "istio_already_installed"
const IstioVersionFieldName = "istio_version"
const IstioIngressEnabledFieldName = "istio_ingress_enabled"
const IstioIngressAnnotationsFieldName = "istio_ingress_annotations"
const EnableApiIntelligenceDASTFieldName = "api_intelligence_dast"
const EnableAutoLabelFieldName = "auto_labeling"
const HoldApplicationUntilProxyStartsFieldName = "hold_application_until_proxy_starts"
const ExternalCAFieldName = "enable_external_ca"
const InternalRegistryFieldName = "internal_registry"
const InternalRegistryFieldNameUrl = "url"
const ServiceDiscoveryIsolationFieldName = "service_discovery_isolation"
const TLSInspectionFieldName = "tls_inspection"
const EnableK8sEventsFieldName = "enable_k8s_events"
const DisableSshMonitorFieldName = "disable_ssh_probing"
const TokenInjectionFieldName = "token_injection"
const SkipReadyCheckFieldName = "skip_ready_check"
const RollbackOnControllerFailureFieldName = "rollback_on_controller_failure"
const ForceRemoveVaultOnDeleteFieldName = "force_remove_vault_on_delete"
const InstallTracingSupportFieldName = "install_tracing_support"
const InstallEnvoyTracingSupportFieldName = "install_envoy_tracing_support"
const SidecarResourcesFieldName = "sidecar_resources"
const SidecarResourcesFieldNameProxyInitLimitsCpu = "proxy_init_limits_cpu"
const SidecarResourcesFieldNameProxyInitLimitsMemory = "proxy_init_limits_memory"
const SidecarResourcesFieldNameProxyInitRequestsCpu = "proxy_init_requests_cpu"
const SidecarResourcesFieldNameProxyInitRequestsMemory = "proxy_init_requests_memory"
const SidecarResourcesFieldNameProxyLimitsCpu = "proxy_limits_cpu"
const SidecarResourcesFieldNameProxyLimitsMemory = "proxy_limits_memory"
const SidecarResourcesFieldNameProxyRequestsCpu = "proxy_requests_cpu"
const SidecarResourcesFieldNameProxyRequestsMemory = "proxy_requests_memory"
const MultiClusterCommunicationSupportFieldName = "multi_cluster_communication_support"
const MultiClusterCommunicationSupportCertsPathFieldName = MultiClusterCommunicationSupportFieldName + "_certs_path"
const InspectIncomingClusterConnectionsFieldName = "inspect_incoming_cluster_connections"
const FailCloseFieldName = "fail_close"
const PersistentStorageFieldName = "persistent_storage"
const ExternalHttpsProxyFieldName = "external_https_proxy"
const OrchestrationTypeFieldName = "orchestration_type"
const MinimumReplicasFieldName = "minimum_replicas"
const CiImageSignatureValidationFieldName = "ci_image_signer_validation_enabled"
const SupportExternalTraceSourceFieldName = "support_external_trace_source"
const AutoUpgradeControllerVersionFieldName = "auto_upgrade_controller_version"

func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,
		Description:   "A Panoptica k8s cluster, Helm v3.8.0 or higher required",
		Schema: map[string]*schema.Schema{
			KubernetesClusterContextFieldName: {Type: schema.TypeString, Required: true, ForceNew: true, Description: "The k8s context name of the cluster", ValidateFunc: validation.StringIsNotEmpty},
			NameFieldName:                     {Type: schema.TypeString, Required: true, Description: "The name of cluster in SecureCN"},
			CiImageValidationFieldName:        {Type: schema.TypeBool, Optional: true, Default: false, Description: "Identify pods only if the image hash matches the value generated by the CI plugin or entered manually in the UI"},
			CdPodTemplateFieldName:            {Type: schema.TypeBool, Optional: true, Default: false, Description: "Identify pod templates only originating from SecureCN CD plugin"},
			RestrictRegistriesFieldName:       {Type: schema.TypeBool, Optional: true, Default: false, Description: "Workload from untrusted registries will be marked as 'unknown'"},
			ConnectionsControlFieldName:       {Type: schema.TypeBool, Optional: true, Default: true, Description: "Enable connections control"},
			KubernetesSecurityFieldName:       {Type: schema.TypeBool, Optional: true, Default: true, Description: "Enable kubernetes security"},
			IstioAlreadyInstalledFieldName:    {Type: schema.TypeBool, Optional: true, Default: false, Description: "if false, istio will be installed, otherwise the controller will use the previously installed istio"},
			IstioVersionFieldName:             {Type: schema.TypeString, Optional: true, Default: nil, Computed: true, Description: "if istio already installed, this specifies its version"},
			IstioIngressEnabledFieldName:      {Type: schema.TypeBool, Optional: true, Computed: true, Description: "If installing Istio, use Istio ingress"},
			IstioIngressAnnotationsFieldName: {Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				}, Optional: true, Description: "when enabling Istio ingress, use these Istio ingress annotations"},
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
			OrchestrationTypeFieldName: {Type: schema.TypeString, Optional: true, Default: "KUBERNETES", Description: "Orchestration type of the kubernetes cluster optional values: GKE, OPENSHIFT, RANCHER, AKS, EKS, KUBERNETES, IKS.",
				ValidateFunc: validation.StringInSlice([]string{
					"GKE", "OPENSHIFT", "RANCHER", "AKS", "EKS", "KUBERNETES", "IKS",
				}, true),
			},
			EnableApiIntelligenceDASTFieldName:       {Type: schema.TypeBool, Optional: true, Default: false, Description: "Enable API Intelligence DAST integration"},
			EnableAutoLabelFieldName:                 {Type: schema.TypeBool, Optional: true, Default: false, Description: "Enable auto labeling of Kubernetes namespaces"},
			HoldApplicationUntilProxyStartsFieldName: {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether the controller should hold the application until the proxy starts"},
			ServiceDiscoveryIsolationFieldName:       {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether the service discovery isolation is enabled"},
			TLSInspectionFieldName:                   {Type: schema.TypeBool, Optional: true, Computed: true, Description: "Indicates whether the TLS inspection is enabled"},
			EnableK8sEventsFieldName:                 {Type: schema.TypeBool, Optional: true, Computed: true, Description: "indicates whether kubernetes events sending is enabled"},
			DisableSshMonitorFieldName:               {Type: schema.TypeBool, Optional: true, Computed: false, Description: "indicates whether SSH monitoring is disabled"},
			TokenInjectionFieldName:                  {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether the token injection is enabled"},
			SkipReadyCheckFieldName:                  {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether the cluster installation should be async"},
			RollbackOnControllerFailureFieldName:     {Type: schema.TypeBool, Optional: true, Default: true, Description: "delete cluster on controller installation failure. default = true"},
			ExternalCAFieldName:                      {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether to use external CA for this cluster"},
			ForceRemoveVaultOnDeleteFieldName:        {Type: schema.TypeBool, Optional: true, Default: false, Description: "delete the vault namespace (that was created for token injection) on delete. default = false"},
			InstallTracingSupportFieldName:           {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether to install tracing support, enable for apiSecurity accounts"},
			InstallEnvoyTracingSupportFieldName:      {Type: schema.TypeBool, Optional: true, Default: false, Description: "Indicates whether to install Envoy tracing support, available when install tracing support is true"},
			InternalRegistryFieldName: {
				Description: "Use an internal container registry for this cluster",
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						InternalRegistryFieldNameUrl: {
							Description: "The InternalRegistryFieldNameUrl of the internal registry",
							Type:        schema.TypeString,
							Optional:    true,
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
						SidecarResourcesFieldNameProxyInitLimitsCpu: {
							Type:     schema.TypeString,
							Optional: true,
						},
						SidecarResourcesFieldNameProxyInitLimitsMemory: {
							Type:     schema.TypeString,
							Optional: true,
						},
						SidecarResourcesFieldNameProxyInitRequestsCpu: {
							Type:     schema.TypeString,
							Optional: true,
						},
						SidecarResourcesFieldNameProxyInitRequestsMemory: {
							Type:     schema.TypeString,
							Optional: true,
						},
						SidecarResourcesFieldNameProxyLimitsCpu: {
							Type:     schema.TypeString,
							Optional: true,
						},
						SidecarResourcesFieldNameProxyLimitsMemory: {
							Type:     schema.TypeString,
							Optional: true,
						},
						SidecarResourcesFieldNameProxyRequestsCpu: {
							Type:     schema.TypeString,
							Optional: true,
						},
						SidecarResourcesFieldNameProxyRequestsMemory: {
							Type:     schema.TypeString,
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
			CiImageSignatureValidationFieldName:   {Type: schema.TypeBool, Optional: true, Default: false, Description: "indicates whether ci image signer validation is Enabled"},
			SupportExternalTraceSourceFieldName:   {Type: schema.TypeBool, Optional: true, Default: false, Description: "indicates whether external trace sources are supported, available when install tracing support is true"},
			AutoUpgradeControllerVersionFieldName: {Type: schema.TypeBool, Optional: true, Default: false, Description: "indicates whether upgrade the controller automatically"},
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
	rollbackOnFailure := d.Get(RollbackOnControllerFailureFieldName).(bool)
	tracingEnabled := d.Get(InstallTracingSupportFieldName).(bool)
	forceRemoveVault := d.Get(ForceRemoveVaultOnDeleteFieldName).(bool)

	err = installAgentWithTimeout(ctx, serviceApi, httpClientWrapper, clusterId, k8sContext, multiClusterFolder, tracingEnabled, tokenInjection, skipReadyCheck)
	if err != nil {
		log.Println("[ERROR] Panoptica controller installation has failed")
		if rollbackOnFailure {
			rollBackOnAgentInstallationFailure(ctx, serviceApi, httpClientWrapper, clusterId, k8sContext, forceRemoveVault)
			return diag.FromErr(err)
		} else {
			log.Println("[ERROR] error while installing Panoptica controller. " +
				"environment remains up for debug according to 'rollback_on_controller_failure' field")
		}
	}

	d.SetId(string(clusterId))
	return resourceClusterRead(ctx, d, m)
}

func rollBackOnAgentInstallationFailure(ctx context.Context, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper, clusterId strfmt.UUID, k8sContext string, forceRemoveVault bool) {
	deleteClusterError := serviceApi.DeleteKubernetesCluster(ctx, httpClientWrapper.HttpClient, clusterId)
	if deleteClusterError != nil {
		log.Println("[WARN] failed to remove cluster from Panoptica:")
		log.Println(deleteClusterError)
	}

	_ = printPortshiftNamespaceBeforeDeletingController(k8sContext)
	deleteAgentError := deleteAgent(k8sContext, forceRemoveVault, ctx, serviceApi, httpClientWrapper, clusterId)
	if deleteAgentError != nil {
		log.Println("[WARN] failed to uninstall controller: ")
		log.Println(deleteAgentError)
	}
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] reading cluster")

	httpClientWrapper := m.(client.HttpClientWrapper)

	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	clusterId := d.Id()

	configErr := validateConfig(d)
	if configErr != nil {
		return diag.FromErr(configErr)
	}

	secureCNCluster, err := serviceApi.GetKubernetesClusterById(ctx, httpClientWrapper.HttpClient, strfmt.UUID(clusterId))
	if err != nil {
		return diag.FromErr(err)
	}

	if secureCNCluster.Payload.ID == "" {
		// Tell terraform the cluster doesn't exist
		d.SetId("")
	} else {
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

	updatedCluster, err := serviceApi.UpdateKubernetesCluster(ctx, httpClientWrapper.HttpClient, kubernetesClusterFromConfig, strfmt.UUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateAgent(ctx, d, updatedCluster.Payload, serviceApi, httpClientWrapper)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(updatedCluster.Payload.ID))
	updateMutableFields(d, updatedCluster.Payload)
	return nil
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[DEBUG] deleting cluster")

	httpClientWrapper := m.(client.HttpClientWrapper)
	serviceApi := utils2.GetServiceApi(&httpClientWrapper)
	clusterId := strfmt.UUID(d.Id())
	k8sContext := d.Get(KubernetesClusterContextFieldName).(string)
	forceRemoveVault := d.Get(ForceRemoveVaultOnDeleteFieldName).(bool)
	err := deleteAgent(k8sContext, forceRemoveVault, ctx, serviceApi, httpClientWrapper, clusterId)
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

func installAgent(ctx context.Context, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper, clusterId strfmt.UUID, context string, multiClusterFolder string, tracingEnabled bool, tokenInjection bool, skipReadyCheck bool) error {
	log.Print("[DEBUG] installing agent")

	rootPath, _ := syscall.Getwd()
	installationDir, kubeconfig, err := setUpInstallation(ctx, serviceApi, httpClientWrapper, clusterId, context)
	if err != nil {
		return err
	}

	defer clearInstallationDir(rootPath, installationDir)

	if tokenInjection {
		err = utils2.MakeExecutable(vaultCertsGenFilePath)
		if err != nil {
			return err
		}
	}

	if tracingEnabled {
		err = utils2.MakeExecutable(tracingCertsFilePath)
		if err != nil {
			return err
		}
	}

	output, err := utils2.ExecuteScript(scriptFilePath, multiClusterFolder, skipReadyCheck, kubeconfig)
	if err != nil {
		log.Print("[DEBUG] controller installation failed")
		return fmt.Errorf("%s:\n%s", err, output)
	}

	return nil
}

func clearInstallationDir(rootPath string, installationDir string) {
	os.Chdir(rootPath)
	removeDirectory(installationDir)
}

func setUpInstallation(ctx context.Context, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper, clusterId strfmt.UUID, k8sContext string) (string, string, error) {
	installationDir := installationDirPrefix + uuid.New().String()
	if err := os.Mkdir(installationDir, os.ModePerm); err != nil {
		return "", "", err
	}

	os.Chdir(installationDir)

	kubeconfig, err := createTempKubeconfig(k8sContext)
	if err != nil {
		return "", "", err
	}

	err = downloadAndExtractBundle(ctx, serviceApi, httpClientWrapper, clusterId)
	if err != nil {
		return "", "", err
	}

	err = utils2.MakeExecutable("./" + scriptFilePath)
	if err != nil {
		return "", "", err
	}

	return installationDir, kubeconfig, err
}

func installAgentWithTimeout(ctx context.Context, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper, clusterId strfmt.UUID, context string, multiClusterFolder string, tracingEnabled bool, tokenInjection bool, skipReadyCheck bool) error {
	err := make(chan error, 1)
	go func() {
		err <- installAgent(ctx, serviceApi, httpClientWrapper, clusterId, context, multiClusterFolder, tracingEnabled, tokenInjection, skipReadyCheck)
	}()
	select {
	case <-time.After(15 * time.Minute):
		return errors.New("timed out during Panoptica controller installation process")
	case err := <-err:
		return err
	}
}

func removeDirectory(installationDir string) error {
	err := os.RemoveAll(installationDir)
	if err != nil {
		log.Print("[DEBUG] failed to delete " + installationDir)
		return err
	}
	return nil
}

func printPortshiftNamespaceBeforeDeletingController(context string) error {
	kubeconfig, err := createTempKubeconfig(context)
	if err != nil {
		return err
	}

	defer os.Remove(kubeconfig)
	getPodsResult, _ := utils2.ExecBashCommand(fmt.Sprintf(getPortshiftPodsFormat, kubeconfig))
	log.Printf("[DEBUG] get pods result: \n" + getPodsResult)

	describePodsResult, _ := utils2.ExecBashCommand(fmt.Sprintf(describePortshiftPodsFormat, kubeconfig))
	log.Printf("[DEBUG] describe pods result: \n" + describePodsResult)

	return nil
}

func deleteAgent(k8sContext string, removeVault bool, ctx context.Context, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper, clusterId strfmt.UUID) error {
	log.Printf("[DEBUG] deleting agent from k8sContext: " + k8sContext)

	rootPath, _ := syscall.Getwd()
	installationDir, kubeconfig, err := setUpInstallation(ctx, serviceApi, httpClientWrapper, clusterId, k8sContext)
	if err != nil {
		return err
	}
	defer clearInstallationDir(rootPath, installationDir)

	output, err := utils2.ExecBashCommand(fmt.Sprintf("KUBECONFIG=%s %s", kubeconfig, uninstallCmd))
	log.Printf("[INFO] " + output)
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

	err := downloadInstallBundle(ctx, serviceApi, httpClientWrapper.HttpClient, clusterId, secureCNBundleFilePath)
	if err != nil {
		return err
	}
	open, err := os.Open(secureCNBundleFilePath)
	if err != nil {
		return err
	}

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
	restrictRegistries := d.Get(RestrictRegistriesFieldName).(bool)
	connectionsControl := d.Get(ConnectionsControlFieldName).(bool)
	kubernetesSecurity := d.Get(KubernetesSecurityFieldName).(bool)
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
	ciImageSignatureValidation := d.Get(CiImageSignatureValidationFieldName).(bool)
	supportExternalTraceSource := d.Get(SupportExternalTraceSourceFieldName).(bool)
	persistentStorage := d.Get(PersistentStorageFieldName).(bool)
	autoUpgradeControllerVersion := d.Get(AutoUpgradeControllerVersionFieldName).(bool)
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
	enableK8sEvents := d.Get(EnableK8sEventsFieldName).(bool)
	disableSshMonitor := d.Get(DisableSshMonitorFieldName).(bool)
	enableTokenInjection := d.Get(TokenInjectionFieldName).(bool)
	installTracingSupport := d.Get(InstallTracingSupportFieldName).(bool)
	installEnvoyTracingSupport := d.Get(InstallEnvoyTracingSupportFieldName).(bool)
	externalCA := d.Get(ExternalCAFieldName).(bool)

	cluster := &model.KubernetesCluster{
		AgentFailClose:                    &failClose,
		APIIntelligenceDAST:               &enableAPIIntelligenceDAST,
		AutoLabelEnabled:                  &enableAutoLabel,
		CiImageValidation:                 &ciImageValidation,
		ClusterPodDefinitionSource:        clusterPodDefinitionSource,
		EnableConnectionsControl:          &connectionsControl,
		KubernetesSecurity:                &kubernetesSecurity,
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
		K8sEventsEnabled:                  &enableK8sEvents,
		SshMonitorDisabled:                &disableSshMonitor,
		TokenInjectionEnabled:             &enableTokenInjection,
		MinimalNumberOfControllerReplicas: minimumReplicas,
		CiImageSignatureValidation:        &ciImageSignatureValidation,
		InstallTracingSupport:             &installTracingSupport,
		InstallEnvoyTracingSupport:        &installEnvoyTracingSupport,
		SupportExternalTraceSource:        &supportExternalTraceSource,
		ExternalCa:                        &externalCA,
		AutoUpgradeControllerVersion:      &autoUpgradeControllerVersion,
	}

	internalRegistryUrl := utils2.ReadNestedStringFromTF(d, InternalRegistryFieldName, InternalRegistryFieldNameUrl, 0)
	if internalRegistryUrl != "" {
		internalRegistryEnabled := true
		cluster.InternalRegistryParameters = &model.InternalRegistryParameters{
			InternalRegistryEnabled: &internalRegistryEnabled,
			InternalRegistry:        internalRegistryUrl,
		}
	}

	proxyInitLimitsCpu := utils2.ReadNestedStringFromTF(d, SidecarResourcesFieldName, "proxy_init_limits_cpu", 0)
	proxyInitLimitsMemory := utils2.ReadNestedStringFromTF(d, SidecarResourcesFieldName, "proxy_init_limits_memory", 0)
	proxyInitRequestsCpu := utils2.ReadNestedStringFromTF(d, SidecarResourcesFieldName, "proxy_init_requests_cpu", 0)
	proxyInitRequestsMemory := utils2.ReadNestedStringFromTF(d, SidecarResourcesFieldName, "proxy_init_requests_memory", 0)
	proxyLimitsCpu := utils2.ReadNestedStringFromTF(d, SidecarResourcesFieldName, "proxy_limits_cpu", 0)
	proxyLimitsMemory := utils2.ReadNestedStringFromTF(d, SidecarResourcesFieldName, "proxy_limits_memory", 0)
	proxyRequestsCpu := utils2.ReadNestedStringFromTF(d, SidecarResourcesFieldName, "proxy_requests_cpu", 0)
	proxyRequestsMemory := utils2.ReadNestedStringFromTF(d, SidecarResourcesFieldName, "proxy_requests_memory", 0)
	cluster.SidecarsResources = &model.SidecarsResource{
		ProxyInitLimitsCPU:      proxyInitLimitsCpu,
		ProxyInitLimitsMemory:   proxyInitLimitsMemory,
		ProxyInitRequestsCPU:    proxyInitRequestsCpu,
		ProxyInitRequestsMemory: proxyInitRequestsMemory,
		ProxyLimitsCPU:          proxyLimitsCpu,
		ProxyLimitsMemory:       proxyLimitsMemory,
		ProxyRequestCPU:         proxyRequestsCpu,
		ProxyRequestMemory:      proxyRequestsMemory,
	}

	return cluster, nil
}

func downloadInstallBundle(ctx context.Context, serviceApi *escherClient.MgmtServiceApiCtx, client *http.Client, clusterId strfmt.UUID, bundlePath string) error {
	log.Print("[DEBUG] downloading file")

	file, err := os.Create(bundlePath)
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
	_ = d.Set(KubernetesSecurityFieldName, secureCNCluster.KubernetesSecurity)
	if secureCNCluster.IstioInstallationParameters == nil {
		_ = d.Set(IstioAlreadyInstalledFieldName, nil)
		_ = d.Set(IstioVersionFieldName, nil)
	} else {
		_ = d.Set(IstioAlreadyInstalledFieldName, secureCNCluster.IstioInstallationParameters.IsIstioAlreadyInstalled)
		_ = d.Set(IstioVersionFieldName, secureCNCluster.IstioInstallationParameters.IstioVersion)
	}
	_ = d.Set(MultiClusterCommunicationSupportFieldName, secureCNCluster.IsMultiCluster)
	_ = d.Set(InspectIncomingClusterConnectionsFieldName, secureCNCluster.PreserveOriginalSourceIP)
	_ = d.Set(FailCloseFieldName, secureCNCluster.AgentFailClose)
	_ = d.Set(CiImageSignatureValidationFieldName, secureCNCluster.CiImageSignatureValidation)
	_ = d.Set(SupportExternalTraceSourceFieldName, secureCNCluster.SupportExternalTraceSource)
	_ = d.Set(PersistentStorageFieldName, secureCNCluster.IsPersistent)
	if secureCNCluster.ProxyConfiguration == nil {
		_ = d.Set(ExternalHttpsProxyFieldName, nil)
	} else {
		_ = d.Set(ExternalHttpsProxyFieldName, secureCNCluster.ProxyConfiguration.HTTPSProxy)
	}
	_ = d.Set(OrchestrationTypeFieldName, secureCNCluster.OrchestrationType)
	_ = d.Set(TLSInspectionFieldName, secureCNCluster.TLSInspectionEnabled)
	_ = d.Set(EnableK8sEventsFieldName, secureCNCluster.K8sEventsEnabled)
	_ = d.Set(DisableSshMonitorFieldName, secureCNCluster.SshMonitorDisabled)
	_ = d.Set(TokenInjectionFieldName, secureCNCluster.TokenInjectionEnabled)
	_ = d.Set(ServiceDiscoveryIsolationFieldName, secureCNCluster.ServiceDiscoveryIsolationEnabled)
	_ = d.Set(RestrictRegistriesFieldName, secureCNCluster.RestrictRegistires)
	_ = d.Set(IstioIngressEnabledFieldName, secureCNCluster.IsIstioIngressEnabled)
	_ = d.Set(IstioIngressAnnotationsFieldName, getIstioAnnotationsMap(secureCNCluster))
	_ = d.Set(EnableApiIntelligenceDASTFieldName, secureCNCluster.APIIntelligenceDAST)
	_ = d.Set(EnableAutoLabelFieldName, secureCNCluster.AutoLabelEnabled)
	_ = d.Set(HoldApplicationUntilProxyStartsFieldName, secureCNCluster.IsHoldApplicationUntilProxyStarts)
	_ = d.Set(InstallTracingSupportFieldName, secureCNCluster.InstallTracingSupport)
	_ = d.Set(InstallEnvoyTracingSupportFieldName, secureCNCluster.InstallEnvoyTracingSupport)
	_ = d.Set(MinimumReplicasFieldName, secureCNCluster.MinimalNumberOfControllerReplicas)
	_ = d.Set(ExternalCAFieldName, secureCNCluster.ExternalCa)
	_ = d.Set(AutoUpgradeControllerVersionFieldName, secureCNCluster.AutoUpgradeControllerVersion)

	if secureCNCluster.InternalRegistryParameters == nil {
		_ = d.Set(InternalRegistryFieldName, nil)
	} else {
		_ = d.Set(InternalRegistryFieldName, utils2.GetTfMapFromKeyValuePairs([]utils2.KeyValue{{
			InternalRegistryFieldNameUrl, secureCNCluster.InternalRegistryParameters.InternalRegistry}}))
	}

	if secureCNCluster.SidecarsResources == nil {
		_ = d.Set(SidecarResourcesFieldName, nil)
	} else {
		_ = d.Set(SidecarResourcesFieldName, utils2.GetTfMapFromKeyValuePairs([]utils2.KeyValue{
			{SidecarResourcesFieldNameProxyInitLimitsCpu, secureCNCluster.SidecarsResources.ProxyInitLimitsCPU},
			{SidecarResourcesFieldNameProxyInitLimitsMemory, secureCNCluster.SidecarsResources.ProxyInitLimitsMemory},
			{SidecarResourcesFieldNameProxyInitRequestsCpu, secureCNCluster.SidecarsResources.ProxyInitRequestsCPU},
			{SidecarResourcesFieldNameProxyInitRequestsMemory, secureCNCluster.SidecarsResources.ProxyInitRequestsMemory},
			{SidecarResourcesFieldNameProxyLimitsCpu, secureCNCluster.SidecarsResources.ProxyLimitsCPU},
			{SidecarResourcesFieldNameProxyLimitsMemory, secureCNCluster.SidecarsResources.ProxyLimitsMemory},
			{SidecarResourcesFieldNameProxyRequestsCpu, secureCNCluster.SidecarsResources.ProxyRequestCPU},
			{SidecarResourcesFieldNameProxyRequestsMemory, secureCNCluster.SidecarsResources.ProxyRequestMemory}}))
	}
}

func getIstioAnnotationsMap(secureCNCluster *model.KubernetesCluster) map[string]string {
	annotationsInSecureCN := secureCNCluster.IstioIngressAnnotations
	annotations := make(map[string]string, len(annotationsInSecureCN))
	for _, annotation := range annotationsInSecureCN {
		keyInSecureCn := annotation.Key
		valueInSecureCn := annotation.Value
		annotations[*keyInSecureCn] = *valueInSecureCn
	}

	return annotations
}

func validateConfig(d *schema.ResourceData) error {
	log.Printf("[DEBUG] validating config")

	isMultiCluster := d.Get(MultiClusterCommunicationSupportFieldName).(bool)
	multiClusterFolder := d.Get(MultiClusterCommunicationSupportCertsPathFieldName).(string)
	connectionsControl := d.Get(ConnectionsControlFieldName).(bool)
	kubernetesSecurity := d.Get(KubernetesSecurityFieldName).(bool)
	inspectIncomingClusterConnections := d.Get(InspectIncomingClusterConnectionsFieldName).(bool)
	installTraceSupport := d.Get(InstallTracingSupportFieldName).(bool)
	installEnvoyTraceSupport := d.Get(InstallEnvoyTracingSupportFieldName).(bool)
	supportExternalTraceSource := d.Get(SupportExternalTraceSourceFieldName).(bool)
	if isMultiCluster && multiClusterFolder == "" {
		return errors.New(fmt.Sprintf("invalid configuration. %s can't be empty when %s is true", MultiClusterCommunicationSupportCertsPathFieldName, MultiClusterCommunicationSupportFieldName))
	}

	if !connectionsControl && isMultiCluster {
		return errors.New(fmt.Sprintf("invalid configuration. %s is off but %s is on ", ConnectionsControlFieldName, MultiClusterCommunicationSupportFieldName))
	}

	if !connectionsControl && inspectIncomingClusterConnections {
		return errors.New(fmt.Sprintf("invalid configuration. %s is off but %s is on", MultiClusterCommunicationSupportCertsPathFieldName, InspectIncomingClusterConnectionsFieldName))
	}

	if !kubernetesSecurity && connectionsControl {
		return errors.New(fmt.Sprintf("invalid configuration. %s is off but %s is on", KubernetesSecurityFieldName, ConnectionsControlFieldName))
	}

	if !installTraceSupport && supportExternalTraceSource {
		return errors.New(fmt.Sprintf("invalid cluster api security config. %s can't be turned on when %s is off", SupportExternalTraceSourceFieldName, InstallTracingSupportFieldName))
	}

	if !installTraceSupport && installEnvoyTraceSupport {
		return errors.New(fmt.Sprintf("invalid cluster api security config. %s can't be turned on when %s is off", InstallEnvoyTracingSupportFieldName, InstallTracingSupportFieldName))
	}

	return nil
}

func updateAgent(ctx context.Context, d *schema.ResourceData, updatedCluster *model.KubernetesCluster, serviceApi *escherClient.MgmtServiceApiCtx, httpClientWrapper client.HttpClientWrapper) error {
	if updatedCluster.ControllerStatus == model.ControllerStatusWAITINGFORUSERUPDATE {
		log.Print("[DEBUG] updating agent")
		context := d.Get(KubernetesClusterContextFieldName).(string)
		forceRemoveVault := d.Get(ForceRemoveVaultOnDeleteFieldName).(bool)
		err := deleteAgent(context, forceRemoveVault, ctx, serviceApi, httpClientWrapper, updatedCluster.ID)
		if err != nil {
			return err
		}
		err = installAgentWithTimeout(ctx, serviceApi, httpClientWrapper, updatedCluster.ID, context, d.Get(MultiClusterCommunicationSupportCertsPathFieldName).(string), d.Get(InstallTracingSupportFieldName).(bool), d.Get(TokenInjectionFieldName).(bool), d.Get(SkipReadyCheckFieldName).(bool))
		if err != nil {
			return err
		}
	}

	return nil
}
