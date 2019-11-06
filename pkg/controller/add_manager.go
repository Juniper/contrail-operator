package controller

import (
	"atom/atom/contrail/operator/pkg/controller/manager"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, manager.Add)
}
