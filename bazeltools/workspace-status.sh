#!/bin/bash
set -eu

# Execute workspace-status.py from the same directory as this bash script
#"${0%.sh}.py"

echo BUILD_SCM_REVISION rev
echo BUILD_SCM_BRANCH branch
#echo BUILD_SCM_REVISION $(git rev-parse --short HEAD)
#echo BUILD_SCM_BRANCH $(git rev-parse --abbrev-ref HEAD)
