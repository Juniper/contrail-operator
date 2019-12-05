package controller

import (
	"github.com/Juniper/contrail-operator/pkg/controller/contrailcommand"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, contrailcommand.Add)
}
