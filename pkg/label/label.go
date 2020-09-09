package label

import (
	"fmt"
)

//New is used to create a new default operators label
func New(instanceType, instanceName string) map[string]string {
	return map[string]string{"contrail_manager": instanceType, instanceType: instanceName}
}

func AsString(instanceType, instanceName string) string {
	return fmt.Sprintf("%s: %s", instanceType, instanceName)
}

//NewLabelSelector is used to create default operator label selector
func NewLabelSelector(instanceType, instanceName string) map[string]string {
	return map[string]string{"contrail_manager": instanceType, instanceType: instanceName}
}
