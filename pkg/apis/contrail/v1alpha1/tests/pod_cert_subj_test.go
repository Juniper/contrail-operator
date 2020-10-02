package contrailtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	crt "github.com/Juniper/contrail-operator/pkg/certificates"
)

func TestUtils(t *testing.T) {
	podItems := []corev1.Pod{}
	podTemplate1 := corev1.Pod{}
	podTemplate1.Name = "pod1"
	podTemplate1.Namespace = "default"
	podTemplate1.Spec.NodeName = "node1"
	podTemplate1.Spec.Hostname = "host1"
	podTemplate1.Status.PodIP = "1.5.5.5"
	podItems = append(podItems, podTemplate1)

	podTemplate2 := corev1.Pod{}
	podTemplate2.Name = "pod2"
	podTemplate2.Namespace = "default"
	podTemplate2.Spec.NodeName = "node2"
	podTemplate2.Spec.Hostname = "host2"
	podTemplate2.Status.PodIP = "1.5.5.6"
	podTemplate2.Annotations = map[string]string{"dataSubnetIP": "172.17.90.2"}
	podItems = append(podItems, podTemplate2)

	podList := corev1.PodList{
		Items: podItems,
	}

	withService := v1alpha1.PodAlternativeIPs{ServiceIP: "2.2.2.2"}
	withServiceAndRetriver := v1alpha1.PodAlternativeIPs{ServiceIP: "2.2.2.2", Retriever: testRetriver}
	var noAltIPs v1alpha1.PodAlternativeIPs
	var emptyList []string
	var hostNetwork *bool
	hostNetworkEnabled := true
	hostNetworkDisabled := false

	t.Run("should return subject with nodename", func(t *testing.T) {
		expected := []crt.CertificateSubject{
			crt.NewSubject("pod1", "node1", "1.5.5.5", emptyList),
			crt.NewSubject("pod2", "node2", "1.5.5.6", emptyList),
		}
		got := v1alpha1.PodsCertSubjects(&podList, hostNetwork, noAltIPs)
		assert.Equal(t, expected, got)
	})
	t.Run("should return subject with hostname", func(t *testing.T) {
		expected := []crt.CertificateSubject{
			crt.NewSubject("pod1", "host1", "1.5.5.5", emptyList),
			crt.NewSubject("pod2", "host2", "1.5.5.6", emptyList),
		}
		got := v1alpha1.PodsCertSubjects(&podList, &hostNetworkDisabled, noAltIPs)
		assert.Equal(t, expected, got)
	})
	t.Run("should return subject with serviceIP", func(t *testing.T) {
		expected := []crt.CertificateSubject{
			crt.NewSubject("pod1", "node1", "1.5.5.5", []string{"2.2.2.2"}),
			crt.NewSubject("pod2", "node2", "1.5.5.6", []string{"2.2.2.2"}),
		}
		got := v1alpha1.PodsCertSubjects(&podList, &hostNetworkEnabled, withService)
		assert.Equal(t, expected, got)
	})
	t.Run("should return subject with serviceIP and dataIP", func(t *testing.T) {
		expected := []crt.CertificateSubject{
			crt.NewSubject("pod1", "node1", "1.5.5.5", []string{"2.2.2.2"}),
			crt.NewSubject("pod2", "node2", "1.5.5.6", []string{"2.2.2.2", "172.17.90.2"}),
		}
		got := v1alpha1.PodsCertSubjects(&podList, &hostNetworkEnabled, withServiceAndRetriver)
		assert.Equal(t, expected, got)
	})

}

func testRetriver(pod corev1.Pod) []string {
	var altIPs []string
	if dataIP, isSet := pod.Annotations["dataSubnetIP"]; isSet {
		altIPs = append(altIPs, dataIP)
	}
	return altIPs
}
