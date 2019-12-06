package owner_test

import (
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/Juniper/contrail-operator/pkg/owner"
)

func TestSetOwnerReference(t *testing.T) {
	trueVal := true
	falseVal := false
	scheme := runtime.NewScheme()
	err := core.SchemeBuilder.AddToScheme(scheme)
	assert.NoError(t, err)
	tests := []struct {
		name      string
		owner     *core.Pod
		dependent *core.Pod
		expected  *core.Pod
	}{
		{
			name: "add owner reference",
			owner: &core.Pod{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "owner", UID: "uuid"},
			},
			dependent: &core.Pod{ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "dependent", UID: "uuid-d"}},
			expected: &core.Pod{ObjectMeta: meta.ObjectMeta{
				Namespace: "default", Name: "dependent", UID: "uuid-d",
				OwnerReferences: []meta.OwnerReference{{"v1", "Pod", "owner", "uuid", &falseVal, &falseVal}},
			}},
		},
		{
			name: "add owner reference when there is a controller",
			owner: &core.Pod{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "owner", UID: "uuid"},
			},
			dependent: &core.Pod{ObjectMeta: meta.ObjectMeta{
				Namespace: "default", Name: "dependent", UID: "uuid-d",
				OwnerReferences: []meta.OwnerReference{{"v1", "Pod", "controller", "uuid-2", &trueVal, &trueVal}},
			}},
			expected: &core.Pod{ObjectMeta: meta.ObjectMeta{
				Namespace: "default", Name: "dependent", UID: "uuid-d",
				OwnerReferences: []meta.OwnerReference{
					{"v1", "Pod", "controller", "uuid-2", &trueVal, &trueVal},
					{"v1", "Pod", "owner", "uuid", &falseVal, &falseVal},
				},
			}},
		},
		{
			name: "update owner reference",
			owner: &core.Pod{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "owner", UID: "uuid-2"},
			},
			dependent: &core.Pod{ObjectMeta: meta.ObjectMeta{
				Namespace: "default", Name: "dependent", UID: "uuid-d",
				OwnerReferences: []meta.OwnerReference{{"v1", "Pod", "owner", "uuid", &falseVal, &falseVal}},
			}},
			expected: &core.Pod{ObjectMeta: meta.ObjectMeta{
				Namespace: "default", Name: "dependent", UID: "uuid-d",
				OwnerReferences: []meta.OwnerReference{{"v1", "Pod", "owner", "uuid-2", &falseVal, &falseVal}},
			}},
		},
		{
			name: "set owner reference is idempotent",
			owner: &core.Pod{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "owner", UID: "uuid"},
			},
			dependent: &core.Pod{ObjectMeta: meta.ObjectMeta{
				Namespace: "default", Name: "dependent", UID: "uuid-d",
				OwnerReferences: []meta.OwnerReference{{"v1", "Pod", "owner", "uuid", &falseVal, &falseVal}},
			}},
			expected: &core.Pod{ObjectMeta: meta.ObjectMeta{
				Namespace: "default", Name: "dependent", UID: "uuid-d",
				OwnerReferences: []meta.OwnerReference{{"v1", "Pod", "owner", "uuid", &falseVal, &falseVal}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, tt.dependent)
			err := owner.EnsureOwnerReference(tt.owner, tt.dependent, cl, scheme)
			assert.NoError(t, err)

			pod := &core.Pod{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      tt.expected.Name,
				Namespace: tt.expected.Namespace,
			}, pod)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, pod)
		})
	}
}
