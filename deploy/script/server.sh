#!/bin/bash

start() {
    sudo ../../docker-compose up -d
    bash user-rpc-test.sh start
}

stop() {
    bash user-rpc-test.sh stop
    sudo ../../docker-compose down
}

restart() {
    stop
    start
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
