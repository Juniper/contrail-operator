package certificates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCertificateSubject(t *testing.T) {
	firstPodName := "first"
	firstPodNodeName := "nodeName1"
	firstPodHostname := "hostName1"
	firstPodIp := "ip1"
	firstPod := core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: firstPodName,
		},
		Spec: core.PodSpec{
			NodeName: firstPodNodeName,
			Hostname: firstPodHostname,
		},
		Status: core.PodStatus{
			PodIP: firstPodIp,
		},
	}

	secondPodName := "second"
	secondPodNodeName := "nodeName2"
	secondPodHostname := "hostName2"
	secondPodIp := "ip2"
	secondPod := core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: secondPodName,
		},
		Spec: core.PodSpec{
			NodeName: secondPodNodeName,
			Hostname: secondPodHostname,
		},
		Status: core.PodStatus{
			PodIP: secondPodIp,
		},
	}

	pods := &core.PodList{Items: []core.Pod{firstPod, secondPod}}

	tests := []struct {
		name             string
		hostNetwork      bool
		podList          *core.PodList
		expectedSubjects []certificateSubject
	}{
		{
			name:    "should create not host network subjects",
			podList: pods,
			expectedSubjects: []certificateSubject{
				certificateSubject{
					name:     firstPodName,
					hostname: firstPodHostname,
					ip:       firstPodIp,
				},
				certificateSubject{
					name:     secondPodName,
					hostname: secondPodHostname,
					ip:       secondPodIp,
				},
			},
		},
		{
			name:        "should create host network subjects",
			hostNetwork: true,
			podList:     pods,
			expectedSubjects: []certificateSubject{
				certificateSubject{
					name:     firstPodName,
					hostname: firstPodNodeName,
					ip:       firstPodIp,
				},
				certificateSubject{
					name:     secondPodName,
					hostname: secondPodNodeName,
					ip:       secondPodIp,
				},
			},
		},
	}

	for _, test := range tests {
		subs := certificateSubjects{test.podList, test.hostNetwork}
		assert.Equal(t, subs.createCertificateSubjects(), test.expectedSubjects)
	}
}
