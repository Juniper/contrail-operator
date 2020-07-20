package label


func New(instanceType, instanceName string) map[string]string {
	return map[string]string{"contrail_manager": instanceType, instanceType: instanceName}
}
