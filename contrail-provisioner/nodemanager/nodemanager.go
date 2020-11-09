package nodemanager

import (
	"gopkg.in/yaml.v2"

	"github.com/Juniper/contrail-operator/contrail-provisioner/analyticsnode"
	"github.com/Juniper/contrail-operator/contrail-provisioner/confignode"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailclient"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailnode"
	"github.com/Juniper/contrail-operator/contrail-provisioner/controlnode"
	"github.com/Juniper/contrail-operator/contrail-provisioner/databasenode"
	"github.com/Juniper/contrail-operator/contrail-provisioner/reconcile"
	"github.com/Juniper/contrail-operator/contrail-provisioner/vrouternode"
)

func getContrailNodesFromBytes(data []byte, nodeType contrailnode.ContrailNodeType) []contrailnode.ContrailNode {
	var contrailNodes []contrailnode.ContrailNode
	switch nodeType {
	case contrailnode.VrouterNode:
		var vrouterNodes []*vrouternode.VrouterNode
		err := yaml.Unmarshal(data, &vrouterNodes)
		if err != nil {
			panic(err)
		}
		for _, v := range vrouterNodes {
			contrailNodes = append(contrailNodes, v)
		}
	case contrailnode.AnalyticsNode:
		var analyticsNodes []*analyticsnode.AnalyticsNode
		err := yaml.Unmarshal(data, &analyticsNodes)
		if err != nil {
			panic(err)
		}
		for _, v := range analyticsNodes {
			contrailNodes = append(contrailNodes, v)
		}
	case contrailnode.ConfigNode:
		var configNodes []*confignode.ConfigNode
		err := yaml.Unmarshal(data, &configNodes)
		if err != nil {
			panic(err)
		}
		for _, v := range configNodes {
			contrailNodes = append(contrailNodes, v)
		}
	case contrailnode.ControlNode:
		var controlNodes []*controlnode.ControlNode
		err := yaml.Unmarshal(data, &controlNodes)
		if err != nil {
			panic(err)
		}
		for _, v := range controlNodes {
			contrailNodes = append(contrailNodes, v)
		}
	case contrailnode.DatabaseNode:
		var databaseNodes []*databasenode.DatabaseNode
		err := yaml.Unmarshal(data, &databaseNodes)
		if err != nil {
			panic(err)
		}
		for _, v := range databaseNodes {
			contrailNodes = append(contrailNodes, v)
		}
	}
	return contrailNodes
}

func getContrailNodesInApiServer(contrailClient contrailclient.ApiClient, nodeType contrailnode.ContrailNodeType) ([]contrailnode.ContrailNode, error) {
	var contrailNodesInApiServer []contrailnode.ContrailNode
	var err error
	switch nodeType {
	case contrailnode.VrouterNode:
		contrailNodesInApiServer, err = vrouternode.GetContrailNodesFromApiServer(contrailClient)
	case contrailnode.AnalyticsNode:
		contrailNodesInApiServer, err = analyticsnode.GetContrailNodesFromApiServer(contrailClient)
	case contrailnode.ConfigNode:
		contrailNodesInApiServer, err = confignode.GetContrailNodesFromApiServer(contrailClient)
	case contrailnode.ControlNode:
		contrailNodesInApiServer, err = controlnode.GetContrailNodesFromApiServer(contrailClient)
	case contrailnode.DatabaseNode:
		contrailNodesInApiServer, err = databasenode.GetContrailNodesFromApiServer(contrailClient)
	}
	return contrailNodesInApiServer, err
}

func ManageNodes(requiredNodesData []byte, requiredAnnotations map[string]string, nodeType contrailnode.ContrailNodeType, contrailClient contrailclient.ApiClient) error {
	requiredNodes := getContrailNodesFromBytes(requiredNodesData, nodeType)
	reconciler := reconcile.NewReconciler(contrailClient, requiredNodes, requiredAnnotations)
	nodesInApiServer, err := getContrailNodesInApiServer(contrailClient, nodeType)
	if err != nil {
		return err
	}
	return reconciler.ReconcileNodes(nodesInApiServer)
}
