package openshift_test
import (
	"testing"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/openshift"
	"github.com/stretchr/testify/assert"
)


func TestCSRSignerCA(t *testing.T) {
	testcases := []struct{
		name string
		configMapName string
		configMapNamespace string
		configMapData map[string] string
		expected string
		expectedErr error
	}{
		{
			name: "CSRSignerCA method retrieves data from the ",
			configMapName: "client-ca",
			configMapNamespace: "openshift-kube-controller-manager",
			configMapData: map[string] string{"ca-bundle.crt": "test"},
			expected: "",
			expectedErr: nil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			inputConfigMap := getConfigMapWithData(test.configMapName, test.configMapNamespace, test.configMapData)
			fakeCoreClient := fake.NewSimpleClientset(inputConfigMap).CoreV1()
			var csrSignerCA contrail.ManagerCSRSignerCA = openshift.CSRSignerCAGetter{Client: fakeCoreClient}
			actual, actualErr := csrSignerCA.CSRSignerCA()
			assert.Equal(t, test.expectedErr, actualErr)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func getConfigMapWithData(name string, namespace string, data map[string] string) *core.ConfigMap {
	cm := &core.ConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		TypeMeta: meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
	}
	cm.Data = data
	return cm
}