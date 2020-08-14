package k8s

import (
	"context"

	"github.com/Juniper/contrail-operator/pkg/label"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Service struct {
	name      string
	servType  core.ServiceType
	port      int32
	ownerType string
	owner     v1.Object
	scheme    *runtime.Scheme
	client    client.Client
	svc       core.Service
}

func (s *Service) EnsureExists() error {
	labels := label.New(s.ownerType, s.owner.GetName())
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
		s.svc.Spec.Type = core.ServiceTypeLoadBalancer
		return controllerutil.SetControllerReference(s.owner, &s.svc, s.scheme)
	})
	return err
}

func (s *Service) ClusterIP() string {
	return s.svc.Spec.ClusterIP
}
