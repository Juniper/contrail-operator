package v1alpha1

// Service is the configuration of the service that exposes a workload
// +k8s:openapi-gen=true
type Service struct {
	ServiceType serviceType       `json:"serviceType,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// +kubebuilder:validation:Enum={"","ClusterIP","NodePort","LoadBalancer","ExternalName"}
type serviceType string

func (s *Service) GetServiceType() serviceType {
	return s.ServiceType
}

func (s *Service) GetAnnotations() map[string]string {
	return s.Annotations
}
