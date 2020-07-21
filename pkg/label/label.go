package label

//New is used to create a new default operators label
func New(instanceType, instanceName string) map[string]string {
	return map[string]string{"contrail_manager": instanceType, instanceType: instanceName}
}
