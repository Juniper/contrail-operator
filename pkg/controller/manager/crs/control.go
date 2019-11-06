package cr
	
import(
	"atom/atom/contrail/operator/pkg/apis/contrail/v1alpha1"

	"github.com/ghodss/yaml"
)

var yamlDataControl= `
apiVersion: contrail.juniper.net/v1alpha1
kind: Control
metadata:
  name: cluster-1
`

func GetControlCr() *v1alpha1.Control{
	cr := v1alpha1.Control{}
	err := yaml.Unmarshal([]byte(yamlDataControl), &cr)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDataControl))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &cr)
	if err != nil {
		panic(err)
	}
	return &cr
}
	
