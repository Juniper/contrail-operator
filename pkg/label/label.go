package label

import (
	"fmt"
	"strings"
)

//New is used to create a new default operators label
func New(instanceType, instanceName string) map[string]string {
	return map[string]string{"contrail_manager": instanceType, instanceType: instanceName}
}

func AsString(instanceType, instanceName string) string {
	var labels []string
	for k, v := range New(instanceType, instanceName) {
		label := fmt.Sprintf("%s: %s", k, v)
		labels = append(labels, label)
	}

	return fmt.Sprintf("{%s}", strings.Join(labels, ", "))
}

//NewLabelSelector is used to create default operator label selector
func NewLabelSelector(instanceType, instanceName string) map[string]string {
	return map[string]string{"contrail_manager": instanceType, instanceType: instanceName}
}
