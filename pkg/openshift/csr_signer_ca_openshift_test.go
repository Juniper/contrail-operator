package openshift_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/Juniper/contrail-operator/pkg/cacertificates"
	"github.com/Juniper/contrail-operator/pkg/openshift"
)

func TestCSRSignerCAOpenshift(t *testing.T) {
	testcases := []struct {
		name               string
		configMapName      string
		configMapNamespace string
		configMapData      map[string]string
		expected           string
		errorExpected      bool
	}{
		{
			name:               "Data retrieved from the the csr-signer-ca ConfigMap",
			configMapName:      "csr-signer-ca",
			configMapNamespace: "openshift-kube-controller-manager-operator",
			configMapData:      map[string]string{"ca-bundle.crt": "test-ca-data"},
			expected:           "test-ca-data",
			errorExpected:      false,
		},
		{
			name:               "Empty string and error returned when the required ConfigMap does not exist",
			configMapName:      "csr-signer-ca",
			configMapNamespace: "other-namespace",
			configMapData:      map[string]string{"ca-bundle.crt": "test-ca-data"},
			errorExpected:      true,
		},
		{
			name:               "Empty string and error returned when the required field in the ConfigMap does not exist",
			configMapName:      "csr-signer-ca",
			configMapNamespace: "openshift-kube-controller-manager-operator",
			configMapData:      map[string]string{"test-field-name": "test-ca-data"},
			errorExpected:      true,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			inputConfigMap := &core.ConfigMap{
				ObjectMeta: meta.ObjectMeta{
					Name:      test.configMapName,
					Namespace: test.configMapNamespace,
				},
				TypeMeta: meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
				Data:     test.configMapData,
			}
			fakeCoreClient := fake.NewSimpleClientset(inputConfigMap).CoreV1()
			var csrSignerCA cacertificates.CSRSignerCA = openshift.CSRSignerCAOpenshift{Client: fakeCoreClient}
			actual, err := csrSignerCA.CSRSignerCA()
			assert.Equal(t, err != nil, test.errorExpected)
			assert.Equal(t, test.expected, actual)
		})
	}
}
