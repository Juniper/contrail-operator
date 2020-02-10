package cr

import (
	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"

	"github.com/ghodss/yaml"
)

var yamlDataWebui = `
apiVersion: contrail.juniper.net/v1alpha1
kind: Webui
metadata:
  name: cluster-1
`

func GetWebuiCr() *v1alpha1.Webui {
	cr := v1alpha1.Webui{}
	err := yaml.Unmarshal([]byte(yamlDataWebui), &cr)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDataWebui))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &cr)
	if err != nil {
		panic(err)
	}
	return &cr
}
