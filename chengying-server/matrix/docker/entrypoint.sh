#!/bin/sh

if [ -z "$STATIC_URL" ];then
        echo "STATIC_URL is empty!"
        exit 1
fi
sed -i "s|^\(\s*api\.static-url\):.*|\1: $STATIC_URL|" example-config.yml

if [ -n "$AGENT_HOST" ];then
        sed -i "s|^\(\s*agent\.host\):.*|\1: $AGENT_HOST|" example-config.yml
fi

if [ -n "$INSTALL_PATH" ];then
        sed -i "s|^\(\s*deploy\.install-path\):.*|\1: $INSTALL_PATH|" example-config.yml
fi

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

if [ -n "$VERIFY_CODE" ];then
        sed -i "s|^\(\s*verify\.code\):.*|\1: $VERIFY_CODE|" example-config.yml
fi

if [ -n "$FLUSH_HOST" ];then
        sed -i "s|^\(\s*autoflush\.etc-host\):.*|\1: $FLUSH_HOST|" example-config.yml
fi

exec ./matrix -c example-config.yml --debug
