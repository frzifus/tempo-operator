package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cfg "sigs.k8s.io/controller-runtime/pkg/config/v1alpha1"
)

// ImagesSpec defines the image for each container.
type ImagesSpec struct {
	// Tempo defines the tempo container image.
	//
	// +optional
	Tempo string `json:"tempo,omitempty"`

	// TempoQuery defines the tempo-query container image.
	//
	// +optional
	TempoQuery string `json:"tempoQuery,omitempty"`

	// TempoGateway defines the tempo-gateway container image.
	//
	// +optional
	TempoGateway string `json:"tempoGateway,omitempty"`

	// TempoGatewayOpa defines the OPA sidecar container for TempoGateway.
	//
	// +optional
	TempoGatewayOpa string `json:"tempoGatewayOpa,omitempty"`
}

// BuiltInCertManagement is the configuration for the built-in facility to generate and rotate
// TLS client and serving certificates for all Tempo services and internal clients. All necessary
// secrets and configmaps for protecting the internal components will be created if this option is enabled.
type BuiltInCertManagement struct {
	// CACertValidity defines the total duration of the CA certificate validity.
	CACertValidity metav1.Duration `json:"caValidity,omitempty"`
	// CACertRefresh defines the duration of the CA certificate validity until a rotation
	// should happen. It can be set up to 80% of CA certificate validity or equal to the
	// CA certificate validity. Latter should be used only for rotating only when expired.
	CACertRefresh metav1.Duration `json:"caRefresh,omitempty"`
	// CertValidity defines the total duration of the validity for all Tempo certificates.
	CertValidity metav1.Duration `json:"certValidity,omitempty"`
	// CertRefresh defines the duration of the certificate validity until a rotation
	// should happen. It can be set up to 80% of certificate validity or equal to the
	// certificate validity. Latter should be used only for rotating only when expired.
	// The refresh is applied to all Tempo certificates at once.
	CertRefresh metav1.Duration `json:"certRefresh,omitempty"`
	// Enabled defines to flag to enable/disable built-in certificate management feature gate.
	Enabled bool `json:"enabled,omitempty"`
}

// OpenShiftFeatureGates is the supported set of all operator features gates on OpenShift.
type OpenShiftFeatureGates struct {
	// ServingCertsService enables OpenShift service-ca annotations on the TempoStack gateway service only
	// to use the in-platform CA and generate a TLS cert/key pair per service for
	// in-cluster data-in-transit encryption.
	// More details: https://docs.openshift.com/container-platform/latest/security/certificate_types_descriptions/service-ca-certificates.html
	ServingCertsService bool `json:"servingCertsService,omitempty"`

	// OpenShiftRoute enables creating OpenShift Route objects.
	// More details: https://docs.openshift.com/container-platform/latest/networking/understanding-networking.html
	OpenShiftRoute bool `json:"openshiftRoute,omitempty"`

	// BaseDomain is used internally for redirect URL in gateway OpenShift auth mode.
	// If empty the operator automatically derives the domain from the cluster.
	BaseDomain string `json:"baseDomain,omitempty"`

	// ClusterTLSPolicy enables usage of TLS policies set in the API Server.
	// More details: https://docs.openshift.com/container-platform/4.11/security/tls-security-profiles.html
	ClusterTLSPolicy bool
}

// TLSProfileType is a TLS security profile based on the Mozilla definitions:
// https://wiki.mozilla.org/Security/Server_Side_TLS
type TLSProfileType string

const (
	// TLSProfileOldType is a TLS security profile based on:
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Old_backward_compatibility
	TLSProfileOldType TLSProfileType = "Old"
	// TLSProfileIntermediateType is a TLS security profile based on:
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Intermediate_compatibility_.28default.29
	TLSProfileIntermediateType TLSProfileType = "Intermediate"
	// TLSProfileModernType is a TLS security profile based on:
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
	TLSProfileModernType TLSProfileType = "Modern"
)

// MetricsFeatureGates configures metrics and alerts of the operator.
type MetricsFeatureGates struct {
	// CreateServiceMonitors defines whether the operator should install ServiceMonitors
	// to scrape metrics of the operator.
	CreateServiceMonitors bool `json:"createServiceMonitors,omitempty"`

	// CreatePrometheusRules defines whether the operator should install PrometheusRules
	// to receive alerts about the operator.
	CreatePrometheusRules bool `json:"createPrometheusRules,omitempty"`
}

// ObservabilityFeatureGates configures observability of the operator.
type ObservabilityFeatureGates struct {
	// Metrics configures metrics of the operator.
	Metrics MetricsFeatureGates `json:"metrics,omitempty"`
}

// FeatureGates is the supported set of all operator feature gates.
type FeatureGates struct {
	// OpenShift contains a set of feature gates supported only on OpenShift.
	OpenShift OpenShiftFeatureGates `json:"openshift,omitempty"`

	// BuiltInCertManagement enables the built-in facility for generating and rotating
	// TLS client and serving certificates for the communication between ingesters and distributors and also between
	// query and query-frontend, In detail all internal Tempo HTTP and GRPC communication is lifted
	// to require mTLS.
	// In addition each service requires a configmap named as the MicroService CR with the
	// suffix `-ca-bundle`, e.g. `tempo-dev-ca-bundle` and the following data:
	// - `service-ca.crt`: The CA signing the service certificate in `tls.crt`.
	// All necessary secrets and configmaps for protecting the internal components will be created if this
	// option is enabled.
	BuiltInCertManagement BuiltInCertManagement `json:"builtInCertManagement,omitempty"`
	// HTTPEncryption enables TLS encryption for all HTTP TempoStack components.
	// Each HTTP component requires a secret, the name should be the name of the component with the
	// suffix `-mtls` and prefix by the TempoStack name e.g `tempo-dev-distributor-mtls`.
	// It should contains the following data:
	// - `tls.crt`: The TLS server side certificate.
	// - `tls.key`: The TLS key for server-side encryption.
	// In addition each service requires a configmap named as the TempoStack CR with the
	// suffix `-ca-bundle`, e.g. `tempo-dev-ca-bundle` and the following data:
	// - `service-ca.crt`: The CA signing the service certificate in `tls.crt`.
	// This will protect all internal communication between the distributors and ingestors and also
	// between ingestor and queriers, and between the queriers and the query-frontend component
	//
	// If BuiltInCertManagement is enabled, you don't need to create this secrets manually.
	//
	// Some considerations when enable mTLS:
	// - If JaegerUI is enabled, it won't be protected by mTLS as it will be considered a public facing
	// component.
	// - If JaegerUI is not enabled, HTTP Tempo API won´t be protected, this will be considered
	// public faced component.
	// - If Gateway is enabled, all comunications between the gateway and the tempo components will be protected
	// by mTLS, and the Gateway itself won´t be, as it will be the only public face component.
	HTTPEncryption bool `json:"httpEncryption,omitempty"`
	// GRPCEncryption enables TLS encryption for all GRPC TempoStack services.
	// Each GRPC component requires a secret, the name should be the name of the component with the
	// suffix `-mtls` and prefix by the TempoStack name e.g `tempo-dev-distributor-mtls`.
	// It should contains the following data:
	// - `tls.crt`: The TLS server side certificate.
	// - `tls.key`: The TLS key for server-side encryption.
	// In addition each service requires a configmap named as the TempoStack CR with the
	// suffix `-ca-bundle`, e.g. `tempo-dev-ca-bundle` and the following data:
	// - `service-ca.crt`: The CA signing the service certificate in `tls.crt`.
	// This will protect all internal communication between the distributors and ingestors and also
	// between ingestor and queriers, and between the queriers and the query-frontend component.
	//
	//
	// If BuiltInCertManagement is enabled, you don't need to create this secrets manually.
	//
	// Some considerations when enable mTLS:
	// - If JaegerUI is enabled, it won´t be protected by mTLS as it will be considered a public face
	// component.
	// - If Gateway is enabled, all comunications between the gateway and the tempo components will be protected
	// by mTLS, and the Gateway itself won´t be, as it will be the only public face component.
	GRPCEncryption bool `json:"grpcEncryption,omitempty"`

	// TLSProfile allows to chose a TLS security profile. Enforced
	// when using HTTPEncryption or GRPCEncryption.
	TLSProfile string `json:"tlsProfile,omitempty"`

	// PrometheusOperator defines whether the Prometheus Operator CRD exists in the cluster.
	// This CRD is part of prometheus-operator.
	PrometheusOperator bool `json:"prometheusOperator,omitempty"`

	// Observability configures observability features of the operator.
	Observability ObservabilityFeatureGates `json:"observability,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ProjectConfig is the Schema for the projectconfigs API.
type ProjectConfig struct {
	metav1.TypeMeta `json:",inline"`
	// ControllerManagerConfigurationSpec returns the configurations for controllers
	cfg.ControllerManagerConfigurationSpec `json:",inline"`

	DefaultImages ImagesSpec `json:"images"`

	Gates FeatureGates `json:"featureGates,omitempty"`

	// Distribution defines the operator distribution name.
	Distribution string `json:"distribution"`
}

func init() {
	SchemeBuilder.Register(&ProjectConfig{})
}
