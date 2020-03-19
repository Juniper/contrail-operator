package certificates_test

import (
	"testing"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
)

func TestCSRSignerCA(t *testing.T) {
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
			expected:           "",
			errorExpected:      true,
		},
		{
			name:               "Empty string and error returned when the required field in the ConfigMap does not exist",
			configMapName:      "csr-signer-ca",
			configMapNamespace: "openshift-kube-controller-manager-operator",
			configMapData:      map[string]string{"test-field-name": "test-ca-data"},
			expected:           "",
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
			var csrSignerCA contrail.ManagerCSRSignerCA = certificates.CSRSignerCAOpenshift{Client: fakeCoreClient}
			actual, err := csrSignerCA.CSRSignerCA()
			if err != nil && !test.errorExpected {
				t.Errorf("Got unexpected error: \"%v\"", err)
			}
			if err == nil && test.errorExpected {
				t.Error("Didn't get expected error")
			}
			if actual != test.expected {
				t.Errorf("Output was incorrect. Expected \"%v\" got \"%v\"", test.expected, actual)
			}
		})
	}
}
