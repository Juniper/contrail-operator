#!/bin/bash
set -eu

# Execute workspace-status.py from the same directory as this bash script
#"${0%.sh}.py"

LOCAL_REGISTRY=${LOCAL_REGISTRY:-"localhost:5000"}

if [ ! -d ".git" ]; then
    rev=${SHORT_SHA}
    branch=${BRANCH_NAME}
else
    rev=`/usr/bin/git rev-parse --short HEAD`
    branch=`/usr/bin/git rev-parse --abbrev-ref HEAD`
fi
echo BUILD_SCM_REVISION ${rev}
echo BUILD_SCM_BRANCH ${branch}
echo LOCAL_REGISTRY ${LOCAL_REGISTRY}
