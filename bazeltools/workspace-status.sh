#!/bin/bash
set -eu

# Execute workspace-status.py from the same directory as this bash script
#"${0%.sh}.py"

LOCAL_REGISTRY=${LOCAL_REGISTRY:-"localhost:5000"}
BRANCH_NAME=${BRANCH_NAME:-"master"}

if [ ! -d ".git" ]; then
    rev=${SHORT_SHA}
else
    rev=`/usr/bin/git rev-parse --short HEAD`
fi
echo BUILD_SCM_REVISION ${rev}
echo BUILD_SCM_BRANCH ${BRANCH_NAME}
echo LOCAL_REGISTRY ${LOCAL_REGISTRY}
