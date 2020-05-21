package vrouter_nodes

import (
	"fmt"

	"github.com/Juniper/contrail-operator/contrail-provisioner/types"

	contrailTypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
)

const nodeType string = "virtual-router"

type Action int

const (
	updateAction Action = iota
	createAction
	deleteAction
	noopAction
)

type NodeWithAction struct {
	node   *types.VrouterNode
	action Action
}

func ReconcileVrouterNodes(contrailClient types.ApiClient, requiredNodes []*types.VrouterNode) error {
	nodesInApiServer, err := getVrouterNodesInApiServer(contrailClient)
	if err != nil {
		return err
	}
	actionMap := createVrouterNodesActionMap(nodesInApiServer, requiredNodes)
	err = executeActionMap(actionMap, contrailClient)
	if err != nil {
		return err
	}
	err = ensureVMIVhost0InterfaceForVirtualRouters(contrailClient)
	if err != nil {
		return err
	}
	return nil
}

func getVrouterNodesInApiServer(contrailClient types.ApiClient) ([]*types.VrouterNode, error) {
	nodesInApiServer := []*types.VrouterNode{}
	vncNodeList, err := contrailClient.List(nodeType)
	if err != nil {
		return nodesInApiServer, err
	}
	for _, vncNode := range vncNodeList {
		obj, err := contrailClient.ReadListResult(nodeType, &vncNode)
		if err != nil {
			return nodesInApiServer, err
		}
		typedNode := obj.(*contrailTypes.VirtualRouter)

		node := &types.VrouterNode{
			IPAddress: typedNode.GetVirtualRouterIpAddress(),
			Hostname:  typedNode.GetName(),
		}
		nodesInApiServer = append(nodesInApiServer, node)
	}

	return nodesInApiServer, nil
}

func createVrouterNodesActionMap(nodesInApiServer []*types.VrouterNode, requiredNodes []*types.VrouterNode) map[string]NodeWithAction {
	var actionMap = make(map[string]NodeWithAction)
	for _, requiredNode := range requiredNodes {
		actionMap[requiredNode.Hostname] = NodeWithAction{node: requiredNode, action: createAction}
	}
	for _, nodeInApiServer := range nodesInApiServer {
		if requiredNodeWithAction, ok := actionMap[nodeInApiServer.Hostname]; ok {
			requiredAction := noopAction
			if requiredNodeWithAction.node.IPAddress != nodeInApiServer.IPAddress {
				requiredAction = updateAction
			}
			actionMap[nodeInApiServer.Hostname] = NodeWithAction{
				node:   requiredNodeWithAction.node,
				action: requiredAction,
			}
		} else {
			actionMap[nodeInApiServer.Hostname] = NodeWithAction{
				node:   nodeInApiServer,
				action: deleteAction,
			}
		}
	}
	return actionMap
}

func executeActionMap(actionMap map[string]NodeWithAction, contrailClient types.ApiClient) error {
	for _, nodeWithAction := range actionMap {
		var err error
		switch nodeWithAction.action {
		case updateAction:
			fmt.Println("updating vrouter node ", nodeWithAction.node.Hostname)
			err = nodeWithAction.node.Update(contrailClient)
		case createAction:
			fmt.Println("creating vrouter node ", nodeWithAction.node.Hostname)
			err = nodeWithAction.node.Create(contrailClient)
		case deleteAction:
			fmt.Println("deleting vrouter node ", nodeWithAction.node.Hostname)
			err = nodeWithAction.node.Delete(contrailClient)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func ensureVMIVhost0InterfaceForVirtualRouters(contrailClient types.ApiClient) error {
	virtualRouterList, err := contrailClient.List(nodeType)
	if err != nil {
		return err
	}
	for _, virtualRouter := range virtualRouterList {
		obj, err := contrailClient.ReadListResult(nodeType, &virtualRouter)
		if err != nil {
			return err
		}
		typedVirtualRouter := obj.(*contrailTypes.VirtualRouter)

		vhost0VMIPresent, err := types.Vhost0VMIPresent(typedVirtualRouter, contrailClient)
		if err != nil {
			return err
		}
		if !vhost0VMIPresent {
			err = types.CreateVhost0VMI(typedVirtualRouter, contrailClient)
			if err != nil {
				return err
			}
		}
	}
	return nil
}