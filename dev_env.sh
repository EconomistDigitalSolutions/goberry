#!/bin/sh

env_script="$(basename "$0")"

cat << EOM
export SERVICE_NAME="goberry"
export SERVICE_REGISTRATION="false"

# Run this command to configure your environment:
# eval "\$(./$env_script)"
EOM
