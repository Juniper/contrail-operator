package cr

import (
	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"

	"github.com/ghodss/yaml"
)

var yamlDataProvisionManager = `
apiVersion: contrail.juniper.net/v1alpha1
kind: ProvisionManager
metadata:
  name: example-ProvisionManager
  labels:
    contrail_manager: ProvisionManager
`

func GetProvisionManagerCr() *v1alpha1.ProvisionManager {
	cr := v1alpha1.ProvisionManager{}
	err := yaml.Unmarshal([]byte(yamlDataProvisionManager), &cr)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDataProvisionManager))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &cr)
	if err != nil {
		panic(err)
	}
	return &cr
}
