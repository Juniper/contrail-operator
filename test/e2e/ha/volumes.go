package ha

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func deleteAllPVs(kubeClient kubernetes.Interface, storageClass string) error {
	pvs, err := kubeClient.CoreV1().PersistentVolumes().List(meta.ListOptions{})
	if err != nil {
		return err
	}

	for _, pv := range pvs.Items {
		if pv.Spec.StorageClassName != storageClass {
			continue
		}
		if err = kubeClient.CoreV1().PersistentVolumes().Delete(pv.GetName(), &meta.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}
