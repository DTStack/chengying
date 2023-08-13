#!/bin/sh

mv /matrix/tmp-easyagent/* /matrix/easyagent

if [ -z "$STATIC_URL" ];then
        echo "STATIC_URL is empty!"
        exit 1
fi
sed -i "s|^\(\s*api\.static-url\):.*|\1: $STATIC_URL|" example-config.yml

if [ -n "$AGENT_HOST" ];then
        sed -i "s|^\(\s*agent\.host\):.*|\1: $AGENT_HOST|" example-config.yml
fi

exec ./matrix -c example-config.yml --debug
