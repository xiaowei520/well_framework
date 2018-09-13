#!/bin/bash
#############################################
## main
## 以托管方式, 启动服务
## control.sh脚本, 必须实现start方法
#############################################
workspace=$(cd $(dirname $0) && pwd -P)
cd $workspace
module=well_framework
app=$module

action=$1
case $action in
    "start" )
        # 获取机房配置
        cluster=test # `cat .deploy/service.cluster.txt`
        clusterconf=conf/global/conf.toml.$cluster
        conf=conf/global/conf.toml

        echo "clusterconf=$clusterconf"
        echo "conf=$conf"

        if [ -f "$clusterconf" ];then
            rm -f $conf
            cp -f $clusterconf $conf
        	else
            echo -n "invalid cluster config file"
            exit 1
        fi
        echo $app
        # 启动服务, 以前台方式启动, 否则无法托管
        exec "./$app" > ./console.log 2>&1
        exit 0
        ;;
    * )
        # 非法命令, 已非0码退出
        echo "unknown command"
        exit 1
        ;;
esac
