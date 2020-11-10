package vrouternode

import (
	"testing"

	"github.com/stretchr/testify/assert"

	contrail "github.com/Juniper/contrail-go-api"

	contrailtypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailnode"
	"github.com/Juniper/contrail-operator/contrail-provisioner/fake"
)

func TestGetVrouterNodesInApiServerCreatesVrouterNodeObjects(t *testing.T) {
	fakeContrailClient := fake.GetDefaultFakeContrailClient()
	fakeContrailClient.ListFake = func(string) ([]contrail.ListResult, error) {
		return []contrail.ListResult{{}, {}}, nil
	}
	virtualRouterOne := contrailtypes.VirtualRouter{}
	virtualRouterOne.SetVirtualRouterIpAddress("1.1.1.1")
	virtualRouterOne.SetName("virtual-router-one")
	fakeContrailClient.ReadListResultFake = func(string, *contrail.ListResult) (contrail.IObject, error) {
		return &virtualRouterOne, nil
	}

	expectedContrailNodes := []contrailnode.ContrailNode{
		&VrouterNode{contrailnode.Node{IPAddress: "1.1.1.1", Hostname: "virtual-router-one", Annotations: map[string]string{}}},
		&VrouterNode{contrailnode.Node{IPAddress: "1.1.1.1", Hostname: "virtual-router-one", Annotations: map[string]string{}}},
	}
	actualContrailNodes, err := GetContrailNodesFromApiServer(fakeContrailClient)

	assert.NoError(t, err)
	assert.Len(t, actualContrailNodes, len(expectedContrailNodes))
	for idx, expectedNode := range expectedContrailNodes {
		assert.Equal(t, expectedNode, actualContrailNodes[idx])
	}
}

func TestGetVrouterNodesInApiServerReturnsEmptyListWhenNoNodesInApiServer(t *testing.T) {
	fakeContrailClient := fake.GetDefaultFakeContrailClient()
	fakeContrailClient.ListFake = func(string) ([]contrail.ListResult, error) {
		return []contrail.ListResult{}, nil
	}
	virtualRouterOne := contrailtypes.VirtualRouter{}
	virtualRouterOne.SetVirtualRouterIpAddress("1.1.1.1")
	virtualRouterOne.SetName("virtual-router-one")
	fakeContrailClient.ReadListResultFake = func(string, *contrail.ListResult) (contrail.IObject, error) {
		return &virtualRouterOne, nil
	}

	actualVirtualRouterNodes, err := GetContrailNodesFromApiServer(fakeContrailClient)

	assert.NoError(t, err)
	assert.Len(t, actualVirtualRouterNodes, 0)
}
