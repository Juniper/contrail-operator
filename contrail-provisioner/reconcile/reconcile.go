package reconcile

import (
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailclient"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailnode"
)

type Action int

const (
	updateAction Action = iota
	createAction
	deleteAction
	noopAction
)

type NodeWithAction struct {
	node   contrailnode.ContrailNode
	action Action
}

type Reconciler struct {
	cl                  contrailclient.ApiClient
	requiredNodes       []contrailnode.ContrailNode
	requiredAnnotations map[string]string
}

func NewReconciler(cl contrailclient.ApiClient, requiredNodes []contrailnode.ContrailNode, requiredAnnotations map[string]string) *Reconciler {
	return &Reconciler{cl: cl, requiredNodes: requiredNodes, requiredAnnotations: requiredAnnotations}
}

func (r *Reconciler) ReconcileNodes(nodesInApiServer []contrailnode.ContrailNode) error {
	r.ensureRequiredAnnotationsOnRequiredNodes()
	managedNodesInApiServer := r.getNodesWithRequiredAnnotations(nodesInApiServer)
	actionMap := r.createContrailNodesActionMap(managedNodesInApiServer)
	if err := r.executeActionMap(actionMap); err != nil {
		return err
	}
	return nil
}

func (r *Reconciler) executeActionMap(actionMap map[string]NodeWithAction) error {
	for _, nodeWithAction := range actionMap {
		var err error
		switch nodeWithAction.action {
		case updateAction:
			err = nodeWithAction.node.Update(r.cl)
		case createAction:
			err = nodeWithAction.node.Create(r.cl)
		case deleteAction:
			err = nodeWithAction.node.Delete(r.cl)
		}
		if err != nil {
			return err
		}
		if nodeWithAction.action != deleteAction {
			if err := nodeWithAction.node.EnsureDependenciesExist(r.cl); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Reconciler) createContrailNodesActionMap(nodesInApiServer []contrailnode.ContrailNode) map[string]NodeWithAction {
	var actionMap = make(map[string]NodeWithAction)
	for _, requiredNode := range r.requiredNodes {
		actionMap[requiredNode.GetHostname()] = NodeWithAction{node: requiredNode, action: createAction}
	}
	for _, nodeInApiServer := range nodesInApiServer {
		nodeToActOn := nodeInApiServer
		requiredAction := deleteAction
		if requiredNodeWithAction, ok := actionMap[nodeInApiServer.GetHostname()]; ok {
			requiredAction = noopAction
			if !nodeInApiServer.Equal(requiredNodeWithAction.node) {
				requiredAction = updateAction
				nodeToActOn = requiredNodeWithAction.node
			}
		}
		actionMap[nodeInApiServer.GetHostname()] = NodeWithAction{
			node:   nodeToActOn,
			action: requiredAction,
		}
	}
	return actionMap
}

func (r *Reconciler) getNodesWithRequiredAnnotations(nodesInApiServer []contrailnode.ContrailNode) []contrailnode.ContrailNode {
	nodesWithRequiredAnnotations := []contrailnode.ContrailNode{}
	for _, node := range nodesInApiServer {
		if contrailclient.HasRequiredAnnotations(node.GetAnnotations(), r.requiredAnnotations) {
			nodesWithRequiredAnnotations = append(nodesWithRequiredAnnotations, node)
		}
	}
	return nodesWithRequiredAnnotations
}

func (r *Reconciler) ensureRequiredAnnotationsOnRequiredNodes() {
	for _, node := range r.requiredNodes {
		ensureRequiredAnnotationsSetOnNode(node, r.requiredAnnotations)
	}
}

func ensureRequiredAnnotationsSetOnNode(node contrailnode.ContrailNode, requiredAnnotations map[string]string) {
	nodeAnnotations := node.GetAnnotations()
	if nodeAnnotations == nil {
		nodeAnnotations = map[string]string{}
	}
	for key, val := range requiredAnnotations {
		nodeAnnotations[key] = val
	}
	node.SetAnnotations(nodeAnnotations)
}
