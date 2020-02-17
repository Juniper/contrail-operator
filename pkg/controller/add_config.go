package controller

import (
	"github.com/Juniper/contrail-operator/pkg/controller/config"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, config.Add)
}
