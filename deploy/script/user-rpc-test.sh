#!/bin/bash

# 停止并且删除镜像,之后重新构建镜像
server_name="user"
server_type="rpc"
container_name="easy-im-user-rpc-test"
image_name="${server_name}-${server_type}-test"

# 1. 没有容器并且需要重新建立容器
start() {
    cd ../..
    make -f ./deploy/mk/user-rpc.mk release-test
    sudo docker run -p 10001:8080 --name=${container_name} --network im-chat_im-chat --link etcd:etcd --link mysql:mysql --link redis:redis -d ${image_name}
}

stop() {
    sudo docker stop ${container_name}
    sudo docker rm ${container_name}
    sudo docker rmi ${image_name}
}

restart() {
    sudo docker stop ${container_name}
    sudo docker rm ${container_name}
    sudo docker rmi ${image_name}
    cd ../mk && make -f user-rpc.mk release_test
    sudo docker run -p 10001:8080 --name=${container_name} --network im-chat_im-chat --link etcd:etcd --link mysql:mysql --link redis:redis -d ${image_name}
}

action="$1"

case "$action" in
start)
    start
    ;;
stop)
    stop
    ;;
restart)
    restart
    ;;
*)
    echo "Usage: $0 {start|stop|restart}"
    exit 1
    ;;
esac
