package certificates

import (
	core "k8s.io/api/core/v1"
)

type certificateSubject struct {
	name     string
	hostname string
	ip       string
}

type certificateSubjects struct {
	pods        *core.PodList
	hostNetwork bool
}

func (c certificateSubjects) createCertificateSubjects() []certificateSubject {
	subjects := []certificateSubject{}
	for _, pod := range c.pods.Items {
		var hostname string
		if c.hostNetwork {
			hostname = pod.Spec.NodeName
		} else {
			hostname = pod.Spec.Hostname
		}
		subjects = append(subjects, certificateSubject{name: pod.Name, hostname: hostname, ip: pod.Status.PodIP})
	}
	return subjects

}
