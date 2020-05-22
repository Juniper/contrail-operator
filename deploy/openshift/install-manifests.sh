#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

get_parameters() {
	while [ ! $# -eq 0 ]
	do
		case "$1" in
			--dir)
				DIRECTORY=$2
				;;
			--operator-dir)
				OPERATOR_DIR=$2
				;;
            --config)
                CONFIG=$2
                ;;
		esac
		shift
	done

	if [ -z $OPERATOR_DIR ]
	then
		usage
		exit 1
	fi

        if [ -z $CONFIG ]
        then
                 CONFIG="${SCRIPT_DIR}/config"
        fi
}

usage ()
{
    echo "usage: $0 [--dir <output-dir>][--config <config file>] --operator-dir <contrail-operator project directory>"
}

copy_manifests() {
	mkdir -p "$DIRECTORY/openshift"
	cp -v ${SCRIPT_DIR}/openshift/* "${DIRECTORY}/openshift"
	mkdir -p "$DIRECTORY/manifests"
	cp -v ${SCRIPT_DIR}/manifests/* "${DIRECTORY}/manifests"
        echo "[INFO] Manifests have been copied to ${DIRECTORY}"
}

copy_and_rename_crds() {
	for f in ${OPERATOR_DIR}/deploy/crds/*_crd.yaml;
	do
		f_filename=$(basename $f)
		cp -v ${f} "${DIRECTORY}/manifests/0000000-contrail-07-${f_filename}"
	done
        echo '[INFO] Manifests CRDs have been properly renamed'
}

read_config() {
    if [ ! -f "$CONFIG" ]; then
        echo "$CONFIG file does not exist. Please provide config file."
        exit 1
    fi

    OPERATOR_IMAGE=$(grep "CONTRAIL_OPERATOR_IMAGE" $CONFIG)
    if [ $? -ne 0 ]; then
        echo "Couldn't find CONTRAIL_OPERATOR_IMAGE parameter. Exiting..."
        exit 1
    fi
    DOCKER_CONFIG=$(grep "DOCKER_CONFIG" $CONFIG)
    if [ $? -ne 0 ]; then
        echo "Couldn't find DOCKER_CONFIG parameter. Exiting..."
        exit 1
    fi
    OPERATOR_IMAGE="${OPERATOR_IMAGE##CONTRAIL_OPERATOR_IMAGE=}"
    DOCKER_CONFIG="${DOCKER_CONFIG##DOCKER_CONFIG=}"
    echo '[INFO] Config properly consumed'
}

apply_config() {
    sed -i 's|<OPERATOR_IMAGE>|'$OPERATOR_IMAGE'|g' ${DIRECTORY}/manifests/0000000-contrail-08-operator.yaml
    sed -i 's|<DOCKER_CONFIG>|'$DOCKER_CONFIG'|g' ${DIRECTORY}/manifests/0000000-contrail-02-registry-secret.yaml
    echo '[INFO] Set proper parameters from config in manifests'
}

DIRECTORY=$(pwd)
get_parameters "$@"
echo '[INFO] Starting setup script'
read_config
copy_manifests
copy_and_rename_crds
apply_config
