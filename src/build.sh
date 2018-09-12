#!/bin/bash
set -e
################ var ################
WORKSPACE=$(cd `dirname $0` && pwd -P)
DATE=$(date "+%s")
#GOPATH=/tmp/go-build${DATE}

# 对应package实际包路径，可以参考glide.yaml package(不包含最后一个目录)，需要RD手动修改
PACKAGE_PATH=Users/admin/Proj/Go/

# 二进制，需要RD手动修改
MODULE_NAME=well_framework

MODULE_PATH=${GOPATH}/src/${PACKAGE_PATH}

################ env ################
export GOPATH
export GOROOT=/usr/local/go
export PATH=${GOROOT}/bin:$GOPATH/bin:${PATH}:$GOBIN

# build
function build() {
    rm -rf $MODULE_PATH/$MODULE_NAME &> /dev/null
    mkdir -p $GOPATH/bin
    mkdir -p $MODULE_PATH
    ln -sf $WORKSPACE ${MODULE_PATH}/${MODULE_NAME}
   # curl https://git.xiaojukeji.com/lego/tools/raw/master/glide/get | sh
    cd $MODULE_PATH/$MODULE_NAME
    echo "Building..."
   # glide cc
   # glide install
    go build -o $MODULE_NAME $PACKAGE_PATH/$MODULE_NAME
    if [[ $? != 0 ]]; then
        echo -e "Build fail !"
        exit 1
    fi
    echo -e "Build success !"
}

# output
function make_output() {
    echo -e "Make output..."
    cd $MODULE_PATH/$MODULE_NAME
    rm -rf output/ &> /dev/null
    mkdir -p output
    cp $MODULE_NAME control.sh output/
    #cp Dockerfile output/
    cp -r conf output/
    #cp...
    echo -e "Make output success !"
}
build
make_output
