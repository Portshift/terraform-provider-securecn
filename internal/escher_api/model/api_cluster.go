// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// KubernetesCluster kubernetes cluster
// swagger:model KubernetesCluster
type KubernetesCluster struct {

	// indicates fail close behavior on SecureCn agent failure
	AgentFailClose *bool `json:"agentFailClose,omitempty"`

	// indicates if the controller is updated. when false, reinstall is needed
	AutoUpdateEnabled *bool `json:"autoUpdateEnabled,omitempty"`

	// indicates whether apiIntelligenceDAST is enabled
	APIIntelligenceDAST *bool `json:"apiIntelligenceDAST,omitempty"`

	// indicates whether auto label is enabled
	AutoLabelEnabled *bool `json:"autoLabelEnabled,omitempty"`

	// Enable pod template images validation
	CiImageValidation *bool `json:"ciImageValidation,omitempty"`

	// cluster pod definition source
	// Required: true
	ClusterPodDefinitionSource ClusterPodDefinitionSource `json:"clusterPodDefinitionSource"`

	// indicates whether SecureCn allows connections actions and detections
	// Required: true
	EnableConnectionsControl *bool `json:"enableConnectionsControl"`

	// external ca
	ExternalCa *ExternalCaDetails `json:"externalCa,omitempty"`

	// Id of the cluster.
	// Read Only: true
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// internal registry parameters
	InternalRegistryParameters *InternalRegistryParameters `json:"internalRegistryParameters,omitempty"`

	// indicates whether the controller should hold the application until the proxy starts
	IsHoldApplicationUntilProxyStarts *bool `json:"isHoldApplicationUntilProxyStarts,omitempty"`

	// indicates whether Istio ingress is enabled
	IsIstioIngressEnabled *bool `json:"isIstioIngressEnabled,omitempty"`

	// indicates whether this cluster should support multi-cluster communication
	IsMultiCluster *bool `json:"isMultiCluster,omitempty"`

	// indicates whether the agent should run in persistent mode
	IsPersistent *bool `json:"isPersistent,omitempty"`

	// annotations for load balancers
	IstioIngressAnnotations []*KubernetesAnnotation `json:"istioIngressAnnotations"`

	// istio installation parameters
	IstioInstallationParameters *IstioInstallationParameters `json:"istioInstallationParameters,omitempty"`

	// name
	// Required: true
	// Min Length: 1
	Name *string `json:"name"`

	// orchestration type
	// Required: true
	// Enum: [GKE OPENSHIFT RANCHER AKS EKS KUBERNETES IKS]
	OrchestrationType *string `json:"orchestrationType"`

	// indicates whether the agent should preserve the original source ip
	PreserveOriginalSourceIP *bool `json:"preserveOriginalSourceIp,omitempty"`

	// indicates whether this cluster should use a proxy server
	ProxyConfiguration *ProxyConfiguration `json:"proxyConfiguration,omitempty"`

	// indicates whether the agent validate the images origin
	RestrictRegistires *bool `json:"restrictRegistires,omitempty"`

	// indicates whether the service discovery isolation is enabled
	ServiceDiscoveryIsolationEnabled *bool `json:"serviceDiscoveryIsolationEnabled,omitempty"`

	// sidecars resources
	SidecarsResources *SidecarsResource `json:"sidecarsResources,omitempty"`

	// indicates whether TLS inspection is enabled
	TLSInspectionEnabled *bool `json:"tlsInspectionEnabled,omitempty"`

	// indicates whether token injection is enabled
	TokenInjectionEnabled *bool `json:"tokenInjectionEnabled,omitempty"`

	// tracing support configuration. enabled for ApiSecurity enabled accounts
	InstallTracingSupport *bool `json:"installTracingSupport,omitempty"`

	InstallEnvoyTracingSupport *bool `json:"installEnvoyTracingSupport,omitempty"`

	// minimum number of controller replicas"
	MinimalNumberOfControllerReplicas int `json:"minimalNumberOfControllerReplicas,omitempty"`

	// indicates whether ci image signer validation is Enabled
	CiImageSignatureValidation *bool `json:"ciImageSignatureValidation,omitempty"`

	SupportExternalTraceSource *bool `json:"supportExternalTraceSource,omitempty"`
}

// Validate validates this kubernetes cluster
func (m *KubernetesCluster) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateClusterPodDefinitionSource(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEnableConnectionsControl(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExternalCa(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInternalRegistryParameters(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIstioIngressAnnotations(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIstioInstallationParameters(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOrchestrationType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProxyConfiguration(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSidecarsResources(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *KubernetesCluster) validateClusterPodDefinitionSource(formats strfmt.Registry) error {

	if err := m.ClusterPodDefinitionSource.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("clusterPodDefinitionSource")
		}
		return err
	}

	return nil
}

func (m *KubernetesCluster) validateEnableConnectionsControl(formats strfmt.Registry) error {

	if err := validate.Required("enableConnectionsControl", "body", m.EnableConnectionsControl); err != nil {
		return err
	}

	return nil
}

func (m *KubernetesCluster) validateExternalCa(formats strfmt.Registry) error {

	if swag.IsZero(m.ExternalCa) { // not required
		return nil
	}

	if m.ExternalCa != nil {
		if err := m.ExternalCa.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("externalCa")
			}
			return err
		}
	}

	return nil
}

func (m *KubernetesCluster) validateID(formats strfmt.Registry) error {

	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *KubernetesCluster) validateInternalRegistryParameters(formats strfmt.Registry) error {

	if swag.IsZero(m.InternalRegistryParameters) { // not required
		return nil
	}

	if m.InternalRegistryParameters != nil {
		if err := m.InternalRegistryParameters.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("internalRegistryParameters")
			}
			return err
		}
	}

	return nil
}

func (m *KubernetesCluster) validateIstioIngressAnnotations(formats strfmt.Registry) error {

	if swag.IsZero(m.IstioIngressAnnotations) { // not required
		return nil
	}

	for i := 0; i < len(m.IstioIngressAnnotations); i++ {
		if swag.IsZero(m.IstioIngressAnnotations[i]) { // not required
			continue
		}

		if m.IstioIngressAnnotations[i] != nil {
			if err := m.IstioIngressAnnotations[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("istioIngressAnnotations" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *KubernetesCluster) validateIstioInstallationParameters(formats strfmt.Registry) error {

	if swag.IsZero(m.IstioInstallationParameters) { // not required
		return nil
	}

	if m.IstioInstallationParameters != nil {
		if err := m.IstioInstallationParameters.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("istioInstallationParameters")
			}
			return err
		}
	}

	return nil
}

func (m *KubernetesCluster) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	if err := validate.MinLength("name", "body", string(*m.Name), 1); err != nil {
		return err
	}

	return nil
}

var kubernetesClusterTypeOrchestrationTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["GKE","OPENSHIFT","RANCHER","AKS","EKS","KUBERNETES","IKS"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		kubernetesClusterTypeOrchestrationTypePropEnum = append(kubernetesClusterTypeOrchestrationTypePropEnum, v)
	}
}

const (

	// KubernetesClusterOrchestrationTypeGKE captures enum value "GKE"
	KubernetesClusterOrchestrationTypeGKE string = "GKE"

	// KubernetesClusterOrchestrationTypeOPENSHIFT captures enum value "OPENSHIFT"
	KubernetesClusterOrchestrationTypeOPENSHIFT string = "OPENSHIFT"

	// KubernetesClusterOrchestrationTypeRANCHER captures enum value "RANCHER"
	KubernetesClusterOrchestrationTypeRANCHER string = "RANCHER"

	// KubernetesClusterOrchestrationTypeAKS captures enum value "AKS"
	KubernetesClusterOrchestrationTypeAKS string = "AKS"

	// KubernetesClusterOrchestrationTypeEKS captures enum value "EKS"
	KubernetesClusterOrchestrationTypeEKS string = "EKS"

	// KubernetesClusterOrchestrationTypeKUBERNETES captures enum value "KUBERNETES"
	KubernetesClusterOrchestrationTypeKUBERNETES string = "KUBERNETES"

	// KubernetesClusterOrchestrationTypeIKS captures enum value "IKS"
	KubernetesClusterOrchestrationTypeIKS string = "IKS"
)

// prop value enum
func (m *KubernetesCluster) validateOrchestrationTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, kubernetesClusterTypeOrchestrationTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *KubernetesCluster) validateOrchestrationType(formats strfmt.Registry) error {

	if err := validate.Required("orchestrationType", "body", m.OrchestrationType); err != nil {
		return err
	}

	// value enum
	if err := m.validateOrchestrationTypeEnum("orchestrationType", "body", *m.OrchestrationType); err != nil {
		return err
	}

	return nil
}

func (m *KubernetesCluster) validateProxyConfiguration(formats strfmt.Registry) error {

	if swag.IsZero(m.ProxyConfiguration) { // not required
		return nil
	}

	if m.ProxyConfiguration != nil {
		if err := m.ProxyConfiguration.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("proxyConfiguration")
			}
			return err
		}
	}

	return nil
}

func (m *KubernetesCluster) validateSidecarsResources(formats strfmt.Registry) error {

	if swag.IsZero(m.SidecarsResources) { // not required
		return nil
	}

	if m.SidecarsResources != nil {
		if err := m.SidecarsResources.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("sidecarsResources")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *KubernetesCluster) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *KubernetesCluster) UnmarshalBinary(b []byte) error {
	var res KubernetesCluster
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
