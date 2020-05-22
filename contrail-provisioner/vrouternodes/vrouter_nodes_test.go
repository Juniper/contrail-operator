package vrouternodes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	contrail "github.com/Juniper/contrail-go-api"
	contrailTypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
	"github.com/Juniper/contrail-operator/contrail-provisioner/fake"
	"github.com/Juniper/contrail-operator/contrail-provisioner/types"
)

func TestGetVrouterNodesInApiServerCreatesVrouterNodeObjects(t *testing.T) {
	fakeContrailClient := fake.GetDefaultFakeContrailClient()
	fakeContrailClient.ListFake = func(string) ([]contrail.ListResult, error) {
		return []contrail.ListResult{{}, {}}, nil
	}
	virtualRouterOne := contrailTypes.VirtualRouter{}
	virtualRouterOne.SetVirtualRouterIpAddress("1.1.1.1")
	virtualRouterOne.SetName("virtual-router-one")
	fakeContrailClient.ReadListResultFake = func(string, *contrail.ListResult) (contrail.IObject, error) {
		return &virtualRouterOne, nil
	}

	expectedVirtualRouterNodes := []*types.VrouterNode{
		{IPAddress: "1.1.1.1", Hostname: "virtual-router-one"},
		{IPAddress: "1.1.1.1", Hostname: "virtual-router-one"},
	}
	actualVirtualRouterNodes, err := getVrouterNodesInApiServer(fakeContrailClient)

	assert.NoError(t, err)
	assert.Len(t, actualVirtualRouterNodes, len(expectedVirtualRouterNodes))
	for i := 0; i < len(expectedVirtualRouterNodes); i++ {
		assert.Equal(t, *expectedVirtualRouterNodes[i], *actualVirtualRouterNodes[i])
	}
}

func TestGetVrouterNodesInApiServerReturnsEmptyListWhenNoNodesInApiServer(t *testing.T) {
	fakeContrailClient := fake.GetDefaultFakeContrailClient()
	fakeContrailClient.ListFake = func(string) ([]contrail.ListResult, error) {
		return []contrail.ListResult{}, nil
	}
	virtualRouterOne := contrailTypes.VirtualRouter{}
	virtualRouterOne.SetVirtualRouterIpAddress("1.1.1.1")
	virtualRouterOne.SetName("virtual-router-one")
	fakeContrailClient.ReadListResultFake = func(string, *contrail.ListResult) (contrail.IObject, error) {
		return &virtualRouterOne, nil
	}

	actualVirtualRouterNodes, err := getVrouterNodesInApiServer(fakeContrailClient)

	assert.NoError(t, err)
	assert.Len(t, actualVirtualRouterNodes, 0)
}

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
