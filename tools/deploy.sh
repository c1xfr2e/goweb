#!/usr/bin/env bash
# Created by paincompiler on 2017/07/13
# © 2017 ZAOSHU All RIGHTS RESERVED.

##########settings##########
set -o errexit
set -o errtrace
set -o pipefail
set -o xtrace
##########settings##########

function finish {
    # cleanup code here
    echo "cleaning up"
}

trap finish EXIT

# script here

source $(dirname $0)/utils.sh

set_basic_env

if [[ "${PROJECT_ENV}" = "production" ]];then
    curl -v http://deployer.zahlgroup.com:9001/deploy/${PROJECT_NAME}/${PROJECT_ENV}/${GIT_COMMIT}
elif [[ "${PROJECT_ENV}" = "debug" ]];then
    curl -v http://deployer.zahlgroup.com:9001/deploy/${PROJECT_NAME}/${PROJECT_ENV}/${GIT_COMMIT}
elif [[ "${PROJECT_ENV}" == "dev" ]];then
    curl -v http://deployer.zahlgroup.com:9001/deploy/${PROJECT_NAME}/${PROJECT_ENV}/${GIT_COMMIT}
else
    echo "invalid branch name ${GIT_BRANCH_BASENAME}"
    exit 1
fi
