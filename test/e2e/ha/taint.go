package ha

import (
	"context"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func untaintNodes(k kubernetes.Interface, nodeLabelSelector string) error {
	nodes, err := k.CoreV1().Nodes().List(context.Background(), meta.ListOptions{
		LabelSelector: nodeLabelSelector,
	})

	if err != nil {
		return err
	}

	for _, n := range nodes.Items {
		for i, tn := range n.Spec.Taints {
			if tn.Key != "e2e.test/failure" {
				continue
			}
			s := n.Spec.Taints
			s[len(s)-1], s[i] = s[i], s[len(s)-1]
			n.Spec.Taints = s[:len(s)-1]
			_, err = k.CoreV1().Nodes().Update(context.Background(), &n, meta.UpdateOptions{})
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func taintWorker(k kubernetes.Interface, nodeLabelSelector string) error {
	nodes, err := k.CoreV1().Nodes().List(context.Background(), meta.ListOptions{
		LabelSelector: nodeLabelSelector,
	})
	if err != nil {
		return err
	}
	for _, node := range nodes.Items {
		if _, ok := node.Labels["node-role.kubernetes.io/master"]; !ok {
			node.Spec.Taints = append(node.Spec.Taints, core.Taint{
				Key:    "e2e.test/failure",
				Effect: core.TaintEffectNoExecute,
			})
			_, err = k.CoreV1().Nodes().Update(context.Background(), &node, meta.UpdateOptions{})
			return err
		}
	}

	return nil
}
