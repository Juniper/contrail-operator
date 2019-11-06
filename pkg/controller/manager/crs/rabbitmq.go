package cr
	
import(
	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"

	"github.com/ghodss/yaml"
)

var yamlDataRabbitmq= `
apiVersion: contrail.juniper.net/v1alpha1
kind: Rabbitmq
metadata:
  name: example-rabbitmq
spec:
  # Add fields here
  size: 3
`

func GetRabbitmqCr() *v1alpha1.Rabbitmq{
	cr := v1alpha1.Rabbitmq{}
	err := yaml.Unmarshal([]byte(yamlDataRabbitmq), &cr)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDataRabbitmq))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &cr)
	if err != nil {
		panic(err)
	}
	return &cr
}
	
