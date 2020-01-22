# Development

## Install Golang 1.13

```
snap install go --classic --channel 1.13/stable
```

## Atom repository can be everywhere because go modules are used (outside GOPATH)

Keep in mind that for Goland you have to have atom folder inside another atom folder.
Open the parent atom folder in Goland.

## Add new API and controller

```
cd github.com/Juniper/contrail-operator
docker run -it -v $(pwd):/contrail-operator -v /var/run/docker.sock:/var/run/docker.sock kaweue/operator-sdk:v.13-go-1.13 bash
$ cd /contrail-operator
$ operator-sdk add api --api-version=contrail.juniper.net/v1alpha1 --kind=Memcached
$ operator-sdk add controller --api-version=contrail.juniper.net/v1alpha1 --kind=Memcached
$ exit

sudo chown -R `id -u`:`id -g` ./**/*
```


## Generate k8s files

```
cd github.com/Juniper/contrail-operator
docker run --rm -it -v $(pwd):/contrail-operator kaweue/operator-sdk:v.13-go-1.13  bash -c "cd /contrail-operator;operator-sdk generate k8s"
docker run --rm -it -v $(pwd):/contrail-operator kaweue/operator-sdk:v.13-go-1.13  bash -c "cd /contrail-operator;operator-sdk generate openapi"
```

## Troubleshooting

* Problem: unsupported type invalid type for invalid type
  Solution: export GOROOT
* Problem: on running operator container `/usr/local/bin/entrypoint: Permission denied`
  Solution:
  ```
  sudo chown -R `id -u`:`id -g` build
  chmod -R 755 build/bin build/_output
  <rebuild operator>
  ```


## Updating Contrail operator
```
cd github.com/Juniper/contrail-operator
docker run -it -v $(pwd):/contrail-operator -v /var/run/docker.sock:/var/run/docker.sock kaweue/operator-sdk:v.13-go-1.13 bash
$ cd /contrail-operator; operator-sdk build contrail-operator
```


## Building CAVA image

* bazel clean
* make cava
   
