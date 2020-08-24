package cr

import (
	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"

	"github.com/ghodss/yaml"
)

var yamlDataContrailCNI = `
apiVersion: contrail.juniper.net/v1alpha1
kind: ContrailCNI
metadata:
  name: cluster-1
`

func GetContrailCNICr() *v1alpha1.ContrailCNI {
	cr := v1alpha1.ContrailCNI{}
	err := yaml.Unmarshal([]byte(yamlDataContrailCNI), &cr)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDataContrailCNI))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &cr)
	if err != nil {
		panic(err)
	}
	return &cr
}
