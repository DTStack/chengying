#!/bin/sh

if [ -z "$SERVER" ];then
        echo "SERVER is empty!"
        exit 1
fi

sed -i "s|host:.*|host: ${SERVER%:*}|" monitor.yml

exec ./easymonitor --config monitor.yml
