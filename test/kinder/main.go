package main

import (
	"time"

	"github.com/Juniper/contrail-operator/test/kinder/kinder"
)

func main() {
	nodeCount := 3
	kindCluster := &kinder.KindCluster{
		Name:  "cluster1",
		Nodes: &nodeCount,
	}

	done := make(chan bool)
	go func() {
		kindCluster.Start()
	}()

	kindCluster.WaitForReady()
	kindCluster.Delete(20 * time.Second)
	close(done)
}
