#!/usr/bin/env bash
# Created by paincompiler on 2017/07/14
# © 2017 ZAOSHU All RIGHTS RESERVED.

##########settings##########
set -o errexit
set -o errtrace
#set -o pipefail
set -o xtrace
##########settings##########

function finish {
    # cleanup code here
    echo "cleaning up"
}

function set_env {
    export GOOS=linux
    export GOARCH=amd64
}

trap finish EXIT

source $(dirname $0)/utils.sh
set_basic_env

set_env

SCRIPT_DIR="$(dirname $0)"
PROJECT_ROOT="${SCRIPT_DIR}/../"
OUTPUT_FILENAME="${PROJECT_NAME}-linux-amd64-${GIT_COMMIT}"

go build -v -o "${OUTPUT_FILENAME}" "${PROJECT_ROOT}"
rm -rf "${OUTPUT_FILENAME}.txz"
tar -cJf "${OUTPUT_FILENAME}.txz" "${OUTPUT_FILENAME}" config
