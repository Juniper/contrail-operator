package v1alpha1

import appsv1 "k8s.io/api/apps/v1"

// Status is the status of the service.
// +k8s:openapi-gen=true
type Status struct {
	Active        bool  `json:"active,omitempty"`
	Replicas      int32 `json:"replicas,omitempty"`
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
}

func (s *Status) FromDeployment(d *appsv1.Deployment) {
	expectedReplicas := int32(1)
	if d.Spec.Replicas != nil {
		expectedReplicas = *d.Spec.Replicas
	}
	s.Replicas = expectedReplicas
	s.ReadyReplicas = d.Status.ReadyReplicas
	if d.Status.ReadyReplicas == expectedReplicas {
		s.Active = true
	} else {
		s.Active = false
	}
}

func (s *Status) FromStatefulSet(d *appsv1.StatefulSet) {
	expectedReplicas := int32(1)
	if d.Spec.Replicas != nil {
		expectedReplicas = *d.Spec.Replicas
	}
	s.Replicas = expectedReplicas
	s.ReadyReplicas = d.Status.ReadyReplicas
	if d.Status.ReadyReplicas == expectedReplicas {
		s.Active = true
	} else {
		s.Active = false
	}
}
