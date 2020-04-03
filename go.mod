module github.com/Juniper/contrail-operator

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/DATA-DOG/go-sqlmock v1.3.3
	github.com/ExpansiveWorlds/instrumentedsql v0.0.0-20171218214018-45abb4b1947d
	github.com/Juniper/asf v0.0.0-20200330081506-d317915f0ee1
	github.com/Juniper/contrail v0.0.0-20200330181744-e78e7561c8fd
	github.com/Juniper/contrail-go-api v1.1.0
	github.com/NYTimes/gziphandler v1.1.1
	github.com/PuerkitoBio/purell v1.1.1
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578
	github.com/apparentlymart/go-cidr v1.0.1
	github.com/bazelbuild/bazel-gazelle v0.20.0 // indirect
	github.com/bazelbuild/rules_docker v0.14.1
	github.com/bazelbuild/rules_go v0.0.0-20190719190356-6dae44dc5cab
	github.com/bitly/go-hostpool v0.0.0-20171023180738-a3a6125de932
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/boltdb/bolt v1.3.1
	github.com/cockroachdb/apd v1.1.0
	github.com/codegangsta/cli v1.20.0
	github.com/containerd/containerd v1.3.0
	github.com/containerd/fifo v0.0.0-20191213151349-ff969a566b00 // indirect
	github.com/containerd/ttrpc v1.0.0
	github.com/containerd/typeurl v0.0.0-20190228175220-2a93cfde8c20
	github.com/coreos/etcd v3.3.15+incompatible
	github.com/databus23/keystone v0.0.0-20180111110916-350fd0e663cd
	github.com/davecgh/go-spew v1.1.1
	github.com/denisenkom/go-mssqldb v0.0.0-20190515213511-eb9f6a1743f3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.4.2-0.20190924003213-a8608b5b67c7
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/ericlagergren/decimal v0.0.0-20191206042408-88212e6cfca9
	github.com/fatih/structs v1.1.0
	github.com/flosch/pongo2 v0.0.0-20190707114632-bbf5a6c351f4
	github.com/friendsofgo/errors v0.9.2
	github.com/fsnotify/fsnotify v1.4.7
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-openapi/jsonpointer v0.19.3
	github.com/go-openapi/jsonreference v0.19.3
	github.com/go-openapi/spec v0.19.4
	github.com/go-openapi/swag v0.19.5
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gocql/gocql v0.0.0-20190301043612-f6df8288f9b4
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/gogo/googleapis v1.3.2 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/dep v0.5.4
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/mock v1.3.1
	github.com/golang/protobuf v1.3.3
	github.com/google/go-cmp v0.4.0
	github.com/google/go-containerregistry v0.0.0-20200320200342-35f57d7d4930
	github.com/gorilla/websocket v1.4.1
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-checkpoint v0.5.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-plugin v1.0.1-0.20190610192547-a1bc61569a26
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/hil v0.0.0-20190212112733-ab17b08d6590
	github.com/hashicorp/terraform v0.12.24
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733
	github.com/jackc/pgx v3.2.0+incompatible
	github.com/jmank88/nuts v0.4.0 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/juju/errors v0.0.0-20181118221551-089d3ea4e4d5
	github.com/juju/testing v0.0.0-20180920084828-472a3e8b2073
	github.com/kat-co/vala v0.0.0-20170210184112-42e1d8b61f12
	github.com/kr/pretty v0.1.0
	github.com/kr/pty v1.1.5
	github.com/kr/text v0.1.0
	github.com/kylelemons/godebug v1.1.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/echo/v4 v4.1.16 // indirect
	github.com/labstack/gommon v0.3.0
	github.com/lib/pq v1.2.0
	github.com/magiconair/properties v1.8.1
	github.com/mailru/easyjson v0.7.0
	github.com/mattn/go-colorable v0.1.6
	github.com/mattn/go-isatty v0.0.12
	github.com/mattn/go-shellwords v1.0.6
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/mitchellh/cli v1.0.0
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db
	github.com/mitchellh/iochan v1.0.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mitchellh/panicwrap v1.0.0
	github.com/mitchellh/prefixedio v0.0.0-20190213213902-5733675afd51
	github.com/mitchellh/reflectwalk v1.0.0
	github.com/nightlyone/lockfile v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/opencontainers/runtime-spec v1.0.0
	github.com/operator-framework/operator-sdk v0.14.1
	github.com/pborman/uuid v1.2.0
	github.com/pelletier/go-toml v1.6.0
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.11.0 // indirect
	github.com/pmylund/go-cache v2.1.0+incompatible
	github.com/pseudomuto/protoc-gen-doc v1.3.1
	github.com/satori/go.uuid v1.2.0
	github.com/sdboyer/constext v0.0.0-20170321163424-836a14457353 // indirect
	github.com/shopspring/decimal v0.0.0-20191009025716-f1972eb1d1f5
	github.com/shurcooL/sanitized_anchor_name v0.0.0-20151028001915-10ef21a441db // indirect
	github.com/sigma/go-inotify v0.0.0-20181102212354-c87b6cf5033d // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.4.0
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/stretchr/testify v1.5.1
	github.com/ulikunitz/xz v0.5.5
	github.com/valyala/bytebufferpool v1.0.0
	github.com/volatiletech/inflect v0.0.0-20170731032912-e7201282ae8d
	github.com/volatiletech/null v8.0.0+incompatible
	github.com/volatiletech/sqlboiler v3.5.0+incompatible
	github.com/yudai/gotty v2.0.0-alpha.3+incompatible
	github.com/yudai/hcl v0.0.0-20151013225006-5fa2393b3552
	go.starlark.net v0.0.0-20200326215636-e8819e807894 // indirect
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae
	golang.org/x/text v0.3.2
	golang.org/x/tools v0.0.0-20200309202150-20ab64c0d93f
	google.golang.org/genproto v0.0.0-20200218151345-dad8c97a84f5
	google.golang.org/grpc v1.27.1
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2
	gopkg.in/inf.v0 v0.9.1
	gopkg.in/volatiletech/null.v6 v6.0.0-20170828023728-0bef4e07ae1b
	gopkg.in/yaml.v2 v2.2.8
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.17.4
	k8s.io/apiextensions-apiserver v0.0.0
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20190918143330-0270cf2f1c1d
	k8s.io/kubernetes v1.16.2
	sigs.k8s.io/controller-runtime v0.4.0
	sigs.k8s.io/kind v0.7.0
)

// Pinned to kubernetes-1.16.2
replace (
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309 // Required by Helm
	k8s.io/api => k8s.io/api v0.0.0-20191016110408-35e52d86657a
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191016113550-5357c4baaf65
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20191016112112-5190913f932d
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20191016114015-74ad18325ed5
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20191016115326-20453efc2458
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20191016115129-c07a134afb42
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20191004115455-8e001e5d1894
	k8s.io/component-base => k8s.io/component-base v0.0.0-20191016111319-039242c015a9
	k8s.io/cri-api => k8s.io/cri-api v0.0.0-20190828162817-608eb1dad4ac
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20191016115521-756ffa5af0bd
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20191016112429-9587704a8ad4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20191016114939-2b2b218dc1df
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20191016114407-2e83b6f20229
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20191016114748-65049c67a58b
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191016120415-2ed914427d51
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20191016114556-7841ed97f1b2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20191016115753-cf0698c3a16b
	k8s.io/metrics => k8s.io/metrics v0.0.0-20191016113814-3b1a734dba6e
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20191016112829-06bb3c9d77c9
)
