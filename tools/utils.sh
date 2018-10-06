#!/usr/bin/env bash
# Created by paincompiler on 2017/07/14
# © 2017 ZAOSHU All RIGHTS RESERVED.

##########settings##########
set -o errexit
set -o errtrace
set -o xtrace
##########settings##########

function set_basic_env {
    export GIT_BRANCH_BASENAME=`git rev-parse --abbrev-ref HEAD`
    export GIT_COMMIT="`git log | head -n 1 | awk '{print $2}'`"

    if [[ "${GIT_BRANCH_BASENAME}" = "master" ]];then
        export PROJECT_ENV="production"
    elif [[ "${GIT_BRANCH_BASENAME}" = "test"* ]];then
        export PROJECT_ENV="debug"
    elif [[ "${GIT_BRANCH_BASENAME}" = "release/"* ]];then
        export PROJECT_ENV="debug"
    elif [[ "${GIT_BRANCH_BASENAME}" = "develop" ]];then
        export PROJECT_ENV="dev"
    elif [[ "${GIT_BRANCH_BASENAME}" = "support/"* ]];then
        export PROJECT_ENV="dev"
    elif [[ "${GIT_BRANCH_BASENAME}" = "feature/"* ]];then
        export PROJECT_ENV="dev"
    elif [[ "${GIT_BRANCH_BASENAME}" = "hotfix/"* ]];then
        export PROJECT_ENV="dev"
    else
        echo "invalid branch name ${GIT_BRANCH_BASENAME}"
        exit 1
    fi

    export PROJECT_NAME=`basename $(git rev-parse --show-toplevel)`

    export SCRIPT_DIR=$(dirname $0)
    export PROJECT_ROOT="${SCRIPT_DIR}/../"
}
