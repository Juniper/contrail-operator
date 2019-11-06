package cr
	
import(
	"atom/atom/contrail/operator/pkg/apis/contrail/v1alpha1"

	"github.com/ghodss/yaml"
)

var yamlDataVrouter= `
apiVersion: contrail.juniper.net/v1alpha1
kind: Vrouter
metadata:
  name: cluster-1
`

func GetVrouterCr() *v1alpha1.Vrouter{
	cr := v1alpha1.Vrouter{}
	err := yaml.Unmarshal([]byte(yamlDataVrouter), &cr)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDataVrouter))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &cr)
	if err != nil {
		panic(err)
	}
	return &cr
}
	
