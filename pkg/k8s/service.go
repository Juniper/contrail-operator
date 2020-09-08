package k8s

import (
	"context"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/Juniper/contrail-operator/pkg/label"
)

// Service is used to create and manage kubernetes services
type Service struct {
	name      string
	servType  core.ServiceType
	port      int32
	ownerType string
	labels    map[string]string
	owner     v1.Object
	scheme    *runtime.Scheme
	client    client.Client
	svc       core.Service
}

// EnsureExists is used to make sure that kubernetes service exists and is correctly configured
func (s *Service) EnsureExists() error {
	labels := s.labels
	if len(labels) == 0 {
		labels = label.New(s.ownerType, s.owner.GetName())
	}
	s.svc = core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      s.name + "-" + s.ownerType,
			Namespace: s.owner.GetNamespace(),
			Labels:    labels,
		},
	}
	_, err := controllerutil.CreateOrUpdate(context.Background(), s.client, &s.svc, func() error {
		nodePort := int32(0)
		for i, p := range s.svc.Spec.Ports {
			if p.Port == s.port {
				nodePort = s.svc.Spec.Ports[i].NodePort
			}
		}

		s.svc.Spec.Ports = []core.ServicePort{
			{Port: s.port, Protocol: "TCP", NodePort: nodePort},
		}
		s.svc.Spec.Selector = labels
		s.svc.Spec.Type = s.servType
		return controllerutil.SetControllerReference(s.owner, &s.svc, s.scheme)
	})
	return err
}

// ClusterIP is used to read clusterIP of service
func (s *Service) ClusterIP() string {
	return s.svc.Spec.ClusterIP
}

// WithLabels is used to set labels on Service
func (s *Service) WithLabels(labels map[string]string) *Service {
	s.labels = labels
	return s
}
