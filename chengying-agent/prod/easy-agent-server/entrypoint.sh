#!/bin/sh

if [ -n "$MYSQL_USERNAME" ];then
    sed -i "s|^\(\s*mysqldb\.user\):.*|\1: $MYSQL_USERNAME|" example-config.yml
fi

if [ -n "$MYSQL_PASSWORD" ];then
    sed -i "s|^\(\s*mysqldb\.password\):.*|\1: $MYSQL_PASSWORD|" example-config.yml
fi

if [ -n "$MATRIX_HOST" ];then
    sed -i "s|^\(\s*publish\.http\.host\):.*|\1: $MATRIX_HOST|" example-config.yml
fi

exec ./easy-agent-server -c example-config.yml --debug
