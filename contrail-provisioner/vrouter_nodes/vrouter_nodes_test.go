package vrouter_nodes

import (
	"testing"

	"github.com/Juniper/contrail-operator/contrail-provisioner/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateVrouterNodesActionMap(t *testing.T) {
	requiredNodeOne := &types.VrouterNode{
		IPAddress: "1.1.1.1",
		Hostname:  "first-node",
	}
	requiredNodeTwo := &types.VrouterNode{
		IPAddress: "2.2.2.2",
		Hostname:  "second-node",
	}
	modifiedRequiredNodeTwo := &types.VrouterNode{
		IPAddress: "2.2.2.3",
		Hostname:  "second-node",
	}
	apiServerNodeOne := &types.VrouterNode{
		IPAddress: "1.1.1.1",
		Hostname:  "first-node",
	}
	apiServerNodeTwo := &types.VrouterNode{
		IPAddress: "2.2.2.2",
		Hostname:  "second-node",
	}

	testCases := []struct {
		name              string
		nodesInApiServer  []*types.VrouterNode
		requiredNodes     []*types.VrouterNode
		expectedActionMap map[string]NodeWithAction
	}{
		{
			name:             "Create action for all required nodes if none in the Api Server",
			nodesInApiServer: []*types.VrouterNode{},
			requiredNodes:    []*types.VrouterNode{requiredNodeOne, requiredNodeTwo},
			expectedActionMap: map[string]NodeWithAction{
				"first-node":  {requiredNodeOne, createAction},
				"second-node": {requiredNodeTwo, createAction},
			},
		},
		{
			name:             "Noop action for nodes in the Api Server",
			nodesInApiServer: []*types.VrouterNode{apiServerNodeTwo},
			requiredNodes:    []*types.VrouterNode{requiredNodeOne, requiredNodeTwo},
			expectedActionMap: map[string]NodeWithAction{
				"first-node":  {requiredNodeOne, createAction},
				"second-node": {requiredNodeTwo, noopAction},
			},
		},
		{
			name:             "Delete action for not required nodes in the Api Server",
			nodesInApiServer: []*types.VrouterNode{apiServerNodeTwo, apiServerNodeOne},
			requiredNodes:    []*types.VrouterNode{requiredNodeTwo},
			expectedActionMap: map[string]NodeWithAction{
				"first-node":  {requiredNodeOne, deleteAction},
				"second-node": {requiredNodeTwo, noopAction},
			},
		},
		{
			name:             "Update action for modified node",
			nodesInApiServer: []*types.VrouterNode{apiServerNodeTwo, apiServerNodeOne},
			requiredNodes:    []*types.VrouterNode{requiredNodeOne, modifiedRequiredNodeTwo},
			expectedActionMap: map[string]NodeWithAction{
				"first-node":  {requiredNodeOne, noopAction},
				"second-node": {modifiedRequiredNodeTwo, updateAction},
			},
		},
		{
			name:              "Empty action map when there are no vrouter nodes",
			nodesInApiServer:  []*types.VrouterNode{},
			requiredNodes:     []*types.VrouterNode{},
			expectedActionMap: map[string]NodeWithAction{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualActionMap := createVrouterNodesActionMap(testCase.nodesInApiServer, testCase.requiredNodes)
			assert.Equal(t, testCase.expectedActionMap, actualActionMap)
		})
	}
}
