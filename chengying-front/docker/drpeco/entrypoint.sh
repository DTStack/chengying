#!/bin/sh

#CONFIG_FILE=/etc/nginx/conf.d/default.conf
CONFIG_FILE=/usr/share/nginx/html/easymanager/config.js

#sed -i s/matrix/${MATRIX_IP}/ ${CONFIG_FILE}
#sed -i s/grafana/${GRAFANA_IP}/ ${CONFIG_FILE}

if [ -n "$GRAFANA_PORT" ];then
    sed -i "s|^\(\s*GRAFANA_PORT\):.*|\1: $GRAFANA_PORT,|" $CONFIG_FILE
fi

nginx -g "daemon off;"

tail -f /dev/null

