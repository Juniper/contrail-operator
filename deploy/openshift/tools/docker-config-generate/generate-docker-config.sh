#!/usr/bin/env bash

if [[ "$OSTYPE" == "darwin"* ]]; then
        BASE64_COMMAND='base64'
else
        BASE64_COMMAND='base64 -w 0'
fi

get_parameters() {
    read -p "Docker registry: "  REGISTRY
    read -p "Docker username: "  USERNAME
    read -s -p "Docker secret: "  SECRET
}

generate_base64() {
    AUTH=$(echo $USERNAME:$SECRET | $BASE64_COMMAND)
    RENDERED_JSON=$(echo {\"auths\":{\"${REGISTRY}\":{\"username\": \"${USERNAME}\",\"password\":\"${SECRET}\",\"auth\":\"${AUTH}\"}}})
    DOCKER_CONFIG=$(echo $RENDERED_JSON | $BASE64_COMMAND)
}

get_parameters
generate_base64

printf "[INFO] Paste this configuration into install manifests config to provide authentication for contrail-operator.\n"
printf "\nDOCKER_CONFIG=$DOCKER_CONFIG\n"
