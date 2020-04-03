package k8s_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestGetServiceAccountTokenSecret(t *testing.T) {
	t.Run("should retrieve the secret of the ServiceAccountTokenSecretType for given ServiceAccount", func(t *testing.T) {
		// given
		serviceAccountTokenSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "default-token-test",
				Namespace: "kube-system",
			},
			Type: core.SecretTypeServiceAccountToken,
			Data: map[string][]byte{
				"ca.crt": []byte("test-ca-value"),
			},
		}
		otherSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-secret",
				Namespace: "kube-system",
			},
			Type: core.SecretTypeOpaque,
			Data: map[string][]byte{
				"ca.crt": []byte("test"),
			},
		}
		defaultServiceAccount := &core.ServiceAccount{
			ObjectMeta: meta.ObjectMeta{
				Name:      "default",
				Namespace: "kube-system",
			},
			Secrets: []core.ObjectReference{
				core.ObjectReference{
					Name:      "default-token-test",
					Namespace: "kube-system",
				},
				core.ObjectReference{
					Name:      "test-secret",
					Namespace: "kube-system",
				},
			},
		}
		fakeCoreClient := fake.NewSimpleClientset(serviceAccountTokenSecret, otherSecret, defaultServiceAccount).CoreV1()
		csrSignerCAK8S := k8s.CSRSignerCA{Client: fakeCoreClient}
		// when
		result, err := csrSignerCAK8S.GetServiceAccountTokenSecret(defaultServiceAccount)
		// then
		if assert.NoError(t, err) {
			assert.Equal(t, serviceAccountTokenSecret, result)
		}

	})

	t.Run("should return error when there is no secret of the ServiceAccountToken type", func(t *testing.T) {
		// given
		serviceAccountTokenSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "default-token-test",
				Namespace: "kube-system",
			},
			Type: core.SecretTypeOpaque,
			Data: map[string][]byte{
				"ca.crt": []byte("test-ca-value"),
			},
		}
		otherSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-secret",
				Namespace: "kube-system",
			},
			Type: core.SecretTypeOpaque,
			Data: map[string][]byte{
				"ca.crt": []byte("test"),
			},
		}
		defaultServiceAccount := &core.ServiceAccount{
			ObjectMeta: meta.ObjectMeta{
				Name:      "default",
				Namespace: "kube-system",
			},
			Secrets: []core.ObjectReference{
				core.ObjectReference{
					Name:      "default-token-test",
					Namespace: "kube-system",
				},
				core.ObjectReference{
					Name:      "test-secret",
					Namespace: "kube-system",
				},
			},
		}
		fakeCoreClient := fake.NewSimpleClientset(serviceAccountTokenSecret, otherSecret, defaultServiceAccount).CoreV1()
		csrSignerCAK8S := k8s.CSRSignerCA{Client: fakeCoreClient}
		// when
		result, err := csrSignerCAK8S.GetServiceAccountTokenSecret(defaultServiceAccount)
		// then
		assert.Error(t, err)
		assert.Equal(t, (*core.Secret)(nil), result)

	})

	t.Run("should return error when secret assigned to the service account cannot be retrieved via corev1 client", func(t *testing.T) {
		// given
		defaultServiceAccount := &core.ServiceAccount{
			ObjectMeta: meta.ObjectMeta{
				Name:      "default",
				Namespace: "kube-system",
			},
			Secrets: []core.ObjectReference{
				core.ObjectReference{
					Name:      "test-secret",
					Namespace: "kube-system",
				},
			},
		}
		fakeCoreClient := fake.NewSimpleClientset(defaultServiceAccount).CoreV1()
		csrSignerCAK8S := k8s.CSRSignerCA{Client: fakeCoreClient}
		// when
		result, err := csrSignerCAK8S.GetServiceAccountTokenSecret(defaultServiceAccount)
		// then
		assert.Error(t, err)
		assert.Equal(t, (*core.Secret)(nil), result)
	})
}

func TestCSRSignerCAK8S(t *testing.T) {
	t.Run("should retrieve the ca.crt value stored in service account token secret", func(t *testing.T) {
		// given
		serviceAccountTokenSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "default-token-test",
				Namespace: "kube-system",
			},
			Type: core.SecretTypeServiceAccountToken,
			Data: map[string][]byte{
				"ca.crt": []byte("test-ca-value"),
			},
		}
		otherSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-secret",
				Namespace: "kube-system",
			},
			Type: core.SecretTypeOpaque,
			Data: map[string][]byte{
				"ca.crt": []byte("test-other-secret-value"),
			},
		}
		defaultServiceAccount := &core.ServiceAccount{
			ObjectMeta: meta.ObjectMeta{
				Name:      "default",
				Namespace: "kube-system",
			},
			Secrets: []core.ObjectReference{
				core.ObjectReference{
					Name:      "default-token-test",
					Namespace: "kube-system",
				},
				core.ObjectReference{
					Name:      "test-secret",
					Namespace: "kube-system",
				},
			},
		}
		fakeCoreClient := fake.NewSimpleClientset(serviceAccountTokenSecret, otherSecret, defaultServiceAccount).CoreV1()
		csrSignerCAK8S := k8s.CSRSignerCA{Client: fakeCoreClient}
		// when
		result, err := csrSignerCAK8S.CACert()
		// then
		if assert.NoError(t, err) {
			assert.Equal(t, "test-ca-value", result)
		}
	})

	t.Run("should return error when no service account token secret is found", func(t *testing.T) {
		// given
		otherSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-secret",
				Namespace: "kube-system",
			},
			Type: core.SecretTypeOpaque,
			Data: map[string][]byte{
				"ca.crt": []byte("test-other-secret-value"),
			},
		}
		defaultServiceAccount := &core.ServiceAccount{
			ObjectMeta: meta.ObjectMeta{
				Name:      "default",
				Namespace: "kube-system",
			},
			Secrets: []core.ObjectReference{
				core.ObjectReference{
					Name:      "test-secret",
					Namespace: "kube-system",
				},
			},
		}
		fakeCoreClient := fake.NewSimpleClientset(otherSecret, defaultServiceAccount).CoreV1()
		csrSignerCAK8S := k8s.CSRSignerCA{Client: fakeCoreClient}
		// when
		result, err := csrSignerCAK8S.CACert()
		// then
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("should return error when service account token secret has no ca.crt field in Data", func(t *testing.T) {
		// given
		serviceAccountTokenSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "default-token-test",
				Namespace: "kube-system",
			},
			Type: core.SecretTypeServiceAccountToken,
		}
		defaultServiceAccount := &core.ServiceAccount{
			ObjectMeta: meta.ObjectMeta{
				Name:      "default",
				Namespace: "kube-system",
			},
			Secrets: []core.ObjectReference{
				core.ObjectReference{
					Name:      "default-token-test",
					Namespace: "kube-system",
				},
			},
		}
		fakeCoreClient := fake.NewSimpleClientset(serviceAccountTokenSecret, defaultServiceAccount).CoreV1()
		csrSignerCAK8S := k8s.CSRSignerCA{Client: fakeCoreClient}
		// when
		result, err := csrSignerCAK8S.CACert()
		// then
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})

}
