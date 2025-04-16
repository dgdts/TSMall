#!/bin/bash
RUN_NAME="ts_mall_service"
BASE_DIR=$(pwd)
OUTPUT=${BASE_DIR}/output
TARGET_BIN=${OUTPUT}/bin
TARGET_CONF=${OUTPUT}/configs
TARGET_DATA=${OUTPUT}/data

DATE=$(date +"%Y%m%d%H%M%S")

rm -rf ${OUTPUT}
mkdir -p ${TARGET_BIN}
mkdir -p ${TARGET_CONF}
mkdir -p ${TARGET_DATA}

# 根据当前系统自动切换编译方式
if [[ `uname` == 'Linux' ]]; then   
    GOOS=linux GOARCH=amd64 go build -o ${TARGET_BIN}/${RUN_NAME} ${BASE_DIR}/cmd/api/main.go
else
    go build -o ${TARGET_BIN}/${RUN_NAME} ${BASE_DIR}/cmd/api/main.go    
fi

cp -rf ${BASE_DIR}/configs/* ${TARGET_CONF}
echo ${RUN_NAME}" Build Success "${DATE}