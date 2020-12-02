module github.com/Juniper/contrail-operator

go 1.13

require (
	github.com/Azure/go-autorest/autorest v0.9.3 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.1 // indirect
	github.com/Juniper/contrail-go-api v1.1.1-0.20200414151206-7bb4264f1da4
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/emicklei/go-restful v2.11.1+incompatible // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-bindata/go-bindata v3.1.2+incompatible // indirect
	github.com/go-openapi/spec v0.19.4
	github.com/go-pg/pg/v10 v10.0.0-beta.2
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golangci/gofmt v0.0.0-20190930125516-244bba706f1a
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/helm/helm-2to3 v0.5.1 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/kylelemons/godebug v1.1.0
	github.com/martinlindhe/base36 v1.0.0 // indirect
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/openshift/api v3.9.1-0.20190924102528-32369d4db2ad+incompatible
	github.com/operator-framework/operator-lifecycle-manager v0.0.0-20200321030439-57b580e57e88 // indirect
	github.com/operator-framework/operator-sdk v0.18.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.5.1
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/tools v0.0.0-20200913032122-97363e29fc9b
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/ini.v1 v1.51.0
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.2
	k8s.io/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/kind v0.7.1-0.20200303021537-981bd80d3802
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.18.2 // Required by prometheus-operator
)
