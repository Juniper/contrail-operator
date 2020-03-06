docker run --rm --mount type=bind,source=${PWD}/config/provision.yaml,target=/provision.yaml  michaelhenkel/contrail-provisioner:latest /contrail-provisioner -file /provision.yaml
