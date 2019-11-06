package cr

import (
	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"

	"github.com/ghodss/yaml"
)

var yamlDataConfig = `
apiVersion: contrail.juniper.net/v1alpha1
kind: Config
metadata:
  name: example-config
spec:
  # Add fields here
  size: 3
  service: 
    activate: true
`

func GetConfigCr() *v1alpha1.Config {
	cr := v1alpha1.Config{}
	err := yaml.Unmarshal([]byte(yamlDataConfig), &cr)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDataConfig))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &cr)
	if err != nil {
		panic(err)
	}
	return &cr
}
