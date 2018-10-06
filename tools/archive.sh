#!/usr/bin/env bash
# Created by paincompiler on 2017/07/14
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

function set_env {
    export PROJECT_PACKAGE_FILENAME="${PROJECT_NAME}-linux-amd64-${GIT_COMMIT}.txz"
    export UPLOAD_URL="https://dl.zaoshu.io/f/beijing/zaoshulibs?filename=${PROJECT_NAME}/${PROJECT_PACKAGE_FILENAME}"
}

trap finish EXIT

# script here

source $(dirname $0)/utils.sh
set_basic_env

set_env

curl -s -XPUT --user "${UPLOAD_AUTHENTICATION}" "${UPLOAD_URL}" -F "file=@${PROJECT_PACKAGE_FILENAME}"
if [ "$?" -ne "0" ];then
  echo "UPLOAD FILE FAILED"
  exit 1
else
  echo "UPLOAD FILE SUCCESSFULLY"
fi

if [[ "$?" -ne "0" ]];then
	echo "upload stage failed"
	exit 1
fi