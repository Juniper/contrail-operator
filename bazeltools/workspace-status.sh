#!/bin/bash
set -eu

# Execute workspace-status.py from the same directory as this bash script
"${0%.sh}.py"
