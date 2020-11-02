#!/usr/bin/env bash

usage()
{
    echo "usage: $0 [--from-cluster <kubeconfig-path>][--to-cluster <kubeconfig-path>][--resource <resource type>][--name <resource-name>][--from-namespace <namespace>][--to-namespace <namespace>]"
}


get_parameters() {
    while [ ! $# -eq 0 ]
    do
        case "$1" in
            --help)
            usage
            exit 0
            ;;
            --from-cluster)
            FROM_CLUSTER=$2
            ;;
            --to-cluster)
            TO_CLUSTER=$2
            ;;
            --resource)
            RESOURCE=$2
            ;;
            --name)
            NAME=$2
            ;;
            --from-namespace)
            FROM_NS=$2
            ;;
            --to-namespace)
            TO_NS=$2
            ;;
        esac
        shift
    done

    if [ -z $TO_CLUSTER ]; then
        if [ -z $FROM_CLUSTER ]; then
            echo "Provide at least one kubeconfig"
            exit 1
        else
            TO_CLUSTER=$FROM_CLUSTER
        fi
    fi

    if [ -z $FROM_CLUSTER ]; then
        if [ -z $TO_CLUSTER ]; then
            echo "Provide at least one kubeconfig"
            exit 1
        else
            FROM_CLUSTER=$TO_CLUSTER
        fi
    fi

    if [ -z $RESOURCE ]; then
        echo "Please provide resource type"
        exit 1
    fi

    if [ -z $NAME ]; then
        echo "Please provide name of the resource"
        exit 1
    fi

    if [ -z $TO_NS ]; then
        if [ -z $FROM_NS ]; then
            echo "Provide at least one namespace name"
            exit 1
        else
            TO_NS=$FROM_NS
        fi
    fi

    if [ -z $FROM_NS ]; then
        if [ -z $TO_NS ]; then
            echo "Provide at least one namespace name"
            exit 1
        else
            FROM_NS=$TO_NS
        fi
    fi
}

copy_resource() {
    export KUBECONFIG=$FROM_CLUSTER
    ORIGINAL_RESOURCE=$(kubectl get $RESOURCE -n $FROM_NS $NAME -o json)
    SCANNED_RESOURCE=$(echo $ORIGINAL_RESOURCE | jq 'del(.metadata.namespace,.metadata.resourceVersion,.metadata.uid) | .metadata.creationTimestamp=null')
    export KUBECONFIG=$TO_CLUSTER
    echo $SCANNED_RESOURCE | kubectl apply -n $TO_NS -f -
}

get_parameters "$@"
echo '[INFO] Arguments consumed properly'
copy_resource

