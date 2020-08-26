package k8s_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestEnsureServiceExists(t *testing.T) {
	scheme := runtime.NewScheme()
	err := core.SchemeBuilder.AddToScheme(scheme)
	require.NoError(t, err)
	tests := []struct {
		name                string
		port                int32
		servType            core.ServiceType
		initDBState         []runtime.Object
		expectedServiceSpec core.ServiceSpec
	}{
		{
			name:     "Create cluster IP service",
			servType: core.ServiceTypeClusterIP,
			port:     5555,
			expectedServiceSpec: core.ServiceSpec{
				Type: core.ServiceTypeClusterIP,
				Ports: []core.ServicePort{{
					Name: "", Protocol: "TCP", Port: 5555,
				}},
				Selector: map[string]string{"contrail_manager": "pod", "pod": "owner"},
			},
		},
		{
			name:     "Create Load balancer IP service",
			servType: core.ServiceTypeLoadBalancer,
			port:     5555,
			expectedServiceSpec: core.ServiceSpec{
				Type: core.ServiceTypeLoadBalancer,
				Ports: []core.ServicePort{{
					Name: "", Protocol: "TCP", Port: 5555,
				}},
				Selector: map[string]string{"contrail_manager": "pod", "pod": "owner"},
			},
		},
		{
			name:     "Don't update node port if it is present",
			servType: core.ServiceTypeLoadBalancer,
			port:     5555,
			expectedServiceSpec: core.ServiceSpec{
				Type: core.ServiceTypeLoadBalancer,
				Ports: []core.ServicePort{{
					Name: "", Protocol: "TCP", Port: 5555, NodePort: 30000,
				}},
				Selector: map[string]string{"contrail_manager": "pod", "pod": "owner"},
			},
			initDBState: []runtime.Object{
				&core.Service{
					ObjectMeta: meta.ObjectMeta{Name: "test-pod"},
					Spec:       core.ServiceSpec{Ports: []core.ServicePort{{Port: 5555, NodePort: 30000}}},
				},
			},
		},
		{
			name:     "Update node port if port changes",
			servType: core.ServiceTypeLoadBalancer,
			port:     5500,
			expectedServiceSpec: core.ServiceSpec{
				Type: core.ServiceTypeLoadBalancer,
				Ports: []core.ServicePort{{
					Name: "", Protocol: "TCP", Port: 5500,
				}},
				Selector: map[string]string{"contrail_manager": "pod", "pod": "owner"},
			},
			initDBState: []runtime.Object{
				&core.Service{
					ObjectMeta: meta.ObjectMeta{Name: "test-pod"},
					Spec:       core.ServiceSpec{Ports: []core.ServicePort{{Port: 5555, NodePort: 30000}}},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			owner := &core.Pod{
				ObjectMeta: meta.ObjectMeta{
					Name: "owner",
					UID:  "uid",
				},
			}
			cl := fake.NewFakeClientWithScheme(scheme, test.initDBState...)
			sc := k8s.New(cl, scheme).Service("test", test.servType, test.port, "pod", owner)
			// When
			err := sc.EnsureExists()
			// Then
			assert.NoError(t, err)
			// And
			service := &core.Service{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name: "test-pod",
			}, service)

			assert.NoError(t, err)
			service.SetResourceVersion("")
			service.TypeMeta = meta.TypeMeta{}
			require.NotNil(t, service)
			trueVal := true
			expectedOwnerRef := []meta.OwnerReference{
				{"v1", "Pod", "owner", "uid", &trueVal, &trueVal},
			}
			assert.Equal(t, expectedOwnerRef, service.OwnerReferences)
			assert.Equal(t, test.expectedServiceSpec, service.Spec)
		})
	}
}

func TestClusterIP(t *testing.T) {
	scheme := runtime.NewScheme()
	err := core.SchemeBuilder.AddToScheme(scheme)
	require.NoError(t, err)
	owner := &core.Pod{ObjectMeta: meta.ObjectMeta{Name: "owner"}}

	t.Run("Get ClusterIP", func(t *testing.T) {
		initService := &core.Service{
			ObjectMeta: meta.ObjectMeta{Name: "test-pod"},
			Spec:       core.ServiceSpec{ClusterIP: "10.0.0.10"},
		}
		cl := fake.NewFakeClientWithScheme(scheme, initService)
		sc := k8s.New(cl, scheme).Service("test", core.ServiceTypeClusterIP, 5555, "pod", owner)
		err := sc.EnsureExists()
		assert.NoError(t, err)
		assert.Equal(t, "10.0.0.10", sc.ClusterIP())
	})
}
