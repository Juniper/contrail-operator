package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"

	//mgr "github.com/Juniper/contrail-operator/pkg/controller/manager"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager) error {
	/*
		if err := mgr.Add(m); err != nil {
			return err
		}
	*/

	for _, f := range AddToManagerFuncs {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}
