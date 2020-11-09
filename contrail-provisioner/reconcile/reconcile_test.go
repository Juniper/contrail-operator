package reconcile

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailnode"
	"github.com/Juniper/contrail-operator/contrail-provisioner/fake"
	"github.com/Juniper/contrail-operator/contrail-provisioner/vrouternode"
)

func TestCreateContrailNodesActionMap(t *testing.T) {
	var requiredNodeOne contrailnode.ContrailNode = &vrouternode.VrouterNode{
		contrailnode.Node{
			IPAddress: "1.1.1.1",
			Hostname:  "first-node",
		},
	}
	var requiredNodeTwo contrailnode.ContrailNode = &vrouternode.VrouterNode{
		contrailnode.Node{
			IPAddress: "2.2.2.2",
			Hostname:  "second-node",
		},
	}
	var modifiedRequiredNodeTwo contrailnode.ContrailNode = &vrouternode.VrouterNode{
		contrailnode.Node{
			IPAddress: "2.2.2.3",
			Hostname:  "second-node",
		},
	}
	var apiServerNodeOne contrailnode.ContrailNode = &vrouternode.VrouterNode{
		contrailnode.Node{
			IPAddress: "1.1.1.1",
			Hostname:  "first-node",
		},
	}
	var apiServerNodeTwo contrailnode.ContrailNode = &vrouternode.VrouterNode{
		contrailnode.Node{
			IPAddress: "2.2.2.2",
			Hostname:  "second-node",
		},
	}

	testCases := []struct {
		name              string
		nodesInApiServer  []contrailnode.ContrailNode
		requiredNodes     []contrailnode.ContrailNode
		expectedActionMap map[string]NodeWithAction
	}{
		{
			name:             "Create action for all required nodes if none in the Api Server",
			nodesInApiServer: []contrailnode.ContrailNode{},
			requiredNodes:    []contrailnode.ContrailNode{requiredNodeOne, requiredNodeTwo},
			expectedActionMap: map[string]NodeWithAction{
				"first-node":  {requiredNodeOne, createAction},
				"second-node": {requiredNodeTwo, createAction},
			},
		},
		{
			name:             "Update action for nodes in the Api Server",
			nodesInApiServer: []contrailnode.ContrailNode{apiServerNodeTwo},
			requiredNodes:    []contrailnode.ContrailNode{requiredNodeOne, requiredNodeTwo},
			expectedActionMap: map[string]NodeWithAction{
				"first-node":  {requiredNodeOne, createAction},
				"second-node": {requiredNodeTwo, updateAction},
			},
		},
		{
			name:             "Delete action for not required nodes in the Api Server",
			nodesInApiServer: []contrailnode.ContrailNode{apiServerNodeTwo, apiServerNodeOne},
			requiredNodes:    []contrailnode.ContrailNode{requiredNodeTwo},
			expectedActionMap: map[string]NodeWithAction{
				"first-node":  {requiredNodeOne, deleteAction},
				"second-node": {requiredNodeTwo, updateAction},
			},
		},
		{
			name:             "Update action for modified node and node in Api Server",
			nodesInApiServer: []contrailnode.ContrailNode{apiServerNodeTwo, apiServerNodeOne},
			requiredNodes:    []contrailnode.ContrailNode{requiredNodeOne, modifiedRequiredNodeTwo},
			expectedActionMap: map[string]NodeWithAction{
				"first-node":  {requiredNodeOne, updateAction},
				"second-node": {modifiedRequiredNodeTwo, updateAction},
			},
		},
		{
			name:              "Empty action map when there are no vrouter nodes",
			nodesInApiServer:  []contrailnode.ContrailNode{},
			requiredNodes:     []contrailnode.ContrailNode{},
			expectedActionMap: map[string]NodeWithAction{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			reconciler := NewReconciler(fake.GetDefaultFakeContrailClient(), testCase.requiredNodes, map[string]string{})
			actualActionMap := reconciler.createContrailNodesActionMap(testCase.nodesInApiServer)
			assert.Equal(t, testCase.expectedActionMap, actualActionMap)
		})
	}
}
