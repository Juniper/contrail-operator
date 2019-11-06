package cr
	
import(
	"atom/atom/contrail/operator/pkg/apis/contrail/v1alpha1"
	
	"github.com/ghodss/yaml"
)

var yamlDataCassandra= `
apiVersion: contrail.juniper.net/v1alpha1
kind: Cassandra
metadata:
  name: example-cassandra
  labels:
    contrail_manager: cassandra
`

func GetCassandraCr() *v1alpha1.Cassandra{
	cr := v1alpha1.Cassandra{}
	err := yaml.Unmarshal([]byte(yamlDataCassandra), &cr)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(yamlDataCassandra))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &cr)
	if err != nil {
		panic(err)
	}
	return &cr
}
	
