# DOCKER_CONFIG generate

## What is it?
This script generates docker config for Contrail Operator install manifests for Openshift encoded in base64 format.

## Encoded format

DOCKER_CONFIG is encoded JSON data in this format:
`{"auths":{"<registry>":{"username":"<username>","password":"<password>",auth":"<base64 encoded username:password>"}}}`

## Usage

Run script and interactively pass necessary data. Afterwards, copy output and paste into config file for install-manifests script.
