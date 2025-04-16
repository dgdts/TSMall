#!/bin/bash

# check if the number of arguments is correct
if [ $# -ne 2 ]; then
    echo "Usage: $0 <module_name> <service_name>"
    echo "Example: $0 MyTest my_test_service"
    exit 1
fi

MODULE_NAME=$1
SERVICE_NAME=$2

hz new --mod=${MODULE_NAME} \
    --service=${SERVICE_NAME} \
    --idl=./idl/goods.proto \
    --customize_layout=./template/layout.yaml \
    --customize_package=./template/package.yaml \
    --handler_dir=biz/handler \
    --router_dir=biz/router \
    --model_dir=hertz_gen \
    -force

go mod tidy
