## Build and push contrail-operator-ringcontroller

Unfortunately this container cannot be run on Mac OS, since scripts are not ready for cross "compilation". Please use following command from the contrail-operator top directory to build and push image.

    docker run --net host -w /workspace -v /var/run/docker.sock:/var/run/docker.sock -v $(pwd):/workspace gcr.io/cloud-builders/bazel run //ringcontroller:contrail-operator-ringcontroller-push-local