# Development

## Install Golang 1.12. 1.13 does not work

```
snap install go --classic --channel 1.12/stable
```

## Atom repository can be everywhere because go modules are used (outside GOPATH)

Keep in mind that for Goland you have to have atom folder inside another atom foleder.
Open the parent atom folder in Goland.

## Generate k8s files

```
cd github.com/Juniper/contrail/operator
# everytime you want to generate files you have to copy go.mod files
# Do not commit those files though and remove them before using Goland
cp ../../go.* .
docker run --rm -it -v $(pwd):/project kaweue/operator-sdk:v.10-go-1.12 bash -c "cd /project;operator-sdk generate k8s"
docker run --rm -it -v $(pwd):/project kaweue/operator-sdk:v.10-go-1.12 bash -c "cd /project;operator-sdk generate openapi"
rm go.*
```

## Troubleshooting

* Problem: `github.com/Juniper/base/go/server/testing/client imports
github.com/Juniper/base/go/server/testing/testserver/testservice: malformed module path "github.com/Juniper/base/go/server/testing/testserver/testservice": missing dot in first path element`
  Solution: use golang 1.12

* Problem: missing imports in Goland
  Solution: remove go.mod and go.sum from atom/contrail/operator

* Problem: unsupported type invalid type for invalid type
  Solution: export GOROOT


## Updating Contrail operator

* make bazel-sync
* make docker-push cmd/manager/contrail_operator_image_base
* docker tag bazel/cmd/manager:contrail_operator_image_base localhost:5000/contrail_operator:latest
* docker push localhost:5000/contrail_operator:latest


## Building CAVA image

* bazel clean
* make cava
   