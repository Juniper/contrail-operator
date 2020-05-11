package certificates

import (
	core "k8s.io/api/core/v1"
)

type certificateSubject struct {
	name     string
	hostname string
	ip       string
}

func createCertificateSubjects(pods *core.PodList, hostNetwork bool) []certificateSubject {
	subject := []certificateSubject{}
	for _, pod := range pods.Items {
		var hostname string
		if hostNetwork {
			hostname = pod.Spec.NodeName
		} else {
			hostname = pod.Spec.Hostname
		}
		subject = append(subject, certificateSubject{name: pod.Name, hostname: hostname, ip: pod.Status.PodIP})
	}
	return subject

}
