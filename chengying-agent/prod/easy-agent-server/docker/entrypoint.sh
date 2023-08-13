#!/bin/sh

if [ -n "$DTSTACK_MYSQL_HOST" ];then
    sed -i "s|^\(\s*mysqldb\.host\):.*|\1: $DTSTACK_MYSQL_HOST|" example-config.yml
fi

if [ -n "$DTSTACK_MYSQL_PORT" ];then
    sed -i "s|^\(\s*mysqldb\.port\):.*|\1: $DTSTACK_MYSQL_PORT|" example-config.yml
fi

if [ -n "$DTSTACK_MYSQL_USER_NAME" ];then
    sed -i "s|^\(\s*mysqldb\.user\):.*|\1: $DTSTACK_MYSQL_USER_NAME|" example-config.yml
fi

if [ -n "$DTSTACK_MYSQL_USER_PWD" ];then
    sed -i "s|^\(\s*mysqldb\.password\):.*|\1: $DTSTACK_MYSQL_USER_PWD|" example-config.yml
fi

exec ./easy-agent-server -c example-config.yml --debug
